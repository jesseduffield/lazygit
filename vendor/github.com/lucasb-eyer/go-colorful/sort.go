// This file provides functions for sorting colors.

package colorful

import (
	"math"
	"sort"
)

// An element represents a single element of a set.  It is used to
// implement a disjoint-set forest.
type element struct {
	parent *element // Parent element
	rank   int      // Rank (approximate depth) of the subtree with this element as root
}

// newElement creates a singleton set and returns its sole element.
func newElement() *element {
	s := &element{}
	s.parent = s
	return s
}

// find returns an arbitrary element of a set when invoked on any element of
// the set, The important feature is that it returns the same value when
// invoked on any element of the set.  Consequently, it can be used to test if
// two elements belong to the same set.
func (e *element) find() *element {
	for e.parent != e {
		e.parent = e.parent.parent
		e = e.parent
	}
	return e
}

// union establishes the union of two sets when given an element from each set.
// Afterwards, the original sets no longer exist as separate entities.
func union(e1, e2 *element) {
	// Ensure the two elements aren't already part of the same union.
	e1Root := e1.find()
	e2Root := e2.find()
	if e1Root == e2Root {
		return
	}

	// Create a union by making the shorter tree point to the root of the
	// larger tree.
	switch {
	case e1Root.rank < e2Root.rank:
		e1Root.parent = e2Root
	case e1Root.rank > e2Root.rank:
		e2Root.parent = e1Root
	default:
		e2Root.parent = e1Root
		e1Root.rank++
	}
}

// An edgeIdxs describes an edge in a graph or tree.  The vertices in the edge
// are indexes into a list of Color values.
type edgeIdxs [2]int

// An edgeDistance is a map from an edge (pair of indices) to a distance
// between the two vertices.
type edgeDistance map[edgeIdxs]float64

// allToAllDistancesCIEDE2000 computes the CIEDE2000 distance between each pair of
// colors.  It returns a map from a pair of indices (u, v) with u < v to a
// distance.
func allToAllDistancesCIEDE2000(cs []Color) edgeDistance {
	nc := len(cs)
	m := make(edgeDistance, nc*nc)
	for u := 0; u < nc-1; u++ {
		for v := u + 1; v < nc; v++ {
			m[edgeIdxs{u, v}] = cs[u].DistanceCIEDE2000(cs[v])
		}
	}
	return m
}

// sortEdges sorts all edges in a distance map by increasing vertex distance.
func sortEdges(m edgeDistance) []edgeIdxs {
	es := make([]edgeIdxs, 0, len(m))
	for uv := range m {
		es = append(es, uv)
	}
	sort.Slice(es, func(i, j int) bool {
		return m[es[i]] < m[es[j]]
	})
	return es
}

// minSpanTree computes a minimum spanning tree from a vertex count and a
// distance-sorted edge list.  It returns the subset of edges that belong to
// the tree, including both (u, v) and (v, u) for each edge.
func minSpanTree(nc int, es []edgeIdxs) map[edgeIdxs]struct{} {
	// Start with each vertex in its own set.
	elts := make([]*element, nc)
	for i := range elts {
		elts[i] = newElement()
	}

	// Run Kruskal's algorithm to construct a minimal spanning tree.
	mst := make(map[edgeIdxs]struct{}, nc)
	for _, uv := range es {
		u, v := uv[0], uv[1]
		if elts[u].find() == elts[v].find() {
			continue // Same set: edge would introduce a cycle.
		}
		mst[uv] = struct{}{}
		mst[edgeIdxs{v, u}] = struct{}{}
		union(elts[u], elts[v])
	}
	return mst
}

// traverseMST walks a minimum spanning tree in prefix order.
func traverseMST(mst map[edgeIdxs]struct{}, root int) []int {
	// Compute a list of neighbors for each vertex.
	neighs := make(map[int][]int, len(mst))
	for uv := range mst {
		u, v := uv[0], uv[1]
		neighs[u] = append(neighs[u], v)
	}
	for u, vs := range neighs {
		sort.Ints(vs)
		copy(neighs[u], vs)
	}

	// Walk the tree from a given vertex.
	order := make([]int, 0, len(neighs))
	visited := make(map[int]bool, len(neighs))
	var walkFrom func(int)
	walkFrom = func(r int) {
		// Visit the starting vertex.
		order = append(order, r)
		visited[r] = true

		// Recursively visit each child in turn.
		for _, c := range neighs[r] {
			if !visited[c] {
				walkFrom(c)
			}
		}
	}
	walkFrom(root)
	return order
}

// Sorted sorts a list of Color values.  Sorting is not a well-defined operation
// for colors so the intention here primarily is to order colors so that the
// transition from one to the next is fairly smooth.
func Sorted(cs []Color) []Color {
	// Do nothing in trivial cases.
	newCs := make([]Color, len(cs))
	if len(cs) < 2 {
		copy(newCs, cs)
		return newCs
	}

	// Compute the distance from each color to every other color.
	dists := allToAllDistancesCIEDE2000(cs)

	// Produce a list of edges in increasing order of the distance between
	// their vertices.
	edges := sortEdges(dists)

	// Construct a minimum spanning tree from the list of edges.
	mst := minSpanTree(len(cs), edges)

	// Find the darkest color in the list.
	var black Color
	var dIdx int             // Index of darkest color
	light := math.MaxFloat64 // Lightness of darkest color (distance from black)
	for i, c := range cs {
		d := black.DistanceCIEDE2000(c)
		if d < light {
			dIdx = i
			light = d
		}
	}

	// Traverse the tree starting from the darkest color.
	idxs := traverseMST(mst, dIdx)

	// Convert the index list to a list of colors, overwriting the input.
	for i, idx := range idxs {
		newCs[i] = cs[idx]
	}
	return newCs
}
