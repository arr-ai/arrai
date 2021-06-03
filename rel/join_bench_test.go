package rel

import (
	"context"
	"testing"

	"github.com/arr-ai/frozen"
	"github.com/arr-ai/wbnf/parser"
	"github.com/stretchr/testify/require"
)

func generateTuples(n int) ([]Tuple, []Tuple) {
	t1, t2 := make([]Tuple, 3*n), make([]Tuple, 3*n)
	for i := 0; i < 3*n; i++ {
		t1[i] = NewTuple(
			NewAttr("a", NewNumber(float64(i))),
			NewAttr("b", NewNumber(float64(i))),
		)
		t2[i] = NewTuple(
			NewAttr("b", NewNumber(float64(i+n))),
			NewAttr("c", NewNumber(float64(i+n))),
		)
	}
	return t1, t2
}

func genericSetTuples(n int) (GenericSet, GenericSet) {
	t1, t2 := generateTuples(n)
	f := func(t []Tuple) GenericSet {
		sb := frozen.NewSetBuilder(3 * n)
		for _, tu := range t {
			sb.Add(tu)
		}
		return GenericSet{sb.Finish()}
	}
	return f(t1), f(t2)
}

func generateRelations(n int) (Relation, Relation) {
	t1, t2 := generateTuples(n)
	f := func(t []Tuple, names []string) Relation {
		sb := newRelationBuilder(names, 3*n)
		for _, tu := range t {
			sb.Add(tu)
		}
		s, err := sb.Finish()
		if err != nil {
			panic(err)
		}
		return s.(Relation)
	}
	return f(t1, []string{"a", "b"}), f(t2, []string{"b", "c"})
}

func BenchmarkGenericSetJoin100(b *testing.B) {
	benchmarkGenericSetJoin(b, 100)
}

func BenchmarkGenericSetJoin1000(b *testing.B) {
	benchmarkGenericSetJoin(b, 1000)
}

func BenchmarkGenericSetJoin10000(b *testing.B) {
	benchmarkGenericSetJoin(b, 10000)
}

// func BenchmarkGenericSetJoin100000(b *testing.B) {
// 	testGenericSetJoin(b, 100000)
// }

// func BenchmarkGenericSetJoin1000000(b *testing.B) {
// 	testGenericSetJoin(b, 1000000)
// }

func BenchmarkRelationSetJoin100(b *testing.B) {
	benchmarkRelationSetJoin(b, 100)
}

func BenchmarkRelationSetJoin1000(b *testing.B) {
	benchmarkRelationSetJoin(b, 1000)
}

func BenchmarkRelationSetJoin10000(b *testing.B) {
	benchmarkRelationSetJoin(b, 10000)
}

// func BenchmarkRelationSetJoin100000(b *testing.B) {
// 	testRelationSetJoin(b, 100000)
// }

// func BenchmarkRelationSetJoin1000000(b *testing.B) {
// 	testRelationSetJoin(b, 1000000)
// }

func benchmarkGenericSetJoin(b *testing.B, n int) {
	g1, g2 := genericSetTuples(n)
	benchmarkSetJoin(b, g1, g2)
}

func benchmarkRelationSetJoin(b *testing.B, n int) {
	r1, r2 := generateRelations(n)
	benchmarkSetJoin(b, r1, r2)
}

func benchmarkSetJoin(b *testing.B, s1, s2 Set) {
	expr := NewJoinExpr(*parser.NewScanner(""), s1, s2)
	_, err := expr.Eval(context.Background(), EmptyScope)
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := expr.Eval(context.Background(), EmptyScope)
		require.NoError(b, err)
	}
}
