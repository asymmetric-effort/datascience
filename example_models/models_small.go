package example_models

// This file contains small Bayesian network models (< 10 nodes) with full CPDs.

import (
	"log"

	"github.com/asymmetric-effort/pgmgo/src/models"
)

// Survey returns the Survey Bayesian network with 6 nodes.
// This network models the relationship between age, sex, education,
// occupation, residence, and travel mode.
// CPD values from the bnlearn repository (Scutari, 2010).
//
// Nodes:
//
//	A (Age): {young, adult, old}
//	S (Sex): {M, F}
//	E (Education): {high, uni}
//	O (Occupation): {emp, self}
//	R (Residence): {small, big}
//	T (Travel): {car, train, other}
//
// Edges: A->E, S->E, E->O, E->R, O->T, R->T
func Survey() *models.BayesianNetwork {
	bn := models.NewBayesianNetwork()

	for _, node := range []string{"A", "E", "O", "R", "S", "T"} {
		must(bn.AddNode(node))
	}
	must(bn.AddEdge("A", "E"))
	must(bn.AddEdge("S", "E"))
	must(bn.AddEdge("E", "O"))
	must(bn.AddEdge("E", "R"))
	must(bn.AddEdge("O", "T"))
	must(bn.AddEdge("R", "T"))

	must(bn.SetStates("A", []string{"young", "adult", "old"}))
	must(bn.SetStates("S", []string{"M", "F"}))
	must(bn.SetStates("E", []string{"high", "uni"}))
	must(bn.SetStates("O", []string{"emp", "self"}))
	must(bn.SetStates("R", []string{"small", "big"}))
	must(bn.SetStates("T", []string{"car", "train", "other"}))

	// P(A)
	must(bn.AddCPD(mustCPD("A", 3, [][]float64{
		{0.3}, // young
		{0.5}, // adult
		{0.2}, // old
	}, nil, nil)))

	// P(S)
	must(bn.AddCPD(mustCPD("S", 2, [][]float64{
		{0.6}, // M
		{0.4}, // F
	}, nil, nil)))

	// P(E | A, S)
	// Columns: A=young,S=M | A=young,S=F | A=adult,S=M | A=adult,S=F | A=old,S=M | A=old,S=F
	must(bn.AddCPD(mustCPD("E", 2, [][]float64{
		{0.75, 0.64, 0.72, 0.70, 0.88, 0.90}, // high
		{0.25, 0.36, 0.28, 0.30, 0.12, 0.10}, // uni
	}, []string{"A", "S"}, []int{3, 2})))

	// P(O | E)
	must(bn.AddCPD(mustCPD("O", 2, [][]float64{
		{0.96, 0.92}, // emp
		{0.04, 0.08}, // self
	}, []string{"E"}, []int{2})))

	// P(R | E)
	must(bn.AddCPD(mustCPD("R", 2, [][]float64{
		{0.25, 0.20}, // small
		{0.75, 0.80}, // big
	}, []string{"E"}, []int{2})))

	// P(T | O, R)
	// Columns: O=emp,R=small | O=emp,R=big | O=self,R=small | O=self,R=big
	must(bn.AddCPD(mustCPD("T", 3, [][]float64{
		{0.48, 0.58, 0.56, 0.70}, // car
		{0.42, 0.24, 0.36, 0.21}, // train
		{0.10, 0.18, 0.08, 0.09}, // other
	}, []string{"O", "R"}, []int{2, 2})))

	if err := bn.CheckModel(); err != nil {
		log.Panicf("example_models: Survey network validation failed: %v", err)
	}
	return bn
}

