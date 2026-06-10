package utils

// Combinations generates all k-element combinations of items.
// The order of elements within each combination matches their order in items.
func Combinations(items []string, k int) [][]string {
	if k < 0 || k > len(items) {
		return nil
	}
	if k == 0 {
		return [][]string{{}}
	}
	var result [][]string
	combo := make([]int, k)
	var generate func(start, depth int)
	generate = func(start, depth int) {
		if depth == k {
			c := make([]string, k)
			for i, idx := range combo {
				c[i] = items[idx]
			}
			result = append(result, c)
			return
		}
		for i := start; i <= len(items)-(k-depth); i++ {
			combo[depth] = i
			generate(i+1, depth+1)
		}
	}
	generate(0, 0)
	return result
}

// Permutations generates all permutations of items using Heap's algorithm.
func Permutations(items []string) [][]string {
	n := len(items)
	if n == 0 {
		return [][]string{{}}
	}
	var result [][]string
	perm := make([]string, n)
	copy(perm, items)

	// snapshot appends a copy of the current permutation.
	snapshot := func() {
		p := make([]string, n)
		copy(p, perm)
		result = append(result, p)
	}

	// Heap's algorithm (iterative).
	c := make([]int, n)
	snapshot()
	i := 0
	for i < n {
		if c[i] < i {
			if i%2 == 0 {
				perm[0], perm[i] = perm[i], perm[0]
			} else {
				perm[c[i]], perm[i] = perm[i], perm[c[i]]
			}
			snapshot()
			c[i]++
			i = 0
		} else {
			c[i] = 0
			i++
		}
	}
	return result
}

// CartesianProduct computes the Cartesian product of integer sets.
// For example, CartesianProduct([][]int{{0,1},{2,3}}) returns
// [[0,2],[0,3],[1,2],[1,3]].
func CartesianProduct(sets [][]int) [][]int {
	if len(sets) == 0 {
		return [][]int{{}}
	}
	// Compute total number of tuples.
	total := 1
	for _, s := range sets {
		if len(s) == 0 {
			return nil
		}
		total *= len(s)
	}
	result := make([][]int, 0, total)
	tuple := make([]int, len(sets))

	var build func(depth int)
	build = func(depth int) {
		if depth == len(sets) {
			t := make([]int, len(sets))
			copy(t, tuple)
			result = append(result, t)
			return
		}
		for _, v := range sets[depth] {
			tuple[depth] = v
			build(depth + 1)
		}
	}
	build(0)
	return result
}
