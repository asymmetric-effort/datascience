package example_models

import (
	"fmt"
	"sort"
	"strings"

	"github.com/asymmetric-effort/pgmgo/src/models"
)

// modelFactory maps lowercase model names to their factory functions.
var modelFactory = map[string]func() *models.BayesianNetwork{
	// Small models with full CPDs
	"student":          Student,
	"asia":             Asia,
	"alarm":            Alarm,
	"cancer":           Cancer,
	"watersprinkler":   WaterSprinkler,
	"survey":           Survey,
	"montyhall":        MontyHall,
	"dogproblem":       DogProblem,
	"frauddetection":   FraudDetection,
	"medicaldiagnosis": MedicalDiagnosis,
	"earthquake":       Earthquake,
	"visitasia":        VisitAsia,
	"cointoss":         CoinToss,

	// Large models — structure only
	"sachs":      Sachs,
	"child":      Child,
	"insurance":  Insurance,
	"alarmfull":  AlarmFull,
	"water":      Water,
	"mildew":     Mildew,
	"barley":     Barley,
	"hailfinder": Hailfinder,
	"hepar2":     Hepar2,
	"win95pts":   Win95pts,
	"pathfinder": Pathfinder,
	"pigs":       Pigs,
}

// List returns the names of all available example models, sorted alphabetically.
func List() []string {
	names := make([]string, 0, len(modelFactory))
	for name := range modelFactory {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Get returns a BayesianNetwork for the given model name (case-insensitive).
// Returns an error if the model name is not recognized.
func Get(name string) (*models.BayesianNetwork, error) {
	factory, ok := modelFactory[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("example_models: unknown model %q; use List() to see available models", name)
	}
	return factory(), nil
}