// MontyHall returns the Monty Hall problem as a Bayesian network with 3 nodes.
//
// Nodes:
//
//	Prize: {Door1, Door2, Door3} — which door hides the prize
//	Guest: {Door1, Door2, Door3} — which door the guest picks
//	Host: {Door1, Door2, Door3} — which door the host opens
//
// Edges: Prize->Host, Guest->Host
func MontyHall() *models.BayesianNetwork {
	bn := models.NewBayesianNetwork()

	for _, node := range []string{"Guest", "Host", "Prize"} {
		must(bn.AddNode(node))
	}
	must(bn.AddEdge("Prize", "Host"))
	must(bn.AddEdge("Guest", "Host"))

	must(bn.SetStates("Prize", []string{"Door1", "Door2", "Door3"}))
	must(bn.SetStates("Guest", []string{"Door1", "Door2", "Door3"}))
	must(bn.SetStates("Host", []string{"Door1", "Door2", "Door3"}))

	// P(Prize) — uniform
	must(bn.AddCPD(mustCPD("Prize", 3, [][]float64{
		{1.0 / 3.0},
		{1.0 / 3.0},
		{1.0 / 3.0},
	}, nil, nil)))

	// P(Guest) — uniform
	must(bn.AddCPD(mustCPD("Guest", 3, [][]float64{
		{1.0 / 3.0},
		{1.0 / 3.0},
		{1.0 / 3.0},
	}, nil, nil)))

	// P(Host | Prize, Guest)
	// Host cannot open the door with the prize or the door the guest picked.
	// If Prize==Guest, host picks uniformly from the other two doors.
	// Columns: Prize=D1,Guest=D1 | Prize=D1,Guest=D2 | Prize=D1,Guest=D3 |
	//          Prize=D2,Guest=D1 | Prize=D2,Guest=D2 | Prize=D2,Guest=D3 |
	//          Prize=D3,Guest=D1 | Prize=D3,Guest=D2 | Prize=D3,Guest=D3
	must(bn.AddCPD(mustCPD("Host", 3, [][]float64{
		{0.0, 0.0, 0.0, 0.0, 0.5, 1.0, 0.0, 1.0, 0.0}, // Door1
		{0.5, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.5}, // Door2
		{0.5, 1.0, 0.0, 1.0, 0.5, 0.0, 0.0, 0.0, 0.5}, // Door3
	}, []string{"Prize", "Guest"}, []int{3, 3})))

	if err := bn.CheckModel(); err != nil {
		log.Panicf("example_models: MontyHall network validation failed: %v", err)
	}
	return bn
}

// DogProblem returns a textbook Bayesian network about a dog's behavior.
// From Charniak (1991) "Bayesian Networks without Tears."
//
// Nodes:
//
//	BowelProblem: {true, false}
//	DogOut: {true, false}
//	FamilyOut: {true, false}
//	HearBark: {true, false}
//	LightOn: {true, false}
//
// Edges: BowelProblem->DogOut, FamilyOut->DogOut, FamilyOut->LightOn, DogOut->HearBark
func DogProblem() *models.BayesianNetwork {
	bn := models.NewBayesianNetwork()

	for _, node := range []string{"BowelProblem", "DogOut", "FamilyOut", "HearBark", "LightOn"} {
		must(bn.AddNode(node))
	}
	must(bn.AddEdge("BowelProblem", "DogOut"))
	must(bn.AddEdge("FamilyOut", "DogOut"))
	must(bn.AddEdge("FamilyOut", "LightOn"))
	must(bn.AddEdge("DogOut", "HearBark"))

	must(bn.SetStates("BowelProblem", []string{"true", "false"}))
	must(bn.SetStates("DogOut", []string{"true", "false"}))
	must(bn.SetStates("FamilyOut", []string{"true", "false"}))
	must(bn.SetStates("HearBark", []string{"true", "false"}))
	must(bn.SetStates("LightOn", []string{"true", "false"}))

	// P(BowelProblem)
	must(bn.AddCPD(mustCPD("BowelProblem", 2, [][]float64{
		{0.01},
		{0.99},
	}, nil, nil)))

	// P(FamilyOut)
	must(bn.AddCPD(mustCPD("FamilyOut", 2, [][]float64{
		{0.15},
		{0.85},
	}, nil, nil)))

	// P(DogOut | BowelProblem, FamilyOut)
	// Columns: BP=true,FO=true | BP=true,FO=false | BP=false,FO=true | BP=false,FO=false
	must(bn.AddCPD(mustCPD("DogOut", 2, [][]float64{
		{0.99, 0.90, 0.97, 0.30}, // true
		{0.01, 0.10, 0.03, 0.70}, // false
	}, []string{"BowelProblem", "FamilyOut"}, []int{2, 2})))

	// P(HearBark | DogOut)
	must(bn.AddCPD(mustCPD("HearBark", 2, [][]float64{
		{0.70, 0.01}, // true
		{0.30, 0.99}, // false
	}, []string{"DogOut"}, []int{2})))

	// P(LightOn | FamilyOut)
	must(bn.AddCPD(mustCPD("LightOn", 2, [][]float64{
		{0.60, 0.05}, // true
		{0.40, 0.95}, // false
	}, []string{"FamilyOut"}, []int{2})))

	if err := bn.CheckModel(); err != nil {
		log.Panicf("example_models: DogProblem network validation failed: %v", err)
	}
	return bn
}

