// Package datasets provides well-known Bayesian network datasets as embedded
// CSV data, comparable to pgmpy's built-in datasets. Each loader function
// returns a *tabgo.DataFrame ready for use in structure learning, parameter
// estimation, or inference examples.
package datasets

import (
	_ "embed"
	"fmt"
	"sort"

	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// ===== Original 8 datasets =====

//go:embed data/asia.csv
var asiaData string

//go:embed data/alarm.csv
var alarmData string

//go:embed data/sachs.csv
var sachsData string

//go:embed data/cancer.csv
var cancerData string

//go:embed data/student.csv
var studentData string

//go:embed data/sprinkler.csv
var sprinklerData string

//go:embed data/survey.csv
var surveyData string

//go:embed data/titanic.csv
var titanicData string

// ===== Well-known ML datasets =====

//go:embed data/adult.csv
var adultData string

//go:embed data/pima_diabetes.csv
var pimaDiabetesData string

//go:embed data/iris.csv
var irisData string

//go:embed data/wine.csv
var wineData string

//go:embed data/heart.csv
var heartData string

//go:embed data/boston.csv
var bostonData string

//go:embed data/breast_cancer.csv
var breastCancerData string

// ===== BN-specific datasets =====

//go:embed data/earthquake.csv
var earthquakeData string

//go:embed data/child.csv
var childData string

//go:embed data/insurance.csv
var insuranceData string

//go:embed data/water.csv
var waterData string

//go:embed data/mildew.csv
var mildewData string

//go:embed data/hailfinder.csv
var hailfinderData string

//go:embed data/hepar2.csv
var hepar2Data string

//go:embed data/lucas.csv
var lucasData string

//go:embed data/andes.csv
var andesData string

//go:embed data/munin.csv
var muninData string

//go:embed data/barley.csv
var barleyData string

//go:embed data/win95pts.csv
var win95ptsData string

// ===== UCI datasets =====

//go:embed data/ecoli.csv
var ecoliData string

//go:embed data/glass.csv
var glassData string

//go:embed data/zoo.csv
var zooData string

//go:embed data/mushroom.csv
var mushroomData string

//go:embed data/nursery.csv
var nurseryData string

//go:embed data/car_evaluation.csv
var carEvaluationData string

//go:embed data/balance_scale.csv
var balanceScaleData string

//go:embed data/monks.csv
var monksData string

//go:embed data/tic_tac_toe.csv
var ticTacToeData string

//go:embed data/vote.csv
var voteData string

//go:embed data/credit_approval.csv
var creditApprovalData string

//go:embed data/hepatitis.csv
var hepatitisData string

//go:embed data/automobile.csv
var automobileData string

// ===== Original 8 loader functions =====

// Asia returns the Asia (Lauritzen-Spiegelhalter) network dataset.
// 1000 rows, 8 binary columns: asia, tub, smoke, lung, bronc, either, xray, dysp.
func Asia() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(asiaData)
}

// Alarm returns the Alarm (Burglary) network dataset.
// 1000 rows, 5 binary columns: Burglary, Earthquake, Alarm, JohnCalls, MaryCalls.
func Alarm() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(alarmData)
}

// Sachs returns the Sachs protein signaling network dataset.
// 500 rows, 11 columns discretized to 3 levels (0=low, 1=medium, 2=high):
// Raf, Mek, Plcg, PIP2, PIP3, Erk, Akt, PKA, PKC, P38, Jnk.
func Sachs() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(sachsData)
}

// Cancer returns the Cancer network dataset.
// 1000 rows, 5 binary columns: Pollution, Smoker, Cancer, Xray, Dyspnoea.
func Cancer() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(cancerData)
}

// Student returns the Student network dataset.
// 1000 rows, 5 columns: D (binary), I (binary), G (ternary 0/1/2), L (binary), S (binary).
func Student() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(studentData)
}

// Sprinkler returns the Water Sprinkler network dataset.
// 1000 rows, 4 binary columns: Cloudy, Sprinkler, Rain, WetGrass.
func Sprinkler() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(sprinklerData)
}

