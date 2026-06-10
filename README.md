<p align="center">
  <img src="docs/img/logo.png" alt="datascience logo" width="150">
</p>

<h1 align="center">datascience</h1>

<p align="center">
  <strong>Data Science and Machine Learning in Go</strong><br>
  A zero-dependency Go library for probabilistic models, deep learning, numerical computing, and quantitative analysis
</p>

<p align="center">
  <a href="LICENSE.txt"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License"></a>
  <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.26+-00ADD8.svg" alt="Go Version"></a>
  <img src="https://img.shields.io/badge/dependencies-zero-brightgreen.svg" alt="Zero Dependencies">
  <a href="https://github.com/asymmetric-effort/datascience/actions"><img src="https://github.com/asymmetric-effort/datascience/actions/workflows/ci.yml/badge.svg" alt="CI Status"></a>
</p>

---

## Overview

datascience is a comprehensive pure Go library for data science and machine learning.
It provides probabilistic graphical models (full pgmpy parity), a TensorFlow-compatible
deep learning framework, BLAS-optimized numerical computing, and quantitative finance
tools — all with **zero third-party dependencies**.

Every numerical, graph-theoretic, statistical, and tabular operation is implemented
from scratch in Go within the project's internal libraries. This eliminates
supply-chain risk and simplifies deployment.

## Installation

```bash
go get github.com/asymmetric-effort/datascience
```

Requires Go 1.26 or later.

## Quick Start

Build a Bayesian network and run inference:

```go
package main

import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/example_models"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
)

func main() {
    bn := example_models.Student()
    if err := bn.CheckModel(); err != nil {
        log.Fatalf("Model validation failed: %v", err)
    }

    markovFactors, _ := bn.ToMarkovFactors()
    ve := inference.NewVariableElimination(markovFactors)
    evidence := map[string]int{"D": 0, "I": 1}
    result, _ := ve.Query([]string{"G"}, evidence)

    for i, state := range bn.GetStates("G") {
        prob := result.GetValue(map[string]int{"G": i})
        fmt.Printf("P(G=%s | D=Easy, I=High) = %.4f\n", state, prob)
    }
}
```

Train a neural network:

```go
package main

import (
    "github.com/asymmetric-effort/datascience/lib/numgo"
    "github.com/asymmetric-effort/datascience/lib/tensorflow/keras"
    "github.com/asymmetric-effort/datascience/lib/tensorflow/nn/layer"
    "github.com/asymmetric-effort/datascience/lib/tensorflow/nn/loss"
)

func main() {
    model := keras.NewSequential()
    model.Add(layer.NewDense(128, "relu"))
    model.Add(layer.NewDense(10, "softmax"))
    model.Compile(loss.CategoricalCrossEntropy, 0.001)

    X := numgo.NewNDArray([]int{100, 784}, nil) // training data
    Y := numgo.NewNDArray([]int{100, 10}, nil)  // labels
    model.Fit(X, Y, 10, 32) // epochs=10, batch=32
}
```

## Libraries

| Library          | Replaces           | Description                                                    |
|------------------|--------------------|----------------------------------------------------------------|
| **lib/numgo**    | NumPy              | N-dimensional arrays, broadcasting, linear algebra, BLAS L1/2/3 |
| **lib/scigo**    | SciPy              | Statistics, distributions, optimization, FFT, SDE solvers, Black-Scholes, portfolio optimization |
| **lib/tabgo**    | Pandas             | DataFrames, CSV I/O, filtering, aggregation, rolling analytics  |
| **lib/graphgo**  | NetworkX           | Graph data structures, algorithms, d-separation, moralization   |
| **lib/gpu**      | PyTorch/Pyro       | Compute backend abstraction (CPU fallback included)             |
| **lib/pgm**      | pgmpy              | Probabilistic graphical models — 13 model types, 7 inference algorithms, 11 learning algorithms |
| **lib/tensorflow**| TensorFlow/Keras  | Neural networks — Dense, Conv2D, LSTM, GRU, Attention, BatchNorm, optimizers, loss functions |

## Probabilistic Graphical Models (lib/pgm)

### Models (13 types)

Bayesian Network, Discrete Bayesian Network, Markov Network, Discrete Markov Network,
Dynamic Bayesian Network, Factor Graph, Cluster Graph, Junction Tree, Naive Bayes,
Markov Chain, Linear Gaussian BN, Functional BN, Structural Equation Model (SEM)

### Inference (7 algorithms)

Variable Elimination, Belief Propagation, MPLP, Approximate Inference,
Causal Inference (do-calculus, backdoor/frontdoor), Dynamic BN Inference, MAP/MPE queries

### Learning (11+ algorithms)

MLE, Bayesian Estimation, EM, Linear Gaussian MLE, SEM Estimation,
Hill Climb, Exhaustive Search, PC, GES, MMHC, Tree Search,
Expert-in-the-Loop, LLM-assisted discovery, IV estimation, Mirror Descent

### File I/O (7 formats)

BIF, XMLBIF, UAI, NET, XBN, XDSL, POMDPX

## Deep Learning (lib/tensorflow)

### Layers

Dense, Conv2D, LSTM, GRU, Attention, BatchNorm, Dropout, Embedding, Flatten, MaxPool2D

### Training

Sequential model, SGD/Adam/RMSProp optimizers, MSE/CrossEntropy/Huber loss,
callbacks, learning rate schedules, regularizers, metrics

### Utilities

GradientTape (automatic differentiation), Variable (trainable state),
model save/load, dataset loading, image processing, weight initializers

## Quantitative Finance (lib/scigo, lib/tabgo)

- Black-Scholes pricing (European calls/puts, Greeks, implied volatility)
- Binomial tree and Monte Carlo pricing
- SDE solvers (Euler-Maruyama, Milstein)
- Ito calculus (Brownian motion, GBM, Ornstein-Uhlenbeck)
- Markowitz mean-variance portfolio optimization
- QP solver with active-set method
- Rolling correlation, beta, alpha, R-squared, PCA
- Rolling Sharpe, Sortino, max drawdown, VaR, CVaR

## Project Structure

```
datascience/
  lib/
    numgo/             Numerical arrays and BLAS
    scigo/             Scientific computing
    graphgo/           Graph algorithms
    tabgo/             Tabular data and analytics
    gpu/               Compute backend
    pgm/               Probabilistic graphical models
      models/            13 model types
      inference/         7 inference algorithms
      learning/          11+ learning algorithms
      sampling/          Forward and Gibbs sampling
      readwrite/         7 file format readers/writers
      factors/           CPD/JPD representations
      metrics/           Scoring and evaluation
      ...
    tensorflow/        Deep learning
      keras/             Sequential model, training
      nn/                Layers, loss, optimizers, activations
      variable/          Trainable variables
      gradtape/          Automatic differentiation
      ...
  examples/            Runnable example programs
  example_models/      Pre-built canonical networks
  tests/               Cross-validation fixtures
  website/             Project website
```

## Documentation

- [Project Website](https://datascience.asymmetric-effort.com)
- [GoDoc](https://pkg.go.dev/github.com/asymmetric-effort/datascience)
- [Examples](examples/)

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) for
guidelines on development workflow, testing, commit conventions, and the
zero-dependency policy.

## Security

To report a security vulnerability, see [SECURITY.md](SECURITY.md). Do not open
public issues for security concerns.

## License

datascience is released under the [MIT License](LICENSE.txt).

Copyright (c) 2026 Asymmetric Effort, LLC.