// FraudDetection returns a simple fraud detection Bayesian network with 6 nodes.
// This is a textbook-style network for credit card fraud detection.
//
// Nodes:
//
//	Fraud: {No, Yes}
//	Age: {Young, Middle, Old}
//	Sex: {Male, Female}
//	ForeignPurchase: {No, Yes}
//	HighAmount: {No, Yes}
//	Alert: {No, Yes}
//
// Edges: Fraud->ForeignPurchase, Fraud->HighAmount, Age->Fraud, Sex->Fraud,
//
//	ForeignPurchase->Alert, HighAmount->Alert
func FraudDetection() *models.BayesianNetwork {
	bn := models.NewBayesianNetwork()

	for _, node := range []string{"Age", "Alert", "ForeignPurchase", "Fraud", "HighAmount", "Sex"} {
		must(bn.AddNode(node))
	}
	must(bn.AddEdge("Age", "Fraud"))
	must(bn.AddEdge("Sex", "Fraud"))
	must(bn.AddEdge("Fraud", "ForeignPurchase"))
	must(bn.AddEdge("Fraud", "HighAmount"))
	must(bn.AddEdge("ForeignPurchase", "Alert"))
	must(bn.AddEdge("HighAmount", "Alert"))

	must(bn.SetStates("Age", []string{"Young", "Middle", "Old"}))
	must(bn.SetStates("Sex", []string{"Male", "Female"}))
	must(bn.SetStates("Fraud", []string{"No", "Yes"}))
	must(bn.SetStates("ForeignPurchase", []string{"No", "Yes"}))
	must(bn.SetStates("HighAmount", []string{"No", "Yes"}))
	must(bn.SetStates("Alert", []string{"No", "Yes"}))

	// P(Age)
	must(bn.AddCPD(mustCPD("Age", 3, [][]float64{
		{0.25}, // Young
		{0.50}, // Middle
		{0.25}, // Old
	}, nil, nil)))

	// P(Sex)
	must(bn.AddCPD(mustCPD("Sex", 2, [][]float64{
		{0.5}, // Male
		{0.5}, // Female
	}, nil, nil)))

	// P(Fraud | Age, Sex)
	// Columns: Age=Young,Sex=Male | Age=Young,Sex=Female | Age=Middle,Sex=Male | Age=Middle,Sex=Female | Age=Old,Sex=Male | Age=Old,Sex=Female
	must(bn.AddCPD(mustCPD("Fraud", 2, [][]float64{
		{0.97, 0.98, 0.99, 0.995, 0.985, 0.99}, // No
		{0.03, 0.02, 0.01, 0.005, 0.015, 0.01}, // Yes
	}, []string{"Age", "Sex"}, []int{3, 2})))

	// P(ForeignPurchase | Fraud)
	must(bn.AddCPD(mustCPD("ForeignPurchase", 2, [][]float64{
		{0.95, 0.20}, // No
		{0.05, 0.80}, // Yes
	}, []string{"Fraud"}, []int{2})))

	// P(HighAmount | Fraud)
	must(bn.AddCPD(mustCPD("HighAmount", 2, [][]float64{
		{0.90, 0.15}, // No
		{0.10, 0.85}, // Yes
	}, []string{"Fraud"}, []int{2})))

	// P(Alert | ForeignPurchase, HighAmount)
	// Columns: FP=No,HA=No | FP=No,HA=Yes | FP=Yes,HA=No | FP=Yes,HA=Yes
	must(bn.AddCPD(mustCPD("Alert", 2, [][]float64{
		{0.99, 0.40, 0.50, 0.02}, // No
		{0.01, 0.60, 0.50, 0.98}, // Yes
	}, []string{"ForeignPurchase", "HighAmount"}, []int{2, 2})))

	if err := bn.CheckModel(); err != nil {
		log.Panicf("example_models: FraudDetection network validation failed: %v", err)
	}
	return bn
}

