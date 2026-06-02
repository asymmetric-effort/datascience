"""
Fixture generator for src/learning package.

Generates test cases by exercising pgmpy's parameter estimation (MLE and
BayesianEstimator) on a simple A->B network, capturing inputs and expected
outputs as fixture data.
"""

from pgmpy.models import DiscreteBayesianNetwork as BayesianNetwork
from pgmpy.factors.discrete import TabularCPD
from pgmpy.estimators import MaximumLikelihoodEstimator, BayesianEstimator
import numpy as np


def generate() -> list[dict]:
    """Generate test cases for the learning package."""
    test_cases = []

    # Build the shared network and data once.
    bn, data = _build_network_and_data()

    test_cases.append(_test_mle(bn, data))
    test_cases.append(_test_bayesian_bdeu(bn, data))

    return test_cases


def _build_network_and_data():
    """Build a simple A->B network with known CPDs and sample data."""
    bn = BayesianNetwork([("A", "B")])

    cpd_a = TabularCPD("A", 2, [[0.4], [0.6]])
    cpd_b = TabularCPD(
        "B", 3,
        [[0.2, 0.3],
         [0.3, 0.5],
         [0.5, 0.2]],
        evidence=["A"],
        evidence_card=[2],
    )
    bn.add_cpds(cpd_a, cpd_b)
    assert bn.check_model()

    data = bn.simulate(10000, seed=42)
    return bn, data


def _cpd_values_as_list(cpd):
    """Extract CPD values as a nested list (rows=states, cols=parent configs)."""
    return cpd.get_values().tolist()


def _test_mle(bn, data):
    """Test MLE parameter estimation on A and B."""
    mle = MaximumLikelihoodEstimator(bn, data)

    cpd_a = mle.estimate_cpd("A")
    cpd_b = mle.estimate_cpd("B")

    return {
        "name": "mle_parameter_estimation",
        "description": "MLE estimation on A->B network with 10000 samples",
        "input": {
            "edges": [["A", "B"]],
            "node_cards": {"A": 2, "B": 3},
            "data_columns": ["A", "B"],
            "data_rows": data[["A", "B"]].values.tolist(),
        },
        "expected": {
            "cpd_A": {
                "variable": "A",
                "variable_card": 2,
                "values": _cpd_values_as_list(cpd_a),
                "evidence": [],
                "evidence_card": [],
            },
            "cpd_B": {
                "variable": "B",
                "variable_card": 3,
                "values": _cpd_values_as_list(cpd_b),
                "evidence": ["A"],
                "evidence_card": [2],
            },
        },
    }


def _test_bayesian_bdeu(bn, data):
    """Test BayesianEstimator with BDeu prior on A and B."""
    be = BayesianEstimator(bn, data)

    cpd_a = be.estimate_cpd("A", prior_type="BDeu", equivalent_sample_size=10)
    cpd_b = be.estimate_cpd("B", prior_type="BDeu", equivalent_sample_size=10)

    return {
        "name": "bayesian_bdeu_parameter_estimation",
        "description": "Bayesian BDeu estimation on A->B network with ESS=10",
        "input": {
            "edges": [["A", "B"]],
            "node_cards": {"A": 2, "B": 3},
            "equivalent_sample_size": 10,
            "data_columns": ["A", "B"],
            "data_rows": data[["A", "B"]].values.tolist(),
        },
        "expected": {
            "cpd_A": {
                "variable": "A",
                "variable_card": 2,
                "values": _cpd_values_as_list(cpd_a),
                "evidence": [],
                "evidence_card": [],
            },
            "cpd_B": {
                "variable": "B",
                "variable_card": 3,
                "values": _cpd_values_as_list(cpd_b),
                "evidence": ["A"],
                "evidence_card": [2],
            },
        },
    }
