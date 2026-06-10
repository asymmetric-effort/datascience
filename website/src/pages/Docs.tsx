import { createElement, useHead, Link } from "@asymmetric-effort/specifyjs";
import { ScrollLink } from "../components/ScrollLink";

export function Docs() {
  useHead({
    title: "Documentation — datascience",
    description: "Comprehensive documentation for datascience, a pure Go library for data science and machine learning. PGMs, TensorFlow, BLAS, financial modeling, and more.",
    canonical: "https://pgmgo.asymmetric-effort.com/#/docs",
  });

  return (
    <div class="page">
      <h1>Documentation</h1>

      <nav class="page-toc">
        <strong>On this page:</strong>{" "}
        <ScrollLink to="overview">Overview</ScrollLink> | <ScrollLink to="installation">Installation</ScrollLink> | <ScrollLink to="getting-started">Getting Started</ScrollLink> | <ScrollLink to="architecture">Architecture</ScrollLink> | <ScrollLink to="library-packages">Library Packages</ScrollLink> | <ScrollLink to="pgm-packages">PGM Packages</ScrollLink> | <ScrollLink to="file-formats">File Formats</ScrollLink> | <ScrollLink to="example-models">Example Models</ScrollLink> | <ScrollLink to="datasets">Datasets</ScrollLink> | <ScrollLink to="testing">Testing</ScrollLink> | <ScrollLink to="configuration">Configuration</ScrollLink> | <ScrollLink to="contributing">Contributing</ScrollLink>
      </nav>

      {/* ============================================================ */}
      {/* OVERVIEW */}
      {/* ============================================================ */}
      <section class="section" id="overview">
        <h2>Overview</h2>
        <p>
          <strong>datascience</strong> is a comprehensive pure Go library for data science and machine learning.
          It provides tools for probabilistic graphical models (PGMs), TensorFlow-compatible deep learning,
          BLAS linear algebra, financial modeling, statistical computing, and more.
        </p>
        <p>
          All numerical primitives -- linear algebra, statistical distributions, graph algorithms, and tabular data
          processing -- are implemented from scratch in pure Go. There are no C bindings, no cgo, no third-party
          module dependencies. The result is a library that compiles to a single static binary with <code>go build</code>,
          cross-compiles trivially, and deploys without runtime dependencies.
        </p>
        <p>
          The current release is <strong>v0.0.37</strong> with approximately 5,000 tests and 392 cross-validation
          fixtures covering inference, learning, sampling, serialization, and cross-validation across the library.
        </p>

        <h3>Key Capabilities</h3>
        <ul>
          <li><strong>Probabilistic Graphical Models</strong>: 13 model types, 7 inference algorithms, 15+ learning algorithms, 16 CI tests, 13 scoring functions (via <code>lib/pgm</code>)</li>
          <li><strong>TensorFlow / Deep Learning</strong>: Neural network construction, training, and inference (via <code>lib/tensorflow</code>)</li>
          <li><strong>BLAS Linear Algebra</strong>: N-dimensional arrays, matrices, vectors, decompositions (via <code>lib/numgo</code>)</li>
          <li><strong>Statistical Computing</strong>: Distributions, optimization, hypothesis tests, special functions (via <code>lib/scigo</code>)</li>
          <li><strong>Graph Algorithms</strong>: Directed/undirected graphs, topological sort, d-separation, connected components (via <code>lib/graphgo</code>)</li>
          <li><strong>Tabular Data</strong>: DataFrames, Series, CSV/Parquet/Excel I/O (via <code>lib/tabgo</code>)</li>
          <li><strong>GPU Acceleration</strong>: Optional GPU compute backend for large-scale operations (via <code>lib/gpu</code>)</li>
          <li><strong>10 file formats</strong>: BIF, XMLBIF, NET, UAI, XDSL, PomdpX, XBN, CSV, JSON, XML</li>
          <li><strong>25 built-in example models</strong> and <strong>40 built-in datasets</strong></li>
          <li><strong>LLM integration</strong> for expert-in-the-loop structure learning</li>
        </ul>

        <h3>Relationship to pgmpy</h3>
        <p>
          The PGM module (<code>lib/pgm</code>) is inspired by <a href="https://pgmpy.org" target="_blank" rel="noopener noreferrer">pgmpy</a> and
          follows similar API patterns where possible. If you have used pgmpy in Python, the concepts and workflows
          will feel familiar. However, datascience is not a direct port -- it is a ground-up reimplementation
          in Go with its own design decisions, particularly around the near-zero-dependency philosophy.
        </p>
        <p>
          Where pgmpy relies on numpy, scipy, networkx, and pandas, datascience provides its own equivalents:
          <code>numgo</code>, <code>scigo</code>, <code>graphgo</code>, and <code>tabgo</code>. These are general-purpose
          libraries that can be used independently of the PGM layer.
        </p>
      </section>

      {/* ============================================================ */}
      {/* INSTALLATION */}
      {/* ============================================================ */}
      <section class="section" id="installation">
        <h2>Installation</h2>

        <h3>Go Module</h3>
        <pre><code>go get github.com/asymmetric-effort/datascience</code></pre>
        <p>
          Requires <strong>Go 1.21</strong> or later. No C compiler needed, no system libraries needed.
          Works on Linux, macOS, and Windows. Cross-compilation works out of the box.
        </p>

        <h3>From Source</h3>
        <pre><code>{`git clone https://github.com/asymmetric-effort/datascience.git
cd datascience
go test ./...`}</code></pre>

        <h3>Verify</h3>
        <pre><code>{`# Run all tests to verify the installation
go test ./...

# Run a quick sanity check
go test ./lib/pgm/models/... -run TestBayesianNetwork -v`}</code></pre>
      </section>

      {/* ============================================================ */}
      {/* GETTING STARTED */}
      {/* ============================================================ */}
      <section class="section" id="getting-started">
        <h2>Getting Started</h2>

        <h3>Your First Bayesian Network</h3>
        <p>
          A Bayesian network is a directed acyclic graph (DAG) where nodes represent random variables and
          edges represent conditional dependencies. Here is a complete, runnable program that creates the
          classic "wet grass" network:
        </p>
        <pre><code>{`package main

import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/lib/pgm/factors"
    "github.com/asymmetric-effort/datascience/lib/pgm/models"
)

func main() {
    // Create an empty Bayesian network
    bn := models.NewBayesianNetwork()

    // Add nodes (random variables)
    bn.AddNode("Cloudy")
    bn.AddNode("Rain")
    bn.AddNode("Sprinkler")
    bn.AddNode("WetGrass")

    // Add directed edges (causal relationships)
    bn.AddEdge("Cloudy", "Rain")
    bn.AddEdge("Cloudy", "Sprinkler")
    bn.AddEdge("Rain", "WetGrass")
    bn.AddEdge("Sprinkler", "WetGrass")

    // Define states for each variable
    bn.SetStates("Cloudy", []string{"clear", "cloudy"})
    bn.SetStates("Rain", []string{"no", "yes"})
    bn.SetStates("Sprinkler", []string{"off", "on"})
    bn.SetStates("WetGrass", []string{"dry", "wet"})

    // P(Cloudy) -- root node, marginal distribution
    bn.SetCPD("Cloudy", factors.NewTabularCPD(
        "Cloudy", 2,
        []float64{0.5, 0.5},
        nil, nil,
    ))

    // P(Rain | Cloudy) -- conditional distribution
    bn.SetCPD("Rain", factors.NewTabularCPD(
        "Rain", 2,
        []float64{0.8, 0.2, 0.2, 0.8},
        []string{"Cloudy"}, []int{2},
    ))

    // P(Sprinkler | Cloudy)
    bn.SetCPD("Sprinkler", factors.NewTabularCPD(
        "Sprinkler", 2,
        []float64{0.5, 0.9, 0.5, 0.1},
        []string{"Cloudy"}, []int{2},
    ))

    // P(WetGrass | Sprinkler, Rain)
    bn.SetCPD("WetGrass", factors.NewTabularCPD(
        "WetGrass", 2,
        []float64{
            1.0, 0.1, 0.1, 0.01,
            0.0, 0.9, 0.9, 0.99,
        },
        []string{"Sprinkler", "Rain"}, []int{2, 2},
    ))

    // Validate: checks DAG, CPD dimensions, probability sums
    if err := bn.CheckModel(); err != nil {
        log.Fatal("Model error:", err)
    }
    fmt.Println("Model is valid!")
    fmt.Printf("Nodes: %d, Edges: %d\\n", len(bn.Nodes()), len(bn.Edges()))
}`}</code></pre>

        <h3>Your First Query</h3>
        <p>
          Once you have a valid model, convert its CPDs to Markov factors and use Variable Elimination
          to compute posterior probabilities:
        </p>
        <pre><code>{`import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/example_models"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
)

func main() {
    // Load a built-in model (13 models have full CPDs)
    bn, _ := example_models.Get("asia")

    // Convert CPDs to Markov factors
    facs, err := bn.ToMarkovFactors()
    if err != nil {
        log.Fatal(err)
    }

    // Create a Variable Elimination engine
    ve := inference.NewVariableElimination(facs)

    // Posterior query: P(Dyspnea | Smoker=yes)
    result, err := ve.Query(
        []string{"Dyspnea"},           // query variables
        map[string]int{"Smoker": 1},   // evidence (1 = "yes")
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("P(Dyspnea | Smoker=yes):", result.Values().Data())

    // MAP query: most likely assignment
    assignment, err := ve.MAP(
        []string{"Lung", "Bronc"},
        map[string]int{"Smoker": 1},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("MAP(Lung, Bronc | Smoker=yes):", assignment)
}`}</code></pre>

        <h3>Your First Structure Learning</h3>
        <p>
          When you have observational data but no known structure, use structure learning to discover the DAG:
        </p>
        <pre><code>{`import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/example_models"
    "github.com/asymmetric-effort/datascience/lib/pgm/learning"
    "github.com/asymmetric-effort/datascience/lib/pgm/sampling"
    "github.com/asymmetric-effort/datascience/lib/pgm/structure_score"
)

func main() {
    // Generate training data from a known model
    bn, _ := example_models.Get("asia")
    bms, _ := sampling.NewBayesianModelSampling(bn, 42)
    data, _ := bms.ForwardSample(5000)

    // Learn structure using hill-climbing with BIC scoring
    scorer := structure_score.NewBIC()
    hc := learning.NewHillClimbSearch(data, scorer.LocalScore)
    learnedBN, err := hc.Estimate()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Learned: %d nodes, %d edges\\n",
        len(learnedBN.Nodes()), len(learnedBN.Edges()))

    // Fit parameters to the learned structure
    mle := learning.NewMLE(learnedBN, data)
    if err := mle.Estimate(); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Parameters fitted successfully")
}`}</code></pre>
      </section>

      {/* ============================================================ */}
      {/* ARCHITECTURE */}
      {/* ============================================================ */}
      <section class="section" id="architecture">
        <h2>Architecture</h2>

        <h3>Project Layout</h3>
        <pre><code>{`datascience/
  lib/                     # All library modules
    numgo/                 # numpy equivalent: NDArray, Matrix, Vector, BLAS
    scigo/                 # scipy equivalent: distributions, optimization, special functions
    graphgo/               # networkx equivalent: DiGraph, PDAG, graph algorithms
    tabgo/                 # pandas equivalent: DataFrame, Series, CSV/Parquet I/O
    gpu/                   # GPU compute backend for large-scale operations
    pgm/                   # Probabilistic Graphical Models (pgmpy-inspired)
      base/                # DAG, PDAG, UndirectedGraph, ADMG, MAG, SimpleCausalModel
      models/              # 13 probabilistic model types
      factors/             # Factor representations: TabularCPD, DiscreteFactor, etc.
      inference/           # 7 inference algorithms (VE, BP, MPLP, Causal, etc.)
      sampling/            # Forward, rejection, likelihood-weighted, Gibbs sampling
      learning/            # 15+ learning algorithms (parameter + structure)
      ci_tests/            # 16 conditional independence tests
      structure_score/     # 13 scoring functions for structure learning
      identification/      # Causal effect identification (back-door, front-door)
      prediction/          # DoubleML, naive adjustment, IV regression
      metrics/             # SHD, confusion matrices, correlation, Fisher's C
      independencies/      # Independence assertion representations
      readwrite/           # 10 file format readers and writers
      config/              # Global configuration
      utils/               # Shared parsing, optimization, compatibility utilities
    tensorflow/            # TensorFlow-compatible deep learning
  example_models/          # 25 built-in Bayesian networks
  examples/
    datasets/              # 40 built-in CSV datasets
  website/                 # Project website (this site)
  docs/                    # Additional documentation`}</code></pre>

        <h3>Dependency Flow</h3>
        <p>
          The dependency flow is strictly layered. Each layer depends only on layers below it:
        </p>
        <pre><code>{`Layer 3: lib/pgm, lib/tensorflow  (domain: PGMs, deep learning, etc.)
         |
Layer 2: lib/numgo, lib/scigo,    (primitives: arrays, stats, graphs, tables, GPU)
         lib/graphgo, lib/tabgo,
         lib/gpu
         |
Layer 1: Go standard library       (only dependency)`}</code></pre>
        <p>
          The <code>lib/</code> packages are general-purpose and can be used independently. For example,
          you could use <code>numgo</code> for matrix math or <code>graphgo</code> for graph algorithms
          without importing any PGM-specific code.
        </p>
        <p>
          The <code>lib/pgm</code> packages build on the primitive libraries to implement PGM-specific functionality.
          They also depend on each other -- for example, <code>inference</code> depends on <code>factors</code>,
          and <code>learning</code> depends on <code>structure_score</code> and <code>ci_tests</code>.
        </p>
        <p>
          The <code>lib/tensorflow</code> package provides TensorFlow-compatible deep learning functionality,
          building on <code>numgo</code> for tensor operations and <code>gpu</code> for acceleration.
        </p>

        <h3>Design Principles</h3>
        <ul>
          <li><strong>Near-zero dependencies:</strong> The entire library compiles with only the Go standard library. No cgo, no system libraries, no third-party modules.</li>
          <li><strong>pgmpy compatibility:</strong> PGM API patterns follow pgmpy where practical, making it easier for Python practitioners to transition to Go.</li>
          <li><strong>Cross-validation:</strong> 392 test fixtures validate results against known-correct outputs, ensuring numerical accuracy across inference, learning, and sampling.</li>
          <li><strong>Layered architecture:</strong> Primitive libraries (numgo, scigo, graphgo, tabgo) are reusable beyond PGMs. Domain libraries (pgm, tensorflow) build on these without polluting them.</li>
        </ul>
      </section>

      {/* ============================================================ */}
      {/* LIBRARY PACKAGES */}
      {/* ============================================================ */}
      <section class="section" id="library-packages">
        <h2>Library Packages (lib/)</h2>
        <p>
          These packages replace common Python scientific computing libraries with pure Go implementations.
          They are general-purpose and can be imported independently of the PGM or deep learning layers.
        </p>

        <h3>numgo -- numpy equivalent</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/numgo</code>
        </p>
        <p>
          N-dimensional arrays, linear algebra, BLAS routines, matrix operations, broadcasting, and element-wise arithmetic.
        </p>
        <table>
          <thead>
            <tr><th>Type</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>NDArray</code></td><td>N-dimensional array with shape, stride, and element-wise operations (add, multiply, etc.)</td></tr>
            <tr><td><code>Matrix</code></td><td>2D matrix with multiply, transpose, inverse, determinant, eigenvalues</td></tr>
            <tr><td><code>Vector</code></td><td>1D vector with dot product, norm, element-wise operations</td></tr>
          </tbody>
        </table>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/lib/numgo"

// Create a matrix
m := numgo.NewMatrix(3, 3)
m.Set(0, 0, 1.0)
m.Set(1, 1, 2.0)
m.Set(2, 2, 3.0)

// Matrix operations
det := m.Det()
inv := m.Inverse()
transposed := m.Transpose()
product := m.Multiply(inv) // should be identity

// NDArray operations
arr := numgo.NewNDArray([]int{2, 3, 4}) // 2x3x4 array
arr.Fill(1.0)
sum := arr.Sum()

// Vector operations
v1 := numgo.NewVector([]float64{1, 2, 3})
v2 := numgo.NewVector([]float64{4, 5, 6})
dot := v1.Dot(v2)
norm := v1.Norm()`}</code></pre>

        <h3>scigo -- scipy equivalent</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/scigo</code>
        </p>
        <p>
          Statistical distributions, optimization routines, special functions, and hypothesis tests.
        </p>
        <table>
          <thead>
            <tr><th>Category</th><th>Key Types / Functions</th></tr>
          </thead>
          <tbody>
            <tr><td>Distributions</td><td><code>Normal</code>, <code>ChiSquared</code>, <code>Beta</code>, <code>Gamma</code>, <code>StudentT</code>, <code>Uniform</code>, <code>Exponential</code></td></tr>
            <tr><td>Optimization</td><td><code>Minimize</code>, <code>GradientDescent</code>, <code>NewtonMethod</code></td></tr>
            <tr><td>Statistics</td><td>Hypothesis tests, p-value computation, quantile functions, CDF/PDF/PPF</td></tr>
            <tr><td>Special Functions</td><td>Gamma function, beta function, incomplete gamma, digamma</td></tr>
          </tbody>
        </table>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/lib/scigo"

// Normal distribution
n := scigo.NewNormal(0, 1) // mean=0, std=1
pdf := n.PDF(1.96)
cdf := n.CDF(1.96)       // ~0.975
ppf := n.PPF(0.975)      // ~1.96

// Chi-squared distribution (used in CI tests)
chi2 := scigo.NewChiSquared(5) // 5 degrees of freedom
pValue := 1.0 - chi2.CDF(11.07)

// Optimization
result := scigo.Minimize(func(x float64) float64 {
    return (x - 3) * (x - 3)
}, 0.0, 10.0)`}</code></pre>

        <h3>graphgo -- networkx equivalent</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/graphgo</code>
        </p>
        <p>
          Directed and undirected graphs with a full suite of graph algorithms.
        </p>
        <table>
          <thead>
            <tr><th>Type</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>DiGraph</code></td><td>Directed graph with adjacency operations, successors, predecessors</td></tr>
            <tr><td><code>Graph</code></td><td>Undirected graph with neighbors, degree, connected components</td></tr>
            <tr><td><code>PDAG</code></td><td>Partially directed acyclic graph (for equivalence classes)</td></tr>
          </tbody>
        </table>
        <p>
          Algorithms: topological sort, d-separation, moral graph, triangulation, maximum cardinality search,
          clique finding, connected components, shortest paths, cycle detection.
        </p>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/lib/graphgo"

// Create a directed graph
g := graphgo.NewDiGraph()
g.AddNode("A")
g.AddNode("B")
g.AddNode("C")
g.AddEdge("A", "B")
g.AddEdge("B", "C")

// Graph queries
parents := g.Predecessors("C")  // ["B"]
children := g.Successors("A")   // ["B"]
sorted := g.TopologicalSort()   // ["A", "B", "C"]

// Undirected graph
ug := graphgo.NewGraph()
ug.AddEdge("X", "Y")
ug.AddEdge("Y", "Z")
neighbors := ug.Neighbors("Y")  // ["X", "Z"]`}</code></pre>

        <h3>tabgo -- pandas equivalent</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/tabgo</code>
        </p>
        <p>
          Tabular data with named columns, row filtering, groupby, and file I/O.
        </p>
        <table>
          <thead>
            <tr><th>Type</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>DataFrame</code></td><td>Tabular data with named columns, row filtering, groupby, merge</td></tr>
            <tr><td><code>Series</code></td><td>Single column with value counts, unique values, statistical summaries</td></tr>
          </tbody>
        </table>
        <p>
          I/O: <code>ReadCSV</code>, <code>WriteCSV</code>, <code>ReadParquet</code>, <code>WriteParquet</code>, <code>ReadXLSX</code>.
        </p>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/lib/tabgo"

// Read CSV data
df, err := tabgo.ReadCSV("observations.csv")
fmt.Printf("Rows: %d, Columns: %d\\n", df.NRows(), len(df.Columns()))

// Access a column as a Series
col := df.Column("Temperature")
fmt.Println("Unique values:", col.Unique())
fmt.Println("Value counts:", col.ValueCounts())

// Filter rows
filtered := df.Filter(func(row map[string]interface{}) bool {
    return row["Temperature"].(int) > 70
})

// Write CSV
tabgo.WriteCSV(filtered, "warm_days.csv")`}</code></pre>

        <h3>gpu -- GPU Compute Backend</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/gpu</code>
        </p>
        <p>
          Optional GPU acceleration for compute-intensive operations on large networks. Provides GPU-backed
          matrix operations and factor computations that can significantly speed up inference, learning,
          and deep learning training.
        </p>

        <h3>tensorflow -- Deep Learning</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/tensorflow</code>
        </p>
        <p>
          TensorFlow-compatible deep learning in pure Go. Provides neural network construction,
          training, and inference capabilities. Builds on <code>numgo</code> for tensor operations
          and <code>gpu</code> for hardware acceleration.
        </p>
      </section>

      {/* ============================================================ */}
      {/* PGM PACKAGES */}
      {/* ============================================================ */}
      <section class="section" id="pgm-packages">
        <h2>PGM Packages (lib/pgm/)</h2>
        <p>
          The PGM module provides a complete probabilistic graphical models toolkit, inspired by pgmpy.
          It is one of several domain modules in the datascience library.
        </p>

        <h3>base -- Foundational Graph Types</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/pgm/base</code>
        </p>
        <p>
          Provides the underlying graph structures that all model types are built on.
        </p>
        <table>
          <thead>
            <tr><th>Type</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>DAG</code></td><td>Directed acyclic graph with cycle detection, topological sort, d-separation</td></tr>
            <tr><td><code>PDAG</code></td><td>Partially directed acyclic graph for Markov equivalence classes</td></tr>
            <tr><td><code>UndirectedGraph</code></td><td>Undirected graph for Markov networks</td></tr>
            <tr><td><code>ADMG</code></td><td>Acyclic directed mixed graph (with bidirected edges for latent confounders)</td></tr>
            <tr><td><code>MAG</code></td><td>Maximal ancestral graph</td></tr>
            <tr><td><code>SimpleCausalModel</code></td><td>Basic causal model with intervention semantics</td></tr>
          </tbody>
        </table>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/lib/pgm/base"

dag := base.NewDAG()
dag.AddNode("X")
dag.AddNode("Y")
dag.AddNode("Z")
dag.AddEdge("X", "Y")
dag.AddEdge("Y", "Z")

// Check d-separation: X _||_ Z | Y?
separated := dag.DSeparation([]string{"X"}, []string{"Z"}, []string{"Y"})
fmt.Println("X _||_ Z | Y:", separated) // true`}</code></pre>

        <h3>models -- Probabilistic Model Types</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/pgm/models</code>
        </p>
        <p>13 model types for different PGM use cases:</p>
        <table>
          <thead>
            <tr><th>Type</th><th>Description</th><th>Constructor</th></tr>
          </thead>
          <tbody>
            <tr><td><code>BayesianNetwork</code></td><td>DAG with TabularCPDs. The primary model type for most use cases.</td><td><code>NewBayesianNetwork()</code></td></tr>
            <tr><td><code>MarkovNetwork</code></td><td>Undirected graphical model with factor potentials.</td><td><code>NewMarkovNetwork()</code></td></tr>
            <tr><td><code>DynamicBayesianNetwork</code></td><td>BN over time slices for temporal modeling.</td><td><code>NewDynamicBayesianNetwork()</code></td></tr>
            <tr><td><code>NaiveBayes</code></td><td>Naive Bayes classifier (BN with single class parent).</td><td><code>NewNaiveBayes()</code></td></tr>
            <tr><td><code>SEM</code></td><td>Structural Equation Model with linear/nonlinear equations.</td><td><code>NewSEM()</code></td></tr>
            <tr><td><code>FactorGraph</code></td><td>Bipartite graph of variable nodes and factor nodes.</td><td><code>NewFactorGraph()</code></td></tr>
            <tr><td><code>JunctionTree</code></td><td>Clique tree for exact inference via message passing.</td><td><code>NewJunctionTree()</code></td></tr>
            <tr><td><code>ClusterGraph</code></td><td>Generalized cluster graph (superset of junction tree).</td><td><code>NewClusterGraph()</code></td></tr>
            <tr><td><code>LinearGaussianBN</code></td><td>BN with continuous, linearly-related Gaussian variables.</td><td><code>NewLinearGaussianBN()</code></td></tr>
            <tr><td><code>FunctionalBN</code></td><td>BN where CPDs are defined by arbitrary functions.</td><td><code>NewFunctionalBN()</code></td></tr>
            <tr><td><code>MarkovChain</code></td><td>First-order Markov chain for sequential data.</td><td><code>NewMarkovChain()</code></td></tr>
            <tr><td><code>DiscreteBayesianNetwork</code></td><td>Specialized discrete-only BN with optimized operations.</td><td><code>NewDiscreteBayesianNetwork()</code></td></tr>
            <tr><td><code>DiscreteMarkovNetwork</code></td><td>Specialized discrete-only Markov network.</td><td><code>NewDiscreteMarkovNetwork()</code></td></tr>
          </tbody>
        </table>
        <p>
          Key methods shared by most model types: <code>AddNode</code>, <code>AddEdge</code>, <code>Nodes()</code>,
          <code>Edges()</code>, <code>SetStates</code>, <code>SetCPD</code>, <code>CheckModel()</code>,
          <code>ToMarkovFactors()</code>, <code>ToJunctionTree()</code>.
        </p>

        <h3>factors -- Factor Representations</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/pgm/factors</code>
        </p>
        <table>
          <thead>
            <tr><th>Type</th><th>Description</th><th>Constructor</th></tr>
          </thead>
          <tbody>
            <tr><td><code>DiscreteFactor</code></td><td>General discrete factor with product, marginalize, reduce, normalize operations.</td><td><code>NewDiscreteFactor()</code></td></tr>
            <tr><td><code>TabularCPD</code></td><td>Conditional probability distribution table. The standard CPD type for BayesianNetwork.</td><td><code>NewTabularCPD()</code></td></tr>
            <tr><td><code>JointProbabilityDistribution</code></td><td>Full joint distribution over a set of variables.</td><td><code>NewJPD()</code></td></tr>
            <tr><td><code>LinearGaussianCPD</code></td><td>Linear Gaussian conditional: child = sum(beta_i * parent_i) + noise.</td><td><code>NewLinearGaussianCPD()</code></td></tr>
            <tr><td><code>NoisyOR</code></td><td>Noisy-OR parameterization. Compact CPD for nodes with many binary parents.</td><td><code>NewNoisyOR()</code></td></tr>
          </tbody>
        </table>

        <h3>inference -- Inference Algorithms</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/pgm/inference</code>
        </p>
        <table>
          <thead>
            <tr><th>Type</th><th>Description</th><th>Constructor</th></tr>
          </thead>
          <tbody>
            <tr><td><code>VariableElimination</code></td><td>Exact inference via factor elimination. Supports <code>Query()</code> for posteriors and <code>MAP()</code> for most-probable assignments.</td><td><code>NewVariableElimination(factors)</code></td></tr>
            <tr><td><code>BeliefPropagation</code></td><td>Message-passing on junction trees. Calibrate once, query multiple times.</td><td><code>NewBeliefPropagation(cliques, separators, factors)</code></td></tr>
            <tr><td><code>MPLP</code></td><td>Max-Product Linear Programming for MAP inference.</td><td><code>NewMPLP(factors)</code></td></tr>
            <tr><td><code>ApproxInference</code></td><td>Sampling-based approximate inference. Uses likelihood-weighted sampling.</td><td><code>NewApproxInference(bn, nSamples)</code></td></tr>
            <tr><td><code>CausalInference</code></td><td>Do-calculus interventional queries. Computes P(Y | do(X=x)).</td><td><code>NewCausalInference(bn)</code></td></tr>
            <tr><td><code>DBNInference</code></td><td>Inference over dynamic Bayesian networks across time slices.</td><td><code>NewDBNInference(dbn)</code></td></tr>
          </tbody>
        </table>

        <h3>sampling -- Sampling Methods</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/pgm/sampling</code>
        </p>
        <table>
          <thead>
            <tr><th>Type</th><th>Methods</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>BayesianModelSampling</code></td><td><code>ForwardSample</code>, <code>RejectionSample</code>, <code>LikelihoodWeightedSample</code></td><td>Exact and weighted sampling from BN joint distribution.</td></tr>
            <tr><td><code>GibbsSampling</code></td><td><code>Sample</code></td><td>MCMC Gibbs sampler with configurable burn-in and thinning.</td></tr>
          </tbody>
        </table>

        <h3>learning -- Learning Algorithms</h3>
        <p>
          Import: <code>github.com/asymmetric-effort/datascience/lib/pgm/learning</code>
        </p>
        <p>15+ algorithms for parameter estimation and structure learning. See the <Link to="/api">API Reference</Link> for full details.</p>

        <h3>Other PGM Packages</h3>
        <table>
          <thead>
            <tr><th>Package</th><th>Import Path</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>ci_tests</code></td><td><code>lib/pgm/ci_tests</code></td><td>16 conditional independence tests (ChiSquare, FisherZ, GCM, etc.)</td></tr>
            <tr><td><code>structure_score</code></td><td><code>lib/pgm/structure_score</code></td><td>13 scoring functions (BIC, AIC, BDeu, K2, etc.)</td></tr>
            <tr><td><code>identification</code></td><td><code>lib/pgm/identification</code></td><td>Causal effect identification (back-door, front-door)</td></tr>
            <tr><td><code>prediction</code></td><td><code>lib/pgm/prediction</code></td><td>DoubleML, naive adjustment, IV regression</td></tr>
            <tr><td><code>metrics</code></td><td><code>lib/pgm/metrics</code></td><td>SHD, confusion matrices, Fisher's C, log-likelihood</td></tr>
            <tr><td><code>independencies</code></td><td><code>lib/pgm/independencies</code></td><td>Independence assertion representations</td></tr>
            <tr><td><code>readwrite</code></td><td><code>lib/pgm/readwrite</code></td><td>10 file format readers and writers</td></tr>
            <tr><td><code>config</code></td><td><code>lib/pgm/config</code></td><td>Global PGM configuration</td></tr>
            <tr><td><code>utils</code></td><td><code>lib/pgm/utils</code></td><td>Shared utilities</td></tr>
          </tbody>
        </table>
      </section>

      {/* ============================================================ */}
      {/* FILE FORMATS */}
      {/* ============================================================ */}
      <section class="section" id="file-formats">
        <h2>File Formats</h2>
        <p>
          datascience supports 10 file formats for reading and writing probabilistic graphical models.
          All formats are accessed through the <code>readwrite</code> package.
        </p>
        <table>
          <thead>
            <tr><th>Format</th><th>Extension</th><th>Read</th><th>Write</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><strong>BIF</strong></td><td>.bif</td><td>Yes</td><td>Yes</td><td>Bayesian Interchange Format. The standard PGM format. Stores structure, states, and CPD tables in a human-readable text format.</td></tr>
            <tr><td><strong>XMLBIF</strong></td><td>.xmlbif</td><td>Yes</td><td>Yes</td><td>XML-based BIF. Same information as BIF but in XML for easier parsing by other tools.</td></tr>
            <tr><td><strong>NET</strong></td><td>.net</td><td>Yes</td><td>Yes</td><td>Hugin NET format. Used by the Hugin BN software.</td></tr>
            <tr><td><strong>UAI</strong></td><td>.uai</td><td>Yes</td><td>Yes</td><td>UAI format. Used by the UAI inference competition. Compact numeric format.</td></tr>
            <tr><td><strong>XDSL</strong></td><td>.xdsl</td><td>Yes</td><td>Yes</td><td>GeNIe XDSL format. Used by the GeNIe/SMILE BN software.</td></tr>
            <tr><td><strong>PomdpX</strong></td><td>.pomdpx</td><td>Yes</td><td>--</td><td>POMDP XML format. Read-only. Used in POMDP planning literature.</td></tr>
            <tr><td><strong>XBN</strong></td><td>.xbn</td><td>Yes</td><td>--</td><td>Microsoft XBN format. Read-only. Legacy Microsoft Research format.</td></tr>
            <tr><td><strong>CSV</strong></td><td>.csv</td><td>Yes</td><td>Yes</td><td>CSV model serialization. Stores structure and parameters in CSV tables.</td></tr>
            <tr><td><strong>JSON</strong></td><td>.json</td><td>Yes</td><td>Yes</td><td>JSON model serialization. Ideal for web applications and REST APIs.</td></tr>
            <tr><td><strong>XML</strong></td><td>.xml</td><td>Yes</td><td>Yes</td><td>XML model serialization. General-purpose XML format.</td></tr>
          </tbody>
        </table>

        <h3>BIF Format Example</h3>
        <pre><code>{`network asia {
}

variable VisitAsia {
  type discrete [ 2 ] { no, yes };
}

variable Tuberculosis {
  type discrete [ 2 ] { no, yes };
}

probability ( VisitAsia ) {
  table 0.99, 0.01;
}

probability ( Tuberculosis | VisitAsia ) {
  (no) 0.99, 0.01;
  (yes) 0.95, 0.05;
}`}</code></pre>

        <h3>JSON Format Example</h3>
        <pre><code>{`{
  "nodes": ["A", "B", "C"],
  "edges": [["A", "B"], ["B", "C"]],
  "states": {
    "A": ["a0", "a1"],
    "B": ["b0", "b1"],
    "C": ["c0", "c1"]
  },
  "cpds": {
    "A": {
      "variable": "A",
      "cardinality": 2,
      "values": [0.4, 0.6],
      "parents": [],
      "parent_cardinalities": []
    }
  }
}`}</code></pre>

        <h3>Reading and Writing</h3>
        <pre><code>{`import (
    "os"
    "github.com/asymmetric-effort/datascience/lib/pgm/readwrite"
)

// Read BIF
f, _ := os.Open("model.bif")
bn, _ := readwrite.ReadBIF(f)
f.Close()

// Write as JSON (for a web API)
out, _ := os.Create("model.json")
readwrite.WriteJSONModel(out, bn)
out.Close()

// Convert: read NET, write XMLBIF
netFile, _ := os.Open("model.net")
bn2, _ := readwrite.ReadNET(netFile)
netFile.Close()

xmlbifFile, _ := os.Create("model.xmlbif")
readwrite.WriteXMLBIF(xmlbifFile, bn2)
xmlbifFile.Close()`}</code></pre>
      </section>

      {/* ============================================================ */}
      {/* EXAMPLE MODELS */}
      {/* ============================================================ */}
      <section class="section" id="example-models">
        <h2>Example Models</h2>
        <p>
          datascience ships with 25 built-in Bayesian networks accessible via the <code>example_models</code> package.
          13 models include full CPDs (conditional probability distributions); 12 are structure-only for use in
          learning and benchmarking.
        </p>

        <h3>Models with Full CPDs</h3>
        <table>
          <thead>
            <tr><th>Name</th><th>Description</th><th>Nodes</th></tr>
          </thead>
          <tbody>
            <tr><td><code>student</code></td><td>Classic Student network (difficulty, intelligence, grade, SAT, letter)</td><td>5</td></tr>
            <tr><td><code>asia</code></td><td>Lung disease diagnosis</td><td>8</td></tr>
            <tr><td><code>alarm</code></td><td>Monitoring system with alarm, burglary, earthquake</td><td>5</td></tr>
            <tr><td><code>cancer</code></td><td>Cancer diagnosis network</td><td>5</td></tr>
            <tr><td><code>watersprinkler</code></td><td>Classic sprinkler/rain/wet grass example</td><td>4</td></tr>
            <tr><td><code>survey</code></td><td>Survey response model</td><td>6</td></tr>
            <tr><td><code>montyhall</code></td><td>Monty Hall problem as a Bayesian network</td><td>3</td></tr>
            <tr><td><code>dogproblem</code></td><td>Dog behavior inference</td><td>5</td></tr>
            <tr><td><code>frauddetection</code></td><td>Financial fraud detection model</td><td>5</td></tr>
            <tr><td><code>medicaldiagnosis</code></td><td>Medical symptom/disease model</td><td>8</td></tr>
            <tr><td><code>earthquake</code></td><td>Earthquake alert network</td><td>5</td></tr>
            <tr><td><code>visitasia</code></td><td>Visit to Asia variant</td><td>8</td></tr>
            <tr><td><code>cointoss</code></td><td>Simple coin toss model</td><td>2</td></tr>
          </tbody>
        </table>

        <h3>Structure-Only Models (Large Networks)</h3>
        <table>
          <thead>
            <tr><th>Name</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>sachs</code></td><td>Protein signaling network (11 nodes)</td></tr>
            <tr><td><code>child</code></td><td>Child health assessment network</td></tr>
            <tr><td><code>insurance</code></td><td>Insurance risk assessment</td></tr>
            <tr><td><code>alarmfull</code></td><td>Full ALARM monitoring network (37 nodes)</td></tr>
            <tr><td><code>water</code></td><td>Water treatment network</td></tr>
            <tr><td><code>mildew</code></td><td>Crop disease model</td></tr>
            <tr><td><code>barley</code></td><td>Barley crop yield model</td></tr>
            <tr><td><code>hailfinder</code></td><td>Severe weather prediction</td></tr>
            <tr><td><code>hepar2</code></td><td>Liver disorder diagnosis</td></tr>
            <tr><td><code>win95pts</code></td><td>Windows 95 printer troubleshooting</td></tr>
            <tr><td><code>pathfinder</code></td><td>Pathology diagnosis</td></tr>
            <tr><td><code>pigs</code></td><td>Pig breeding network</td></tr>
          </tbody>
        </table>

        <h3>Usage</h3>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/example_models"

