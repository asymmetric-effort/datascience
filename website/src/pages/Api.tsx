import { createElement, useHead } from "@asymmetric-effort/specifyjs";

export function Api() {
  useHead({
    title: "API Reference — pgmgo",
    description: "Go API reference for pgmgo probabilistic graphical models library.",
    canonical: "https://pgmgo.asymmetric-effort.com/#/api",
  });

  return (
    <div class="page">
      <h1>API Reference</h1>

      <section class="section">
        <h2>Import</h2>
        <pre><code>{`import (
    "github.com/asymmetric-effort/pgmgo/src/models"
    "github.com/asymmetric-effort/pgmgo/src/factors"
    "github.com/asymmetric-effort/pgmgo/src/inference"
)`}</code></pre>
      </section>

      <section class="section">
        <h2>Library Packages (src/)</h2>
        <table>
          <thead>
            <tr><th>Package</th><th>Key Types</th></tr>
          </thead>
          <tbody>
            <tr><td><code>base</code></td><td>DAG, PDAG, UndirectedGraph, ADMG, MAG, SimpleCausalModel</td></tr>
            <tr><td><code>models</code></td><td>BayesianNetwork, MarkovNetwork, DynamicBayesianNetwork, NaiveBayes, SEM, FactorGraph, JunctionTree</td></tr>
            <tr><td><code>factors</code></td><td>DiscreteFactor, TabularCPD, JointProbabilityDistribution, LinearGaussianCPD, NoisyOR</td></tr>
            <tr><td><code>inference</code></td><td>VariableElimination, BeliefPropagation, MPLP, CausalInference, DBNInference</td></tr>
            <tr><td><code>sampling</code></td><td>BayesianModelSampling, GibbsSampling</td></tr>
            <tr><td><code>learning</code></td><td>MLE, BayesianEstimator, EM, HillClimbSearch, PC, GES, ExhaustiveSearch</td></tr>
            <tr><td><code>ci_tests</code></td><td>ChiSquare, GSq, FisherZ, Pearsonr, GCM, HotellingLawley</td></tr>
            <tr><td><code>structure_score</code></td><td>BIC, AIC, BDeu, BDs, K2, LogLikelihood</td></tr>
            <tr><td><code>identification</code></td><td>Adjustment, Frontdoor</td></tr>
            <tr><td><code>prediction</code></td><td>DoubleMLRegressor, NaiveAdjustmentRegressor, NaiveIVRegressor</td></tr>
            <tr><td><code>metrics</code></td><td>SHD, AdjacencyConfusionMatrix, OrientationConfusionMatrix</td></tr>
            <tr><td><code>readwrite</code></td><td>BIFReader/Writer, XMLBIFReader/Writer, NETReader/Writer, UAIReader/Writer</td></tr>
          </tbody>
        </table>
      </section>

      <section class="section">
        <h2>Primitive Modules (lib/)</h2>
        <table>
          <thead>
            <tr><th>Module</th><th>Key Types</th></tr>
          </thead>
          <tbody>
            <tr><td><code>numgo</code></td><td>NDArray, Matrix, Vector, linear algebra operations</td></tr>
            <tr><td><code>scigo</code></td><td>Distributions, optimization, statistical tests</td></tr>
            <tr><td><code>graphgo</code></td><td>Graph, DiGraph, topological sort, d-separation, clique finding</td></tr>
            <tr><td><code>tabgo</code></td><td>DataFrame, Series, groupby, filtering, CSV I/O</td></tr>
          </tbody>
        </table>
      </section>

      <section class="section">
        <h2>Example: Bayesian Network</h2>
        <pre><code>{`bn := models.NewBayesianNetwork()
bn.AddNode("A")
bn.AddNode("B")
bn.AddNode("C")
bn.AddEdge("A", "B")
bn.AddEdge("B", "C")

// Assign CPDs
bn.SetCPD("A", factors.NewDiscreteFactor(
    []string{"A"},
    []float64{0.4, 0.6},
))

// Run inference
ve := inference.NewVariableElimination(bn)
result := ve.Query([]string{"C"}, map[string]int{"A": 0})`}</code></pre>
      </section>

      <section class="section">
        <h2>Example: Structure Learning</h2>
        <pre><code>{`data := tabgo.ReadCSV("observations.csv")
learner := learning.NewHillClimbSearch(data)
model := learner.Estimate()`}</code></pre>
      </section>
    </div>
  );
}
