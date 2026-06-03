//go:build ignore
// +build ignore

// Generator for dataset CSV files. Run with: go run examples/datasets/gen/main.go
package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	dir := filepath.Join("examples", "datasets", "data")
	os.MkdirAll(dir, 0755)

	rng := rand.New(rand.NewSource(42))

	// Original 8 datasets
	genAsia(rng, dir)
	genAlarm(rng, dir)
	genSachs(rng, dir)
	genCancer(rng, dir)
	genStudent(rng, dir)
	genSprinkler(rng, dir)
	genSurvey(rng, dir)
	genTitanic(rng, dir)

	// Well-known ML datasets
	genAdult(rng, dir)
	genPimaDiabetes(rng, dir)
	genIris(rng, dir)
	genWine(rng, dir)
	genHeart(rng, dir)
	genBoston(rng, dir)
	genBreastCancer(rng, dir)

	// BN-specific datasets
	genEarthquake(rng, dir)
	genChild(rng, dir)
	genInsurance(rng, dir)
	genWater(rng, dir)
	genMildew(rng, dir)
	genHailfinder(rng, dir)
	genHepar2(rng, dir)
	genLucas(rng, dir)
	genAndes(rng, dir)
	genMunin(rng, dir)
	genBarley(rng, dir)
	genWin95pts(rng, dir)

	// UCI datasets
	genEcoli(rng, dir)
	genGlass(rng, dir)
	genZoo(rng, dir)
	genMushroom(rng, dir)
	genNursery(rng, dir)
	genCarEvaluation(rng, dir)
	genBalanceScale(rng, dir)
	genMonks(rng, dir)
	genTicTacToe(rng, dir)
	genVote(rng, dir)
	genCreditApproval(rng, dir)
	genHepatitis(rng, dir)
	genAutomobile(rng, dir)

	fmt.Println("All datasets generated.")
}

func writeCSV(path string, headers []string, rows [][]string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Write(headers)
	for _, r := range rows {
		w.Write(r)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		panic(err)
	}
	fmt.Printf("Wrote %s (%d rows)\n", path, len(rows))
}

func itoa(v int) string { return strconv.Itoa(v) }

func ftoa(v float64) string { return strconv.FormatFloat(v, 'f', 2, 64) }

func bern(rng *rand.Rand, p float64) int {
	if rng.Float64() < p {
		return 1
	}
	return 0
}

func categorical(rng *rand.Rand, probs []float64) int {
	r := rng.Float64()
	cum := 0.0
	for i, p := range probs {
		cum += p
		if r < cum {
			return i
		}
	}
	return len(probs) - 1
}

// clampF clamps a float64 to [lo, hi].
func clampF(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// normalSample returns a sample from N(mean, stddev).
func normalSample(rng *rand.Rand, mean, stddev float64) float64 {
	return rng.NormFloat64()*stddev + mean
}

// categoricalStr picks a string from choices with given probs.
func categoricalStr(rng *rand.Rand, choices []string, probs []float64) string {
	return choices[categorical(rng, probs)]
}

// ===== Original 8 Datasets =====

func genAsia(rng *rand.Rand, dir string) {
	n := 1000
	headers := []string{"asia", "tub", "smoke", "lung", "bronc", "either", "xray", "dysp"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		asia := bern(rng, 0.01)
		tub := 0
		if asia == 1 {
			tub = bern(rng, 0.05)
		} else {
			tub = bern(rng, 0.01)
		}
		smoke := bern(rng, 0.50)
		lung := 0
		if smoke == 1 {
			lung = bern(rng, 0.10)
		} else {
			lung = bern(rng, 0.01)
		}
		bronc := 0
		if smoke == 1 {
			bronc = bern(rng, 0.60)
		} else {
			bronc = bern(rng, 0.30)
		}
		either := 0
		if tub == 1 || lung == 1 {
			either = 1
		}
		xray := 0
		if either == 1 {
			xray = bern(rng, 0.98)
		} else {
			xray = bern(rng, 0.05)
		}
		dysp := 0
		if either == 1 && bronc == 1 {
			dysp = bern(rng, 0.90)
		} else if either == 1 && bronc == 0 {
			dysp = bern(rng, 0.70)
		} else if either == 0 && bronc == 1 {
			dysp = bern(rng, 0.80)
		} else {
			dysp = bern(rng, 0.10)
		}
		rows[i] = []string{itoa(asia), itoa(tub), itoa(smoke), itoa(lung), itoa(bronc), itoa(either), itoa(xray), itoa(dysp)}
	}
	writeCSV(filepath.Join(dir, "asia.csv"), headers, rows)
}

func genAlarm(rng *rand.Rand, dir string) {
	n := 1000
	headers := []string{"Burglary", "Earthquake", "Alarm", "JohnCalls", "MaryCalls"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		burg := bern(rng, 0.001)
		eq := bern(rng, 0.002)
		alarm := 0
		if burg == 1 && eq == 1 {
			alarm = bern(rng, 0.95)
		} else if burg == 1 {
			alarm = bern(rng, 0.94)
		} else if eq == 1 {
			alarm = bern(rng, 0.29)
		} else {
			alarm = bern(rng, 0.001)
		}
		john := 0
		if alarm == 1 {
			john = bern(rng, 0.90)
		} else {
			john = bern(rng, 0.05)
		}
		mary := 0
		if alarm == 1 {
			mary = bern(rng, 0.70)
		} else {
			mary = bern(rng, 0.01)
		}
		rows[i] = []string{itoa(burg), itoa(eq), itoa(alarm), itoa(john), itoa(mary)}
	}
	writeCSV(filepath.Join(dir, "alarm.csv"), headers, rows)
}

func genSachs(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Raf", "Mek", "Plcg", "PIP2", "PIP3", "Erk", "Akt", "PKA", "PKC", "P38", "Jnk"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		plcg := categorical(rng, []float64{0.30, 0.40, 0.30})
		pkc := categorical(rng, []float64{0.40, 0.35, 0.25})
		pip3Probs := [][]float64{
			{0.60, 0.30, 0.10},
			{0.25, 0.50, 0.25},
			{0.10, 0.30, 0.60},
		}
		pip3 := categorical(rng, pip3Probs[plcg])
		pip2Base := (plcg + pip3)
		pip2 := 0
		if pip2Base <= 1 {
			pip2 = categorical(rng, []float64{0.60, 0.30, 0.10})
		} else if pip2Base <= 3 {
			pip2 = categorical(rng, []float64{0.20, 0.50, 0.30})
		} else {
			pip2 = categorical(rng, []float64{0.10, 0.30, 0.60})
		}
		pkaProbs := [][]float64{
			{0.50, 0.35, 0.15},
			{0.25, 0.45, 0.30},
			{0.10, 0.30, 0.60},
		}
		pka := categorical(rng, pkaProbs[pkc])
		rafBase := (pkc + pka)
		raf := 0
		if rafBase <= 1 {
			raf = categorical(rng, []float64{0.55, 0.30, 0.15})
		} else if rafBase <= 3 {
			raf = categorical(rng, []float64{0.20, 0.50, 0.30})
		} else {
			raf = categorical(rng, []float64{0.10, 0.25, 0.65})
		}
		mekBase := (raf + pkc + pka)
		mek := 0
		if mekBase <= 2 {
			mek = categorical(rng, []float64{0.55, 0.30, 0.15})
		} else if mekBase <= 4 {
			mek = categorical(rng, []float64{0.20, 0.50, 0.30})
		} else {
			mek = categorical(rng, []float64{0.10, 0.25, 0.65})
		}
		erkBase := (mek + pka)
		erk := 0
		if erkBase <= 1 {
			erk = categorical(rng, []float64{0.55, 0.30, 0.15})
		} else if erkBase <= 3 {
			erk = categorical(rng, []float64{0.20, 0.50, 0.30})
		} else {
			erk = categorical(rng, []float64{0.10, 0.25, 0.65})
		}
		aktBase := (pka + erk)
		akt := 0
		if aktBase <= 1 {
			akt = categorical(rng, []float64{0.55, 0.30, 0.15})
		} else if aktBase <= 3 {
			akt = categorical(rng, []float64{0.20, 0.50, 0.30})
		} else {
			akt = categorical(rng, []float64{0.10, 0.25, 0.65})
		}
		p38Base := (pkc + pka)
		p38 := 0
		if p38Base <= 1 {
			p38 = categorical(rng, []float64{0.55, 0.30, 0.15})
		} else if p38Base <= 3 {
			p38 = categorical(rng, []float64{0.20, 0.50, 0.30})
		} else {
			p38 = categorical(rng, []float64{0.10, 0.25, 0.65})
		}
		jnkBase := (pkc + pka)
		jnk := 0
		if jnkBase <= 1 {
			jnk = categorical(rng, []float64{0.55, 0.30, 0.15})
		} else if jnkBase <= 3 {
			jnk = categorical(rng, []float64{0.20, 0.50, 0.30})
		} else {
			jnk = categorical(rng, []float64{0.10, 0.25, 0.65})
		}
		rows[i] = []string{itoa(raf), itoa(mek), itoa(plcg), itoa(pip2), itoa(pip3), itoa(erk), itoa(akt), itoa(pka), itoa(pkc), itoa(p38), itoa(jnk)}
	}
	writeCSV(filepath.Join(dir, "sachs.csv"), headers, rows)
}

