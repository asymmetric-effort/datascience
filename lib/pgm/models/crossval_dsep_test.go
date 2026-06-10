//go:build unit

package models_test

import (
	"testing"

	"github.com/asymmetric-effort/datascience/lib/graphgo"
	"github.com/asymmetric-effort/datascience/tests/testutil"
)

func TestCrossval_DSeparation(t *testing.T) {
	ff := testutil.LoadFixtures(t, "dseparation/fixtures.json")

	for _, tc := range ff.TestCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			var input struct {
				Edges [][]string `json:"edges"`
				X     []string   `json:"x"`
				Y     []string   `json:"y"`
				Z     []string   `json:"z"`
			}
			tc.UnmarshalInput(t, &input)

			var expected struct {
				DSeparated bool `json:"d_separated"`
			}
			tc.UnmarshalExpected(t, &expected)

			// Build a DiGraph from the edges.
			g := graphgo.NewDiGraph()
			for _, edge := range input.Edges {
				g.AddEdge(edge[0], edge[1])
			}

			// Build set arguments for DSeparation.
			xSet := make(map[string]bool, len(input.X))
			for _, n := range input.X {
				xSet[n] = true
			}
			ySet := make(map[string]bool, len(input.Y))
			for _, n := range input.Y {
				ySet[n] = true
			}
			zSet := make(map[string]bool, len(input.Z))
			for _, n := range input.Z {
				zSet[n] = true
			}

			got := graphgo.DSeparation(g, xSet, ySet, zSet)
			if got != expected.DSeparated {
				t.Errorf("DSeparation(x=%v, y=%v, z=%v): expected %v, got %v",
					input.X, input.Y, input.Z, expected.DSeparated, got)
			}
		})
	}
}
