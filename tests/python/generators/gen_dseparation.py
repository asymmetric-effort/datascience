"""
Fixture generator for d-separation cross-validation.

Uses the student BN (D->G, I->G, G->L, I->S) to test d-separation queries
via pgmpy's active_trail_nodes method. D-separation holds when the target
node is NOT in the active trail from the source.
"""

from pgmpy.models import DiscreteBayesianNetwork


def generate() -> list[dict]:
    """Generate test cases for d-separation."""
    test_cases = []

    bn = DiscreteBayesianNetwork([("D", "G"), ("I", "G"), ("G", "L"), ("I", "S")])

    # Each query: (x, y, z, description)
    queries = [
        # Independent causes: D and I are d-separated with no observations.
        {
            "x": ["D"],
            "y": ["I"],
            "z": [],
            "description": "D _|_ I | {} (independent causes, d-separated)",
        },
        # V-structure activated: conditioning on G opens the path D-G-I.
        {
            "x": ["D"],
            "y": ["I"],
            "z": ["G"],
            "description": "D _|_ I | {G} (v-structure active, NOT d-separated)",
        },
        # D and S have no undirected path not through I or G.
        {
            "x": ["D"],
            "y": ["S"],
            "z": [],
            "description": "D _|_ S | {} (no active trail, d-separated)",
        },
        # D -> G -> L is an active trail with no observations.
        {
            "x": ["D"],
            "y": ["L"],
            "z": [],
            "description": "D _|_ L | {} (active trail D->G->L, NOT d-separated)",
        },
        # Conditioning on G blocks D -> G -> L.
        {
            "x": ["L"],
            "y": ["D"],
            "z": ["G"],
            "description": "L _|_ D | {G} (G blocks path, d-separated)",
        },
        # Conditioning on G blocks L-G-I path.
        {
            "x": ["L"],
            "y": ["S"],
            "z": ["G"],
            "description": "L _|_ S | {G} (G blocks all paths, d-separated)",
        },
        # I -> S is an active trail.
        {
            "x": ["I"],
            "y": ["S"],
            "z": [],
            "description": "I _|_ S | {} (active trail I->S, NOT d-separated)",
        },
        # Conditioning on I blocks the only path I->S.
        {
            "x": ["D"],
            "y": ["S"],
            "z": ["I"],
            "description": "D _|_ S | {I} (I blocks, d-separated)",
        },
    ]

    for q in queries:
        x_node = q["x"][0]
        z_nodes = q["z"]
        y_node = q["y"][0]

        trail = bn.active_trail_nodes(x_node, observed=z_nodes)
        is_reachable = y_node in trail[x_node]
        d_separated = not is_reachable

        test_cases.append({
            "name": f"dsep_{'_'.join(q['x'])}_{'_'.join(q['y'])}_given_{'_'.join(q['z']) if q['z'] else 'empty'}",
            "description": q["description"],
            "input": {
                "edges": [["D", "G"], ["I", "G"], ["G", "L"], ["I", "S"]],
                "x": q["x"],
                "y": q["y"],
                "z": q["z"],
            },
            "expected": {
                "d_separated": d_separated,
            },
        })

    return test_cases