func genCancer(rng *rand.Rand, dir string) {
	n := 1000
	headers := []string{"Pollution", "Smoker", "Cancer", "Xray", "Dyspnoea"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		poll := bern(rng, 0.10)
		smoker := bern(rng, 0.30)
		cancer := 0
		if poll == 1 && smoker == 1 {
			cancer = bern(rng, 0.05)
		} else if poll == 1 {
			cancer = bern(rng, 0.03)
		} else if smoker == 1 {
			cancer = bern(rng, 0.02)
		} else {
			cancer = bern(rng, 0.001)
		}
		xray := 0
		if cancer == 1 {
			xray = bern(rng, 0.90)
		} else {
			xray = bern(rng, 0.20)
		}
		dysp := 0
		if cancer == 1 {
			dysp = bern(rng, 0.65)
		} else {
			dysp = bern(rng, 0.30)
		}
		rows[i] = []string{itoa(poll), itoa(smoker), itoa(cancer), itoa(xray), itoa(dysp)}
	}
	writeCSV(filepath.Join(dir, "cancer.csv"), headers, rows)
}

func genStudent(rng *rand.Rand, dir string) {
	n := 1000
	headers := []string{"D", "I", "G", "L", "S"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		d := bern(rng, 0.40)
		intl := bern(rng, 0.30)
		g := 0
		if d == 0 && intl == 0 {
			g = categorical(rng, []float64{0.30, 0.40, 0.30})
		} else if d == 0 && intl == 1 {
			g = categorical(rng, []float64{0.05, 0.25, 0.70})
		} else if d == 1 && intl == 0 {
			g = categorical(rng, []float64{0.50, 0.35, 0.15})
		} else {
			g = categorical(rng, []float64{0.20, 0.40, 0.40})
		}
		l := 0
		if g == 2 {
			l = bern(rng, 0.90)
		} else if g == 1 {
			l = bern(rng, 0.40)
		} else {
			l = bern(rng, 0.10)
		}
		s := 0
		if intl == 1 {
			s = bern(rng, 0.80)
		} else {
			s = bern(rng, 0.40)
		}
		rows[i] = []string{itoa(d), itoa(intl), itoa(g), itoa(l), itoa(s)}
	}
	writeCSV(filepath.Join(dir, "student.csv"), headers, rows)
}

func genSprinkler(rng *rand.Rand, dir string) {
	n := 1000
	headers := []string{"Cloudy", "Sprinkler", "Rain", "WetGrass"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		cloudy := bern(rng, 0.50)
		sprinkler := 0
		if cloudy == 1 {
			sprinkler = bern(rng, 0.10)
		} else {
			sprinkler = bern(rng, 0.50)
		}
		rain := 0
		if cloudy == 1 {
			rain = bern(rng, 0.80)
		} else {
			rain = bern(rng, 0.20)
		}
		wet := 0
		if sprinkler == 1 && rain == 1 {
			wet = bern(rng, 0.99)
		} else if sprinkler == 1 {
			wet = bern(rng, 0.90)
		} else if rain == 1 {
			wet = bern(rng, 0.90)
		} else {
			wet = bern(rng, 0.01)
		}
		rows[i] = []string{itoa(cloudy), itoa(sprinkler), itoa(rain), itoa(wet)}
	}
	writeCSV(filepath.Join(dir, "sprinkler.csv"), headers, rows)
}

func genSurvey(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Age", "Education", "Occupation", "Residence", "Transportation"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		age := categorical(rng, []float64{0.30, 0.50, 0.20})
		edu := 0
		if age == 0 {
			edu = bern(rng, 0.60)
		} else if age == 1 {
			edu = bern(rng, 0.40)
		} else {
			edu = bern(rng, 0.20)
		}
		occ := 0
		if edu == 1 {
			occ = bern(rng, 0.30)
		} else {
			occ = bern(rng, 0.50)
		}
		res := 0
		if edu == 1 {
			res = bern(rng, 0.60)
		} else {
			res = bern(rng, 0.40)
		}
		trans := 0
		if occ == 0 && res == 1 {
			trans = categorical(rng, []float64{0.40, 0.45, 0.15})
		} else if occ == 0 {
			trans = categorical(rng, []float64{0.60, 0.20, 0.20})
		} else if res == 1 {
			trans = categorical(rng, []float64{0.50, 0.30, 0.20})
		} else {
			trans = categorical(rng, []float64{0.70, 0.15, 0.15})
		}
		rows[i] = []string{itoa(age), itoa(edu), itoa(occ), itoa(res), itoa(trans)}
	}
	writeCSV(filepath.Join(dir, "survey.csv"), headers, rows)
}

func genTitanic(rng *rand.Rand, dir string) {
	n := 800
	headers := []string{"Class", "Sex", "Age", "Survived"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		class := categorical(rng, []float64{0.15, 0.15, 0.35, 0.35})
		sex := bern(rng, 0.35)
		age := bern(rng, 0.95)
		p := 0.20
		if class == 0 {
			p += 0.30
		} else if class == 1 {
			p += 0.10
		} else if class == 2 {
			p -= 0.05
		} else {
			p -= 0.02
		}
		if sex == 1 {
			p += 0.35
		}
		if age == 0 {
			p += 0.20
		}
		if p > 0.95 {
			p = 0.95
		}
		if p < 0.05 {
			p = 0.05
		}
		survived := bern(rng, p)
		rows[i] = []string{itoa(class), itoa(sex), itoa(age), itoa(survived)}
	}
	writeCSV(filepath.Join(dir, "titanic.csv"), headers, rows)
}

// ===== Well-known ML datasets =====

// Adult (Census Income) - 14 features + income label, discretized
// Age, Workclass, Education, MaritalStatus, Occupation, Relationship, Race, Sex, HoursPerWeek, NativeCountry, Income
func genAdult(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Age", "Workclass", "Education", "MaritalStatus", "Occupation",
		"Relationship", "Race", "Sex", "HoursPerWeek", "NativeCountry", "Income"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		age := categorical(rng, []float64{0.15, 0.35, 0.30, 0.20}) // young,mid,senior,old
		workclass := categorical(rng, []float64{0.70, 0.10, 0.08, 0.07, 0.05})
		edu := categorical(rng, []float64{0.10, 0.25, 0.35, 0.20, 0.10})
		marital := categorical(rng, []float64{0.45, 0.35, 0.20})
		occ := categorical(rng, []float64{0.15, 0.20, 0.15, 0.15, 0.10, 0.10, 0.15})
		rel := categorical(rng, []float64{0.40, 0.25, 0.15, 0.10, 0.10})
		race := categorical(rng, []float64{0.80, 0.10, 0.05, 0.03, 0.02})
		sex := bern(rng, 0.33) // 33% female
		hours := categorical(rng, []float64{0.10, 0.60, 0.20, 0.10})
		country := categorical(rng, []float64{0.90, 0.10})
		// Income depends on education, age, hours
		p := 0.15
		if edu >= 3 {
			p += 0.25
		}
		if age >= 1 && age <= 2 {
			p += 0.10
		}
		if hours >= 2 {
			p += 0.10
		}
		if marital == 0 {
			p += 0.10
		}
		income := bern(rng, clampF(p, 0.05, 0.85))
		rows[i] = []string{itoa(age), itoa(workclass), itoa(edu), itoa(marital), itoa(occ),
			itoa(rel), itoa(race), itoa(sex), itoa(hours), itoa(country), itoa(income)}
	}
	writeCSV(filepath.Join(dir, "adult.csv"), headers, rows)
}

