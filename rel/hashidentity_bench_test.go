package rel

import "testing"

func BenchmarkSetEqual(b *testing.B) {
	sb := NewSetBuilder()
	for i := 0; i < 2000; i++ {
		sb.Add(NewNumber(float64(i)))
	}
	s, err := sb.Finish()
	if err != nil {
		b.Fatal(err)
	}
	t := s
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !s.Equal(t) {
			b.Fatal("expected equal")
		}
	}
}

func BenchmarkHashIdentityRelationHas(b *testing.B) {
	r := benchRelation(b)
	probe := NewTuple(
		NewAttr("id", NewNumber(0)),
		NewAttr("cust", NewNumber(0)),
		NewAttr("sku", NewNumber(0)),
		NewAttr("qty", NewNumber(1)),
		NewAttr("region", NewString([]rune("r0"))),
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.Has(probe)
	}
}