// MedicalDiagnosis returns a simple medical diagnosis Bayesian network with 8 nodes.
// A textbook model for reasoning about patient symptoms and diseases.
//
// Nodes:
//
//	Smoking: {No, Yes}
//	Pollution: {Low, High}
//	LungCancer: {No, Yes}
//	Bronchitis: {No, Yes}
//	Fatigue: {No, Yes}
//	ChestPain: {No, Yes}
//	Cough: {No, Yes}
//	ShortnessOfBreath: {No, Yes}
//
// Edges: Smoking->LungCancer, Smoking->Bronchitis, Pollution->LungCancer,
//
//	LungCancer->Fatigue, LungCancer->ChestPain, Bronchitis->Cough,
//	Bronchitis->ShortnessOfBreath, LungCancer->ShortnessOfBreath
func MedicalDiagnosis() *models.BayesianNetwork {
	bn := models.NewBayesianNetwork()

	for _, node := range []string{"Bronchitis", "ChestPain", "Cough", "Fatigue", "LungCancer", "Pollution", "ShortnessOfBreath", "Smoking"} {
		must(bn.AddNode(node))
	}
	must(bn.AddEdge("Smoking", "LungCancer"))
	must(bn.AddEdge("Smoking", "Bronchitis"))
	must(bn.AddEdge("Pollution", "LungCancer"))
	must(bn.AddEdge("LungCancer", "Fatigue"))
	must(bn.AddEdge("LungCancer", "ChestPain"))
	must(bn.AddEdge("Bronchitis", "Cough"))
	must(bn.AddEdge("Bronchitis", "ShortnessOfBreath"))
	must(bn.AddEdge("LungCancer", "ShortnessOfBreath"))

	must(bn.SetStates("Smoking", []string{"No", "Yes"}))
	must(bn.SetStates("Pollution", []string{"Low", "High"}))
	must(bn.SetStates("LungCancer", []string{"No", "Yes"}))
	must(bn.SetStates("Bronchitis", []string{"No", "Yes"}))
	must(bn.SetStates("Fatigue", []string{"No", "Yes"}))
	must(bn.SetStates("ChestPain", []string{"No", "Yes"}))
	must(bn.SetStates("Cough", []string{"No", "Yes"}))
	must(bn.SetStates("ShortnessOfBreath", []string{"No", "Yes"}))

	// P(Smoking)
	must(bn.AddCPD(mustCPD("Smoking", 2, [][]float64{
		{0.70}, // No
		{0.30}, // Yes
	}, nil, nil)))

	// P(Pollution)
	must(bn.AddCPD(mustCPD("Pollution", 2, [][]float64{
		{0.90}, // Low
		{0.10}, // High
	}, nil, nil)))

	// P(LungCancer | Smoking, Pollution)
	// Columns: S=No,P=Low | S=No,P=High | S=Yes,P=Low | S=Yes,P=High
	must(bn.AddCPD(mustCPD("LungCancer", 2, [][]float64{
		{0.999, 0.95, 0.97, 0.92}, // No
		{0.001, 0.05, 0.03, 0.08}, // Yes
	}, []string{"Smoking", "Pollution"}, []int{2, 2})))

	// P(Bronchitis | Smoking)
	must(bn.AddCPD(mustCPD("Bronchitis", 2, [][]float64{
		{0.70, 0.40}, // No
		{0.30, 0.60}, // Yes
	}, []string{"Smoking"}, []int{2})))

	// P(Fatigue | LungCancer)
	must(bn.AddCPD(mustCPD("Fatigue", 2, [][]float64{
		{0.65, 0.10}, // No
		{0.35, 0.90}, // Yes
	}, []string{"LungCancer"}, []int{2})))

	// P(ChestPain | LungCancer)
	must(bn.AddCPD(mustCPD("ChestPain", 2, [][]float64{
		{0.98, 0.30}, // No
		{0.02, 0.70}, // Yes
	}, []string{"LungCancer"}, []int{2})))

	// P(Cough | Bronchitis)
	must(bn.AddCPD(mustCPD("Cough", 2, [][]float64{
		{0.70, 0.10}, // No
		{0.30, 0.90}, // Yes
	}, []string{"Bronchitis"}, []int{2})))

	// P(ShortnessOfBreath | Bronchitis, LungCancer)
	// Columns: B=No,LC=No | B=No,LC=Yes | B=Yes,LC=No | B=Yes,LC=Yes
	must(bn.AddCPD(mustCPD("ShortnessOfBreath", 2, [][]float64{
		{0.90, 0.20, 0.30, 0.05}, // No
		{0.10, 0.80, 0.70, 0.95}, // Yes
	}, []string{"Bronchitis", "LungCancer"}, []int{2, 2})))

	if err := bn.CheckModel(); err != nil {
		log.Panicf("example_models: MedicalDiagnosis network validation failed: %v", err)
	}
	return bn
}