// Pima Indians Diabetes - 8 features + outcome, discretized to categories
func genPimaDiabetes(rng *rand.Rand, dir string) {
	n := 300
	headers := []string{"Pregnancies", "Glucose", "BloodPressure", "SkinThickness",
		"Insulin", "BMI", "DiabetesPedigree", "Age", "Outcome"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		preg := categorical(rng, []float64{0.30, 0.35, 0.25, 0.10})    // 0-2,3-5,6-9,10+
		glucose := categorical(rng, []float64{0.15, 0.40, 0.30, 0.15}) // low,normal,high,vhigh
		bp := categorical(rng, []float64{0.15, 0.50, 0.25, 0.10})      // low,normal,high,vhigh
		skin := categorical(rng, []float64{0.25, 0.40, 0.25, 0.10})    // thin,normal,thick,vthick
		insulin := categorical(rng, []float64{0.30, 0.35, 0.25, 0.10}) // low,normal,high,vhigh
		bmi := categorical(rng, []float64{0.15, 0.35, 0.30, 0.20})     // under,normal,over,obese
		pedigree := categorical(rng, []float64{0.40, 0.35, 0.25})      // low,medium,high
		age := categorical(rng, []float64{0.35, 0.35, 0.20, 0.10})     // young,mid,senior,old
		// Outcome depends on glucose, bmi, age, pedigree
		p := 0.15
		if glucose >= 2 {
			p += 0.20
		}
		if bmi >= 2 {
			p += 0.15
		}
		if age >= 2 {
			p += 0.10
		}
		if pedigree >= 2 {
			p += 0.10
		}
		outcome := bern(rng, clampF(p, 0.05, 0.85))
		rows[i] = []string{itoa(preg), itoa(glucose), itoa(bp), itoa(skin),
			itoa(insulin), itoa(bmi), itoa(pedigree), itoa(age), itoa(outcome)}
	}
	writeCSV(filepath.Join(dir, "pima_diabetes.csv"), headers, rows)
}

// Iris - 4 continuous features + species, discretized
func genIris(rng *rand.Rand, dir string) {
	n := 150
	headers := []string{"SepalLength", "SepalWidth", "PetalLength", "PetalWidth", "Species"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		// 3 species equally distributed
		species := i / 50
		if species > 2 {
			species = 2
		}
		var sl, sw, pl, pw float64
		switch species {
		case 0: // setosa
			sl = clampF(normalSample(rng, 5.0, 0.35), 4.0, 6.0)
			sw = clampF(normalSample(rng, 3.4, 0.38), 2.0, 4.5)
			pl = clampF(normalSample(rng, 1.5, 0.17), 1.0, 2.0)
			pw = clampF(normalSample(rng, 0.2, 0.10), 0.1, 0.6)
		case 1: // versicolor
			sl = clampF(normalSample(rng, 5.9, 0.52), 4.5, 7.5)
			sw = clampF(normalSample(rng, 2.8, 0.31), 2.0, 3.5)
			pl = clampF(normalSample(rng, 4.3, 0.47), 3.0, 5.5)
			pw = clampF(normalSample(rng, 1.3, 0.20), 0.8, 1.8)
		case 2: // virginica
			sl = clampF(normalSample(rng, 6.6, 0.64), 5.0, 8.0)
			sw = clampF(normalSample(rng, 3.0, 0.32), 2.0, 4.0)
			pl = clampF(normalSample(rng, 5.6, 0.55), 4.5, 7.0)
			pw = clampF(normalSample(rng, 2.0, 0.27), 1.4, 2.5)
		}
		rows[i] = []string{ftoa(sl), ftoa(sw), ftoa(pl), ftoa(pw), itoa(species)}
	}
	// Shuffle
	rng.Shuffle(n, func(i, j int) { rows[i], rows[j] = rows[j], rows[i] })
	writeCSV(filepath.Join(dir, "iris.csv"), headers, rows)
}

// Wine Quality - 11 features + quality, discretized
func genWine(rng *rand.Rand, dir string) {
	n := 300
	headers := []string{"FixedAcidity", "VolatileAcidity", "CitricAcid", "ResidualSugar",
		"Chlorides", "FreeSO2", "TotalSO2", "Density", "pH", "Sulphates", "Alcohol", "Quality"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		fa := categorical(rng, []float64{0.20, 0.40, 0.30, 0.10})
		va := categorical(rng, []float64{0.25, 0.40, 0.25, 0.10})
		ca := categorical(rng, []float64{0.25, 0.35, 0.25, 0.15})
		rs := categorical(rng, []float64{0.30, 0.35, 0.25, 0.10})
		cl := categorical(rng, []float64{0.20, 0.45, 0.25, 0.10})
		fso2 := categorical(rng, []float64{0.20, 0.40, 0.25, 0.15})
		tso2 := categorical(rng, []float64{0.20, 0.35, 0.30, 0.15})
		dens := categorical(rng, []float64{0.20, 0.40, 0.30, 0.10})
		ph := categorical(rng, []float64{0.15, 0.45, 0.30, 0.10})
		sulph := categorical(rng, []float64{0.20, 0.40, 0.25, 0.15})
		alc := categorical(rng, []float64{0.15, 0.35, 0.30, 0.20})
		// Quality 0=low, 1=medium, 2=high
		qp := float64(alc+sulph+ca-va-cl) / 10.0
		quality := categorical(rng, []float64{clampF(0.30-qp, 0.10, 0.60),
			clampF(0.45, 0.20, 0.60), clampF(0.25+qp, 0.10, 0.60)})
		_ = fa
		_ = rs
		_ = fso2
		_ = tso2
		_ = dens
		_ = ph
		rows[i] = []string{itoa(fa), itoa(va), itoa(ca), itoa(rs), itoa(cl), itoa(fso2),
			itoa(tso2), itoa(dens), itoa(ph), itoa(sulph), itoa(alc), itoa(quality)}
	}
	writeCSV(filepath.Join(dir, "wine.csv"), headers, rows)
}

// Heart disease - 13 features + target
func genHeart(rng *rand.Rand, dir string) {
	n := 300
	headers := []string{"Age", "Sex", "ChestPain", "RestBP", "Cholesterol",
		"FastingBS", "RestECG", "MaxHR", "ExerciseAngina", "Oldpeak",
		"Slope", "NumVessels", "Thal", "Target"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		age := categorical(rng, []float64{0.10, 0.30, 0.35, 0.25})
		sex := bern(rng, 0.32)
		cp := categorical(rng, []float64{0.45, 0.15, 0.25, 0.15})
		rbp := categorical(rng, []float64{0.15, 0.50, 0.25, 0.10})
		chol := categorical(rng, []float64{0.15, 0.40, 0.30, 0.15})
		fbs := bern(rng, 0.15)
		recg := categorical(rng, []float64{0.50, 0.35, 0.15})
		maxhr := categorical(rng, []float64{0.15, 0.35, 0.35, 0.15})
		exang := bern(rng, 0.33)
		oldpeak := categorical(rng, []float64{0.45, 0.30, 0.15, 0.10})
		slope := categorical(rng, []float64{0.35, 0.45, 0.20})
		nv := categorical(rng, []float64{0.55, 0.25, 0.15, 0.05})
		thal := categorical(rng, []float64{0.05, 0.55, 0.35, 0.05})
		// Target depends on cp, age, exang, oldpeak
		p := 0.25
		if cp >= 2 {
			p += 0.15
		}
		if age >= 2 {
			p += 0.10
		}
		if exang == 1 {
			p += 0.15
		}
		if oldpeak >= 2 {
			p += 0.10
		}
		if thal >= 2 {
			p += 0.10
		}
		target := bern(rng, clampF(p, 0.05, 0.90))
		rows[i] = []string{itoa(age), itoa(sex), itoa(cp), itoa(rbp), itoa(chol),
			itoa(fbs), itoa(recg), itoa(maxhr), itoa(exang), itoa(oldpeak),
			itoa(slope), itoa(nv), itoa(thal), itoa(target)}
	}
	writeCSV(filepath.Join(dir, "heart.csv"), headers, rows)
}

// Boston Housing - 13 features + MEDV target, discretized
func genBoston(rng *rand.Rand, dir string) {
	n := 300
	headers := []string{"CRIM", "ZN", "INDUS", "CHAS", "NOX", "RM", "AGE",
		"DIS", "RAD", "TAX", "PTRATIO", "B", "LSTAT", "MEDV"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		crim := categorical(rng, []float64{0.40, 0.30, 0.20, 0.10})
		zn := categorical(rng, []float64{0.70, 0.15, 0.10, 0.05})
		indus := categorical(rng, []float64{0.30, 0.35, 0.25, 0.10})
		chas := bern(rng, 0.07)
		nox := categorical(rng, []float64{0.25, 0.35, 0.25, 0.15})
		rm := categorical(rng, []float64{0.10, 0.35, 0.35, 0.20})
		age := categorical(rng, []float64{0.20, 0.25, 0.30, 0.25})
		dis := categorical(rng, []float64{0.25, 0.35, 0.25, 0.15})
		rad := categorical(rng, []float64{0.40, 0.25, 0.20, 0.15})
		tax := categorical(rng, []float64{0.30, 0.30, 0.25, 0.15})
		ptratio := categorical(rng, []float64{0.20, 0.35, 0.30, 0.15})
		b := categorical(rng, []float64{0.10, 0.15, 0.25, 0.50})
		lstat := categorical(rng, []float64{0.25, 0.35, 0.25, 0.15})
		// MEDV depends on RM, LSTAT, CRIM
		mp := float64(rm) - float64(lstat)*0.5 - float64(crim)*0.3
		medv := 0
		if mp >= 1.5 {
			medv = categorical(rng, []float64{0.05, 0.15, 0.40, 0.40})
		} else if mp >= 0.5 {
			medv = categorical(rng, []float64{0.10, 0.35, 0.35, 0.20})
		} else if mp >= -0.5 {
			medv = categorical(rng, []float64{0.25, 0.40, 0.25, 0.10})
		} else {
			medv = categorical(rng, []float64{0.45, 0.30, 0.15, 0.10})
		}
		_ = zn
		_ = indus
		_ = chas
		_ = nox
		_ = age
		_ = dis
		_ = rad
		_ = tax
		_ = ptratio
		_ = b
		rows[i] = []string{itoa(crim), itoa(zn), itoa(indus), itoa(chas), itoa(nox), itoa(rm),
			itoa(age), itoa(dis), itoa(rad), itoa(tax), itoa(ptratio), itoa(b), itoa(lstat), itoa(medv)}
	}
	writeCSV(filepath.Join(dir, "boston.csv"), headers, rows)
}

