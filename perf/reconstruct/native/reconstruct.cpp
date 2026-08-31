// A handcrafted native implementation of the reconstruct scenario: load a
// Sysl model from binary protobuf, normalise it, and render it back to Sysl
// source, emitting the same arr.ai-formatted dict that
// `arrai run vendor/run.arrai model.sysl.pb` prints. It exists to put a
// native-speed floor under the arr.ai evaluator's number for the same
// workload: output is byte-identical to expected.arrai.
//
// This is a faithful port of the pipeline in
// vendor/sysl/pkg/arrai/{sysl,reconstruct}.arrai — the same normalisation,
// ordering, quoting and template rules, implemented generally over the
// protobuf features that pipeline renders — not a shortcut that exploits the
// synthetic model's regularity. Features the arr.ai pipeline itself leaves
// unrendered (endpoint params and tags, events, aliases, views) are dropped
// here for the same reason: the port matches the pipeline, quirks included.
//
// Build:  clang++ -std=c++20 -O2 -o reconstruct reconstruct.cpp
// Run:    ./reconstruct ../model.sysl.pb

#include <algorithm>
#include <cstdint>
#include <cstdio>
#include <cstring>
#include <map>
#include <memory>
#include <set>
#include <string>
#include <string_view>
#include <vector>

using std::string;
using std::string_view;

// ---------------------------------------------------------------------------
// Protobuf wire-format reader (no libprotobuf; unknown fields are skipped).

struct Reader {
    const uint8_t* p;
    const uint8_t* end;

    bool done() const { return p >= end; }

    uint64_t varint() {
        uint64_t x = 0;
        int s = 0;
        while (p < end) {
            uint8_t c = *p++;
            x |= uint64_t(c & 0x7f) << s;
            if (!(c & 0x80)) return x;
            s += 7;
        }
        return x;
    }

    // next reads one field's tag; value access depends on wire type.
    struct Field {
        uint32_t num;
        uint32_t wire;
        uint64_t ival;    // wire 0/1/5
        string_view data; // wire 2
    };

    Field next() {
        uint64_t tag = varint();
        Field f{uint32_t(tag >> 3), uint32_t(tag & 7), 0, {}};
        switch (f.wire) {
        case 0: f.ival = varint(); break;
        case 1: std::memcpy(&f.ival, p, 8); p += 8; break;
        case 5: { uint32_t v; std::memcpy(&v, p, 4); p += 4; f.ival = v; break; }
        case 2: {
            uint64_t n = varint();
            f.data = string_view(reinterpret_cast<const char*>(p), n);
            p += n;
            break;
        }
        default:
            std::fprintf(stderr, "unsupported wire type %u\n", f.wire);
            std::exit(1);
        }
        return f;
    }
};

static Reader sub(string_view s) {
    auto* b = reinterpret_cast<const uint8_t*>(s.data());
    return Reader{b, b + s.size()};
}

// ---------------------------------------------------------------------------
// The slice of the sysl proto model the pipeline consumes.

// Attribute: s | i | n | a{elt}. Annotation values in the rendered output are
// strings or arrays of strings.
struct AttrValue {
    string s;
    std::vector<AttrValue> elems;
    bool isArray = false;
};

using AttrMap = std::vector<std::pair<string, AttrValue>>; // insertion order

struct SrcCtx {
    int64_t startLine = -1;
    bool present = false;
};

struct FieldType {
    // Mirrors parseFieldType results: primitive | tuple | ref | set/sequence.
    enum Kind { None, Primitive, Tuple, Ref, Set, Sequence } kind = None;
    string primitive;                  // lowercased primitive name
    std::vector<string> appName;       // ref
    std::vector<string> typePath;      // ref
    std::unique_ptr<FieldType> inner;  // set/sequence
};

struct FieldDef {
    string name;
    FieldType type;
    bool opt = false;
    AttrMap attrs;
};

struct TypeDef {
    string name;
    bool isTuple = false;
    bool isRelation = false;
    bool isEnum = false;
    std::vector<std::pair<string, int64_t>> enumItems;
    std::vector<FieldDef> fields; // tuple or relation attr_defs
    AttrMap attrs;
    FieldType selfType; // for alias inference (unused in rendering)
};

struct Stmt {
    enum Kind { None, Action, Call, Cond, Loop, LoopN, Foreach, Alt, Group, Ret } kind = None;
    string action;
    std::vector<string> callTarget;
    string callEp;
    string retPayload;
    AttrMap attrs;
};

struct Endpoint {
    string name;
    string longName;
    bool isPubsub = false;
    bool hasRest = false;
    string restPath;
    int restMethod = 0;
    std::vector<Stmt> stmts;
    AttrMap attrs;
};

struct App {
    string key; // map key in Module.apps
    std::vector<string> nameParts;
    string longName;
    AttrMap attrs;
    std::vector<Endpoint> eps;
    std::vector<TypeDef> types;
    SrcCtx src;
    string srcFile;
};

// ---------------------------------------------------------------------------
// Decoding.

