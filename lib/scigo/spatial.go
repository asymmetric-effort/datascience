package scigo

import (
	"errors"
	"math"
	"sort"
)

// ---------------------------------------------------------------------------
// Pairwise Distance Functions
// ---------------------------------------------------------------------------

// Cdist computes the pairwise distance between each pair of rows in xa and xb.
// Supported metrics: "euclidean" (default), "manhattan", "chebyshev".
func Cdist(xa, xb [][]float64, metric string) [][]float64 {
	na := len(xa)
	nb := len(xb)
	if na == 0 || nb == 0 {
		return nil
	}

	distFn := getDistFunc(metric)
	result := make([][]float64, na)
	for i := 0; i < na; i++ {
		result[i] = make([]float64, nb)
		for j := 0; j < nb; j++ {
			result[i][j] = distFn(xa[i], xb[j])
		}
	}
	return result
}

// Pdist computes the pairwise distance between each pair of rows in x,
// returning a condensed distance vector.
func Pdist(x [][]float64, metric string) []float64 {
	n := len(x)
	if n < 2 {
		return nil
	}

	distFn := getDistFunc(metric)
	size := n * (n - 1) / 2
	result := make([]float64, size)
	idx := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			result[idx] = distFn(x[i], x[j])
			idx++
		}
	}
	return result
}