// Breast Cancer Wisconsin - 9 features + diagnosis
func genBreastCancer(rng *rand.Rand, dir string) {
	n := 300
	headers := []string{"ClumpThickness", "UniformCellSize", "UniformCellShape",
		"MarginalAdhesion", "SingleEpithelialCellSize", "BareNuclei",
		"BlandChromatin", "NormalNucleoli", "Mitoses", "Diagnosis"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		// Benign samples have features skewed to low values
		malignant := bern(rng, 0.35)
		var probs []float64
		if malignant == 1 {
			probs = []float64{0.05, 0.10, 0.15, 0.20, 0.15, 0.15, 0.10, 0.05, 0.03, 0.02}
		} else {
			probs = []float64{0.35, 0.25, 0.15, 0.10, 0.05, 0.04, 0.03, 0.02, 0.005, 0.005}
		}
		ct := categorical(rng, probs) + 1
		ucs := categorical(rng, probs) + 1
		ucsh := categorical(rng, probs) + 1
		ma := categorical(rng, probs) + 1
		sec := categorical(rng, probs) + 1
		bn := categorical(rng, probs) + 1
		bc := categorical(rng, probs) + 1
		nn := categorical(rng, probs) + 1
		mit := categorical(rng, probs) + 1
		rows[i] = []string{itoa(ct), itoa(ucs), itoa(ucsh), itoa(ma), itoa(sec),
			itoa(bn), itoa(bc), itoa(nn), itoa(mit), itoa(malignant)}
	}
	writeCSV(filepath.Join(dir, "breast_cancer.csv"), headers, rows)
}

// ===== BN-specific datasets =====

// Earthquake BN: Burglary -> Alarm; Earthquake -> Alarm; Alarm -> {Radio, Calls}
func genEarthquake(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Burglary", "Earthquake", "Alarm", "Radio", "Calls"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		burg := bern(rng, 0.01)
		eq := bern(rng, 0.02)
		alarm := 0
		if burg == 1 && eq == 1 {
			alarm = bern(rng, 0.95)
		} else if burg == 1 {
			alarm = bern(rng, 0.90)
		} else if eq == 1 {
			alarm = bern(rng, 0.40)
		} else {
			alarm = bern(rng, 0.01)
		}
		radio := 0
		if eq == 1 {
			radio = bern(rng, 0.95)
		} else {
			radio = bern(rng, 0.02)
		}
		calls := 0
		if alarm == 1 {
			calls = bern(rng, 0.85)
		} else {
			calls = bern(rng, 0.03)
		}
		rows[i] = []string{itoa(burg), itoa(eq), itoa(alarm), itoa(radio), itoa(calls)}
	}
	writeCSV(filepath.Join(dir, "earthquake.csv"), headers, rows)
}

// Child BN: ~20 nodes about childhood diseases
func genChild(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"BirthAsphyxia", "Disease", "Age", "LVH", "DuctFlow",
		"CardiacMixing", "LungParench", "LungFlow", "Sick", "HypDistrib",
		"HypoxiaInO2", "CO2", "CO2Report", "XrayReport", "Grunting",
		"LVHReport", "LowerBodyO2", "RUQO2", "GruntingReport", "GruntStatus"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		ba := bern(rng, 0.10)
		disease := categorical(rng, []float64{0.15, 0.20, 0.20, 0.15, 0.15, 0.15})
		age := categorical(rng, []float64{0.30, 0.40, 0.30})
		lvh := bern(rng, 0.15)
		ductFlow := categorical(rng, []float64{0.35, 0.40, 0.25})
		cm := categorical(rng, []float64{0.40, 0.35, 0.15, 0.10})
		lp := categorical(rng, []float64{0.30, 0.40, 0.30})
		lf := categorical(rng, []float64{0.35, 0.40, 0.25})
		sick := 0
		if ba == 1 || disease >= 3 {
			sick = bern(rng, 0.80)
		} else {
			sick = bern(rng, 0.20)
		}
		hd := categorical(rng, []float64{0.35, 0.40, 0.25})
		hio := categorical(rng, []float64{0.40, 0.35, 0.25})
		co2 := categorical(rng, []float64{0.35, 0.40, 0.25})
		co2r := categorical(rng, []float64{0.35, 0.40, 0.25})
		xray := categorical(rng, []float64{0.20, 0.30, 0.30, 0.20})
		grunt := bern(rng, 0.20)
		lvhr := categorical(rng, []float64{0.40, 0.35, 0.25})
		lbo := categorical(rng, []float64{0.35, 0.40, 0.25})
		ruq := categorical(rng, []float64{0.35, 0.40, 0.25})
		gruntr := bern(rng, 0.20)
		grunts := bern(rng, 0.20)
		rows[i] = []string{itoa(ba), itoa(disease), itoa(age), itoa(lvh), itoa(ductFlow),
			itoa(cm), itoa(lp), itoa(lf), itoa(sick), itoa(hd),
			itoa(hio), itoa(co2), itoa(co2r), itoa(xray), itoa(grunt),
			itoa(lvhr), itoa(lbo), itoa(ruq), itoa(gruntr), itoa(grunts)}
	}
	writeCSV(filepath.Join(dir, "child.csv"), headers, rows)
}

// Insurance BN: ~27 nodes
func genInsurance(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Age", "GoodStudent", "SocioEcon", "RiskAversion", "VehicleYear",
		"MakeModel", "DrivingSkill", "SeniorTrain", "ThisCarDam", "RuggedAuto",
		"Accident", "ThisCarCost", "Theft", "AntiTheft", "HomeBase",
		"CarValue", "OtherCarCost", "PropCost", "MedCost", "ILiCost",
		"OtherCost", "Cushioning", "Airbag", "DrivQuality", "Mileage",
		"Antilock", "PolicyType"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		age := categorical(rng, []float64{0.20, 0.30, 0.30, 0.20})
		gs := bern(rng, 0.25)
		se := categorical(rng, []float64{0.20, 0.35, 0.30, 0.15})
		ra := categorical(rng, []float64{0.20, 0.35, 0.30, 0.15})
		vy := categorical(rng, []float64{0.30, 0.40, 0.30})
		mm := categorical(rng, []float64{0.20, 0.30, 0.30, 0.20})
		ds := categorical(rng, []float64{0.20, 0.40, 0.30, 0.10})
		st := bern(rng, 0.15)
		tcd := categorical(rng, []float64{0.50, 0.30, 0.15, 0.05})
		ra2 := categorical(rng, []float64{0.30, 0.40, 0.30})
		acc := categorical(rng, []float64{0.60, 0.25, 0.10, 0.05})
		tcc := categorical(rng, []float64{0.40, 0.30, 0.20, 0.10})
		theft := bern(rng, 0.05)
		at := bern(rng, 0.50)
		hb := categorical(rng, []float64{0.30, 0.30, 0.20, 0.20})
		cv := categorical(rng, []float64{0.15, 0.30, 0.30, 0.15, 0.10})
		occ := categorical(rng, []float64{0.50, 0.30, 0.15, 0.05})
		pc := categorical(rng, []float64{0.45, 0.30, 0.15, 0.10})
		mc := categorical(rng, []float64{0.50, 0.30, 0.15, 0.05})
		ilc := categorical(rng, []float64{0.50, 0.30, 0.15, 0.05})
		oc := categorical(rng, []float64{0.50, 0.30, 0.15, 0.05})
		cush := categorical(rng, []float64{0.15, 0.35, 0.35, 0.15})
		ab := bern(rng, 0.60)
		dq := categorical(rng, []float64{0.15, 0.35, 0.35, 0.15})
		mil := categorical(rng, []float64{0.20, 0.30, 0.30, 0.20})
		al := bern(rng, 0.65)
		pt := categorical(rng, []float64{0.35, 0.35, 0.30})
		_ = gs
		_ = se
		_ = ra
		_ = st
		_ = ra2
		_ = at
		_ = hb
		rows[i] = []string{itoa(age), itoa(gs), itoa(se), itoa(ra), itoa(vy),
			itoa(mm), itoa(ds), itoa(st), itoa(tcd), itoa(ra2),
			itoa(acc), itoa(tcc), itoa(theft), itoa(at), itoa(hb),
			itoa(cv), itoa(occ), itoa(pc), itoa(mc), itoa(ilc),
			itoa(oc), itoa(cush), itoa(ab), itoa(dq), itoa(mil),
			itoa(al), itoa(pt)}
	}
	writeCSV(filepath.Join(dir, "insurance.csv"), headers, rows)
}

