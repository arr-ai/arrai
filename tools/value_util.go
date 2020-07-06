package tools

import "github.com/arr-ai/arrai/rel"

func ValueAsString(v rel.Value) (string, bool) {
	if v == nil {
		return "", false
	}

	switch v := v.(type) {
	case rel.String:
		return v.String(), true
	case rel.GenericSet:
		return "", !v.IsTrue()
	}
	return "", false
}

func ValueAsBytes(v rel.Value) ([]byte, bool) {
	if v == nil {
		return nil, false
	}

	switch v := v.(type) {
	case rel.Bytes:
		return v.Bytes(), true
	case rel.GenericSet:
		return nil, !v.IsTrue()
	}
	return nil, false
}