static SrcCtx decodeSrc(string_view s, string* file) {
    SrcCtx src;
    src.present = true;
    Reader r = sub(s);
    while (!r.done()) {
        auto f = r.next();
        if (f.num == 1 && file) {
            *file = string(f.data);
        } else if (f.num == 2) { // start: Location{line}
            Reader lr = sub(f.data);
            while (!lr.done()) {
                auto lf = lr.next();
                if (lf.num == 1) src.startLine = int64_t(lf.ival);
            }
        }
    }
    return src;
}

static AttrValue decodeAttr(string_view s) {
    AttrValue v;
    Reader r = sub(s);
    while (!r.done()) {
        auto f = r.next();
        switch (f.num) {
        case 4: v.s = string(f.data); break;
        case 5: v.s = std::to_string(int64_t(f.ival)); break; // i: unused here
        case 7: { // a{elt*}
            v.isArray = true;
            Reader ar = sub(f.data);
            while (!ar.done()) {
                auto af = ar.next();
                if (af.num == 1) v.elems.push_back(decodeAttr(af.data));
            }
            break;
        }
        default: break;
        }
    }
    return v;
}

static void decodeAttrsEntry(string_view s, AttrMap* attrs) {
    Reader r = sub(s);
    string key;
    AttrValue val;
    while (!r.done()) {
        auto f = r.next();
        if (f.num == 1) key = string(f.data);
        else if (f.num == 2) val = decodeAttr(f.data);
    }
    attrs->emplace_back(std::move(key), std::move(val));
}

static const char* primitiveName(uint64_t v) {
    switch (v) {
    case 1: return "empty";
    case 2: return "any";
    case 3: return "bool";
    case 4: return "int";
    case 5: return "float";
    case 6: return "string";
    case 7: return "bytes";
    case 8: return "string_8";
    case 9: return "date";
    case 10: return "datetime";
    case 11: return "xml";
    case 12: return "decimal";
    case 13: return "uuid";
    default: return "";
    }
}

static std::vector<string> decodeAppName(string_view s) {
    std::vector<string> parts;
    Reader r = sub(s);
    while (!r.done()) {
        auto f = r.next();
        if (f.num == 1) parts.emplace_back(f.data);
    }
    return parts;
}

// decodeFieldType mirrors parseFieldType: the salient details of a Type when
// used as a field/alias/return type. appName is the owning app's name parts,
// used when a ref has no explicit appname.
static FieldType decodeFieldType(string_view s, const std::vector<string>& appName);

static FieldType decodeScopedRef(string_view s, const std::vector<string>& appName) {
    FieldType t;
    Reader r = sub(s);
    while (!r.done()) {
        auto f = r.next();
        if (f.num == 2) { // ref: Scope{appname, path}
            t.kind = FieldType::Ref;
            t.appName = appName;
            Reader sr = sub(f.data);
            while (!sr.done()) {
                auto sf = sr.next();
                if (sf.num == 1) t.appName = decodeAppName(sf.data);
                else if (sf.num == 2) t.typePath.emplace_back(sf.data);
            }
        }
    }
    return t;
}

static FieldType decodeFieldType(string_view s, const std::vector<string>& appName) {
    FieldType t;
    Reader r = sub(s);
    while (!r.done()) {
        auto f = r.next();
        switch (f.num) {
        case 1:
            t.kind = FieldType::Primitive;
            t.primitive = primitiveName(f.ival);
            break;
        case 3:
            t.kind = FieldType::Tuple;
            break;
        case 9:
            t = decodeScopedRef(f.data, appName);
            break;
        case 13: {
            FieldType inner = decodeFieldType(f.data, appName);
            t.kind = FieldType::Set;
            t.inner = std::make_unique<FieldType>(std::move(inner));
            break;
        }
        case 15: {
            FieldType inner = decodeFieldType(f.data, appName);
            t.kind = FieldType::Sequence;
            t.inner = std::make_unique<FieldType>(std::move(inner));
            break;
        }
        default: break;
        }
    }
    return t;
}