// Water BN: water treatment plant with ~32 nodes
func genWater(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"CODO_I", "CODB_I", "SS_I", "SVS_I", "SED_I",
		"COND_I", "PH_I", "CODO_D", "CODB_D", "SS_D",
		"SED_D", "COND_D", "PH_D", "CODO_S", "CODB_S",
		"SS_S", "SED_S", "COND_S", "PH_S", "CODO_E",
		"CODB_E", "SS_E", "SED_E", "COND_E", "PH_E",
		"RD_DBO_P", "RD_SS_P", "RD_SED_P", "RD_DBO_S", "RD_SS_S",
		"RD_SED_S", "RD_DBO_G"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		row := make([]string, 32)
		for j := 0; j < 32; j++ {
			row[j] = itoa(categorical(rng, []float64{0.25, 0.35, 0.25, 0.15}))
		}
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "water.csv"), headers, rows)
}

// Mildew BN: crop disease with ~35 nodes
func genMildew(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Mite_0", "Mite_1", "Mite_2", "Mildew_0", "Mildew_1",
		"Mildew_2", "Micro_0", "Micro_1", "Micro_2", "Photo_0",
		"Photo_1", "Photo_2", "Dew_0", "Dew_1", "Dew_2",
		"Lai_0", "Lai_1", "Lai_2", "Dm_0", "Dm_1",
		"Dm_2", "Nmin_0", "Nmin_1", "Nmin_2", "Rain_0",
		"Rain_1", "Rain_2", "Temp_0", "Temp_1", "Temp_2",
		"Soil_0", "Soil_1", "Soil_2", "Fert_0", "Fert_1"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		row := make([]string, 35)
		for j := 0; j < 35; j++ {
			row[j] = itoa(categorical(rng, []float64{0.30, 0.40, 0.30}))
		}
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "mildew.csv"), headers, rows)
}

// Hailfinder BN: weather prediction ~56 nodes
func genHailfinder(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"N07muVerworworworworworwor", "SubjVertMo", "QGVertMotion", "CombVerMo",
		"AreaMesoALS", "SatContMoist", "RaosContrAM", "AMInstability",
		"Lifted_Index", "K_Index", "TotalTotals", "SWEAT_Index",
		"MorningCIN", "MorningSFC", "AMDewpoint", "AMTempDiff",
		"WindFieldMT", "WindFieldPln", "CldShadeOth", "CldShadeConv",
		"CompPlFcst", "CapChange", "LoLevlRel", "MidLevRel",
		"MvmtFeatures", "LargeScale", "Boundaries", "InitCapCh",
		"ScenRelAMCIN", "ScenRelSFC", "AreaPrecipA", "Outflow",
		"IRWindspr", "ConvInhibit", "CurPropConv", "CldPrecip",
		"ScenPrecip", "Visibility", "InsChange", "R5FcstH",
		"Hail", "PlainsFcst", "CombMoisture", "AreaMoDryAir",
		"VISCloudCov", "MorningBound", "WindAloft", "MountainFcst",
		"Date", "ScenRel34", "LatestCIN", "SfcWndShift",
		"TempDis", "MtVerWorRel", "CapInScen", "Scenario"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		row := make([]string, 56)
		for j := 0; j < 56; j++ {
			row[j] = itoa(categorical(rng, []float64{0.25, 0.30, 0.25, 0.20}))
		}
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "hailfinder.csv"), headers, rows)
}

// Hepar2 BN: liver disease diagnosis ~70 nodes
func genHepar2(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Age", "Sex", "Alcoholism", "Fibrosis", "Diabetes",
		"Obesity", "Hepatotoxic", "TH_Viremial", "TH_Viremial2", "Steatosis",
		"Cirrhosis", "Gallstones", "ChHepatitis", "NonAlcFatty", "Fatigue",
		"Itching", "UpperAbPain", "Jaundice", "Bleeding", "Ascites",
		"Hepatomegaly", "Splenomegaly", "Encephalopathy", "Edema", "PH_Albumin",
		"PH_INR", "PH_Bilirubin", "PH_ALP", "PH_ALT", "PH_AST",
		"PH_GGT", "PH_Cholesterol", "PH_Platelets", "PH_AFP", "Imaging_US",
		"Imaging_CT", "BiopsyResult", "HistologyPat", "HistoryDrugs", "HistoryAlcohol"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		row := make([]string, 40)
		age := categorical(rng, []float64{0.15, 0.30, 0.30, 0.25})
		sex := bern(rng, 0.45)
		alc := bern(rng, 0.20)
		row[0] = itoa(age)
		row[1] = itoa(sex)
		row[2] = itoa(alc)
		for j := 3; j < 40; j++ {
			if j < 12 {
				// Disease nodes - lower probability
				row[j] = itoa(categorical(rng, []float64{0.60, 0.25, 0.15}))
			} else if j < 24 {
				// Symptom nodes
				row[j] = itoa(bern(rng, 0.25))
			} else {
				// Lab/test nodes
				row[j] = itoa(categorical(rng, []float64{0.30, 0.40, 0.30}))
			}
		}
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "hepar2.csv"), headers, rows)
}

// Lucas (Lung Cancer Simple) BN
func genLucas(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Smoking", "YellowFingers", "Anxiety", "PeerPressure",
		"Genetics", "AttentionDisorder", "BornEvenDay", "CarAccident",
		"Fatigue", "Allergy", "Coughing", "LungCancer"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		smoking := bern(rng, 0.35)
		pp := bern(rng, 0.40)
		genetics := bern(rng, 0.10)
		bed := bern(rng, 0.50)
		anxiety := bern(rng, 0.30)
		yf := 0
		if smoking == 1 {
			yf = bern(rng, 0.70)
		} else {
			yf = bern(rng, 0.05)
		}
		ad := 0
		if anxiety == 1 {
			ad = bern(rng, 0.40)
		} else {
			ad = bern(rng, 0.10)
		}
		ca := bern(rng, 0.05)
		fatigue := 0
		if smoking == 1 || anxiety == 1 {
			fatigue = bern(rng, 0.60)
		} else {
			fatigue = bern(rng, 0.20)
		}
		allergy := bern(rng, 0.25)
		coughing := 0
		if smoking == 1 || allergy == 1 {
			coughing = bern(rng, 0.65)
		} else {
			coughing = bern(rng, 0.10)
		}
		lc := 0
		if smoking == 1 && genetics == 1 {
			lc = bern(rng, 0.40)
		} else if smoking == 1 {
			lc = bern(rng, 0.15)
		} else if genetics == 1 {
			lc = bern(rng, 0.10)
		} else {
			lc = bern(rng, 0.02)
		}
		_ = pp
		_ = bed
		rows[i] = []string{itoa(smoking), itoa(yf), itoa(anxiety), itoa(pp),
			itoa(genetics), itoa(ad), itoa(bed), itoa(ca),
			itoa(fatigue), itoa(allergy), itoa(coughing), itoa(lc)}
	}
	writeCSV(filepath.Join(dir, "lucas.csv"), headers, rows)
}

// Andes BN: educational tutoring system ~223 nodes (we use 30 representative)
func genAndes(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Ability", "Motivation", "Knowledge_1", "Knowledge_2", "Knowledge_3",
		"Knowledge_4", "Knowledge_5", "Skill_1", "Skill_2", "Skill_3",
		"Skill_4", "Skill_5", "Task_1", "Task_2", "Task_3",
		"Task_4", "Task_5", "Hint_1", "Hint_2", "Hint_3",
		"Error_1", "Error_2", "Error_3", "Score_1", "Score_2",
		"Score_3", "TimeSpent_1", "TimeSpent_2", "TimeSpent_3", "Overall"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		ability := categorical(rng, []float64{0.20, 0.40, 0.30, 0.10})
		motivation := categorical(rng, []float64{0.15, 0.35, 0.35, 0.15})
		row := make([]string, 30)
		row[0] = itoa(ability)
		row[1] = itoa(motivation)
		for j := 2; j < 29; j++ {
			base := []float64{0.25, 0.35, 0.25, 0.15}
			if ability >= 2 {
				base = []float64{0.10, 0.25, 0.40, 0.25}
			}
			row[j] = itoa(categorical(rng, base))
		}
		// Overall depends on ability+motivation
		op := float64(ability+motivation) / 6.0
		overall := categorical(rng, []float64{clampF(0.30-op, 0.05, 0.50),
			0.35, clampF(0.35+op, 0.15, 0.50)})
		row[29] = itoa(overall)
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "andes.csv"), headers, rows)
}