// List all available models
names := example_models.List()
for _, name := range names {
    fmt.Println(name)
}

// Load a specific model
bn, err := example_models.Get("asia")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Nodes: %d, Edges: %d\\n", len(bn.Nodes()), len(bn.Edges()))

// Use the model for inference
facs, _ := bn.ToMarkovFactors()
ve := inference.NewVariableElimination(facs)
result, _ := ve.Query([]string{"Dyspnea"}, map[string]int{"Smoker": 1})
fmt.Println(result.Values().Data())`}</code></pre>
      </section>

      {/* ============================================================ */}
      {/* DATASETS */}
      {/* ============================================================ */}
      <section class="section" id="datasets">
        <h2>Datasets</h2>
        <p>
          datascience includes 40 built-in datasets accessible via the <code>examples/datasets</code> package.
          These datasets are embedded in the binary using Go's <code>embed</code> package, so they are
          always available without external file dependencies.
        </p>

        <h3>BN-Specific Datasets</h3>
        <table>
          <thead>
            <tr><th>Name</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>asia</code></td><td>Sampled data from the Asia (lung disease) network</td></tr>
            <tr><td><code>alarm</code></td><td>Sampled data from the ALARM monitoring network</td></tr>
            <tr><td><code>sachs</code></td><td>Protein signaling data (Sachs et al.)</td></tr>
            <tr><td><code>cancer</code></td><td>Cancer diagnosis observations</td></tr>
            <tr><td><code>student</code></td><td>Student performance data</td></tr>
            <tr><td><code>sprinkler</code></td><td>Sprinkler/rain/wet grass observations</td></tr>
            <tr><td><code>survey</code></td><td>Survey response data</td></tr>
            <tr><td><code>earthquake</code></td><td>Earthquake alert observations</td></tr>
            <tr><td><code>child</code></td><td>Child health assessment data</td></tr>
            <tr><td><code>insurance</code></td><td>Insurance risk data</td></tr>
            <tr><td><code>water</code></td><td>Water treatment data</td></tr>
            <tr><td><code>mildew</code></td><td>Crop disease data</td></tr>
            <tr><td><code>hailfinder</code></td><td>Severe weather observations</td></tr>
            <tr><td><code>hepar2</code></td><td>Liver disorder data</td></tr>
            <tr><td><code>barley</code></td><td>Barley crop data</td></tr>
            <tr><td><code>win95pts</code></td><td>Windows 95 troubleshooting data</td></tr>
            <tr><td><code>andes</code></td><td>ANDES intelligent tutoring system data</td></tr>
            <tr><td><code>munin</code></td><td>MUNIN neural network data</td></tr>
            <tr><td><code>lucas</code></td><td>LUCAS causal discovery benchmark</td></tr>
          </tbody>
        </table>

        <h3>Classic ML Datasets</h3>
        <table>
          <thead>
            <tr><th>Name</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>titanic</code></td><td>Titanic survival data</td></tr>
            <tr><td><code>iris</code></td><td>Fisher's Iris flower dataset</td></tr>
            <tr><td><code>heart</code></td><td>Heart disease prediction</td></tr>
            <tr><td><code>wine</code></td><td>Wine quality classification</td></tr>
            <tr><td><code>boston</code></td><td>Boston housing prices</td></tr>
            <tr><td><code>pima_diabetes</code></td><td>Pima Indians diabetes</td></tr>
            <tr><td><code>adult</code></td><td>Adult income prediction (Census)</td></tr>
            <tr><td><code>breast_cancer</code></td><td>Wisconsin breast cancer diagnosis</td></tr>
          </tbody>
        </table>

        <h3>UCI Repository Datasets</h3>
        <table>
          <thead>
            <tr><th>Name</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>zoo</code></td><td>Zoo animal classification</td></tr>
            <tr><td><code>glass</code></td><td>Glass identification</td></tr>
            <tr><td><code>ecoli</code></td><td>E. coli protein localization</td></tr>
            <tr><td><code>monks</code></td><td>MONKS problem</td></tr>
            <tr><td><code>nursery</code></td><td>Nursery school evaluation</td></tr>
            <tr><td><code>credit_approval</code></td><td>Credit card approval</td></tr>
            <tr><td><code>balance_scale</code></td><td>Balance scale weight/distance</td></tr>
            <tr><td><code>automobile</code></td><td>Automobile price prediction</td></tr>
            <tr><td><code>mushroom</code></td><td>Mushroom edibility classification</td></tr>
            <tr><td><code>car_evaluation</code></td><td>Car evaluation</td></tr>
            <tr><td><code>hepatitis</code></td><td>Hepatitis prognosis</td></tr>
            <tr><td><code>vote</code></td><td>Congressional voting records</td></tr>
            <tr><td><code>tic_tac_toe</code></td><td>Tic-tac-toe endgame</td></tr>
          </tbody>
        </table>

        <h3>Usage</h3>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/examples/datasets"

// List all available datasets
names := datasets.List()

// Load a dataset as a DataFrame
df, err := datasets.Load("asia")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Rows: %d, Columns: %d\\n", df.NRows(), len(df.Columns()))

// Use the dataset for structure learning
scorer := structure_score.NewBIC()
hc := learning.NewHillClimbSearch(df, scorer.LocalScore)
bn, _ := hc.Estimate()`}</code></pre>
      </section>

      {/* ============================================================ */}
      {/* TESTING */}
      {/* ============================================================ */}
      <section class="section" id="testing">
        <h2>Testing</h2>
        <p>
          datascience has approximately 5,000 tests and 392 cross-validation fixtures spanning unit tests,
          integration tests, and cross-validation tests across the library.
        </p>

        <h3>Running Tests</h3>
        <pre><code>{`# Run all tests
go test ./...

# Run tests for a specific package
go test ./lib/pgm/inference/...

# Run with verbose output
go test -v ./lib/pgm/models/...

# Run with race detector
go test -race ./...

# Run a specific test function
go test -run TestVariableElimination ./lib/pgm/inference/...

# Run cross-validation tests only
go test -run CrossVal ./...

# Run benchmarks
go test -bench=. ./lib/pgm/inference/...`}</code></pre>

        <h3>Cross-Validation System</h3>
        <p>
          Many packages include <code>crossval_*_test.go</code> files that validate algorithms against
          known-correct results. These tests load built-in example models, run computations, and compare
          outputs against pre-computed reference values. Examples:
        </p>
        <ul>
          <li><code>lib/pgm/models/crossval_dsep_test.go</code> -- validates d-separation queries</li>
          <li><code>lib/pgm/inference/crossval_causal_test.go</code> -- validates causal inference results</li>
          <li><code>lib/pgm/inference/crossval_ve_test.go</code> -- validates Variable Elimination posteriors</li>
          <li><code>lib/pgm/learning/crossval_hillclimb_test.go</code> -- validates structure learning output</li>
          <li><code>lib/pgm/sampling/crossval_forward_test.go</code> -- validates sampling distributions</li>
        </ul>
        <p>
          Cross-validation fixtures are generated by running the equivalent pgmpy code in Python and
          storing the results. This ensures datascience produces the same numerical outputs as the reference
          implementation.
        </p>

        <h3>Test Fixture Generation</h3>
        <p>
          Test fixtures use the built-in example models from the <code>example_models</code> package.
          This ensures reproducible test data without external file dependencies. When adding new tests,
          use existing models or create new ones in the <code>example_models</code> package.
        </p>

        <h3>Writing Tests</h3>
        <pre><code>{`func TestMyFeature(t *testing.T) {
    // Load a known model
    bn, err := example_models.Get("asia")
    if err != nil {
        t.Fatal(err)
    }

    // Perform computation
    facs, _ := bn.ToMarkovFactors()
    ve := inference.NewVariableElimination(facs)
    result, err := ve.Query([]string{"Dyspnea"}, map[string]int{"Smoker": 1})
    if err != nil {
        t.Fatal(err)
    }

    // Compare against expected values
    values := result.Values().Data()
    if math.Abs(values[0] - 0.304) > 0.01 {
        t.Errorf("expected P(Dyspnea=0|Smoker=1) ~ 0.304, got %f", values[0])
    }
}`}</code></pre>
      </section>

      {/* ============================================================ */}
      {/* CONFIGURATION */}
      {/* ============================================================ */}
      <section class="section" id="configuration">
        <h2>Configuration</h2>
        <p>
          The <code>config</code> package provides global configuration options that control
          default behavior across the PGM module.
        </p>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/lib/pgm/config"

