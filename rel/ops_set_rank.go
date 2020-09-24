package rel

import (
	"sort"
)

// Rank ...
func Rank(s Set, rankerf func(v Tuple) (Tuple, error)) (Set, error) {
	if !s.IsTrue() {
		return None, nil
	}
	entries := []rankerEntry{}
	for e := s.Enumerator(); e.MoveNext(); {
		input := e.Current().(Tuple)
		rankers, err := rankerf(input)
		if err != nil {
			return nil, err
		}
		entries = append(entries, rankerEntry{input, rankers})
	}
	ranker := newRanker(entries)
	for _, attr := range entries[0].ranker.Names().Names() {
		ranker.attr = attr
		sort.Sort(ranker)
		current := ranker.entries[0].ranker.MustGet(attr)
		rank := 0
		for i, r := range ranker.entries {
			a := r.ranker.MustGet(attr)
			if !a.Equal(current) {
				rank = i
				current = a
			}
			r.input = r.input.With(attr, NewNumber(float64(rank)))
		}
	}

	values := make([]Value, 0, len(entries))
	for _, entry := range entries {
		values = append(values, entry.input)
	}
	return NewSet(values...)
}

type rankerEntry struct {
	input, ranker Tuple
}

type rankerSlice struct {
	entries []*rankerEntry
	attr    string
}

func newRanker(entries []rankerEntry) rankerSlice {
	var r rankerSlice
	for i := range entries {
		r.entries = append(r.entries, &entries[i])
	}
	return r
}

// Len is the number of elements in the collection.
func (o rankerSlice) Len() int {
	return len(o.entries)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (o rankerSlice) Less(i, j int) bool {
	a := o.entries[i].ranker.MustGet(o.attr)
	b := o.entries[j].ranker.MustGet(o.attr)
	return a.Less(b)
}

// Swap swaps the elements with indexes i and j.
func (o rankerSlice) Swap(i, j int) {
	o.entries[i], o.entries[j] = o.entries[j], o.entries[i]
}