// Munin BN: neuromuscular disorders ~1041 nodes (we use 30 representative)
func genMunin(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Age", "Sex", "NerveDamage", "MuscleTone", "Reflexes",
		"Sensation_1", "Sensation_2", "Sensation_3", "Motor_1", "Motor_2",
		"Motor_3", "EMG_1", "EMG_2", "EMG_3", "NCV_1",
		"NCV_2", "NCV_3", "Temp_1", "Temp_2", "Temp_3",
		"FWave_1", "FWave_2", "HReflex_1", "HReflex_2", "Diagnosis_1",
		"Diagnosis_2", "Diagnosis_3", "Treatment_1", "Treatment_2", "Prognosis"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		age := categorical(rng, []float64{0.15, 0.30, 0.30, 0.25})
		sex := bern(rng, 0.48)
		nd := categorical(rng, []float64{0.50, 0.30, 0.15, 0.05})
		row := make([]string, 30)
		row[0] = itoa(age)
		row[1] = itoa(sex)
		row[2] = itoa(nd)
		for j := 3; j < 30; j++ {
			if nd >= 2 {
				row[j] = itoa(categorical(rng, []float64{0.10, 0.25, 0.35, 0.30}))
			} else {
				row[j] = itoa(categorical(rng, []float64{0.35, 0.35, 0.20, 0.10}))
			}
		}
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "munin.csv"), headers, rows)
}

// Barley BN: crop yield ~48 nodes (we use 30 representative)
func genBarley(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Soil", "Site", "Variety", "Sowing", "Fertilizer",
		"Fungicide", "Herbicide", "Irrigation", "Rainfall", "Temperature",
		"SunHours", "Wind", "Humidity", "Protein", "StarchContent",
		"KernelWeight", "PlantHeight", "TillerCount", "EarLength", "GrainPerEar",
		"DiseaseLevel", "WeedLevel", "Lodging", "MaturityDate", "HarvestDate",
		"Yield", "MaltQuality", "NitrogenUptake", "WaterUse", "Profit"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		soil := categorical(rng, []float64{0.25, 0.35, 0.25, 0.15})
		site := categorical(rng, []float64{0.30, 0.40, 0.30})
		fert := categorical(rng, []float64{0.20, 0.35, 0.30, 0.15})
		rain := categorical(rng, []float64{0.20, 0.35, 0.30, 0.15})
		row := make([]string, 30)
		row[0] = itoa(soil)
		row[1] = itoa(site)
		row[2] = itoa(categorical(rng, []float64{0.20, 0.25, 0.30, 0.25}))
		row[3] = itoa(categorical(rng, []float64{0.25, 0.40, 0.35}))
		row[4] = itoa(fert)
		for j := 5; j < 25; j++ {
			row[j] = itoa(categorical(rng, []float64{0.20, 0.35, 0.30, 0.15}))
		}
		// Yield depends on soil, fert, rain
		yp := float64(soil+fert+rain) / 9.0
		row[25] = itoa(categorical(rng, []float64{clampF(0.30-yp, 0.05, 0.50),
			0.35, clampF(0.35+yp, 0.15, 0.50)}))
		for j := 26; j < 30; j++ {
			row[j] = itoa(categorical(rng, []float64{0.25, 0.35, 0.25, 0.15}))
		}
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "barley.csv"), headers, rows)
}

// Win95pts BN: Windows 95 printer troubleshooting ~76 nodes (30 representative)
func genWin95pts(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Problem", "ErrorCode", "PrinterType", "DriverVersion", "OSVersion",
		"ConnectionType", "PortConfig", "SpoolerStatus", "MemoryAvail", "DiskSpace",
		"NetworkStatus", "PermissionLevel", "PaperJam", "InkLevel", "PrintQuality",
		"PageOrientation", "FontIssue", "ColorMode", "Resolution", "Duplex",
		"TraySelection", "PaperSize", "Timeout", "RetryCount", "ErrorLog",
		"UserAction", "AutoFix", "ManualFix", "Escalation", "Resolution2"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		problem := categorical(rng, []float64{0.20, 0.25, 0.20, 0.15, 0.10, 0.10})
		row := make([]string, 30)
		row[0] = itoa(problem)
		for j := 1; j < 30; j++ {
			if j < 12 {
				row[j] = itoa(categorical(rng, []float64{0.30, 0.35, 0.25, 0.10}))
			} else {
				row[j] = itoa(bern(rng, 0.30))
			}
		}
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "win95pts.csv"), headers, rows)
}

// ===== UCI datasets =====

// Ecoli - protein localization sites, 7 features + class
func genEcoli(rng *rand.Rand, dir string) {
	n := 336
	headers := []string{"Mcg", "Gvh", "Lip", "Chg", "Aac", "Alm1", "Alm2", "Site"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		site := categorical(rng, []float64{0.42, 0.21, 0.15, 0.06, 0.06, 0.05, 0.03, 0.02})
		var mcgP, gvhP, lipP, chgP, aacP, alm1P, alm2P []float64
		switch site {
		case 0: // cp (cytoplasm)
			mcgP = []float64{0.40, 0.35, 0.20, 0.05}
			gvhP = []float64{0.35, 0.35, 0.20, 0.10}
			lipP = []float64{0.80, 0.20}
			chgP = []float64{0.85, 0.15}
			aacP = []float64{0.35, 0.35, 0.20, 0.10}
			alm1P = []float64{0.30, 0.35, 0.25, 0.10}
			alm2P = []float64{0.30, 0.35, 0.25, 0.10}
		case 1: // im (inner membrane)
			mcgP = []float64{0.10, 0.25, 0.35, 0.30}
			gvhP = []float64{0.15, 0.30, 0.35, 0.20}
			lipP = []float64{0.80, 0.20}
			chgP = []float64{0.85, 0.15}
			aacP = []float64{0.15, 0.30, 0.35, 0.20}
			alm1P = []float64{0.10, 0.20, 0.35, 0.35}
			alm2P = []float64{0.10, 0.20, 0.35, 0.35}
		default:
			mcgP = []float64{0.25, 0.30, 0.25, 0.20}
			gvhP = []float64{0.25, 0.30, 0.25, 0.20}
			lipP = []float64{0.70, 0.30}
			chgP = []float64{0.75, 0.25}
			aacP = []float64{0.25, 0.30, 0.25, 0.20}
			alm1P = []float64{0.25, 0.30, 0.25, 0.20}
			alm2P = []float64{0.25, 0.30, 0.25, 0.20}
		}
		mcg := categorical(rng, mcgP)
		gvh := categorical(rng, gvhP)
		lip := categorical(rng, lipP)
		chg := categorical(rng, chgP)
		aac := categorical(rng, aacP)
		alm1 := categorical(rng, alm1P)
		alm2 := categorical(rng, alm2P)
		rows[i] = []string{itoa(mcg), itoa(gvh), itoa(lip), itoa(chg), itoa(aac), itoa(alm1), itoa(alm2), itoa(site)}
	}
	writeCSV(filepath.Join(dir, "ecoli.csv"), headers, rows)
}

// Glass identification - 9 features + type
func genGlass(rng *rand.Rand, dir string) {
	n := 214
	headers := []string{"RI", "Na", "Mg", "Al", "Si", "K", "Ca", "Ba", "Fe", "Type"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		glassType := categorical(rng, []float64{0.33, 0.35, 0.08, 0.06, 0.04, 0.14})
		row := make([]string, 10)
		for j := 0; j < 9; j++ {
			row[j] = itoa(categorical(rng, []float64{0.20, 0.30, 0.30, 0.20}))
		}
		row[9] = itoa(glassType)
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "glass.csv"), headers, rows)
}

