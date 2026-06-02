//go:build unit

package scigo

import (
	"math"
	"sort"
	"testing"
)

// ---------------------------------------------------------------------------
// Cdist Tests
// ---------------------------------------------------------------------------

func TestCdist(t *testing.T) {
	xa := [][]float64{{0, 0}, {1, 0}}
	xb := [][]float64{{0, 1}, {1, 1}}

	result := Cdist(xa, xb, "euclidean")
	if len(result) != 2 || len(result[0]) != 2 {
		t.Fatal("Cdist: wrong shape")
	}
	// d(0,0)-(0,1) = 1, d(0,0)-(1,1) = sqrt(2), d(1,0)-(0,1) = sqrt(2), d(1,0)-(1,1) = 1
	if !approxEqual(result[0][0], 1, 1e-14) {
		t.Errorf("Cdist[0][0] = %v, want 1", result[0][0])
	}
	if !approxEqual(result[0][1], math.Sqrt(2), 1e-14) {
		t.Errorf("Cdist[0][1] = %v, want sqrt(2)", result[0][1])
	}
	if !approxEqual(result[1][0], math.Sqrt(2), 1e-14) {
		t.Errorf("Cdist[1][0] = %v, want sqrt(2)", result[1][0])
	}
	if !approxEqual(result[1][1], 1, 1e-14) {
		t.Errorf("Cdist[1][1] = %v, want 1", result[1][1])
	}
}

func TestCdistManhattan(t *testing.T) {
	xa := [][]float64{{0, 0}}
	xb := [][]float64{{3, 4}}
	result := Cdist(xa, xb, "manhattan")
	if !approxEqual(result[0][0], 7, 1e-14) {
		t.Errorf("Cdist manhattan = %v, want 7", result[0][0])
	}
}

func TestCdistChebyshev(t *testing.T) {
	xa := [][]float64{{0, 0}}
	xb := [][]float64{{3, 5}}
	result := Cdist(xa, xb, "chebyshev")
	if !approxEqual(result[0][0], 5, 1e-14) {
		t.Errorf("Cdist chebyshev = %v, want 5", result[0][0])
	}
}

// ---------------------------------------------------------------------------
// Pdist Tests
// ---------------------------------------------------------------------------

func TestPdist(t *testing.T) {
	x := [][]float64{{0, 0}, {1, 0}, {0, 1}}
	result := Pdist(x, "euclidean")
	// Pairs: (0,1)=1, (0,2)=1, (1,2)=sqrt(2)
	if len(result) != 3 {
		t.Fatalf("Pdist: length = %d, want 3", len(result))
	}
	if !approxEqual(result[0], 1, 1e-14) {
		t.Errorf("Pdist[0] = %v, want 1", result[0])
	}
	if !approxEqual(result[1], 1, 1e-14) {
		t.Errorf("Pdist[1] = %v, want 1", result[1])
	}
	if !approxEqual(result[2], math.Sqrt(2), 1e-14) {
		t.Errorf("Pdist[2] = %v, want sqrt(2)", result[2])
	}
}

// ---------------------------------------------------------------------------
// Squareform Tests
// ---------------------------------------------------------------------------

func TestSquareform(t *testing.T) {
	d := []float64{1, 2, 3}
	result := Squareform(d, 3)
	if result == nil {
		t.Fatal("Squareform: nil result")
	}
	expected := [][]float64{
		{0, 1, 2},
		{1, 0, 3},
		{2, 3, 0},
	}
	for i := range expected {
		for j := range expected[i] {
			if result[i][j] != expected[i][j] {
				t.Errorf("Squareform[%d][%d] = %v, want %v", i, j, result[i][j], expected[i][j])
			}
		}
	}
}

func TestSquareformInvalidSize(t *testing.T) {
	result := Squareform([]float64{1, 2}, 3) // 3 points need 3 distances
	if result != nil {
		t.Error("Squareform: expected nil for wrong size")
	}
}

// ---------------------------------------------------------------------------
// KDTree Tests
// ---------------------------------------------------------------------------

func TestKDTreeBasic(t *testing.T) {
	points := [][]float64{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
		{0.5, 0.5},
	}
	tree := NewKDTree(points)

	// Query nearest neighbor to (0.1, 0.1) — should be (0, 0).
	indices, dists := tree.Query([]float64{0.1, 0.1}, 1)
	if len(indices) != 1 {
		t.Fatalf("KDTree.Query: got %d results, want 1", len(indices))
	}
	if indices[0] != 0 {
		t.Errorf("KDTree.Query: nearest = %d, want 0", indices[0])
	}
	expectedDist := math.Sqrt(0.01 + 0.01)
	if !approxEqual(dists[0], expectedDist, 1e-10) {
		t.Errorf("KDTree.Query: dist = %v, want %v", dists[0], expectedDist)
	}
}

func TestKDTreeKNN(t *testing.T) {
	points := [][]float64{
		{0, 0},
		{1, 0},
		{0, 1},
		{1, 1},
		{0.5, 0.5},
	}
	tree := NewKDTree(points)

	indices, dists := tree.Query([]float64{0.5, 0.5}, 3)
	if len(indices) != 3 {
		t.Fatalf("KDTree.Query: got %d results, want 3", len(indices))
	}

	// The nearest should be (0.5, 0.5) itself (index 4).
	if indices[0] != 4 {
		t.Errorf("KDTree.Query: nearest = %d, want 4", indices[0])
	}
	if !approxEqual(dists[0], 0, 1e-14) {
		t.Errorf("KDTree.Query: dist[0] = %v, want 0", dists[0])
	}

	// Next should be equidistant corners.
	for i := 1; i < 3; i++ {
		if !approxEqual(dists[i], math.Sqrt(0.5), 1e-10) {
			t.Errorf("KDTree.Query: dist[%d] = %v, want %v", i, dists[i], math.Sqrt(0.5))
		}
	}
}

