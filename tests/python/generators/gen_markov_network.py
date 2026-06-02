"""
Fixture generator for MarkovNetwork (MRF) cross-validation.

Builds a simple 3-node triangle MRF (A-B-C) with factors, computes the
partition function, and captures structure and factor data as fixtures.
"""

from pgmpy.models import DiscreteMarkovNetwork
from pgmpy.factors.discrete import DiscreteFactor
import numpy as np


def generate() -> list[dict]:
    """Generate test cases for the Markov network package."""
    test_cases = []

    test_cases.append(_test_markov_network_triangle())
    test_cases.append(_test_markov_network_chain())

    return test_cases


def _test_markov_network_triangle():
    """Test a 3-node triangle MRF with pairwise factors."""
    mn = DiscreteMarkovNetwork([("A", "B"), ("B", "C"), ("A", "C")])

    # Pairwise factors over binary variables.
    f_ab = DiscreteFactor(["A", "B"], [2, 2], [10, 1, 1, 10])
    f_bc = DiscreteFactor(["B", "C"], [2, 2], [5, 1, 1, 5])
    f_ac = DiscreteFactor(["A", "C"], [2, 2], [3, 2, 2, 3])

    mn.add_factors(f_ab, f_bc, f_ac)
    is_valid = mn.check_model()
    Z = mn.get_partition_function()

    return {
        "name": "markov_network_triangle",
        "description": "3-node triangle MRF (A-B-C) with pairwise factors",
        "input": {
            "edges": [["A", "B"], ["B", "C"], ["A", "C"]],
            "factors": [
                {
                    "variables": ["A", "B"],
                    "cardinality": [2, 2],
                    "values": [10.0, 1.0, 1.0, 10.0],
                },
                {
                    "variables": ["B", "C"],
                    "cardinality": [2, 2],
                    "values": [5.0, 1.0, 1.0, 5.0],
                },
                {
                    "variables": ["A", "C"],
                    "cardinality": [2, 2],
                    "values": [3.0, 2.0, 2.0, 3.0],
                },
            ],
        },
        "expected": {
            "nodes": sorted(mn.nodes()),
            "edges": sorted([sorted(list(e)) for e in mn.edges()]),
            "num_nodes": len(mn.nodes()),
            "num_edges": len(mn.edges()),
            "is_valid": is_valid,
            "partition_function": float(Z),
            "num_factors": 3,
        },
    }


def _test_markov_network_chain():
    """Test a 3-node chain MRF (X-Y-Z) with pairwise factors."""
    mn = DiscreteMarkovNetwork([("X", "Y"), ("Y", "Z")])

    f_xy = DiscreteFactor(["X", "Y"], [2, 2], [8, 2, 2, 8])
    f_yz = DiscreteFactor(["Y", "Z"], [2, 2], [6, 1, 1, 6])

    mn.add_factors(f_xy, f_yz)
    is_valid = mn.check_model()
    Z = mn.get_partition_function()

    return {
        "name": "markov_network_chain",
        "description": "3-node chain MRF (X-Y-Z) with pairwise factors",
        "input": {
            "edges": [["X", "Y"], ["Y", "Z"]],
            "factors": [
                {
                    "variables": ["X", "Y"],
                    "cardinality": [2, 2],
                    "values": [8.0, 2.0, 2.0, 8.0],
                },
                {
                    "variables": ["Y", "Z"],
                    "cardinality": [2, 2],
                    "values": [6.0, 1.0, 1.0, 6.0],
                },
            ],
        },
        "expected": {
            "nodes": sorted(mn.nodes()),
            "edges": sorted([sorted(list(e)) for e in mn.edges()]),
            "num_nodes": len(mn.nodes()),
            "num_edges": len(mn.edges()),
            "is_valid": is_valid,
            "partition_function": float(Z),
            "num_factors": 2,
        },
    }