// Zoo - 16 boolean features + type
func genZoo(rng *rand.Rand, dir string) {
	n := 101
	headers := []string{"Hair", "Feathers", "Eggs", "Milk", "Airborne", "Aquatic",
		"Predator", "Toothed", "Backbone", "Breathes", "Venomous", "Fins",
		"Legs", "Tail", "Domestic", "Catsize", "Type"}
	rows := make([][]string, n)
	typeProbs := []float64{0.40, 0.20, 0.05, 0.13, 0.04, 0.08, 0.10}
	for i := 0; i < n; i++ {
		anType := categorical(rng, typeProbs)
		row := make([]string, 17)
		switch anType {
		case 0: // mammal
			row[0] = itoa(bern(rng, 0.90))                                // hair
			row[1] = "0"                                                  // feathers
			row[2] = itoa(bern(rng, 0.10))                                // eggs
			row[3] = "1"                                                  // milk
			row[4] = itoa(bern(rng, 0.10))                                // airborne
			row[5] = itoa(bern(rng, 0.10))                                // aquatic
			row[6] = itoa(bern(rng, 0.55))                                // predator
			row[7] = "1"                                                  // toothed
			row[8] = "1"                                                  // backbone
			row[9] = "1"                                                  // breathes
			row[10] = "0"                                                 // venomous
			row[11] = "0"                                                 // fins
			row[12] = itoa(categorical(rng, []float64{0.05, 0.90, 0.05})) // legs 0/1/2 -> 0,4,2
			row[13] = itoa(bern(rng, 0.80))                               // tail
			row[14] = itoa(bern(rng, 0.20))                               // domestic
			row[15] = itoa(bern(rng, 0.60))                               // catsize
		case 1: // bird
			row[0] = "0"
			row[1] = "1"
			row[2] = "1"
			row[3] = "0"
			row[4] = itoa(bern(rng, 0.70))
			row[5] = itoa(bern(rng, 0.20))
			row[6] = itoa(bern(rng, 0.40))
			row[7] = "0"
			row[8] = "1"
			row[9] = "1"
			row[10] = "0"
			row[11] = "0"
			row[12] = "1"
			row[13] = "1"
			row[14] = itoa(bern(rng, 0.15))
			row[15] = itoa(bern(rng, 0.30))
		default: // other types
			for j := 0; j < 16; j++ {
				row[j] = itoa(bern(rng, 0.35))
			}
		}
		row[16] = itoa(anType)
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "zoo.csv"), headers, rows)
}

// Mushroom - 22 categorical features + edible/poisonous
func genMushroom(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"CapShape", "CapSurface", "CapColor", "Bruises", "Odor",
		"GillAttach", "GillSpacing", "GillSize", "GillColor", "StalkShape",
		"StalkRoot", "StalkSurfAboveRing", "StalkSurfBelowRing", "StalkColorAbove",
		"StalkColorBelow", "VeilType", "VeilColor", "RingNumber", "RingType",
		"SporePrintColor", "Population", "Habitat", "Class"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		edible := bern(rng, 0.52)
		row := make([]string, 23)
		if edible == 1 {
			row[0] = itoa(categorical(rng, []float64{0.20, 0.30, 0.25, 0.15, 0.10}))
			row[4] = itoa(categorical(rng, []float64{0.50, 0.20, 0.15, 0.10, 0.05})) // odor: none mostly
		} else {
			row[0] = itoa(categorical(rng, []float64{0.15, 0.20, 0.25, 0.25, 0.15}))
			row[4] = itoa(categorical(rng, []float64{0.10, 0.20, 0.25, 0.25, 0.20})) // odor: foul more common
		}
		for j := 1; j < 22; j++ {
			if j == 4 {
				continue
			}
			switch j {
			case 3: // bruises binary
				row[j] = itoa(bern(rng, 0.40))
			case 15: // veil type (almost always 0)
				row[j] = "0"
			default:
				nCats := 3 + rng.Intn(4) // 3-6 categories
				probs := make([]float64, nCats)
				sum := 0.0
				for k := range probs {
					probs[k] = rng.Float64() + 0.1
					sum += probs[k]
				}
				for k := range probs {
					probs[k] /= sum
				}
				row[j] = itoa(categorical(rng, probs))
			}
		}
		row[22] = itoa(edible)
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "mushroom.csv"), headers, rows)
}

// Nursery - 8 categorical features + class
func genNursery(rng *rand.Rand, dir string) {
	n := 500
	headers := []string{"Parents", "HasNurs", "Form", "Children", "Housing",
		"Finance", "Social", "Health", "Class"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		parents := categorical(rng, []float64{0.33, 0.34, 0.33})
		hasNurs := categorical(rng, []float64{0.20, 0.20, 0.20, 0.20, 0.20})
		form := categorical(rng, []float64{0.25, 0.25, 0.25, 0.25})
		children := categorical(rng, []float64{0.33, 0.34, 0.33})
		housing := categorical(rng, []float64{0.33, 0.34, 0.33})
		finance := bern(rng, 0.50)
		social := categorical(rng, []float64{0.33, 0.34, 0.33})
		health := categorical(rng, []float64{0.33, 0.34, 0.33})
		// Class depends on multiple factors
		score := parents + hasNurs/2 + housing + finance + social + health
		class := 0
		if score >= 8 {
			class = categorical(rng, []float64{0.05, 0.10, 0.30, 0.40, 0.15})
		} else if score >= 5 {
			class = categorical(rng, []float64{0.10, 0.25, 0.35, 0.25, 0.05})
		} else {
			class = categorical(rng, []float64{0.35, 0.30, 0.20, 0.10, 0.05})
		}
		rows[i] = []string{itoa(parents), itoa(hasNurs), itoa(form), itoa(children),
			itoa(housing), itoa(finance), itoa(social), itoa(health), itoa(class)}
	}
	writeCSV(filepath.Join(dir, "nursery.csv"), headers, rows)
}

// Car Evaluation - 6 categorical features + class
func genCarEvaluation(rng *rand.Rand, dir string) {
	n := 400
	headers := []string{"Buying", "Maint", "Doors", "Persons", "LugBoot", "Safety", "Class"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		buying := categorical(rng, []float64{0.25, 0.25, 0.25, 0.25})
		maint := categorical(rng, []float64{0.25, 0.25, 0.25, 0.25})
		doors := categorical(rng, []float64{0.25, 0.25, 0.25, 0.25})
		persons := categorical(rng, []float64{0.33, 0.34, 0.33})
		lug := categorical(rng, []float64{0.33, 0.34, 0.33})
		safety := categorical(rng, []float64{0.33, 0.34, 0.33})
		// Class: 0=unacc, 1=acc, 2=good, 3=vgood
		score := (3 - buying) + (3 - maint) + persons + lug + safety
		class := 0
		if score >= 12 {
			class = categorical(rng, []float64{0.05, 0.20, 0.35, 0.40})
		} else if score >= 8 {
			class = categorical(rng, []float64{0.15, 0.45, 0.30, 0.10})
		} else if score >= 5 {
			class = categorical(rng, []float64{0.50, 0.35, 0.10, 0.05})
		} else {
			class = categorical(rng, []float64{0.80, 0.15, 0.04, 0.01})
		}
		rows[i] = []string{itoa(buying), itoa(maint), itoa(doors), itoa(persons),
			itoa(lug), itoa(safety), itoa(class)}
	}
	writeCSV(filepath.Join(dir, "car_evaluation.csv"), headers, rows)
}

// Balance Scale - 4 features + class
func genBalanceScale(rng *rand.Rand, dir string) {
	n := 300
	headers := []string{"LeftWeight", "LeftDist", "RightWeight", "RightDist", "Class"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		lw := rng.Intn(5) + 1
		ld := rng.Intn(5) + 1
		rw := rng.Intn(5) + 1
		rd := rng.Intn(5) + 1
		class := 0 // balanced
		lm := lw * ld
		rm := rw * rd
		if lm > rm {
			class = 1 // left
		} else if rm > lm {
			class = 2 // right
		}
		rows[i] = []string{itoa(lw), itoa(ld), itoa(rw), itoa(rd), itoa(class)}
	}
	writeCSV(filepath.Join(dir, "balance_scale.csv"), headers, rows)
}

// MONKS - 6 features + class
func genMonks(rng *rand.Rand, dir string) {
	n := 300
	headers := []string{"A1", "A2", "A3", "A4", "A5", "A6", "Class"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		a1 := rng.Intn(3) + 1
		a2 := rng.Intn(3) + 1
		a3 := rng.Intn(2) + 1
		a4 := rng.Intn(3) + 1
		a5 := rng.Intn(4) + 1
		a6 := rng.Intn(2) + 1
		// MONKS-1 rule: (a1 == a2) or (a5 == 1)
		class := 0
		if a1 == a2 || a5 == 1 {
			class = 1
		}
		rows[i] = []string{itoa(a1), itoa(a2), itoa(a3), itoa(a4), itoa(a5), itoa(a6), itoa(class)}
	}
	writeCSV(filepath.Join(dir, "monks.csv"), headers, rows)
}

// Tic-Tac-Toe endgame - 9 positions + class
func genTicTacToe(rng *rand.Rand, dir string) {
	n := 400
	headers := []string{"TL", "TM", "TR", "ML", "MM", "MR", "BL", "BM", "BR", "Class"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		// 0=x, 1=o, 2=blank
		board := make([]int, 9)
		for j := range board {
			board[j] = categorical(rng, []float64{0.40, 0.40, 0.20})
		}
		// Check if X wins (simplistic)
		xWins := 0
		wins := [][3]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}}
		for _, w := range wins {
			if board[w[0]] == 0 && board[w[1]] == 0 && board[w[2]] == 0 {
				xWins = 1
				break
			}
		}
		row := make([]string, 10)
		for j := 0; j < 9; j++ {
			row[j] = itoa(board[j])
		}
		row[9] = itoa(xWins)
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "tic_tac_toe.csv"), headers, rows)
}