// Get the global config
cfg := config.Global()

// Configuration is used internally by various packages
// to control default inference methods, scoring functions,
// numerical tolerances, and other global settings.`}</code></pre>
        <p>
          Configuration options include numerical tolerances for probability comparisons,
          default inference methods, default scoring functions for structure learning,
          convergence thresholds for iterative algorithms (EM, Belief Propagation),
          and logging verbosity.
        </p>
      </section>

      {/* ============================================================ */}
      {/* CONTRIBUTING */}
      {/* ============================================================ */}
      <section class="section" id="contributing">
        <h2>Contributing</h2>
        <p>
          Contributions are welcome. See the full{" "}
          <a href="https://github.com/asymmetric-effort/datascience/blob/main/CONTRIBUTING.md" target="_blank" rel="noopener noreferrer">
            CONTRIBUTING.md
          </a>{" "}
          for details.
        </p>

        <h3>Development Workflow</h3>
        <ol>
          <li>Fork the repository and clone your fork</li>
          <li>Create a feature branch: <code>git checkout -b feature/my-feature</code></li>
          <li>Make changes and add tests</li>
          <li>Run <code>go test ./...</code> to verify all tests pass</li>
          <li>Run <code>go vet ./...</code> for static analysis</li>
          <li>Commit with a clear message and submit a pull request</li>
        </ol>

        <h3>Guidelines</h3>
        <ul>
          <li><strong>Near-zero dependencies:</strong> Do not add third-party modules. All functionality must be implemented in pure Go using only the standard library.</li>
          <li><strong>Tests required:</strong> All new functionality must include unit tests. Cross-validation tests against pgmpy are strongly encouraged.</li>
          <li><strong>Backward compatibility:</strong> Public API changes require discussion in an issue before implementation.</li>
          <li><strong>Documentation:</strong> Exported types and functions must have godoc comments.</li>
        </ul>
      </section>
    </div>
  );
}
