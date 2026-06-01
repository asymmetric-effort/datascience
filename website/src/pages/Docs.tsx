import { createElement, useHead } from "@asymmetric-effort/specifyjs";

export function Docs() {
  useHead({
    title: "Documentation — pgmgo",
    description: "Documentation for pgmgo, a Go library for probabilistic graphical models.",
    canonical: "https://pgmgo.asymmetric-effort.com/#/docs",
  });

  return (
    <div class="page">
      <h1>Documentation</h1>

      <section class="section">
        <h2>Overview</h2>
        <p>
          pgmgo is a Go library for working with probabilistic graphical models
          with 100% feature parity with <a href="https://pgmpy.org" target="_blank" rel="noopener noreferrer">pgmpy</a>.
          It provides tools for creating, parameterizing, and performing inference
          on Bayesian networks, Markov networks, and related structures with zero
          third-party runtime dependencies.
        </p>
      </section>

      <section class="section">
        <h2>Library Packages (src/)</h2>
        <table>
          <thead>
            <tr><th>Package</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>base</code></td><td>Foundational graph types: DAG, PDAG, UndirectedGraph, ADMG, MAG</td></tr>
            <tr><td><code>models</code></td><td>BayesianNetwork, MarkovNetwork, DynamicBN, NaiveBayes, SEM, and more</td></tr>
            <tr><td><code>factors</code></td><td>DiscreteFactor, TabularCPD, LinearGaussianCPD, NoisyOR, factor operations</td></tr>
            <tr><td><code>inference</code></td><td>VariableElimination, BeliefPropagation, MPLP, CausalInference</td></tr>
            <tr><td><code>sampling</code></td><td>Forward, rejection, likelihood-weighted, and Gibbs sampling</td></tr>
            <tr><td><code>learning</code></td><td>MLE, BayesianEstimator, EM, structure learning (PC, GES, HillClimb)</td></tr>
            <tr><td><code>ci_tests</code></td><td>Conditional independence tests: ChiSquare, FisherZ, GCM, and more</td></tr>
            <tr><td><code>structure_score</code></td><td>BIC, AIC, BDeu, BDs, K2, log-likelihood scoring functions</td></tr>
            <tr><td><code>identification</code></td><td>Causal effect identification: back-door, front-door</td></tr>
            <tr><td><code>prediction</code></td><td>DoubleML, naive adjustment, instrumental variable regression</td></tr>
            <tr><td><code>metrics</code></td><td>SHD, confusion matrices, correlation scores, Fisher's C</td></tr>
            <tr><td><code>independencies</code></td><td>Independence assertion representations</td></tr>
            <tr><td><code>readwrite</code></td><td>BIF, XMLBIF, NET, UAI, XDSL, PomdpX, XBN format I/O</td></tr>
            <tr><td><code>utils</code></td><td>Shared parsing, optimization, and compatibility utilities</td></tr>
          </tbody>
        </table>
      </section>

      <section class="section">
        <h2>Internal Primitive Modules (lib/)</h2>
        <table>
          <thead>
            <tr><th>Module</th><th>Replaces</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>numgo</code></td><td>numpy</td><td>N-dimensional arrays, linear algebra, numerical primitives</td></tr>
            <tr><td><code>scigo</code></td><td>scipy</td><td>Statistical distributions, optimization, special functions</td></tr>
            <tr><td><code>graphgo</code></td><td>networkx</td><td>Graph data structures, algorithms, d-separation, topological sort</td></tr>
            <tr><td><code>tabgo</code></td><td>pandas</td><td>Tabular data structures, loading, filtering, grouping</td></tr>
          </tbody>
        </table>
      </section>

      <section class="section">
        <h2>Project Structure</h2>
        <pre><code>{`pgmgo/
  cmd/             # CLI tools
    pgmgo/         # Main CLI entry point
  lib/             # Internal primitive modules
    numgo/         # numpy equivalent
    scigo/         # scipy equivalent
    graphgo/       # networkx equivalent
    tabgo/         # pandas equivalent
  src/             # pgmgo library code
    base/          # Foundational graph types
    models/        # Probabilistic model classes
    factors/       # Factor representations
    inference/     # Inference algorithms
    sampling/      # Sampling methods
    learning/      # Parameter and structure learning
    ci_tests/      # Conditional independence tests
    structure_score/ # Scoring functions
    identification/  # Causal identification
    prediction/    # Causal prediction
    metrics/       # Model evaluation
    independencies/  # Independence relations
    readwrite/     # File format I/O
    utils/         # Shared utilities
  docs/            # Documentation
  website/         # Project website`}</code></pre>
      </section>

      <section class="section">
        <h2>Datasets and Example Models</h2>
        <p>
          pgmgo does not bundle datasets or example networks. For reference data, use the
          datasets and example models provided by pgmpy:
        </p>
        <ul>
          <li><a href="https://pgmpy.org/models/bayesiannetwork.html#module-pgmpy.utils" target="_blank" rel="noopener noreferrer">pgmpy Datasets</a></li>
          <li><a href="https://pgmpy.org/examples.html" target="_blank" rel="noopener noreferrer">pgmpy Example Models</a></li>
        </ul>
        <p>
          pgmgo's readwrite package can load models in BIF, XMLBIF, NET, UAI, XDSL,
          PomdpX, and XBN formats exported from pgmpy or other tools.
        </p>
      </section>
    </div>
  );
}
