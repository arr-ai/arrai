package syntax

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/rel"
	"github.com/spf13/afero"
)

const compiledPlanPath = "/plan.bin"

func init() {
	rel.RegisterPlanLower(lowerSyntaxExpr)
	rel.RegisterPlanLift("import", liftImport)
	rel.RegisterPlanLift("package", liftPackage)
	rel.RegisterPlanLift("xstr", liftXstr)
}

// portablePath guards against a machine-specific absolute path ending up in
// the compiled plan. compile.go already resolves ImportExpr.path to a
// portable, bundle-relative form (see portableBundlePath) before it gets
// here, so this should be a no-op in practice; it's a cheap safety net so any
// future path-construction bug drops the path (losing debug info) rather
// than silently making the bundle non-reproducible. ModuleDir/NoModuleDir
// prefixes are allow-listed since they're virtual bundle-zip paths, not real
// filesystem locations -- filepath.IsAbs would otherwise flag them too, since
// they start with "/" like any genuine OS-absolute path on Unix.
func portablePath(p string) string {
	if strings.HasPrefix(p, ModuleDir+"/") || strings.HasPrefix(p, NoModuleDir+"/") {
		return p
	}
	if filepath.IsAbs(p) {
		return ""
	}
	return p
}

func lowerChild(e rel.Expr) (rel.PlanNode, error) {
	p, err := rel.LowerPlan(e)
	if err != nil {
		return rel.PlanNode{}, err
	}
	return p.Root, nil
}

func lowerSyntaxExpr(e rel.Expr) (rel.PlanNode, bool, error) {
	switch e := e.(type) {
	case ImportExpr:
		inner, err := lowerChild(e.importedExpr)
		if err != nil {
			return rel.PlanNode{}, true, err
		}
		return rel.PlanNode{K: "import", Str: portablePath(e.path), Kids: []rel.PlanNode{inner}}, true, nil
	case PackageExpr:
		inner, err := lowerChild(e.a)
		if err != nil {
			return rel.PlanNode{}, true, err
		}
		return rel.PlanNode{K: "package", Kids: []rel.PlanNode{inner}}, true, nil
	case *xstrExpr:
		kids := make([]rel.PlanNode, len(e.parts))
		for i, p := range e.parts {
			part := rel.PlanNode{K: "xpart", Str: p.literal, Op: p.format, Attr: p.delim, Bytes: []byte(p.tail)}
			if p.expr != nil {
				n, err := lowerChild(p.expr)
				if err != nil {
					return rel.PlanNode{}, true, err
				}
				part.Kids = []rel.PlanNode{n}
			}
			kids[i] = part
		}
		return rel.PlanNode{K: "xstr", Kids: kids}, true, nil
	default:
		return rel.PlanNode{}, false, nil
	}
}

func liftImport(n rel.PlanNode) (rel.Expr, error) {
	inner, err := rel.LiftPlan(&rel.Plan{Root: n.Kids[0]})
	if err != nil {
		return nil, err
	}
	return NewImportExpr(rel.PlanScanner(), inner, n.Str), nil
}

func liftPackage(n rel.PlanNode) (rel.Expr, error) {
	inner, err := rel.LiftPlan(&rel.Plan{Root: n.Kids[0]})
	if err != nil {
		return nil, err
	}
	return NewPackageExpr(rel.PlanScanner(), inner), nil
}

func liftXstr(n rel.PlanNode) (rel.Expr, error) {
	parts := make([]xstrPart, len(n.Kids))
	for i, k := range n.Kids {
		p := xstrPart{literal: k.Str, format: k.Op, delim: k.Attr, tail: string(k.Bytes)}
		if len(k.Kids) > 0 {
			e, err := rel.LiftPlan(&rel.Plan{Root: k.Kids[0]})
			if err != nil {
				return nil, err
			}
			p.expr = e
		}
		parts[i] = p
	}
	return &xstrExpr{ExprScanner: rel.ExprScanner{Src: rel.PlanScanner()}, parts: parts}, nil
}

// WriteCompiledPlan adds /plan.bin to the bundle zip (🎯T25).
func WriteCompiledPlan(ctx context.Context, expr rel.Expr) error {
	if !isBundling(ctx) {
		return nil
	}
	p, err := rel.LowerPlan(expr)
	if err != nil {
		return err
	}
	b, err := rel.EncodePlan(p)
	if err != nil {
		return err
	}
	return ctxfs.ZipCreate(ctx, bundleFsKey, compiledPlanPath, b)
}

// LoadCompiledPlan reads /plan.bin from a running bundle, or nil if absent.
func LoadCompiledPlan(ctx context.Context) (*rel.Plan, error) {
	fs := ctxfs.SourceFsFrom(ctx)
	if fs == nil {
		return nil, nil
	}
	buf, err := afero.ReadFile(fs, compiledPlanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return rel.DecodePlan(buf)
}
