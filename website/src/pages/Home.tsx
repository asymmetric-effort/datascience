import { createElement, useHead, Link } from "@asymmetric-effort/specifyjs";

export function Home() {
  useHead({
    title: "datascience — Data Science and Machine Learning in Go",
    description: "A comprehensive pure Go library for data science and machine learning. PGMs, TensorFlow/deep learning, BLAS, financial modeling, and more.",
    canonical: "https://datascience.asymmetric-effort.com/",
    og: {
      title: "datascience — Data Science and Machine Learning in Go",
      description: "A comprehensive pure Go library for data science and machine learning.",
      url: "https://datascience.asymmetric-effort.com/",
    },
  });

  return (
    <div class="page">
      <section class="hero">
        <img src="/docs/img/logo.png" alt="datascience logo" class="hero-logo" />
        <h1>datascience</h1>
        <p class="hero-subtitle">Data Science and Machine Learning in Go</p>
        <p class="hero-description">
          A comprehensive pure Go library for probabilistic graphical models, TensorFlow-compatible
          deep learning, BLAS linear algebra, financial modeling, and more. Built entirely in Go
          with near-zero dependencies.
        </p>
        <div class="badges">
          <span class="badge">v0.0.37</span>
          <span class="badge">Near-Zero Dependencies</span>
          <span class="badge">Go 1.21+</span>
          <span class="badge">MIT License</span>
          <span class="badge">~5,000 Tests</span>
          <span class="badge">392 Cross-Validation Fixtures</span>
          <span class="badge">PGMs</span>
          <span class="badge">TensorFlow / Deep Learning</span>
          <span class="badge">BLAS</span>
          <span class="badge">Financial Modeling</span>
          <span class="badge">40 Built-in Datasets</span>
        </div>
        <div class="hero-actions">
          <Link to="/docs" class="btn btn-primary">Get Started</Link>
          <a href="https://github.com/asymmetric-effort/datascience" target="_blank" rel="noopener noreferrer" class="btn btn-secondary">GitHub</a>
        </div>
      </section>

      <section class="section">
        <h2>Installation</h2>
        <pre><code>go get github.com/asymmetric-effort/datascience</code></pre>
        <p>Requires Go 1.21 or later. No C dependencies, no cgo, no third-party modules.</p>
      </section>

      <section class="section">
        <h2>Quick Start</h2>
        <pre><code>{`package main

import (
    "fmt"
    "github.com/asymmetric-effort/datascience/example_models"
    "github.com/asymmetric-effort/datascience/lib/pgm/factors"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
    "github.com/asymmetric-effort/datascience/lib/pgm/models"
)

func main() {
    // Option 1: Load a built-in example model
    asia, _ := example_models.Get("asia")
    fmt.Println("Asia model:", len(asia.Nodes()), "nodes")

    // Option 2: Build a Bayesian Network from scratch
    bn := models.NewBayesianNetwork()
    bn.AddNode("Rain")
    bn.AddNode("Sprinkler")
    bn.AddNode("WetGrass")
    bn.AddEdge("Rain", "WetGrass")
    bn.AddEdge("Sprinkler", "WetGrass")

    // Add CPDs
    bn.SetStates("Rain", []string{"no", "yes"})
    bn.SetStates("Sprinkler", []string{"off", "on"})
    bn.SetStates("WetGrass", []string{"dry", "wet"})

    bn.SetCPD("Rain", factors.NewTabularCPD(
        "Rain", 2, []float64{0.8, 0.2}, nil, nil,
    ))

    // Run inference with Variable Elimination
    facs, _ := bn.ToMarkovFactors()
    ve := inference.NewVariableElimination(facs)
    result, _ := ve.Query(
        []string{"WetGrass"},
        map[string]int{"Rain": 1},
    )
    fmt.Println(result)
}`}</code></pre>
        <p>See the <Link to="/docs">documentation</Link> for complete getting-started guides, or jump to <Link to="/tutorials">tutorials</Link> for step-by-step walkthroughs.</p>
      </section>

      <section class="section">
        <h2>Features</h2>
        <div class="features-grid">
          <div class="feature-card">
            <h3>Probabilistic Graphical Models</h3>
            <p>13 model types including BayesianNetwork, MarkovNetwork, DynamicBN, NaiveBayes, SEM, FactorGraph, JunctionTree, and more. 7 inference algorithms, 15+ learning algorithms.</p>
          </div>
          <div class="feature-card">
            <h3>TensorFlow / Deep Learning</h3>
            <p>TensorFlow-compatible deep learning via the <code>lib/tensorflow</code> package. Neural network construction, training, and inference in pure Go.</p>
          </div>
          <div class="feature-card">
            <h3>BLAS Linear Algebra</h3>
            <p>High-performance BLAS routines for matrix multiplication, decomposition, eigenvalue computation, and other linear algebra operations via <code>lib/numgo</code>.</p>
          </div>
          <div class="feature-card">
            <h3>Financial Modeling</h3>
            <p>Tools for quantitative finance including time series analysis, risk modeling, portfolio optimization, and statistical forecasting.</p>
          </div>
          <div class="feature-card">
            <h3>Statistical Computing</h3>
            <p>Comprehensive statistical distributions, hypothesis tests, optimization routines, and special functions via <code>lib/scigo</code> (scipy equivalent).</p>
          </div>
          <div class="feature-card">
            <h3>Graph Algorithms</h3>
            <p>Directed and undirected graphs with topological sort, d-separation, connected components, shortest paths, and more via <code>lib/graphgo</code> (networkx equivalent).</p>
          </div>
          <div class="feature-card">
            <h3>Tabular Data</h3>
            <p>DataFrames, Series, CSV/Parquet/Excel I/O, filtering, groupby, and merge operations via <code>lib/tabgo</code> (pandas equivalent).</p>
          </div>
          <div class="feature-card">
            <h3>GPU Compute Backend</h3>
            <p>Optional GPU acceleration via the <code>lib/gpu</code> package for compute-intensive operations on large networks and deep learning models.</p>
          </div>
          <div class="feature-card">
            <h3>Causal Inference</h3>
            <p>Do-calculus interventional queries, back-door and front-door identification, ATE estimation with DoubleML, naive adjustment, and IV regression.</p>
          </div>
          <div class="feature-card">
            <h3>LLM Integration</h3>
            <p>Expert-in-the-loop structure learning with LLM client support for AI-assisted model construction and knowledge elicitation.</p>
          </div>
          <div class="feature-card">
            <h3>40 Built-in Datasets</h3>
            <p>Ready-to-use CSV datasets for structure learning, parameter estimation, benchmarking, and general machine learning tasks.</p>
          </div>
          <div class="feature-card">
            <h3>Near-Zero Dependencies</h3>
            <p>Built entirely in Go with custom implementations of numpy (numgo), scipy (scigo), networkx (graphgo), and pandas (tabgo). No cgo, no external libraries.</p>
          </div>
        </div>
      </section>

      <section class="section">
        <h2>Why datascience?</h2>
        <p>
          datascience brings comprehensive data science and machine learning capabilities to the Go ecosystem.
          Whether you need probabilistic graphical models (inspired by pgmpy), deep learning (TensorFlow-compatible),
          BLAS linear algebra, or financial modeling -- this library provides a unified, pure Go solution.
        </p>
        <p>
          Every numerical primitive -- linear algebra, statistical distributions, graph algorithms, tabular data
          processing -- is implemented from scratch in pure Go. This means <code>go build</code> just works,
          cross-compilation just works, and static binaries just work.
        </p>
        <p>
          The layered architecture means you can use individual packages independently.
          Need just matrix math? Import <code>lib/numgo</code>. Just graphs? Import <code>lib/graphgo</code>.
          Just PGMs? Import <code>lib/pgm</code>. Just deep learning? Import <code>lib/tensorflow</code>.
        </p>
      </section>
    </div>
  );
}