static TypeDef decodeType(string_view s, const std::vector<string>& appParts) {
    TypeDef t;
    Reader r = sub(s);
    while (!r.done()) {
        auto f = r.next();
        switch (f.num) {
        case 2: { // enum
            t.isEnum = true;
            Reader er = sub(f.data);
            while (!er.done()) {
                auto ef = er.next();
                if (ef.num == 1) { // ItemsEntry
                    Reader ir = sub(ef.data);
                    string k;
                    int64_t v = 0;
                    while (!ir.done()) {
                        auto kf = ir.next();
                        if (kf.num == 1) k = string(kf.data);
                        else if (kf.num == 2) v = int64_t(kf.ival);
                    }
                    t.enumItems.emplace_back(std::move(k), v);
                }
            }
            break;
        }
        case 3:   // tuple
        case 7: { // relation
            (f.num == 3 ? t.isTuple : t.isRelation) = true;
            Reader tr = sub(f.data);
            while (!tr.done()) {
                auto tf = tr.next();
                if (tf.num == 1) { // AttrDefsEntry
                    Reader fr = sub(tf.data);
                    FieldDef fd;
                    while (!fr.done()) {
                        auto ff = fr.next();
                        if (ff.num == 1) {
                            fd.name = string(ff.data);
                        } else if (ff.num == 2) {
                            Reader vr = sub(ff.data);
                            fd.type = decodeFieldType(ff.data, appParts);
                            while (!vr.done()) {
                                auto vf = vr.next();
                                if (vf.num == 8) decodeAttrsEntry(vf.data, &fd.attrs);
                                else if (vf.num == 12) fd.opt = vf.ival != 0;
                            }
                        }
                    }
                    t.fields.push_back(std::move(fd));
                }
            }
            break;
        }
        case 8: decodeAttrsEntry(f.data, &t.attrs); break;
        default: break;
        }
    }
    if (!t.isTuple) t.selfType = decodeFieldType(s, appParts);
    return t;
}

static Stmt decodeStmt(string_view s) {
    Stmt st;
    Reader r = sub(s);
    while (!r.done()) {
        auto f = r.next();
        switch (f.num) {
        case 1: { // action
            st.kind = Stmt::Action;
            Reader ar = sub(f.data);
            while (!ar.done()) {
                auto af = ar.next();
                if (af.num == 2) st.action = string(af.data);
            }
            break;
        }
        case 2: { // call
            st.kind = Stmt::Call;
            Reader cr = sub(f.data);
            while (!cr.done()) {
                auto cf = cr.next();
                if (cf.num == 1) st.callTarget = decodeAppName(cf.data);
                else if (cf.num == 2) st.callEp = string(cf.data);
            }
            break;
        }
        case 3: st.kind = Stmt::Cond; break;
        case 4: st.kind = Stmt::Loop; break;
        case 5: st.kind = Stmt::LoopN; break;
        case 6: st.kind = Stmt::Alt; break;
        case 7: st.kind = Stmt::Group; break;
        case 8: { // ret
            st.kind = Stmt::Ret;
            Reader rr = sub(f.data);
            while (!rr.done()) {
                auto rf = rr.next();
                if (rf.num == 1) st.retPayload = string(rf.data);
            }
            break;
        }
        case 9: st.kind = Stmt::Foreach; break;
        case 10: decodeAttrsEntry(f.data, &st.attrs); break;
        default: break;
        }
    }
    return st;
}

static Endpoint decodeEndpoint(string_view s) {
    Endpoint ep;
    Reader r = sub(s);
    while (!r.done()) {
        auto f = r.next();
        switch (f.num) {
        case 1: ep.name = string(f.data); break;
        case 2: ep.longName = string(f.data); break;
        case 4: decodeAttrsEntry(f.data, &ep.attrs); break;
        case 6: ep.isPubsub = f.ival != 0; break;
        case 7: ep.stmts.push_back(decodeStmt(f.data)); break;
        case 8: { // rest_params
            ep.hasRest = true;
            Reader rr = sub(f.data);
            while (!rr.done()) {
                auto rf = rr.next();
                if (rf.num == 1) ep.restMethod = int(rf.ival);
                else if (rf.num == 2) ep.restPath = string(rf.data);
            }
            break;
        }
        default: break;
        }
    }
    return ep;
}

static App decodeApp(string_view s) {
    App app;
    Reader r = sub(s);
    while (!r.done()) {
        auto f = r.next();
        switch (f.num) {
        case 1: app.nameParts = decodeAppName(f.data); break;
        case 2: app.longName = string(f.data); break;
        case 4: decodeAttrsEntry(f.data, &app.attrs); break;
        case 5: { // EndpointsEntry
            Reader er = sub(f.data);
            string key;
            Endpoint ep;
            while (!er.done()) {
                auto ef = er.next();
                if (ef.num == 1) key = string(ef.data);
                else if (ef.num == 2) ep = decodeEndpoint(ef.data);
            }
            ep.name = key; // the map key is the endpoint's identity
            app.eps.push_back(std::move(ep));
            break;
        }
        case 6: { // TypesEntry
            Reader tr = sub(f.data);
            string key;
            string_view body;
            while (!tr.done()) {
                auto tf = tr.next();
                if (tf.num == 1) key = string(tf.data);
                else if (tf.num == 2) body = tf.data;
            }
            TypeDef t = decodeType(body, app.nameParts);
            t.name = std::move(key);
            app.types.push_back(std::move(t));
            break;
        }
        case 99: app.src = decodeSrc(f.data, &app.srcFile); break;
        default: break;
        }
    }
    return app;
}

// ---------------------------------------------------------------------------
// Name utilities: syslSafeName and friends, ported from
// vendor/sysl/pkg/importer/utils.arrai.