func TestKDTreeEmpty(t *testing.T) {
	tree := NewKDTree(nil)
	indices, dists := tree.Query([]float64{0, 0}, 1)
	if len(indices) != 0 || len(dists) != 0 {
		t.Error("KDTree.Query on empty tree should return empty")
	}
}

func TestKDTree3D(t *testing.T) {
	points := [][]float64{
		{0, 0, 0},
		{1, 1, 1},
		{2, 2, 2},
	}
	tree := NewKDTree(points)
	indices, _ := tree.Query([]float64{0.9, 0.9, 0.9}, 1)
	if len(indices) != 1 || indices[0] != 1 {
		t.Errorf("KDTree 3D: nearest = %v, want [1]", indices)
	}
}

// ---------------------------------------------------------------------------
// ConvexHull Tests
// ---------------------------------------------------------------------------

func TestConvexHullSquare(t *testing.T) {
	points := [][]float64{
		{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0.5, 0.5},
	}
	hull, err := ConvexHull(points)
	if err != nil {
		t.Fatalf("ConvexHull: %v", err)
	}

	// Should contain exactly 4 hull vertices (the corners).
	if len(hull) != 4 {
		t.Errorf("ConvexHull: got %d vertices, want 4", len(hull))
	}

	// The interior point (0.5, 0.5) at index 4 should not be in the hull.
	for _, idx := range hull {
		if idx == 4 {
			t.Error("ConvexHull: interior point should not be on hull")
		}
	}
}

func TestConvexHullTriangle(t *testing.T) {
	points := [][]float64{
		{0, 0}, {2, 0}, {1, 2},
	}
	hull, err := ConvexHull(points)
	if err != nil {
		t.Fatalf("ConvexHull: %v", err)
	}
	if len(hull) != 3 {
		t.Errorf("ConvexHull triangle: got %d vertices, want 3", len(hull))
	}
}

func TestConvexHullCollinear(t *testing.T) {
	points := [][]float64{
		{0, 0}, {1, 0}, {2, 0},
	}
	hull, err := ConvexHull(points)
	if err != nil {
		t.Fatalf("ConvexHull: %v", err)
	}
	// Collinear points: hull should have 2 endpoints.
	if len(hull) < 2 {
		t.Error("ConvexHull collinear: expected at least 2 points")
	}
}

func TestConvexHullEmpty(t *testing.T) {
	_, err := ConvexHull(nil)
	if err == nil {
		t.Error("ConvexHull: expected error for empty input")
	}
}

// ---------------------------------------------------------------------------
// Delaunay Tests
// ---------------------------------------------------------------------------

func TestDelaunayTriangle(t *testing.T) {
	points := [][]float64{
		{0, 0}, {1, 0}, {0, 1},
	}
	tris, err := Delaunay(points)
	if err != nil {
		t.Fatalf("Delaunay: %v", err)
	}
	if len(tris) != 1 {
		t.Errorf("Delaunay triangle: got %d triangles, want 1", len(tris))
	}

	// The single triangle should contain all three points.
	if len(tris) > 0 {
		tri := tris[0]
		sort.Ints(tri)
		if tri[0] != 0 || tri[1] != 1 || tri[2] != 2 {
			t.Errorf("Delaunay triangle: vertices = %v, want [0,1,2]", tri)
		}
	}
}

func TestDelaunaySquare(t *testing.T) {
	points := [][]float64{
		{0, 0}, {1, 0}, {1, 1}, {0, 1},
	}
	tris, err := Delaunay(points)
	if err != nil {
		t.Fatalf("Delaunay: %v", err)
	}
	if len(tris) != 2 {
		t.Errorf("Delaunay square: got %d triangles, want 2", len(tris))
	}
}

func TestDelaunayFivePoints(t *testing.T) {
	points := [][]float64{
		{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0.5, 0.5},
	}
	tris, err := Delaunay(points)
	if err != nil {
		t.Fatalf("Delaunay: %v", err)
	}
	if len(tris) < 4 {
		t.Errorf("Delaunay 5 points: got %d triangles, expected >= 4", len(tris))
	}
}

func TestDelaunayTooFew(t *testing.T) {
	_, err := Delaunay([][]float64{{0, 0}, {1, 0}})
	if err == nil {
		t.Error("Delaunay: expected error for 2 points")
	}
}

// ---------------------------------------------------------------------------
// Voronoi Tests
// ---------------------------------------------------------------------------

func TestVoronoi(t *testing.T) {
	points := [][]float64{
		{0, 0}, {1, 0}, {0.5, 1},
	}
	vertices, regions, err := Voronoi(points)
	if err != nil {
		t.Fatalf("Voronoi: %v", err)
	}
	if len(vertices) == 0 {
		t.Error("Voronoi: no vertices")
	}
	if len(regions) != len(points) {
		t.Errorf("Voronoi: %d regions, want %d", len(regions), len(points))
	}
}

func TestVoronoiSquare(t *testing.T) {
	points := [][]float64{
		{0, 0}, {2, 0}, {2, 2}, {0, 2},
	}
	vertices, regions, err := Voronoi(points)
	if err != nil {
		t.Fatalf("Voronoi: %v", err)
	}
	_ = vertices
	// Each corner point should be in some region.
	for i, r := range regions {
		if len(r) == 0 {
			t.Errorf("Voronoi: region %d is empty", i)
		}
	}
}
