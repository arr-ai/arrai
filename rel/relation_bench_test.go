package rel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Benchmarks for the core relational operations at a realistic size. These
// exist to make relation work measurable on its own, apart from the
// interpreter overhead that dominates end-to-end scripts.

const benchRelRows = 12000

func benchRelation(b *testing.B) Relation {
	b.Helper()
	sb := NewSetBuilder()
	for i := 0; i < benchRelRows; i++ {
		sb.Add(NewTuple(
			NewAttr("id", NewNumber(float64(i))),
			NewAttr("cust", NewNumber(float64(i%900))),
			NewAttr("sku", NewNumber(float64(i%320))),
			NewAttr("qty", NewNumber(float64(i%7+1))),
			NewAttr("region", NewString([]rune(fmt.Sprintf("r%d", i%12)))),
		))
	}
	s, err := sb.Finish()
	if err != nil {
		b.Fatal(err)
	}
	return s.(Relation)
}

func BenchmarkRelationBuild(b *testing.B) {
	tuples := make([]Value, benchRelRows)
	for i := range tuples {
		tuples[i] = NewTuple(
			NewAttr("id", NewNumber(float64(i))),
			NewAttr("cust", NewNumber(float64(i%900))),
			NewAttr("qty", NewNumber(float64(i%7+1))),
		)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sb := NewSetBuilder()
		for _, t := range tuples {
			sb.Add(t)
		}
		if _, err := sb.Finish(); err != nil {
			b.Fatal(err)
		}
	}
}

// Where keeps a subset of an existing relation: no new rows, no duplicates
// possible, yet the row set is rebuilt.
func BenchmarkRelationWhere(b *testing.B) {
	r := benchRelation(b)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := r.Where(func(v Value) (bool, error) {
			return v.(Tuple).MustGet("qty").(Number).Float64() > 4, nil
		}); err != nil {
			b.Fatal(err)
		}
	}
}

// Enumerating every row, the common denominator of most operations.
func BenchmarkRelationEnumerate(b *testing.B) {
	r := benchRelation(b)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		n := 0
		for e := r.Enumerator(); e.MoveNext(); {
			if e.Current().(Tuple).MustGet("qty") != nil {
				n++
			}
		}
	}
}

func BenchmarkRelationJoin(b *testing.B) {
	r := benchRelation(b)
	sb := NewSetBuilder()
	for i := 0; i < 900; i++ {
		sb.Add(NewTuple(
			NewAttr("cust", NewNumber(float64(i))),
			NewAttr("tier", NewString([]rune(fmt.Sprintf("t%d", i%4)))),
		))
	}
	c, err := sb.Finish()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = r.Join(c.(Relation), NamesSlice{"cust"},
			NamesSlice{"id", "cust", "sku", "qty", "region"}, NamesSlice{"tier"})
	}
}

// The memoised groupBy index that backs joins and indexed `where`.
func BenchmarkRelationGroupByCold(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := benchRelation(b)
		b.StartTimer()
		_ = r.rows.groupBy(valueProjector{r.attrMap["cust"]})
	}
}

func BenchmarkRelationHas(b *testing.B) {
	r := benchRelation(b)
	probes := make([]Value, 500)
	for i := range probes {
		probes[i] = NewTuple(
			NewAttr("id", NewNumber(float64(i*7))),
			NewAttr("cust", NewNumber(float64((i*7)%900))),
			NewAttr("sku", NewNumber(float64((i*7)%320))),
			NewAttr("qty", NewNumber(float64((i*7)%7+1))),
			NewAttr("region", NewString([]rune(fmt.Sprintf("r%d", (i*7)%12)))),
		)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, p := range probes {
			_ = r.Has(p)
		}
	}
}

// A group-by key must round-trip through the index: the encoding used to
// build it and the encoding used to probe it have to agree, including for
// values whose Go type is not comparable (Strings, sets, arrays).
func TestGroupByKeyRoundTrip(t *testing.T) {
	t.Parallel()

	rows := []Tuple{
		NewTuple(NewAttr("a", NewNumber(1)), NewAttr("b", NewString([]rune("x")))),
		NewTuple(NewAttr("a", NewNumber(2)), NewAttr("b", NewString([]rune("y")))),
		NewTuple(NewAttr("a", NewNumber(3)), NewAttr("b", NewString([]rune("x")))),
	}
	sb := NewSetBuilder()
	for _, r := range rows {
		sb.Add(r)
	}
	s, err := sb.Finish()
	require.NoError(t, err)
	r := s.(Relation)

	for _, attr := range []string{"a", "b"} {
		idx := r.attrMap[attr]
		p := valueProjector{idx}
		group := r.rows.groupBy(p)
		for e := r.Enumerator(); e.MoveNext(); {
			row := r.tupleToValues(e.Current().(Tuple))
			got, has := group.getKey(row[idx])
			assert.True(t, has, "probing %s with %v must find its own row", attr, row)
			found := false
			for _, id := range got {
				if group.row(id).equalValues(row) {
					found = true
					break
				}
			}
			assert.True(t, found, "row %v missing from its group for %s", row, attr)
		}
	}
}