// Survey returns the Survey network dataset.
// 500 rows, 5 categorical columns: Age, Education, Occupation, Residence, Transportation.
func Survey() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(surveyData)
}

// Titanic returns the Titanic survival dataset.
// 800 rows, 4 categorical columns: Class, Sex, Age, Survived.
func Titanic() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(titanicData)
}

// ===== Well-known ML dataset loaders =====

// Adult returns the Adult (Census Income) dataset.
// 500 rows, 11 discretized columns including Age, Workclass, Education, etc.
func Adult() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(adultData)
}

// PimaDiabetes returns the Pima Indians Diabetes dataset.
// 300 rows, 9 discretized columns: Pregnancies, Glucose, BloodPressure, etc.
func PimaDiabetes() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(pimaDiabetesData)
}

// Iris returns Fisher's Iris dataset.
// 150 rows, 5 columns: SepalLength, SepalWidth, PetalLength, PetalWidth, Species.
func Iris() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(irisData)
}

// Wine returns the Wine Quality dataset.
// 300 rows, 12 discretized columns including acidity, sugar, alcohol, and quality.
func Wine() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(wineData)
}

// Heart returns the Heart Disease dataset.
// 300 rows, 14 discretized columns including Age, ChestPain, Cholesterol, Target.
func Heart() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(heartData)
}

// Boston returns the Boston Housing dataset.
// 300 rows, 14 discretized columns including CRIM, RM, LSTAT, MEDV.
func Boston() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(bostonData)
}

// BreastCancer returns the Breast Cancer Wisconsin dataset.
// 300 rows, 10 columns: 9 cell features + Diagnosis.
func BreastCancer() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(breastCancerData)
}

// ===== BN-specific dataset loaders =====

// Earthquake returns the Earthquake Bayesian network dataset.
// 500 rows, 5 binary columns: Burglary, Earthquake, Alarm, Radio, Calls.
func Earthquake() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(earthquakeData)
}

// Child returns the Child Bayesian network dataset.
// 500 rows, 20 columns representing childhood disease diagnosis variables.
func Child() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(childData)
}

// Insurance returns the Insurance Bayesian network dataset.
// 500 rows, 27 columns representing automobile insurance risk variables.
func Insurance() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(insuranceData)
}

// Water returns the Water treatment plant Bayesian network dataset.
// 500 rows, 32 columns representing water quality measurements.
func Water() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(waterData)
}

// Mildew returns the Mildew crop disease Bayesian network dataset.
// 500 rows, 35 columns representing crop disease and environmental variables.
func Mildew() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(mildewData)
}

// Hailfinder returns the Hailfinder weather prediction Bayesian network dataset.
// 500 rows, 56 columns representing meteorological variables.
func Hailfinder() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(hailfinderData)
}

// Hepar2 returns the Hepar2 liver disease Bayesian network dataset.
// 500 rows, 40 columns representing liver disease diagnosis variables.
func Hepar2() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(hepar2Data)
}

// Lucas returns the Lucas (Lung Cancer Simple) Bayesian network dataset.
// 500 rows, 12 binary columns: Smoking, YellowFingers, Anxiety, etc.
func Lucas() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(lucasData)
}

// Andes returns the Andes educational tutoring Bayesian network dataset.
// 500 rows, 30 columns representing student knowledge and skill variables.
func Andes() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(andesData)
}

// Munin returns the Munin neuromuscular disorder Bayesian network dataset.
// 500 rows, 30 columns representing nerve and muscle diagnostic variables.
func Munin() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(muninData)
}

// Barley returns the Barley crop yield Bayesian network dataset.
// 500 rows, 30 columns representing crop, soil, and weather variables.
func Barley() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(barleyData)
}

// Win95pts returns the Win95pts printer troubleshooting Bayesian network dataset.
// 500 rows, 30 columns representing printer diagnostic variables.
func Win95pts() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(win95ptsData)
}

// ===== UCI dataset loaders =====