// Congressional Voting Records - 16 votes + party
func genVote(rng *rand.Rand, dir string) {
	n := 435
	headers := []string{"HandicappedInfants", "WaterProjectCost", "AdoptionBudget",
		"PhysicianFeeFreeze", "ElSalvadorAid", "ReligiousGroupsSchools",
		"AntiSatelliteTestBan", "NicaraguanContras", "MXMissile",
		"Immigration", "SynfuelsCutback", "EducationSpending",
		"SuperfundSueRight", "Crime", "DutyFreeExports", "ExportAdminActSA", "Party"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		party := bern(rng, 0.45) // 0=republican, 1=democrat
		row := make([]string, 17)
		for j := 0; j < 16; j++ {
			// Democrats and Republicans vote differently on issues
			p := 0.50
			if party == 1 {
				// Democrats more likely to vote 1 on social issues
				if j == 0 || j == 2 || j == 8 || j == 11 || j == 14 {
					p = 0.70
				} else {
					p = 0.35
				}
			} else {
				if j == 3 || j == 4 || j == 5 || j == 13 {
					p = 0.75
				} else {
					p = 0.40
				}
			}
			// 0=no, 1=yes, 2=abstain
			r := rng.Float64()
			if r < 0.05 {
				row[j] = "2" // abstain
			} else if rng.Float64() < p {
				row[j] = "1"
			} else {
				row[j] = "0"
			}
		}
		row[16] = itoa(party)
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "vote.csv"), headers, rows)
}

// Credit Approval - 15 features + approved
func genCreditApproval(rng *rand.Rand, dir string) {
	n := 300
	headers := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8",
		"A9", "A10", "A11", "A12", "A13", "A14", "A15", "Approved"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		row := make([]string, 16)
		score := 0.0
		for j := 0; j < 15; j++ {
			if j == 0 || j == 3 || j == 4 || j == 5 || j == 6 || j == 12 {
				// Categorical
				v := categorical(rng, []float64{0.30, 0.40, 0.30})
				row[j] = itoa(v)
				score += float64(v) * 0.1
			} else {
				// Discretized continuous
				v := categorical(rng, []float64{0.20, 0.30, 0.30, 0.20})
				row[j] = itoa(v)
				score += float64(v) * 0.08
			}
		}
		// Approval probability based on accumulated score
		p := clampF(0.20+score*0.5, 0.10, 0.80)
		row[15] = itoa(bern(rng, p))
		rows[i] = row
	}
	writeCSV(filepath.Join(dir, "credit_approval.csv"), headers, rows)
}

// Hepatitis - 19 features + class
func genHepatitis(rng *rand.Rand, dir string) {
	n := 155
	headers := []string{"Age", "Sex", "Steroid", "Antivirals", "Fatigue", "Malaise",
		"Anorexia", "LiverBig", "LiverFirm", "SpleenPalpable", "Spiders",
		"Ascites", "Varices", "Bilirubin", "AlkPhosphatase", "Sgot",
		"Albumin", "Protime", "Histology", "Class"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		die := bern(rng, 0.21) // ~21% mortality
		age := categorical(rng, []float64{0.10, 0.30, 0.35, 0.25})
		sex := bern(rng, 0.20)
		steroid := bern(rng, 0.50)
		antiv := bern(rng, 0.25)
		fatigue := bern(rng, 0.65)
		malaise := bern(rng, 0.45)
		anorexia := bern(rng, 0.30)
		liverBig := bern(rng, 0.60)
		liverFirm := bern(rng, 0.50)
		spleen := bern(rng, 0.30)
		spiders := bern(rng, 0.35)
		ascites := 0
		if die == 1 {
			ascites = bern(rng, 0.50)
		} else {
			ascites = bern(rng, 0.10)
		}
		varices := 0
		if die == 1 {
			varices = bern(rng, 0.40)
		} else {
			varices = bern(rng, 0.08)
		}
		bili := categorical(rng, []float64{0.30, 0.35, 0.25, 0.10})
		alkPhos := categorical(rng, []float64{0.25, 0.40, 0.25, 0.10})
		sgot := categorical(rng, []float64{0.25, 0.35, 0.25, 0.15})
		albumin := categorical(rng, []float64{0.15, 0.30, 0.35, 0.20})
		protime := categorical(rng, []float64{0.15, 0.30, 0.35, 0.20})
		histology := bern(rng, 0.55)
		rows[i] = []string{itoa(age), itoa(sex), itoa(steroid), itoa(antiv), itoa(fatigue),
			itoa(malaise), itoa(anorexia), itoa(liverBig), itoa(liverFirm), itoa(spleen),
			itoa(spiders), itoa(ascites), itoa(varices), itoa(bili), itoa(alkPhos),
			itoa(sgot), itoa(albumin), itoa(protime), itoa(histology), itoa(die)}
	}
	writeCSV(filepath.Join(dir, "hepatitis.csv"), headers, rows)
}

// Automobile - 25 features + symboling (risk rating)
func genAutomobile(rng *rand.Rand, dir string) {
	n := 205
	headers := []string{"Symboling", "NormLosses", "Make", "FuelType", "Aspiration",
		"NumDoors", "BodyStyle", "DriveWheels", "EngineLocation", "WheelBase",
		"Length", "Width", "Height", "CurbWeight", "EngineType",
		"NumCylinders", "EngineSize", "FuelSystem", "Bore", "Stroke",
		"CompressionRatio", "Horsepower", "PeakRPM", "CityMPG", "HighwayMPG",
		"Price"}
	rows := make([][]string, n)
	for i := 0; i < n; i++ {
		symboling := categorical(rng, []float64{0.05, 0.15, 0.30, 0.30, 0.15, 0.05})
		normLoss := categorical(rng, []float64{0.20, 0.35, 0.30, 0.15})
		make_ := categorical(rng, []float64{0.10, 0.10, 0.10, 0.10, 0.10, 0.10, 0.10, 0.10, 0.10, 0.10})
		fuelType := bern(rng, 0.20)   // gas/diesel
		aspiration := bern(rng, 0.18) // std/turbo
		numDoors := bern(rng, 0.55)   // two/four
		bodyStyle := categorical(rng, []float64{0.15, 0.35, 0.20, 0.15, 0.15})
		driveWheels := categorical(rng, []float64{0.45, 0.35, 0.20})
		engLoc := bern(rng, 0.03) // front/rear
		wheelBase := categorical(rng, []float64{0.20, 0.40, 0.30, 0.10})
		length := categorical(rng, []float64{0.15, 0.35, 0.35, 0.15})
		width := categorical(rng, []float64{0.15, 0.40, 0.30, 0.15})
		height := categorical(rng, []float64{0.15, 0.35, 0.35, 0.15})
		curbWeight := categorical(rng, []float64{0.15, 0.35, 0.30, 0.20})
		engType := categorical(rng, []float64{0.40, 0.30, 0.15, 0.10, 0.05})
		numCyl := categorical(rng, []float64{0.55, 0.25, 0.10, 0.05, 0.05})
		engSize := categorical(rng, []float64{0.20, 0.35, 0.30, 0.15})
		fuelSys := categorical(rng, []float64{0.30, 0.30, 0.20, 0.20})
		bore := categorical(rng, []float64{0.20, 0.35, 0.30, 0.15})
		stroke := categorical(rng, []float64{0.20, 0.35, 0.30, 0.15})
		compRatio := categorical(rng, []float64{0.20, 0.40, 0.30, 0.10})
		hp := categorical(rng, []float64{0.20, 0.35, 0.30, 0.15})
		peakRPM := categorical(rng, []float64{0.15, 0.40, 0.30, 0.15})
		cityMPG := categorical(rng, []float64{0.15, 0.35, 0.30, 0.20})
		hwyMPG := categorical(rng, []float64{0.15, 0.35, 0.30, 0.20})
		price := categorical(rng, []float64{0.20, 0.30, 0.30, 0.20})
		_ = normLoss
		_ = make_
		_ = fuelType
		_ = aspiration
		_ = engLoc
		rows[i] = []string{itoa(symboling), itoa(normLoss), itoa(make_), itoa(fuelType), itoa(aspiration),
			itoa(numDoors), itoa(bodyStyle), itoa(driveWheels), itoa(engLoc), itoa(wheelBase),
			itoa(length), itoa(width), itoa(height), itoa(curbWeight), itoa(engType),
			itoa(numCyl), itoa(engSize), itoa(fuelSys), itoa(bore), itoa(stroke),
			itoa(compRatio), itoa(hp), itoa(peakRPM), itoa(cityMPG), itoa(hwyMPG),
			itoa(price)}
	}
	writeCSV(filepath.Join(dir, "automobile.csv"), headers, rows)
}

// Suppress unused import warning
var _ = math.Abs
