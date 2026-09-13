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
	for b.Loop() {
		if !s.Equal(t) {
			b.Fatal("expected equal")
		}
	}
}
