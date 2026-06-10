import { createElement, useHead, Link } from "@asymmetric-effort/specifyjs";
import { ScrollLink } from "../components/ScrollLink";

export function Libraries() {
  useHead({
    title: "Libraries — datascience",
    description: "Overview of all library packages in datascience: numgo, scigo, graphgo, tabgo, gpu, pgm, and tensorflow.",
    canonical: "https://pgmgo.asymmetric-effort.com/#/libraries",
  });

  return (
    <div class="page">
      <h1>Libraries</h1>

      <nav class="page-toc">
        <strong>Packages:</strong>{" "}
        <ScrollLink to="lib-overview">Overview</ScrollLink> | <ScrollLink to="lib-numgo">numgo</ScrollLink> | <ScrollLink to="lib-scigo">scigo</ScrollLink> | <ScrollLink to="lib-graphgo">graphgo</ScrollLink> | <ScrollLink to="lib-tabgo">tabgo</ScrollLink> | <ScrollLink to="lib-gpu">gpu</ScrollLink> | <ScrollLink to="lib-pgm">pgm</ScrollLink> | <ScrollLink to="lib-tensorflow">tensorflow</ScrollLink>
      </nav>

      <section class="section" id="lib-overview">
        <h2>Overview</h2>
        <p>
          The datascience library is organized into independent, composable packages under <code>lib/</code>.
          Each package can be imported and used on its own. The primitive packages (numgo, scigo, graphgo, tabgo)
          provide general-purpose data science foundations, while the domain packages (pgm, tensorflow)
          build on them for specific use cases.
        </p>
        <pre><code>{`go get github.com/asymmetric-effort/datascience`}</code></pre>
        <table>
          <thead>
            <tr><th>Package</th><th>Import Path</th><th>Description</th><th>Python Equivalent</th></tr>
          </thead>
          <tbody>
            <tr><td><strong>numgo</strong></td><td><code>lib/numgo</code></td><td>N-dimensional arrays, matrices, vectors, BLAS</td><td>numpy</td></tr>
            <tr><td><strong>scigo</strong></td><td><code>lib/scigo</code></td><td>Distributions, optimization, special functions</td><td>scipy</td></tr>
            <tr><td><strong>graphgo</strong></td><td><code>lib/graphgo</code></td><td>Directed/undirected graphs, algorithms</td><td>networkx</td></tr>
            <tr><td><strong>tabgo</strong></td><td><code>lib/tabgo</code></td><td>DataFrames, Series, CSV/Parquet I/O</td><td>pandas</td></tr>
            <tr><td><strong>gpu</strong></td><td><code>lib/gpu</code></td><td>GPU compute backend</td><td>cupy / CUDA</td></tr>
            <tr><td><strong>pgm</strong></td><td><code>lib/pgm/...</code></td><td>Probabilistic graphical models</td><td>pgmpy</td></tr>
            <tr><td><strong>tensorflow</strong></td><td><code>lib/tensorflow</code></td><td>Deep learning</td><td>tensorflow</td></tr>
          </tbody>
        </table>
      </section>

      <section class="section" id="lib-numgo">
        <h2>numgo -- Arrays, Matrices, and BLAS</h2>
        <p>Import: <code>github.com/asymmetric-effort/datascience/lib/numgo</code></p>
        <p>
          Pure Go implementation of N-dimensional arrays, 2D matrices, and 1D vectors with
          BLAS-level linear algebra operations. Provides the numerical foundation for all
          other packages.
        </p>

        <h3>Key Types</h3>
        <table>
          <thead>
            <tr><th>Type</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>NDArray</code></td><td>N-dimensional array with shape, stride, broadcasting, element-wise operations</td></tr>
            <tr><td><code>Matrix</code></td><td>2D matrix with multiply, transpose, inverse, determinant, eigenvalues, decompositions</td></tr>
            <tr><td><code>Vector</code></td><td>1D vector with dot product, norm, cross product, element-wise operations</td></tr>
          </tbody>
        </table>

        <h3>Example</h3>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/lib/numgo"

// Matrix operations
m := numgo.NewMatrixFromData(3, 3, []float64{
    1, 2, 3,
    0, 1, 4,
    5, 6, 0,
})
det := m.Det()           // determinant
inv := m.Inverse()       // matrix inverse
t := m.Transpose()       // transpose

// Vector operations
v1 := numgo.NewVector([]float64{1, 2, 3})
v2 := numgo.NewVector([]float64{4, 5, 6})
dot := v1.Dot(v2)       // 32
norm := v1.Norm()        // sqrt(14)

// NDArray operations
arr := numgo.NewNDArray([]int{2, 3, 4})
arr.Fill(1.0)
sum := arr.Sum()         // 24`}</code></pre>
      </section>

      <section class="section" id="lib-scigo">
        <h2>scigo -- Statistical Computing</h2>
        <p>Import: <code>github.com/asymmetric-effort/datascience/lib/scigo</code></p>
        <p>
          Statistical distributions (Normal, Chi-Squared, Beta, Gamma, Student-t, Uniform, Exponential),
          optimization routines (gradient descent, Newton's method, minimization), special functions
          (gamma, beta, digamma), and hypothesis tests.
        </p>

        <h3>Key Capabilities</h3>
        <table>
          <thead>
            <tr><th>Category</th><th>Functions / Types</th></tr>
          </thead>
          <tbody>
            <tr><td>Distributions</td><td><code>Normal</code>, <code>ChiSquared</code>, <code>Beta</code>, <code>Gamma</code>, <code>StudentT</code>, <code>Uniform</code>, <code>Exponential</code> -- each with PDF, CDF, PPF, Sample</td></tr>
            <tr><td>Optimization</td><td><code>Minimize</code>, <code>GradientDescent</code>, <code>NewtonMethod</code></td></tr>
            <tr><td>Statistics</td><td><code>Mean</code>, <code>Std</code>, <code>PearsonCorrelation</code>, p-value computation</td></tr>
            <tr><td>Special Functions</td><td>Gamma, beta, incomplete gamma, digamma, log-gamma</td></tr>
          </tbody>
        </table>

        <h3>Example</h3>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/lib/scigo"

// Normal distribution
n := scigo.NewNormal(0, 1)
cdf := n.CDF(1.96)       // ~0.975
ppf := n.PPF(0.975)      // ~1.96

// Chi-squared test
chi2 := scigo.NewChiSquared(5)
pValue := 1.0 - chi2.CDF(11.07)

// Optimization
result := scigo.Minimize(func(x float64) float64 {
    return (x - 3) * (x - 3)
}, 0.0, 10.0)  // ~3.0`}</code></pre>
      </section>

      <section class="section" id="lib-graphgo">
        <h2>graphgo -- Graph Algorithms</h2>
        <p>Import: <code>github.com/asymmetric-effort/datascience/lib/graphgo</code></p>
        <p>
          Directed and undirected graphs with a full suite of algorithms: topological sort,
          d-separation, moral graph, triangulation, maximum cardinality search, clique finding,
          connected components, shortest paths, and cycle detection.
        </p>

        <h3>Key Types</h3>
        <table>
          <thead>
            <tr><th>Type</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>DiGraph</code></td><td>Directed graph with successors, predecessors, topological sort, d-separation</td></tr>
            <tr><td><code>Graph</code></td><td>Undirected graph with neighbors, degree, connected components</td></tr>
            <tr><td><code>PDAG</code></td><td>Partially directed acyclic graph for equivalence classes</td></tr>
          </tbody>
        </table>

        <h3>Example</h3>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/lib/graphgo"

g := graphgo.NewDiGraph()
g.AddNode("A")
g.AddNode("B")
g.AddNode("C")
g.AddEdge("A", "B")
g.AddEdge("B", "C")

sorted := g.TopologicalSort()  // [A, B, C]
hasCycle := g.HasCycle()       // false`}</code></pre>
      </section>

      <section class="section" id="lib-tabgo">
        <h2>tabgo -- Tabular Data</h2>
        <p>Import: <code>github.com/asymmetric-effort/datascience/lib/tabgo</code></p>
        <p>
          DataFrames and Series for tabular data manipulation. Supports CSV, Parquet, and Excel I/O,
          row filtering, groupby, merge, column operations, and statistical summaries.
        </p>

        <h3>Key Types</h3>
        <table>
          <thead>
            <tr><th>Type</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>DataFrame</code></td><td>Tabular data with named columns, filtering, groupby, merge</td></tr>
            <tr><td><code>Series</code></td><td>Single column with value counts, unique values, statistics</td></tr>
          </tbody>
        </table>

        <h3>I/O Functions</h3>
        <table>
          <thead>
            <tr><th>Function</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>ReadCSV / WriteCSV</code></td><td>CSV file I/O</td></tr>
            <tr><td><code>ReadParquet / WriteParquet</code></td><td>Apache Parquet I/O</td></tr>
            <tr><td><code>ReadXLSX</code></td><td>Excel file reading</td></tr>
          </tbody>
        </table>

        <h3>Example</h3>
        <pre><code>{`import "github.com/asymmetric-effort/datascience/lib/tabgo"

df, _ := tabgo.ReadCSV("data.csv")
fmt.Printf("%d rows, %d columns\\n", df.NRows(), len(df.Columns()))

col := df.Column("Temperature")
fmt.Println("Unique:", col.Unique())

filtered := df.Filter(func(row map[string]interface{}) bool {
    return row["Temperature"].(int) > 70
})
tabgo.WriteCSV(filtered, "warm_days.csv")`}</code></pre>
      </section>

      <section class="section" id="lib-gpu">
        <h2>gpu -- GPU Compute Backend</h2>
        <p>Import: <code>github.com/asymmetric-effort/datascience/lib/gpu</code></p>
        <p>
          Optional GPU acceleration for compute-intensive operations. Provides GPU-backed
          matrix operations, tensor computations, and factor operations that can significantly
          speed up inference, learning, and deep learning training on large models.
        </p>
        <p>
          The GPU package is designed to be a drop-in accelerator. Other packages detect GPU
          availability and automatically offload heavy computations when a GPU is present.
        </p>
      </section>

      <section class="section" id="lib-pgm">
        <h2>pgm -- Probabilistic Graphical Models</h2>
        <p>Import: <code>github.com/asymmetric-effort/datascience/lib/pgm/...</code></p>
        <p>
          A complete PGM toolkit inspired by <a href="https://pgmpy.org" target="_blank" rel="noopener noreferrer">pgmpy</a>.
          Provides 13 model types, 7 inference algorithms, 15+ learning algorithms, 16 conditional
          independence tests, 13 scoring functions, 10 file formats, causal inference, and more.
        </p>

        <h3>Sub-packages</h3>
        <table>
          <thead>
            <tr><th>Package</th><th>Description</th></tr>
          </thead>
          <tbody>
            <tr><td><code>lib/pgm/base</code></td><td>DAG, PDAG, UndirectedGraph, ADMG, MAG, SimpleCausalModel</td></tr>
            <tr><td><code>lib/pgm/models</code></td><td>13 probabilistic model types (BayesianNetwork, MarkovNetwork, etc.)</td></tr>
            <tr><td><code>lib/pgm/factors</code></td><td>TabularCPD, DiscreteFactor, LinearGaussianCPD, NoisyOR, JPD</td></tr>
            <tr><td><code>lib/pgm/inference</code></td><td>Variable Elimination, Belief Propagation, MPLP, Causal, Approximate, DBN</td></tr>
            <tr><td><code>lib/pgm/sampling</code></td><td>Forward, rejection, likelihood-weighted, Gibbs sampling</td></tr>
            <tr><td><code>lib/pgm/learning</code></td><td>MLE, Bayesian, EM, HillClimb, PC, GES, ExhaustiveSearch, TreeSearch, MMHC</td></tr>
            <tr><td><code>lib/pgm/ci_tests</code></td><td>ChiSquare, FisherZ, GCM, and 13 more CI tests</td></tr>
            <tr><td><code>lib/pgm/structure_score</code></td><td>BIC, AIC, BDeu, BDs, K2, LogLikelihood, Gaussian scores</td></tr>
            <tr><td><code>lib/pgm/identification</code></td><td>Back-door and front-door causal identification</td></tr>
            <tr><td><code>lib/pgm/prediction</code></td><td>DoubleML, naive adjustment, IV regression</td></tr>
            <tr><td><code>lib/pgm/metrics</code></td><td>SHD, confusion matrices, Fisher's C, log-likelihood</td></tr>
            <tr><td><code>lib/pgm/readwrite</code></td><td>BIF, XMLBIF, NET, UAI, XDSL, JSON, CSV, XML I/O</td></tr>
            <tr><td><code>lib/pgm/independencies</code></td><td>Independence assertion representations</td></tr>
            <tr><td><code>lib/pgm/config</code></td><td>Global configuration</td></tr>
            <tr><td><code>lib/pgm/utils</code></td><td>Shared utilities</td></tr>
          </tbody>
        </table>

        <h3>Quick Example</h3>
        <pre><code>{`import (
    "github.com/asymmetric-effort/datascience/lib/pgm/models"
    "github.com/asymmetric-effort/datascience/lib/pgm/factors"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
)

bn := models.NewBayesianNetwork()
bn.AddNode("A")
bn.AddNode("B")
bn.AddEdge("A", "B")
bn.SetStates("A", []string{"a0", "a1"})
bn.SetStates("B", []string{"b0", "b1"})
bn.SetCPD("A", factors.NewTabularCPD("A", 2, []float64{0.6, 0.4}, nil, nil))
bn.SetCPD("B", factors.NewTabularCPD("B", 2, []float64{0.2, 0.8, 0.75, 0.25}, []string{"A"}, []int{2}))

facs, _ := bn.ToMarkovFactors()
ve := inference.NewVariableElimination(facs)
result, _ := ve.Query([]string{"B"}, map[string]int{"A": 1})
fmt.Println("P(B | A=1):", result.Values().Data())`}</code></pre>
        <p>See the <Link to="/docs">Documentation</Link> and <Link to="/api">API Reference</Link> for complete details.</p>
      </section>

      <section class="section" id="lib-tensorflow">
        <h2>tensorflow -- Deep Learning</h2>
        <p>Import: <code>github.com/asymmetric-effort/datascience/lib/tensorflow</code></p>
        <p>
          TensorFlow-compatible deep learning in pure Go. Provides neural network construction,
          training, and inference capabilities. Builds on <code>numgo</code> for tensor operations
          and <code>gpu</code> for optional hardware acceleration.
        </p>
        <p>
          The tensorflow package brings deep learning capabilities to the datascience library,
          complementing the probabilistic graphical models in <code>lib/pgm</code>. Together,
          they cover both classical statistical modeling and modern neural network approaches.
        </p>
      </section>
    </div>
  );
}