static bool isWordChar(char c) {
    return c == '-' || c == '_' || (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') ||
           (c >= 'A' && c <= 'Z');
}

static int hexVal(char c) {
    if (c >= '0' && c <= '9') return c - '0';
    if (c >= 'a' && c <= 'f') return c - 'a' + 10;
    if (c >= 'A' && c <= 'F') return c - 'A' + 10;
    return -1;
}

static const std::set<string>& syslKeywords() {
    static const std::set<string> kw = {
        "int",  "int32", "int64", "float", "float32", "float64", "decimal",
        "string", "date", "datetime", "bool", "bytes", "any",
    };
    return kw;
}

static string toLower(string_view s) {
    string out(s);
    for (char& c : out)
        if (c >= 'A' && c <= 'Z') c += 32;
    return out;
}

static string syslSafeName(string_view name) {
    // Collapse runs of [/\{} ] to '_'.
    string collapsed;
    for (size_t i = 0; i < name.size();) {
        char c = name[i];
        if (c == '/' || c == '\\' || c == '{' || c == '}' || c == ' ') {
            collapsed += '_';
            while (i < name.size() && (name[i] == '/' || name[i] == '\\' || name[i] == '{' ||
                                       name[i] == '}' || name[i] == ' '))
                i++;
        } else {
            collapsed += c;
            i++;
        }
    }
    // shallowEscape: unescape %XX, then percent-encode anything not a word char.
    string esc;
    for (size_t i = 0; i < collapsed.size(); i++) {
        char c = collapsed[i];
        if (c == '%' && i + 2 < collapsed.size() && hexVal(collapsed[i + 1]) >= 0 &&
            hexVal(collapsed[i + 2]) >= 0) {
            c = char(hexVal(collapsed[i + 1]) * 16 + hexVal(collapsed[i + 2]));
            i += 2;
        }
        if (isWordChar(c)) {
            esc += c;
        } else {
            static const char* hex = "0123456789ABCDEF";
            esc += '%';
            esc += hex[(uint8_t(c) >> 4) & 0xf];
            esc += hex[uint8_t(c) & 0xf];
        }
    }
    if (!esc.empty()) {
        char c0 = esc[0];
        bool safeStart = c0 == '_' || (c0 >= 'a' && c0 <= 'z') || (c0 >= 'A' && c0 <= 'Z');
        if (!safeStart) esc.insert(esc.begin(), '_');
    }
    if (syslKeywords().count(toLower(esc))) esc += '_';
    return esc;
}

// syslSafeType keeps native type names as-is; everything else is safenamed.
static string syslSafeType(string_view name) {
    static const std::set<string> native = {
        "int",  "int32", "int64", "string", "any",      "float", "float32",
        "float64", "date", "bool",  "decimal", "datetime", "bytes",
    };
    if (native.count(string(name))) return string(name);
    return syslSafeName(name);
}

static string joinParts(const std::vector<string>& parts, const char* sep, bool safe) {
    string out;
    for (size_t i = 0; i < parts.size(); i++) {
        if (i) out += sep;
        out += safe ? syslSafeName(parts[i]) : parts[i];
    }
    return out;
}

// ---------------------------------------------------------------------------
// Annotations and tags, following sysl.arrai's annotations/tags.

struct Anno {
    string name;
    AttrValue value;
};

static std::vector<Anno> annotationsOf(const AttrMap& attrs) {
    std::vector<Anno> out;
    for (const auto& [k, v] : attrs)
        if (k != "patterns") out.push_back({k, v});
    std::sort(out.begin(), out.end(), [](const Anno& a, const Anno& b) { return a.name < b.name; });
    return out;
}

static std::vector<string> tagsOf(const AttrMap& attrs) {
    std::vector<string> out;
    for (const auto& [k, v] : attrs)
        if (k == "patterns")
            for (const auto& e : v.elems) out.push_back(e.s);
    return out;
}

// tagsString: {'b','a'} -> " [~a, ~b]"
static string tagsString(std::vector<string> tags) {
    if (tags.empty()) return "";
    for (auto& t : tags) t = "~" + t;
    std::sort(tags.begin(), tags.end());
    string out = " [";
    for (size_t i = 0; i < tags.size(); i++) {
        if (i) out += ", ";
        out += tags[i];
    }
    return out + "]";
}

// ---------------------------------------------------------------------------
// Return-payload parsing, following payloadMacro's grammar:
//   payload -> (status ("<:" type)? | (status "<:")? type) attr?
// with status = ok | error | [1-5][0-9][0-9], type up to '[' or newline, and
// attr a bracketed list of ~modifiers and name=value pairs.

struct RetPayload {
    string status = "ok";
    string type;
    std::vector<string> modifiers;                  // tag names
    std::vector<std::pair<string, string>> nvps;    // name -> raw string value
};

static string_view trimSpaces(string_view s) {
    while (!s.empty() && s.front() == ' ') s.remove_prefix(1);
    while (!s.empty() && s.back() == ' ') s.remove_suffix(1);
    return s;
}

static string_view trimWs(string_view s) {
    auto isws = [](char c) { return c == ' ' || c == '\t' || c == '\n' || c == '\r'; };
    while (!s.empty() && isws(s.front())) s.remove_prefix(1);
    while (!s.empty() && isws(s.back())) s.remove_suffix(1);
    return s;
}

static RetPayload parseRetPayload(string_view payload) {
    RetPayload ret;
    string_view s = trimWs(payload);

    // Optional bracketed attr list at the end.
    if (!s.empty() && s.back() == ']') {
        size_t open = s.rfind('[');
        if (open != string_view::npos) {
            string_view attrs = s.substr(open + 1, s.size() - open - 2);
            s = trimWs(s.substr(0, open));
            size_t i = 0;
            while (i < attrs.size()) {
                while (i < attrs.size() && (attrs[i] == ' ' || attrs[i] == ',')) i++;
                if (i >= attrs.size()) break;
                if (attrs[i] == '~') {
                    size_t j = ++i;
                    while (j < attrs.size() && (isWordChar(attrs[j]) || attrs[j] == '+')) j++;
                    ret.modifiers.emplace_back(attrs.substr(i, j - i));
                    i = j;
                } else {
                    size_t j = i;
                    while (j < attrs.size() && isWordChar(attrs[j])) j++;
                    string name(attrs.substr(i, j - i));
                    i = j;
                    while (i < attrs.size() && attrs[i] == ' ') i++;
                    if (i < attrs.size() && attrs[i] == '=') i++;
                    while (i < attrs.size() && attrs[i] == ' ') i++;
                    if (i < attrs.size() && (attrs[i] == '"' || attrs[i] == '\'')) {
                        char q = attrs[i++];
                        string val;
                        while (i < attrs.size() && attrs[i] != q) {
                            if (q == '"' && attrs[i] == '\\' && i + 1 < attrs.size()) {
                                char e = attrs[++i];
                                switch (e) {
                                case 'n': val += '\n'; break;
                                case 't': val += '\t'; break;
                                case 'r': val += '\r'; break;
                                case 'b': val += '\b'; break;
                                default: val += e;
                                }
                            } else {
                                val += attrs[i];
                            }
                            i++;
                        }
                        if (i < attrs.size()) i++; // closing quote
                        ret.nvps.emplace_back(std::move(name), std::move(val));
                    }
                }
            }
            std::sort(ret.modifiers.begin(), ret.modifiers.end());
            std::sort(ret.nvps.begin(), ret.nvps.end());
        }
    }

    // status ("<:" type)? | (status "<:")? type
    auto isStatus = [](string_view w) {
        if (w == "ok" || w == "error") return true;
        return w.size() == 3 && w[0] >= '1' && w[0] <= '5' && w[1] >= '0' && w[1] <= '9' &&
               w[2] >= '0' && w[2] <= '9';
    };
    size_t sep = s.find("<:");
    if (sep != string_view::npos) {
        string_view st = trimWs(s.substr(0, sep));
        if (isStatus(st)) ret.status = string(st);
        ret.type = string(trimWs(s.substr(sep + 2)));
    } else if (isStatus(s)) {
        ret.status = string(s);
    } else {
        ret.type = string(s);
    }
    return ret;
}

// ---------------------------------------------------------------------------
// Type resolution: fixType + parseFieldType + resolvedType.

// unpackType splits "App.Type.Field" / "Type" / "App.Type" on '.'.
static FieldType fixType(const std::set<string, std::less<>>& appKeys,
                         const std::vector<string>& appParts, string_view payload) {
    string_view s = trimSpaces(payload);
    std::vector<string> segs;
    size_t start = 0;
    string str(s);
    while (true) {
        size_t dot = str.find('.', start);
        if (dot == string::npos) {
            segs.push_back(string(trimSpaces(string_view(str).substr(start))));
            break;
        }
        segs.push_back(string(trimSpaces(string_view(str).substr(start, dot - start))));
        start = dot + 1;
    }

    FieldType t;
    if (segs.size() == 1 && syslKeywords().count(toLower(segs[0]))) {
        t.kind = FieldType::Primitive;
        t.primitive = segs[0]; // resolvedType lowercases
        return t;
    }
    t.kind = FieldType::Ref;
    // unpackType: app = first segment iff there are 2+ segments; the last of
    // 2+ remaining segments is the field.
    string app;
    std::vector<string> path;
    if (segs.size() == 1) {
        path = {segs[0]};
    } else {
        app = segs[0];
        path.assign(segs.begin() + 1, segs.end());
    }
    auto appToParts = [](const string& name) {
        std::vector<string> parts;
        size_t at = 0;
        while (true) {
            size_t sep2 = name.find("::", at);
            if (sep2 == string::npos) {
                parts.push_back(string(trimSpaces(string_view(name).substr(at))));
                break;
            }
            parts.push_back(string(trimSpaces(string_view(name).substr(at, sep2 - at))));
            at = sep2 + 2;
        }
        return parts;
    };
    std::vector<string> fullPath;
    if (app.empty()) {
        t.appName = appParts;
        fullPath = path;
    } else if (appKeys.count(app)) {
        t.appName = appToParts(app);
        fullPath = path;
    } else {
        t.appName = appParts;
        fullPath.push_back(app);
        fullPath.insert(fullPath.end(), path.begin(), path.end());
    }
    for (auto& p : fullPath)
        if (!p.empty()) t.typePath.push_back(std::move(p));
    return t;
}

static string resolvedType(const FieldType& t) {
    switch (t.kind) {
    case FieldType::Primitive: return toLower(t.primitive);
    case FieldType::Ref: {
        string out;
        if (!t.appName.empty()) {
            for (size_t i = 0; i < t.appName.size(); i++) {
                if (i) out += "::";
                out += syslSafeName(t.appName[i]);
            }
            out += '.';
        }
        for (size_t i = 0; i < t.typePath.size(); i++) {
            if (i) out += '.';
            out += syslSafeName(t.typePath[i]);
        }
        return out;
    }
    case FieldType::Set: return "set of " + resolvedType(*t.inner);
    case FieldType::Sequence: return "sequence of " + resolvedType(*t.inner);
    case FieldType::Tuple: return "";
    default: return "";
    }
}

// ---------------------------------------------------------------------------
// Rendering, following reconstruct.arrai's template.

// goQuote renders a string as arr.ai's `:q` format does (Go %q).
static string goQuote(string_view s) {
    string out = "\"";
    for (char c : s) {
        switch (c) {
        case '"': out += "\\\""; break;
        case '\\': out += "\\\\"; break;
        case '\n': out += "\\n"; break;
        case '\t': out += "\\t"; break;
        case '\r': out += "\\r"; break;
        default:
            if (uint8_t(c) < 32) {
                char buf[8];
                std::snprintf(buf, sizeof buf, "\\x%02x", c);
                out += buf;
            } else {
                out += c;
            }
        }
    }
    return out + "\"";
}

// Emit appends lines at a given indent, skipping the indent on empty lines
// (trimLines strips trailing whitespace from every line).
struct Out {
    string buf;

    void line(int indent, string_view text) {
        if (!text.empty()) buf.append(indent, ' ');
        buf.append(text);
        buf += '\n';
    }

    void blank() { buf += '\n'; }
};

// renderAnnoValue follows resolvedAnnotations.renderValue for string and
// array-of-string annotation values.
static void renderAnnos(Out& out, int indent, const std::vector<Anno>& annos) {
    for (const auto& a : annos) {
        const AttrValue& v = a.value;
        if (v.isArray) {
            string parts = "[";
            for (size_t i = 0; i < v.elems.size(); i++) {
                if (i) parts += ", ";
                parts += goQuote(v.elems[i].s);
            }
            parts += "]";
            out.line(indent, "@" + a.name + " = " + parts);
        } else if (v.s.find('\n') != string::npos) {
            out.line(indent, "@" + a.name + " =:");
            // Trim trailing whitespace per line and any trailing newline.
            string_view s = v.s;
            while (!s.empty() && (s.back() == '\n' || s.back() == ' ' || s.back() == '\t'))
                s.remove_suffix(1);
            size_t start = 0;
            while (start <= s.size()) {
                size_t nl = s.find('\n', start);
                string_view ln = s.substr(start, nl == string_view::npos ? nl : nl - start);
                while (!ln.empty() && (ln.back() == ' ' || ln.back() == '\t')) ln.remove_suffix(1);
                out.line(indent + 4, "| " + string(ln));
                if (nl == string_view::npos) break;
                start = nl + 1;
            }
        } else {
            out.line(indent, "@" + a.name + " = " + goQuote(v.s));
        }
    }
}

// renderInlineAnnoAndTags: " [~tag, name=\"value\"]" or "".
static string inlineAnnoAndTags(const std::vector<std::pair<string, string>>& nvps,
                                std::vector<string> tags) {
    std::vector<string> items;
    for (auto& t : tags) items.push_back("~" + t);
    std::sort(items.begin(), items.end());
    for (const auto& [name, val] : nvps) items.push_back(name + "=\"" + val + "\"");
    if (items.empty()) return "";
    string out = " [";
    for (size_t i = 0; i < items.size(); i++) {
        if (i) out += ", ";
        out += items[i];
    }
    return out + "]";
}

struct Model {
    std::vector<App> apps;
    std::set<string, std::less<>> appKeys;
};

// Statement-level annotations render inline as name=value pairs.
static std::vector<std::pair<string, string>> annotationsToNvps(const AttrMap& attrs);

// renderStmt returns one statement's text, or "" for kinds reconstruct
// leaves unrendered (cond/loop/foreach/group and '...' actions).
static string renderStmt(const Model& m, const App& app, const Stmt& st) {
    string body;
    switch (st.kind) {
    case Stmt::Action:
        if (st.action == "...") return "";
        body = st.action;
        break;
    case Stmt::Call:
        body = joinParts(st.callTarget, " :: ", false) + " <- " + st.callEp;
        break;
    case Stmt::Ret: {
        RetPayload ret = parseRetPayload(st.retPayload);
        FieldType t = fixType(m.appKeys, app.nameParts, ret.type);
        body = "return " + ret.status + " <: " + resolvedType(t) +
               inlineAnnoAndTags(ret.nvps, ret.modifiers);
        break;
    }
    default:
        return ""; // cond/loop/loopN/foreach/alt/group: TODO in the original too
    }
    return body + inlineAnnoAndTags(annotationsToNvps(st.attrs), tagsOf(st.attrs));
}

static std::vector<std::pair<string, string>> annotationsToNvps(const AttrMap& attrs) {
    std::vector<std::pair<string, string>> out;
    for (const auto& [k, v] : attrs)
        if (k != "patterns") out.emplace_back(k, v.s);
    std::sort(out.begin(), out.end());
    return out;
}

static void renderApp(Out& out, const Model& m, const App& app) {
    // Header: joined safe name, optional long name, sorted tags.
    string header = joinParts(app.nameParts, " :: ", true);
    if (!app.longName.empty()) header += " \"" + app.longName + "\"";
    header += tagsString(tagsOf(app.attrs));
    header += ":";
    out.line(0, header);

    std::vector<Anno> appAnnos = annotationsOf(app.attrs);

    // Partition types: enums render separately; aliases and events are not
    // rendered (matching the original's TODOs).
    std::vector<const TypeDef*> enums, types;
    for (const auto& t : app.types) {
        if (t.isEnum) enums.push_back(&t);
        else if (t.isTuple || t.isRelation) types.push_back(&t);
    }
    auto byName = [](const TypeDef* a, const TypeDef* b) { return a->name < b->name; };
    std::sort(enums.begin(), enums.end(), byName);
    std::sort(types.begin(), types.end(), byName);

    std::vector<const Endpoint*> plainEps, restEps;
    for (const auto& ep : app.eps) {
        if (ep.isPubsub || ep.name == "...") continue;
        (ep.hasRest ? restEps : plainEps).push_back(&ep);
    }
    std::sort(plainEps.begin(), plainEps.end(),
              [](const Endpoint* a, const Endpoint* b) { return a->name < b->name; });

    bool isEmpty = appAnnos.empty() && enums.empty() && types.empty() && plainEps.empty() &&
                   restEps.empty();
    if (isEmpty) {
        out.line(4, "...");
        return;
    }

    renderAnnos(out, 4, appAnnos);

    // Sections separated by blank lines, per the template's ::\i\n:\n joins.
    bool needBlank = false;
    auto sectionBlock = [&]() {
        if (needBlank) out.blank();
        needBlank = true;
    };

    for (const TypeDef* e : enums) {
        sectionBlock();
        out.line(4, "!enum " + syslSafeName(e->name) + ":");
        if (e->enumItems.empty()) {
            out.line(8, "...");
        } else {
            auto items = e->enumItems;
            std::sort(items.begin(), items.end(),
                      [](const auto& a, const auto& b) { return a.second < b.second; });
            for (const auto& [k, v] : items) out.line(8, k + ": " + std::to_string(v));
        }
    }

    for (const TypeDef* t : types) {
        sectionBlock();
        out.line(4, string(t->isRelation ? "!table " : "!type ") + syslSafeName(t->name) +
                        tagsString(tagsOf(t->attrs)) + ":");
        renderAnnos(out, 8, annotationsOf(t->attrs));
        std::vector<const FieldDef*> fields;
        for (const auto& f : t->fields) fields.push_back(&f);
        std::sort(fields.begin(), fields.end(),
                  [](const FieldDef* a, const FieldDef* b) { return a->name < b->name; });
        if (fields.empty()) {
            out.line(8, "...");
        } else {
            for (const FieldDef* fp : fields) {
                const FieldDef& f = *fp;
                string ln = syslSafeName(f.name) + " <: " + resolvedType(f.type);
                if (f.opt) ln += "?";
                ln += tagsString(tagsOf(f.attrs));
                std::vector<Anno> fieldAnnos = annotationsOf(f.attrs);
                if (!fieldAnnos.empty()) {
                    ln += ":";
                    out.line(8, ln);
                    renderAnnos(out, 12, fieldAnnos);
                } else {
                    out.line(8, ln);
                }
            }
        }
    }

    static const char* methodNames[] = {"", "GET", "HEAD", "PUT", "POST", "DELETE", "PATCH", ""};
    if (!restEps.empty()) {
        // Group by path, order by path; methods ordered by name.
        std::map<string, std::vector<const Endpoint*>> byPath;
        for (const Endpoint* ep : restEps) byPath[ep->restPath].push_back(ep);
        for (auto& [path, eps] : byPath) {
            sectionBlock();
            out.line(4, path + ":");
            std::sort(eps.begin(), eps.end(), [](const Endpoint* a, const Endpoint* b) {
                return string(methodNames[a->restMethod]) < methodNames[b->restMethod];
            });
            for (const Endpoint* ep : eps) {
                out.line(8, string(methodNames[ep->restMethod]) + ":");
                bool any = false;
                for (const auto& st : ep->stmts) {
                    string s = renderStmt(m, app, st);
                    if (!s.empty()) {
                        out.line(12, s);
                        any = true;
                    }
                }
                if (!any) out.line(12, "...");
            }
        }
    }

    for (const Endpoint* ep : plainEps) {
        sectionBlock();
        out.line(4, syslSafeName(ep->name) + ":");
        bool any = false;
        for (const auto& st : ep->stmts) {
            string s = renderStmt(m, app, st);
            if (!s.empty()) {
                out.line(8, s);
                any = true;
            }
        }
        if (!any) out.line(8, "...");
    }
}

// ---------------------------------------------------------------------------
// arr.ai value repr for the final dict-of-files value.

static string arraiRepr(string_view s) {
    char delim = s.find('\'') == string_view::npos ? '\'' : '"';
    string out(1, delim);
    for (char c : s) {
        if (c == '\\' || c == delim) {
            out += '\\';
            out += c;
        } else if (uint8_t(c) >= 32) {
            out += c;
        } else if (c == '\n') {
            out += "\\n";
        } else if (c == '\t') {
            out += "\\t";
        } else if (c == '\r') {
            out += "\\r";
        } else {
            char buf[8];
            std::snprintf(buf, sizeof buf, "\\x%02x", c);
            out += buf;
        }
    }
    out += delim;
    return out;
}

int main(int argc, char** argv) {
    if (argc != 2) {
        std::fprintf(stderr, "usage: %s <model.sysl.pb>\n", argv[0]);
        return 2;
    }
    FILE* f = std::fopen(argv[1], "rb");
    if (!f) {
        std::perror(argv[1]);
        return 1;
    }
    std::fseek(f, 0, SEEK_END);
    long size = std::ftell(f);
    std::fseek(f, 0, SEEK_SET);
    string data(size_t(size), 0);
    if (std::fread(data.data(), 1, size_t(size), f) != size_t(size)) {
        std::fprintf(stderr, "short read\n");
        return 1;
    }
    std::fclose(f);

    // Decode the Module.
    Model m;
    Reader r = sub(data);
    while (!r.done()) {
        auto fld = r.next();
        if (fld.num == 2) { // AppsEntry
            Reader ar = sub(fld.data);
            string key;
            string_view body;
            while (!ar.done()) {
                auto af = ar.next();
                if (af.num == 1) key = string(af.data);
                else if (af.num == 2) body = af.data;
            }
            App app = decodeApp(body);
            app.key = std::move(key);
            m.appKeys.insert(app.key);
            m.apps.push_back(std::move(app));
        }
    }

    // Group apps into files by their source file (missing -> default.sysl),
    // ordered within each file by source start line, falling back to name
    // (numbers order before names, as in arr.ai).
    std::map<string, std::vector<const App*>> files;
    for (const auto& app : m.apps)
        files[app.src.present ? app.srcFile : "default.sysl"].push_back(&app);
    for (auto& [file, apps] : files) {
        std::stable_sort(apps.begin(), apps.end(), [](const App* a, const App* b) {
            bool la = a->src.present && a->src.startLine >= 0;
            bool lb = b->src.present && b->src.startLine >= 0;
            if (la != lb) return la; // numbers sort before name arrays
            if (la) return a->src.startLine < b->src.startLine;
            return a->nameParts < b->nameParts;
        });
    }

    // Render each file's contents and wrap in the (possibly nested) dict the
    // arr.ai pipeline builds from path segments.
    string result = "{";
    bool firstFile = true;
    for (const auto& [file, apps] : files) {
        Out out;
        bool firstApp = true;
        for (const App* app : apps) {
            if (!firstApp) out.blank();
            firstApp = false;
            renderApp(out, m, *app);
        }
        // trimLines leaves exactly one trailing newline.
        while (!out.buf.empty() && out.buf.back() == '\n') out.buf.pop_back();
        out.buf += '\n';

        if (!firstFile) result += ", ";
        firstFile = false;
        // fileMap splits the path on '/' into nested dicts.
        string_view path = file;
        while (!path.empty() && path.front() == '/') path.remove_prefix(1);
        std::vector<string_view> segs;
        size_t at = 0;
        while (true) {
            size_t slash = path.find('/', at);
            if (slash == string_view::npos) {
                segs.push_back(path.substr(at));
                break;
            }
            segs.push_back(path.substr(at, slash - at));
            at = slash + 1;
        }
        string open, close;
        for (size_t i = 0; i < segs.size(); i++) {
            open += arraiRepr(segs[i]);
            open += ": ";
            if (i + 1 < segs.size()) {
                open += "{";
                close = "}" + close;
            }
        }
        result += open + arraiRepr(out.buf) + close;
    }
    result += "}\n";

    std::fwrite(result.data(), 1, result.size(), stdout);
    return 0;
}