// Ecoli returns the E. coli protein localization dataset.
// 336 rows, 8 columns: 7 features + Site.
func Ecoli() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(ecoliData)
}

// Glass returns the Glass Identification dataset.
// 214 rows, 10 columns: 9 features + Type.
func Glass() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(glassData)
}

// Zoo returns the Zoo animal classification dataset.
// 101 rows, 17 columns: 16 boolean features + Type.
func Zoo() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(zooData)
}

// Mushroom returns the Mushroom classification dataset.
// 500 rows, 23 columns: 22 categorical features + Class.
func Mushroom() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(mushroomData)
}

// Nursery returns the Nursery school evaluation dataset.
// 500 rows, 9 columns: 8 categorical features + Class.
func Nursery() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(nurseryData)
}

// CarEvaluation returns the Car Evaluation dataset.
// 400 rows, 7 columns: 6 categorical features + Class.
func CarEvaluation() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(carEvaluationData)
}

// BalanceScale returns the Balance Scale dataset.
// 300 rows, 5 columns: LeftWeight, LeftDist, RightWeight, RightDist, Class.
func BalanceScale() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(balanceScaleData)
}

// Monks returns the MONKS problem dataset.
// 300 rows, 7 columns: 6 features + Class.
func Monks() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(monksData)
}

// TicTacToe returns the Tic-Tac-Toe endgame dataset.
// 400 rows, 10 columns: 9 board positions + Class.
func TicTacToe() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(ticTacToeData)
}

// Vote returns the Congressional Voting Records dataset.
// 435 rows, 17 columns: 16 vote features + Party.
func Vote() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(voteData)
}

// CreditApproval returns the Credit Approval dataset.
// 300 rows, 16 columns: 15 features + Approved.
func CreditApproval() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(creditApprovalData)
}

// Hepatitis returns the Hepatitis dataset.
// 155 rows, 20 columns: 19 features + Class.
func Hepatitis() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(hepatitisData)
}

// Automobile returns the Automobile dataset.
// 205 rows, 26 columns: 25 features + Price.
func Automobile() (*tabgo.DataFrame, error) {
	return tabgo.ReadCSVFromString(automobileData)
}

// registry maps dataset names to loader functions.
var registry = map[string]func() (*tabgo.DataFrame, error){
	// Original 8
	"asia":      Asia,
	"alarm":     Alarm,
	"sachs":     Sachs,
	"cancer":    Cancer,
	"student":   Student,
	"sprinkler": Sprinkler,
	"survey":    Survey,
	"titanic":   Titanic,
	// Well-known ML
	"adult":         Adult,
	"pima_diabetes": PimaDiabetes,
	"iris":          Iris,
	"wine":          Wine,
	"heart":         Heart,
	"boston":        Boston,
	"breast_cancer": BreastCancer,
	// BN-specific
	"earthquake": Earthquake,
	"child":      Child,
	"insurance":  Insurance,
	"water":      Water,
	"mildew":     Mildew,
	"hailfinder": Hailfinder,
	"hepar2":     Hepar2,
	"lucas":      Lucas,
	"andes":      Andes,
	"munin":      Munin,
	"barley":     Barley,
	"win95pts":   Win95pts,
	// UCI
	"ecoli":           Ecoli,
	"glass":           Glass,
	"zoo":             Zoo,
	"mushroom":        Mushroom,
	"nursery":         Nursery,
	"car_evaluation":  CarEvaluation,
	"balance_scale":   BalanceScale,
	"monks":           Monks,
	"tic_tac_toe":     TicTacToe,
	"vote":            Vote,
	"credit_approval": CreditApproval,
	"hepatitis":       Hepatitis,
	"automobile":      Automobile,
}

// List returns the names of all available datasets in sorted order.
func List() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Load loads a dataset by name. The name is case-sensitive and must match
// one of the names returned by List.
func Load(name string) (*tabgo.DataFrame, error) {
	fn, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("datasets: unknown dataset %q; available: %v", name, List())
	}
	return fn()
}