// Earthquake returns the classic 5-node Earthquake/Burglary/Alarm network.
// This is an alias for the same model as Alarm() but with node names matching
// the bnlearn "earthquake" dataset. The structure is identical to the classic
// Pearl (1988) example.
//
// This function is provided for compatibility with the bnlearn repository naming.
func Earthquake() *models.BayesianNetwork {
	return Alarm()
}

// VisitAsia returns the Asia network under an alternative name for compatibility.
// Some libraries refer to this model as "VisitAsia" rather than "Asia."
func VisitAsia() *models.BayesianNetwork {
	return Asia()
}

// CoinToss returns a simple 3-node Bayesian network modeling a biased coin toss.
//
// Nodes:
//
//	Bias: {Fair, Biased} — whether the coin is biased
//	FirstToss: {Heads, Tails}
//	SecondToss: {Heads, Tails}
//
// Edges: Bias->FirstToss, Bias->SecondToss
func CoinToss() *models.BayesianNetwork {
	bn := models.NewBayesianNetwork()

	for _, node := range []string{"Bias", "FirstToss", "SecondToss"} {
		must(bn.AddNode(node))
	}
	must(bn.AddEdge("Bias", "FirstToss"))
	must(bn.AddEdge("Bias", "SecondToss"))

	must(bn.SetStates("Bias", []string{"Fair", "Biased"}))
	must(bn.SetStates("FirstToss", []string{"Heads", "Tails"}))
	must(bn.SetStates("SecondToss", []string{"Heads", "Tails"}))

	// P(Bias)
	must(bn.AddCPD(mustCPD("Bias", 2, [][]float64{
		{0.9}, // Fair
		{0.1}, // Biased
	}, nil, nil)))

	// P(FirstToss | Bias)
	must(bn.AddCPD(mustCPD("FirstToss", 2, [][]float64{
		{0.5, 0.8}, // Heads
		{0.5, 0.2}, // Tails
	}, []string{"Bias"}, []int{2})))

	// P(SecondToss | Bias)
	must(bn.AddCPD(mustCPD("SecondToss", 2, [][]float64{
		{0.5, 0.8}, // Heads
		{0.5, 0.2}, // Tails
	}, []string{"Bias"}, []int{2})))

	if err := bn.CheckModel(); err != nil {
		log.Panicf("example_models: CoinToss network validation failed: %v", err)
	}
	return bn
}
