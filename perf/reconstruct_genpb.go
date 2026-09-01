package perf

import (
	"strconv"

	"google.golang.org/protobuf/encoding/protowire"
)

// generateSyslModelPB emits a Module protobuf matching gen_model.py's synthetic
// apps, so the reconstruct scenario can be rescaled without the sysl CLI.
func generateSyslModelPB(n int) []byte {
	out := make([]byte, 0, n*1200)
	for i := 0; i < n; i++ {
		out = appendBytes(out, 2, encodeAppEntry(i, n))
	}
	return out
}

func encodeAppEntry(i, n int) []byte {
	name := appName(i)
	return concat(pbString(1, name), pbBytes(2, encodeApp(i, n)))
}

func encodeApp(i, n int) []byte {
	name := appName(i)
	line := 4 + i*12
	b := pbBytes(1, pbString(1, name))
	b = appendBytes(b, 4, attrString("owner", "team"+strconv.Itoa(i%20)))
	b = appendBytes(b, 4, attrPattern("service"))
	b = appendBytes(b, 4, attrString("description", "Synthetic service "+strconv.Itoa(i)))
	b = appendBytes(b, 5, encodeEndpointEntry(i, n))
	b = appendBytes(b, 6, encodeTypeEntry())
	return appendBytes(b, 99, sourceContext("model.sysl", line))
}

func encodeEndpointEntry(i, n int) []byte {
	nxt := (i + 1) % n
	ep := pbString(1, "Get")
	ep = appendBytes(ep, 4, attrPattern("rest"))
	call := concat(pbBytes(1, pbString(1, appName(nxt))), pbString(2, "Get"))
	ep = appendBytes(ep, 7, pbBytes(2, call))
	ret := pbString(1, "ok <: "+appName(i)+".Data")
	ep = appendBytes(ep, 7, pbBytes(8, ret))
	paramType := pbVarint(1, 4) // primitive INT
	param := concat(pbString(1, "id"), pbBytes(2, paramType))
	ep = appendBytes(ep, 9, param)
	return concat(pbString(1, "Get"), pbBytes(2, ep))
}

func encodeTypeEntry() []byte {
	idType := concat(pbBytes(8, attrPattern("pk")), pbVarint(1, 4))
	id := concat(pbString(1, "id"), pbBytes(2, idType))
	nameType := concat(pbBytes(8, attrString("description", "Synthetic field")), pbVarint(1, 6))
	name := concat(pbString(1, "name"), pbBytes(2, nameType))
	valType := concat(pbVarint(12, 1), pbVarint(1, 4)) // opt int
	val := concat(pbString(1, "value"), pbBytes(2, valType))
	tuple := concat(pbBytes(1, id), pbBytes(1, name), pbBytes(1, val))
	return concat(pbString(1, "Data"), pbBytes(2, pbBytes(3, tuple)))
}

func attrString(key, val string) []byte {
	return concat(pbString(1, key), pbBytes(2, pbString(4, val)))
}

func attrPattern(name string) []byte {
	return concat(pbString(1, "patterns"), pbBytes(2, pbBytes(7, pbBytes(1, pbString(4, name)))))
}

func sourceContext(file string, startLine int) []byte {
	return concat(pbString(1, file), pbBytes(2, pbVarint(1, uint64(startLine))))
}

func appName(i int) string { return "App" + strconv.Itoa(i) }

func pbString(num int, s string) []byte {
	b := protowire.AppendTag(nil, protowire.Number(num), protowire.BytesType)
	return protowire.AppendString(b, s)
}

func pbBytes(num int, msg []byte) []byte {
	return appendBytes(nil, num, msg)
}

func pbVarint(num int, v uint64) []byte {
	b := protowire.AppendTag(nil, protowire.Number(num), protowire.VarintType)
	return protowire.AppendVarint(b, v)
}

func appendBytes(b []byte, num int, msg []byte) []byte {
	b = protowire.AppendTag(b, protowire.Number(num), protowire.BytesType)
	return protowire.AppendBytes(b, msg)
}

func concat(parts ...[]byte) []byte {
	n := 0
	for _, p := range parts {
		n += len(p)
	}
	out := make([]byte, 0, n)
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}
