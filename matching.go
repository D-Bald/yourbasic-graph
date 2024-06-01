package graph

// Matching calculates a maximal matching on a bipartite graph and returns the solution
// as an array `a` such that each row `i` is matched to column `a[i]`.
// If g isn't bipartite, it returns an empty slice and sets ok to false.
func Matching(g Iterator) (matching []int, ok bool) {
	// TODO: Implement hungarian method:
	// Ressources:
	//	- https://github.com/oddg/hungarian-algorithm/blob/master/label.go
	//	- https://cse.hkust.edu.hk/~golin/COMP572/Notes/Matching.pdf

	n := g.Order()

	part, ok := Bipartition(g)
	if !ok {
		return []int{}, false
	}

	// Create Sets for each partition
	x := make(set, len(part))
	x.insertSlice(part)
	y := make(set, n-len(part))
	for i := range n {
		if !x.contains(i) {
			y.insert(i)
		}
	}

	// 1. Generate initial labelling ℓ and matching M in Eℓ.
	// ∀y ∈ Y, ℓ(y) = 0, ∀x ∈ X, ℓ(x) = max y∈Y {w(x, y)}
	labels := make([]int64, n)
	for v := range x {
		g.Visit(v, func(w int, c int64) (skip bool) {
			if c > labels[v] {
				labels[v] = c
			}
			return
		})
	}

	// 2. f M perfect, stop.
	// Otherwise pick free vertex u ∈ X.
	// Set S = {u}, T = ∅

	// 3. If Nℓ(S) = T, update labels (forcing Nℓ(S) 6 = T
	// 		αℓ = mins∈S, y6 ∈T {ℓ(x) + ℓ(y) − w(x, y)}
	// 				ℓ(v) − αℓ if v ∈ S
	// 		ℓ′(v) = ℓ(v) + αℓ if v ∈ T
	// 				ℓ(v) otherwise

	// 4. If Nℓ(S) 6 = T , pick y ∈ Nℓ(S) − T
	// 4 a) If y free, u − y is augmenting path.
	// 		Augment M and go to 2.

	// 4 b) If y matched, say to z, extend alternating tree:
	//		S = S ∪ {z}, T = T ∪ {y}. Go to 3

	return []int{}, true
}

// Set implementation close to https://github.com/hashicorp/go-set/blob/v0.1.14/set.go

type nothing struct{}

var sentinel = nothing{}

type set map[int]nothing

// Insert item into s.
//
// Return true if s was modified (item was not already in s), false otherwise.
func (s set) insert(item int) bool {
	if _, exists := s[item]; exists {
		return false
	}
	s[item] = sentinel
	return true
}

// InsertSlice will insert each item in items into s.
//
// Return true if s was modified (at least one item was not already in s), false otherwise.
func (s set) insertSlice(items []int) bool {
	modified := false
	for _, item := range items {
		if s.insert(item) {
			modified = true
		}
	}
	return modified
}

// Remove will remove item from s.
//
// Return true if s was modified (item was present), false otherwise.
func (s set) remove(item int) bool {
	if _, exists := s[item]; !exists {
		return false
	}
	delete(s, item)
	return true
}

// RemoveSlice will remove each item in items from s.
//
// Return true if s was modified (any item was present), false otherwise.
func (s set) removeSlice(items []int) bool {
	modified := false
	for _, item := range items {
		if s.remove(item) {
			modified = true
		}
	}
	return modified
}

// Contains returns whether item is present in s.
func (s set) contains(item int) bool {
	_, exists := s[item]
	return exists
}

// Equal returns whether s and o contain the same elements.
func (s set) equal(o set) bool {
	if len(s) != len(o) {
		return false
	}

	for item := range s {
		if !o.contains(item) {
			return false
		}
	}

	return true
}