// Squareform converts a condensed distance vector to a square distance matrix.
// n is the number of original observations.
func Squareform(d []float64, n int) [][]float64 {
	expected := n * (n - 1) / 2
	if len(d) != expected {
		return nil
	}

	result := make([][]float64, n)
	for i := range result {
		result[i] = make([]float64, n)
	}

	idx := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			result[i][j] = d[idx]
			result[j][i] = d[idx]
			idx++
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// KDTree
// ---------------------------------------------------------------------------

// KDTree implements a k-d tree for efficient nearest-neighbor queries.
type KDTree struct {
	root   *kdNode
	points [][]float64
	dim    int
}

type kdNode struct {
	index int
	axis  int
	left  *kdNode
	right *kdNode
}

// NewKDTree builds a k-d tree from the given points.
func NewKDTree(points [][]float64) *KDTree {
	if len(points) == 0 {
		return &KDTree{}
	}
	dim := len(points[0])
	indices := make([]int, len(points))
	for i := range indices {
		indices[i] = i
	}
	root := buildKD(points, indices, 0, dim)
	return &KDTree{root: root, points: points, dim: dim}
}

func buildKD(points [][]float64, indices []int, depth, dim int) *kdNode {
	if len(indices) == 0 {
		return nil
	}
	axis := depth % dim

	// Sort indices by the axis coordinate.
	sort.Slice(indices, func(i, j int) bool {
		return points[indices[i]][axis] < points[indices[j]][axis]
	})

	mid := len(indices) / 2
	return &kdNode{
		index: indices[mid],
		axis:  axis,
		left:  buildKD(points, indices[:mid], depth+1, dim),
		right: buildKD(points, indices[mid+1:], depth+1, dim),
	}
}

// Query finds the k nearest neighbors to the given point.
// Returns the indices and distances of the k nearest neighbors.
func (t *KDTree) Query(point []float64, k int) ([]int, []float64) {
	if t.root == nil || k <= 0 {
		return nil, nil
	}
	if k > len(t.points) {
		k = len(t.points)
	}

	// Use a max-heap of size k.
	type neighbor struct {
		index int
		dist  float64
	}
	best := make([]neighbor, 0, k)
	worstDist := math.Inf(1)

	var search func(node *kdNode)
	search = func(node *kdNode) {
		if node == nil {
			return
		}

		d := euclidean(point, t.points[node.index])

		if len(best) < k {
			best = append(best, neighbor{node.index, d})
			if len(best) == k {
				// Find worst distance.
				worstDist = best[0].dist
				for _, b := range best {
					if b.dist > worstDist {
						worstDist = b.dist
					}
				}
			}
		} else if d < worstDist {
			// Replace the worst.
			wi := 0
			for i, b := range best {
				if b.dist > best[wi].dist {
					wi = i
				}
			}
			best[wi] = neighbor{node.index, d}
			worstDist = best[0].dist
			for _, b := range best {
				if b.dist > worstDist {
					worstDist = b.dist
				}
			}
		}

		// Decide which subtree to search first.
		diff := point[node.axis] - t.points[node.index][node.axis]
		var first, second *kdNode
		if diff < 0 {
			first, second = node.left, node.right
		} else {
			first, second = node.right, node.left
		}

		search(first)
		if len(best) < k || math.Abs(diff) < worstDist {
			search(second)
		}
	}

	search(t.root)

	// Sort results by distance.
	sort.Slice(best, func(i, j int) bool {
		return best[i].dist < best[j].dist
	})

	indices := make([]int, len(best))
	dists := make([]float64, len(best))
	for i, b := range best {
		indices[i] = b.index
		dists[i] = b.dist
	}
	return indices, dists
}

// ---------------------------------------------------------------------------
// Convex Hull (Gift Wrapping for 2D)
// ---------------------------------------------------------------------------

// ConvexHull computes the convex hull of a set of 2D points using the gift
// wrapping (Jarvis march) algorithm. Returns the indices of the hull vertices
// in counter-clockwise order.
func ConvexHull(points [][]float64) ([]int, error) {
	n := len(points)
	if n < 3 {
		if n == 0 {
			return nil, errors.New("scigo.ConvexHull: empty point set")
		}
		indices := make([]int, n)
		for i := range indices {
			indices[i] = i
		}
		return indices, nil
	}

	for _, p := range points {
		if len(p) < 2 {
			return nil, errors.New("scigo.ConvexHull: points must be 2D")
		}
	}

	// Find leftmost point.
	leftmost := 0
	for i := 1; i < n; i++ {
		if points[i][0] < points[leftmost][0] ||
			(points[i][0] == points[leftmost][0] && points[i][1] < points[leftmost][1]) {
			leftmost = i
		}
	}

	hull := []int{}
	p := leftmost
	for {
		hull = append(hull, p)
		q := (p + 1) % n
		for r := 0; r < n; r++ {
			if cross2D(points[p], points[q], points[r]) < 0 {
				q = r
			}
		}
		p = q
		if p == leftmost {
			break
		}
		if len(hull) > n {
			break // Safety.
		}
	}

	return hull, nil
}

// cross2D returns the cross product of vectors (q-p) and (r-p).
// Negative means r is to the left of p->q (counter-clockwise).
func cross2D(p, q, r []float64) float64 {
	return (q[0]-p[0])*(r[1]-p[1]) - (q[1]-p[1])*(r[0]-p[0])
}

// ---------------------------------------------------------------------------
// Delaunay Triangulation (Bowyer-Watson for 2D)
// ---------------------------------------------------------------------------

// Delaunay computes a Delaunay triangulation of a set of 2D points.
// Returns a list of triangles, where each triangle is a slice of 3 point indices.
func Delaunay(points [][]float64) ([][]int, error) {
	n := len(points)
	if n < 3 {
		return nil, errors.New("scigo.Delaunay: need at least 3 points")
	}
	for _, p := range points {
		if len(p) < 2 {
			return nil, errors.New("scigo.Delaunay: points must be 2D")
		}
	}

	// Create a super-triangle that contains all points.
	minX, minY := points[0][0], points[0][1]
	maxX, maxY := minX, minY
	for _, p := range points {
		if p[0] < minX {
			minX = p[0]
		}
		if p[0] > maxX {
			maxX = p[0]
		}
		if p[1] < minY {
			minY = p[1]
		}
		if p[1] > maxY {
			maxY = p[1]
		}
	}

	dx := maxX - minX
	dy := maxY - minY
	dmax := dx
	if dy > dmax {
		dmax = dy
	}
	midX := (minX + maxX) / 2
	midY := (minY + maxY) / 2

	// Super-triangle vertices (indices n, n+1, n+2).
	superPts := [][]float64{
		{midX - 2*dmax, midY - dmax},
		{midX + 2*dmax, midY - dmax},
		{midX, midY + 2*dmax},
	}

	allPts := make([][]float64, n+3)
	copy(allPts, points)
	allPts[n] = superPts[0]
	allPts[n+1] = superPts[1]
	allPts[n+2] = superPts[2]

	type triangle struct {
		a, b, c int
	}

	triangles := []triangle{{n, n + 1, n + 2}}

	for i := 0; i < n; i++ {
		px, py := allPts[i][0], allPts[i][1]

		// Find triangles whose circumcircle contains the new point.
		type edge struct{ a, b int }
		var badTris []int
		for ti, tri := range triangles {
			if inCircumcircle(allPts, tri.a, tri.b, tri.c, px, py) {
				badTris = append(badTris, ti)
			}
		}

		// Find boundary edges of the bad triangles.
		edgeCount := make(map[edge]int)
		for _, ti := range badTris {
			tri := triangles[ti]
			edges := []edge{
				{tri.a, tri.b}, {tri.b, tri.c}, {tri.c, tri.a},
			}
			for _, e := range edges {
				// Normalize edge direction.
				ne := e
				if ne.a > ne.b {
					ne.a, ne.b = ne.b, ne.a
				}
				edgeCount[ne]++
			}
		}

		// Remove bad triangles (in reverse order to keep indices valid).
		sort.Sort(sort.Reverse(sort.IntSlice(badTris)))
		for _, ti := range badTris {
			triangles[ti] = triangles[len(triangles)-1]
			triangles = triangles[:len(triangles)-1]
		}

		// Add new triangles from boundary edges to the new point.
		for e, count := range edgeCount {
			if count == 1 {
				triangles = append(triangles, triangle{e.a, e.b, i})
			}
		}
	}

	// Remove triangles that share vertices with the super-triangle.
	var result [][]int
	for _, tri := range triangles {
		if tri.a >= n || tri.b >= n || tri.c >= n {
			continue
		}
		result = append(result, []int{tri.a, tri.b, tri.c})
	}

	if len(result) == 0 {
		return nil, errors.New("scigo.Delaunay: degenerate triangulation")
	}

	return result, nil
}

func inCircumcircle(pts [][]float64, a, b, c int, px, py float64) bool {
	ax, ay := pts[a][0]-px, pts[a][1]-py
	bx, by := pts[b][0]-px, pts[b][1]-py
	cx, cy := pts[c][0]-px, pts[c][1]-py

	det := ax*(by*((cx*cx)+(cy*cy))-cy*((bx*bx)+(by*by))) -
		bx*(ay*((cx*cx)+(cy*cy))-cy*((ax*ax)+(ay*ay))) +
		cx*(ay*((bx*bx)+(by*by))-by*((ax*ax)+(ay*ay)))

	// The sign depends on the orientation of the triangle.
	// For CCW orientation, det > 0 means inside.
	orient := (pts[b][0]-pts[a][0])*(pts[c][1]-pts[a][1]) -
		(pts[b][1]-pts[a][1])*(pts[c][0]-pts[a][0])
	if orient < 0 {
		det = -det
	}
	return det > 0
}

// ---------------------------------------------------------------------------
// Voronoi (from Delaunay dual)
// ---------------------------------------------------------------------------

// Voronoi computes a Voronoi diagram from a set of 2D points.
// Returns the Voronoi vertices and a list of regions (each region is a list of vertex indices).
// This is a simplified implementation that derives the Voronoi from the Delaunay triangulation.
func Voronoi(points [][]float64) ([][]float64, [][]int, error) {
	triangles, err := Delaunay(points)
	if err != nil {
		return nil, nil, err
	}

	// Compute circumcenters of each triangle — these are the Voronoi vertices.
	vertices := make([][]float64, len(triangles))
	for i, tri := range triangles {
		vertices[i] = circumcenter(points[tri[0]], points[tri[1]], points[tri[2]])
	}

	// Build adjacency: for each input point, find which triangles it belongs to.
	n := len(points)
	pointTris := make([][]int, n)
	for ti, tri := range triangles {
		for _, pi := range tri {
			pointTris[pi] = append(pointTris[pi], ti)
		}
	}

	// For each input point, the Voronoi region is the polygon formed by
	// the circumcenters of its adjacent triangles, ordered angularly.
	regions := make([][]int, n)
	for pi := 0; pi < n; pi++ {
		tris := pointTris[pi]
		if len(tris) == 0 {
			regions[pi] = nil
			continue
		}

		// Sort triangles by angle of circumcenter relative to the point.
		cx, cy := points[pi][0], points[pi][1]
		sort.Slice(tris, func(i, j int) bool {
			ai := math.Atan2(vertices[tris[i]][1]-cy, vertices[tris[i]][0]-cx)
			aj := math.Atan2(vertices[tris[j]][1]-cy, vertices[tris[j]][0]-cx)
			return ai < aj
		})

		region := make([]int, len(tris))
		for i, ti := range tris {
			region[i] = ti
		}
		regions[pi] = region
	}

	return vertices, regions, nil
}

func circumcenter(a, b, c []float64) []float64 {
	ax, ay := a[0], a[1]
	bx, by := b[0], b[1]
	cx, cy := c[0], c[1]

	d := 2 * (ax*(by-cy) + bx*(cy-ay) + cx*(ay-by))
	if math.Abs(d) < 1e-14 {
		// Degenerate: return midpoint.
		return []float64{(ax + bx + cx) / 3, (ay + by + cy) / 3}
	}

	ux := ((ax*ax+ay*ay)*(by-cy) + (bx*bx+by*by)*(cy-ay) + (cx*cx+cy*cy)*(ay-by)) / d
	uy := ((ax*ax+ay*ay)*(cx-bx) + (bx*bx+by*by)*(ax-cx) + (cx*cx+cy*cy)*(bx-ax)) / d

	return []float64{ux, uy}
}

// ---------------------------------------------------------------------------
// Distance helpers
// ---------------------------------------------------------------------------

func euclidean(a, b []float64) float64 {
	s := 0.0
	for i := range a {
		if i < len(b) {
			d := a[i] - b[i]
			s += d * d
		}
	}
	return math.Sqrt(s)
}

func manhattan(a, b []float64) float64 {
	s := 0.0
	for i := range a {
		if i < len(b) {
			s += math.Abs(a[i] - b[i])
		}
	}
	return s
}

func chebyshev(a, b []float64) float64 {
	m := 0.0
	for i := range a {
		if i < len(b) {
			d := math.Abs(a[i] - b[i])
			if d > m {
				m = d
			}
		}
	}
	return m
}

func getDistFunc(metric string) func([]float64, []float64) float64 {
	switch metric {
	case "manhattan":
		return manhattan
	case "chebyshev":
		return chebyshev
	default:
		return euclidean
	}
}
