import { createElement, useHead } from "@asymmetric-effort/specifyjs";

export function Home() {
  useHead({
    title: "pgmgo — Probabilistic Graphical Models in Go",
    description: "A zero-dependency Go library for probabilistic graphical models, similar to pgmpy.",
    canonical: "https://pgmgo.asymmetric-effort.com/",
    og: {
      title: "pgmgo — Probabilistic Graphical Models in Go",
      description: "A zero-dependency Go library for probabilistic graphical models, similar to pgmpy.",
      url: "https://pgmgo.asymmetric-effort.com/",
    },
  });

  return (
    <div class="page">
      <section class="hero">
        <img src="/docs/img/logo.png" alt="pgmgo logo" class="hero-logo" />
        <h1>pgmgo</h1>
        <p class="hero-subtitle">Probabilistic Graphical Models in Go</p>
        <div class="badges">
          <span class="badge">Zero Dependencies</span>
          <span class="badge">Go</span>
          <span class="badge">MIT License</span>
          <span class="badge">pgmpy-Inspired</span>
        </div>
      </section>

      <section class="section">
        <h2>Installation</h2>
        <pre><code>go get github.com/asymmetric-effort/pgmgo</code></pre>
      </section>

      <section class="section">
        <h2>Quick Start</h2>
        <pre><code>{`package main

import (
    "fmt"
    "github.com/asymmetric-effort/pgmgo/src/models"
)

func main() {
    // Create a Bayesian Network
    bn := models.NewBayesianNetwork()
    bn.AddNode("Rain")
    bn.AddNode("Sprinkler")
    bn.AddNode("WetGrass")
    bn.AddEdge("Rain", "WetGrass")
    bn.AddEdge("Sprinkler", "WetGrass")
    fmt.Println(bn)
}`}</code></pre>
      </section>

      <section class="section">
        <h2>Features</h2>
        <div class="features-grid">
          <div class="feature-card">
            <h3>Bayesian Networks</h3>
            <p>Create, parameterize, and query Bayesian networks with discrete and continuous variables.</p>
          </div>
          <div class="feature-card">
            <h3>Markov Networks</h3>
            <p>Undirected graphical models with factor-based parameterization.</p>
          </div>
          <div class="feature-card">
            <h3>Exact Inference</h3>
            <p>Variable elimination and belief propagation algorithms.</p>
          </div>
          <div class="feature-card">
            <h3>Approximate Inference</h3>
            <p>Sampling-based methods including MCMC and importance sampling.</p>
          </div>
          <div class="feature-card">
            <h3>Structure Learning</h3>
            <p>Learn graph structure from data using score-based and constraint-based methods.</p>
          </div>
          <div class="feature-card">
            <h3>Zero Dependencies</h3>
            <p>Built entirely in Go with no third-party runtime dependencies.</p>
          </div>
        </div>
      </section>
    </div>
  );
}
