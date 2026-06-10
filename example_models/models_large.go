package example_models

// This file contains large Bayesian network models (10+ nodes) with
// structure only (nodes + edges). CPDs are omitted because the full
// parameterizations are too large to hardcode.

import "github.com/asymmetric-effort/datascience/lib/pgm/models"

// structureOnly is a helper that builds a structure-only BayesianNetwork
// from a list of node names and directed edges. State names are set to
// binary {State0, State1} as placeholders since no CPDs are provided.
func structureOnly(nodes []string, edges [][2]string) *models.BayesianNetwork {
	bn := models.NewBayesianNetwork()
	for _, n := range nodes {
		must(bn.AddNode(n))
	}
	for _, e := range edges {
		must(bn.AddEdge(e[0], e[1]))
	}
	return bn
}

// Child returns the Child congenital heart disease network with 20 nodes and 25 edges.
// From Spiegelhalter et al. (1992).
// Structure only — CPDs too large to hardcode.
func Child() *models.BayesianNetwork {
	return structureOnly(
		[]string{
			"BirthAsphyxia", "HypDistrib", "HypoxiaInO2", "CO2", "ChestXray",
			"Grunting", "LVHreport", "LowerBodyO2", "RUQO2", "CO2Report",
			"XrayReport", "Disease", "GruntingReport", "Age", "LVH",
			"DuctFlow", "CardiacMixing", "LungParench", "LungFlow", "Sick",
		},
		[][2]string{
			{"DuctFlow", "HypDistrib"}, {"CardiacMixing", "HypDistrib"},
			{"CardiacMixing", "HypoxiaInO2"}, {"LungParench", "HypoxiaInO2"},
			{"LungParench", "CO2"},
			{"LungParench", "ChestXray"}, {"LungFlow", "ChestXray"},
			{"LungParench", "Grunting"}, {"Sick", "Grunting"},
			{"LVH", "LVHreport"},
			{"HypDistrib", "LowerBodyO2"}, {"HypoxiaInO2", "LowerBodyO2"},
			{"HypoxiaInO2", "RUQO2"},
			{"CO2", "CO2Report"},
			{"ChestXray", "XrayReport"},
			{"BirthAsphyxia", "Disease"},
			{"Grunting", "GruntingReport"},
			{"Disease", "Age"}, {"Sick", "Age"},
			{"Disease", "LVH"},
			{"Disease", "DuctFlow"},
			{"Disease", "CardiacMixing"},
			{"Disease", "LungParench"},
			{"Disease", "LungFlow"},
			{"Disease", "Sick"},
		},
	)
}

// Insurance returns the Insurance evaluation network with 27 nodes and 52 edges.
// From Binder et al. (1997).
// Structure only — CPDs too large to hardcode.
func Insurance() *models.BayesianNetwork {
	return structureOnly(
		[]string{
			"GoodStudent", "Age", "SocioEcon", "RiskAversion", "VehicleYear",
			"ThisCarDam", "RuggedAuto", "Accident", "MakeModel", "DrivQuality",
			"Mileage", "Antilock", "DrivingSkill", "SeniorTrain", "ThisCarCost",
			"Theft", "CarValue", "HomeBase", "AntiTheft", "PropCost",
			"OtherCarCost", "OtherCar", "MedCost", "Cushioning", "Airbag",
			"ILiCost", "DrivHist",
		},
		[][2]string{
			{"SocioEcon", "GoodStudent"}, {"Age", "GoodStudent"},
			{"Age", "SocioEcon"},
			{"Age", "RiskAversion"}, {"SocioEcon", "RiskAversion"},
			{"SocioEcon", "VehicleYear"}, {"RiskAversion", "VehicleYear"},
			{"Accident", "ThisCarDam"}, {"RuggedAuto", "ThisCarDam"},
			{"MakeModel", "RuggedAuto"}, {"VehicleYear", "RuggedAuto"},
			{"Antilock", "Accident"}, {"Mileage", "Accident"}, {"DrivQuality", "Accident"},
			{"SocioEcon", "MakeModel"}, {"RiskAversion", "MakeModel"},
			{"DrivingSkill", "DrivQuality"}, {"RiskAversion", "DrivQuality"},
			{"MakeModel", "Antilock"}, {"VehicleYear", "Antilock"},
			{"Age", "DrivingSkill"}, {"SeniorTrain", "DrivingSkill"},
			{"Age", "SeniorTrain"}, {"RiskAversion", "SeniorTrain"},
			{"ThisCarDam", "ThisCarCost"}, {"CarValue", "ThisCarCost"}, {"Theft", "ThisCarCost"},
			{"AntiTheft", "Theft"}, {"HomeBase", "Theft"}, {"CarValue", "Theft"},
			{"MakeModel", "CarValue"}, {"VehicleYear", "CarValue"}, {"Mileage", "CarValue"},
			{"RiskAversion", "HomeBase"}, {"SocioEcon", "HomeBase"},
			{"RiskAversion", "AntiTheft"}, {"SocioEcon", "AntiTheft"},
			{"OtherCarCost", "PropCost"}, {"ThisCarCost", "PropCost"},
			{"Accident", "OtherCarCost"}, {"RuggedAuto", "OtherCarCost"},
			{"SocioEcon", "OtherCar"},
			{"Accident", "MedCost"}, {"Age", "MedCost"}, {"Cushioning", "MedCost"},
			{"RuggedAuto", "Cushioning"}, {"Airbag", "Cushioning"},
			{"MakeModel", "Airbag"}, {"VehicleYear", "Airbag"},
			{"Accident", "ILiCost"},
			{"DrivingSkill", "DrivHist"}, {"RiskAversion", "DrivHist"},
		},
	)
}

// AlarmFull returns the full ALARM (A Logical Alarm Reduction Mechanism) network
// with 37 nodes and 46 edges. From Beinlich et al. (1989).
// Structure only — CPDs too large to hardcode.
func AlarmFull() *models.BayesianNetwork {
	return structureOnly(
		[]string{
			"HISTORY", "CVP", "PCWP", "HYPOVOLEMIA", "LVEDVOLUME",
			"LVFAILURE", "STROKEVOLUME", "ERRLOWOUTPUT", "HRBP", "HREKG",
			"ERRCAUTER", "HRSAT", "INSUFFANESTH", "ANAPHYLAXIS", "TPR",
			"EXPCO2", "KINKEDTUBE", "MINVOL", "FIO2", "PVSAT",
			"SAO2", "PAP", "PULMEMBOLUS", "SHUNT", "INTUBATION",
			"PRESS", "DISCONNECT", "MINVOLSET", "VENTMACH", "VENTTUBE",
			"VENTLUNG", "VENTALV", "ARTCO2", "CATECHOL", "HR",
			"CO", "BP",
		},
		[][2]string{
			{"LVFAILURE", "HISTORY"},
			{"LVEDVOLUME", "CVP"},
			{"LVEDVOLUME", "PCWP"},
			{"HYPOVOLEMIA", "LVEDVOLUME"}, {"LVFAILURE", "LVEDVOLUME"},
			{"HYPOVOLEMIA", "STROKEVOLUME"}, {"LVFAILURE", "STROKEVOLUME"},
			{"ERRLOWOUTPUT", "HRBP"}, {"HR", "HRBP"},
			{"ERRCAUTER", "HREKG"}, {"HR", "HREKG"},
			{"ERRCAUTER", "HRSAT"}, {"HR", "HRSAT"},
			{"ANAPHYLAXIS", "TPR"},
			{"ARTCO2", "EXPCO2"}, {"VENTLUNG", "EXPCO2"},
			{"INTUBATION", "MINVOL"}, {"VENTLUNG", "MINVOL"},
			{"FIO2", "PVSAT"}, {"VENTALV", "PVSAT"},
			{"PVSAT", "SAO2"}, {"SHUNT", "SAO2"},
			{"PULMEMBOLUS", "PAP"},
			{"INTUBATION", "SHUNT"}, {"PULMEMBOLUS", "SHUNT"},
			{"INTUBATION", "PRESS"}, {"KINKEDTUBE", "PRESS"}, {"VENTTUBE", "PRESS"},
			{"MINVOLSET", "VENTMACH"},
			{"DISCONNECT", "VENTTUBE"}, {"VENTMACH", "VENTTUBE"},
			{"INTUBATION", "VENTLUNG"}, {"KINKEDTUBE", "VENTLUNG"}, {"VENTTUBE", "VENTLUNG"},
			{"INTUBATION", "VENTALV"}, {"VENTLUNG", "VENTALV"},
			{"VENTALV", "ARTCO2"},
			{"ARTCO2", "CATECHOL"}, {"INSUFFANESTH", "CATECHOL"}, {"SAO2", "CATECHOL"}, {"TPR", "CATECHOL"},
			{"CATECHOL", "HR"},
			{"HR", "CO"}, {"STROKEVOLUME", "CO"},
			{"CO", "BP"}, {"TPR", "BP"},
		},
	)
}

// Water returns the Water purification network with 32 nodes and 66 edges.
// From Jensen et al. (1989).
// Structure only — CPDs too large to hardcode.
func Water() *models.BayesianNetwork {
	return structureOnly(
		[]string{
			"C_NI_12_00", "CKNI_12_00", "CBODD_12_00", "CKND_12_00",
			"CNOD_12_00", "CBODN_12_00", "CKNN_12_00", "CNON_12_00",
			"C_NI_12_15", "CKNI_12_15", "CBODD_12_15", "CKND_12_15",
			"CNOD_12_15", "CBODN_12_15", "CKNN_12_15", "CNON_12_15",
			"C_NI_12_30", "CKNI_12_30", "CBODD_12_30", "CKND_12_30",
			"CNOD_12_30", "CBODN_12_30", "CKNN_12_30", "CNON_12_30",
			"C_NI_12_45", "CKNI_12_45", "CBODD_12_45", "CKND_12_45",
			"CNOD_12_45", "CBODN_12_45", "CKNN_12_45", "CNON_12_45",
		},
		[][2]string{
			{"C_NI_12_00", "C_NI_12_15"}, {"CKNI_12_00", "CKNI_12_15"},
			{"C_NI_12_00", "CBODD_12_15"}, {"CKNI_12_00", "CBODD_12_15"}, {"CBODD_12_00", "CBODD_12_15"}, {"CNOD_12_00", "CBODD_12_15"}, {"CBODN_12_00", "CBODD_12_15"},
			{"CKNI_12_00", "CKND_12_15"}, {"CKND_12_00", "CKND_12_15"}, {"CKNN_12_00", "CKND_12_15"},
			{"CBODD_12_00", "CNOD_12_15"}, {"CNOD_12_00", "CNOD_12_15"}, {"CNON_12_00", "CNOD_12_15"},
			{"CBODD_12_00", "CBODN_12_15"}, {"CBODN_12_00", "CBODN_12_15"}, {"CNON_12_00", "CBODN_12_15"},
			{"CKND_12_00", "CKNN_12_15"}, {"CKNN_12_00", "CKNN_12_15"},
			{"CNOD_12_00", "CNON_12_15"}, {"CBODN_12_00", "CNON_12_15"}, {"CKNN_12_00", "CNON_12_15"}, {"CNON_12_00", "CNON_12_15"},
			{"C_NI_12_15", "C_NI_12_30"}, {"CKNI_12_15", "CKNI_12_30"},
			{"C_NI_12_15", "CBODD_12_30"}, {"CKNI_12_15", "CBODD_12_30"}, {"CBODD_12_15", "CBODD_12_30"}, {"CNOD_12_15", "CBODD_12_30"}, {"CBODN_12_15", "CBODD_12_30"},
			{"CKNI_12_15", "CKND_12_30"}, {"CKND_12_15", "CKND_12_30"}, {"CKNN_12_15", "CKND_12_30"},
			{"CBODD_12_15", "CNOD_12_30"}, {"CNOD_12_15", "CNOD_12_30"}, {"CNON_12_15", "CNOD_12_30"},
			{"CBODD_12_15", "CBODN_12_30"}, {"CBODN_12_15", "CBODN_12_30"}, {"CNON_12_15", "CBODN_12_30"},
			{"CKND_12_15", "CKNN_12_30"}, {"CKNN_12_15", "CKNN_12_30"},
			{"CNOD_12_15", "CNON_12_30"}, {"CBODN_12_15", "CNON_12_30"}, {"CKNN_12_15", "CNON_12_30"}, {"CNON_12_15", "CNON_12_30"},
			{"C_NI_12_30", "C_NI_12_45"}, {"CKNI_12_30", "CKNI_12_45"},
			{"C_NI_12_30", "CBODD_12_45"}, {"CKNI_12_30", "CBODD_12_45"}, {"CBODD_12_30", "CBODD_12_45"}, {"CNOD_12_30", "CBODD_12_45"}, {"CBODN_12_30", "CBODD_12_45"},
			{"CKNI_12_30", "CKND_12_45"}, {"CKND_12_30", "CKND_12_45"}, {"CKNN_12_30", "CKND_12_45"},
			{"CBODD_12_30", "CNOD_12_45"}, {"CNOD_12_30", "CNOD_12_45"}, {"CNON_12_30", "CNOD_12_45"},
			{"CBODD_12_30", "CBODN_12_45"}, {"CBODN_12_30", "CBODN_12_45"}, {"CNON_12_30", "CBODN_12_45"},
			{"CKND_12_30", "CKNN_12_45"}, {"CKNN_12_30", "CKNN_12_45"},
			{"CNOD_12_30", "CNON_12_45"}, {"CBODN_12_30", "CNON_12_45"}, {"CKNN_12_30", "CNON_12_45"}, {"CNON_12_30", "CNON_12_45"},
		},
	)
}

// Mildew returns the Mildew crop disease network with 35 nodes and 46 edges.
// From Kjærulff (1992).
// Structure only — CPDs too large to hardcode.
func Mildew() *models.BayesianNetwork {
	return structureOnly(
		[]string{
			"dm_1", "foto_1", "straaling_1", "temp_1", "lai_1", "meldug_1", "lai_0",
			"dm_2", "foto_2", "straaling_2", "temp_2", "lai_2", "meldug_2",
			"dm_3", "foto_3", "straaling_3", "temp_3", "lai_3", "meldug_3",
			"dm_4", "foto_4", "straaling_4", "temp_4", "lai_4", "meldug_4",
			"mikro_1", "middel_1", "mikro_2", "middel_2", "mikro_3", "middel_3",
			"nedboer_1", "nedboer_2", "nedboer_3", "udbytte",
		},
		[][2]string{
			{"foto_1", "dm_1"}, {"lai_1", "foto_1"}, {"temp_1", "foto_1"}, {"straaling_1", "foto_1"},
			{"lai_0", "lai_1"}, {"meldug_1", "lai_1"},
			{"foto_2", "dm_2"}, {"dm_1", "dm_2"},
			{"lai_2", "foto_2"}, {"temp_2", "foto_2"}, {"straaling_2", "foto_2"},
			{"lai_1", "lai_2"}, {"meldug_2", "lai_2"},
			{"middel_1", "meldug_2"}, {"mikro_1", "meldug_2"}, {"meldug_1", "meldug_2"},
			{"foto_3", "dm_3"}, {"dm_2", "dm_3"},
			{"lai_3", "foto_3"}, {"temp_3", "foto_3"}, {"straaling_3", "foto_3"},
			{"lai_2", "lai_3"}, {"meldug_3", "lai_3"},
			{"middel_2", "meldug_3"}, {"mikro_2", "meldug_3"}, {"meldug_2", "meldug_3"},
			{"foto_4", "dm_4"}, {"dm_3", "dm_4"},
			{"lai_4", "foto_4"}, {"temp_4", "foto_4"}, {"straaling_4", "foto_4"},
			{"lai_3", "lai_4"}, {"meldug_4", "lai_4"},
			{"middel_3", "meldug_4"}, {"mikro_3", "meldug_4"}, {"meldug_3", "meldug_4"},
			{"lai_1", "mikro_1"}, {"temp_1", "mikro_1"}, {"nedboer_1", "mikro_1"},
			{"lai_2", "mikro_2"}, {"temp_2", "mikro_2"}, {"nedboer_2", "mikro_2"},
			{"lai_3", "mikro_3"}, {"temp_3", "mikro_3"}, {"nedboer_3", "mikro_3"},
			{"dm_4", "udbytte"},
		},
	)
}

// Barley returns the Barley crop yield network with 48 nodes and 84 edges.
// From Kristensen & Rasmussen (2002).
// Structure only — CPDs too large to hardcode.
func Barley() *models.BayesianNetwork {
	return structureOnly(
		[]string{
			"jordtype", "komm", "nedbarea", "nmin", "aar_mod", "forfrugt",
			"potnmin", "jordn", "pesticid", "exptgens", "mod_nmin", "ngodnt",
			"nopt", "ngodnn", "ngodn", "nprot", "saatid", "rokap",
			"dgv1059", "sort", "srtprot", "nplac", "dg25", "ngtilg",
			"ntilg", "saamng", "tkvs", "saakern", "partigerm", "frspdag",
			"jordinf", "markgrm", "antplnt", "sorttkv", "aks_m2", "keraks",
			"dgv5980", "aks_vgt", "srtsize", "ksort", "protein", "udb",
			"spndx", "tkv", "slt22", "s2225", "s2528", "bgbyg",
		},
		[][2]string{
			{"komm", "nedbarea"},
			{"jordtype", "nmin"}, {"nedbarea", "nmin"},
			{"komm", "aar_mod"}, {"jordtype", "aar_mod"},
			{"jordtype", "potnmin"}, {"forfrugt", "potnmin"},
			{"nmin", "jordn"}, {"aar_mod", "jordn"}, {"potnmin", "jordn"},
			{"jordtype", "exptgens"}, {"forfrugt", "exptgens"}, {"pesticid", "exptgens"},
			{"nmin", "mod_nmin"}, {"aar_mod", "mod_nmin"},
			{"forfrugt", "ngodnt"}, {"exptgens", "ngodnt"}, {"mod_nmin", "ngodnt"},
			{"exptgens", "nopt"}, {"pesticid", "nopt"},
			{"nopt", "ngodnn"}, {"jordn", "ngodnn"},
			{"ngodnt", "ngodn"}, {"ngodnn", "ngodn"},
			{"jordn", "nprot"}, {"ngodn", "nprot"},
			{"jordtype", "rokap"},
			{"saatid", "dgv1059"}, {"rokap", "dgv1059"},
			{"sort", "srtprot"},
			{"saatid", "dg25"},
			{"ngodn", "ngtilg"}, {"nplac", "ngtilg"}, {"dg25", "ngtilg"},
			{"ngtilg", "ntilg"}, {"jordn", "ntilg"},
			{"saamng", "saakern"}, {"tkvs", "saakern"},
			{"saatid", "frspdag"},
			{"frspdag", "jordinf"},
			{"partigerm", "markgrm"}, {"jordinf", "markgrm"},
			{"saakern", "antplnt"}, {"markgrm", "antplnt"},
			{"sort", "sorttkv"},
			{"antplnt", "aks_m2"}, {"ntilg", "aks_m2"}, {"dgv1059", "aks_m2"}, {"sorttkv", "aks_m2"},
			{"ntilg", "keraks"}, {"dgv1059", "keraks"}, {"aks_m2", "keraks"},
			{"rokap", "dgv5980"},
			{"ntilg", "aks_vgt"}, {"dgv5980", "aks_vgt"}, {"aks_m2", "aks_vgt"},
			{"sort", "srtsize"},
			{"keraks", "ksort"}, {"aks_vgt", "ksort"}, {"srtsize", "ksort"},
			{"nprot", "protein"}, {"dgv1059", "protein"}, {"srtprot", "protein"}, {"ksort", "protein"},
			{"aks_m2", "udb"}, {"aks_vgt", "udb"},
			{"ntilg", "spndx"}, {"dgv5980", "spndx"}, {"ksort", "spndx"},
			{"aks_m2", "tkv"}, {"keraks", "tkv"}, {"ntilg", "tkv"}, {"sorttkv", "tkv"},
			{"keraks", "slt22"}, {"aks_vgt", "slt22"}, {"srtsize", "slt22"},
			{"keraks", "s2225"}, {"aks_vgt", "s2225"}, {"srtsize", "s2225"},
			{"keraks", "s2528"}, {"aks_vgt", "s2528"}, {"srtsize", "s2528"},
			{"dgv1059", "bgbyg"}, {"dgv5980", "bgbyg"},
		},
	)
}

// Hailfinder returns the Hailfinder weather network with 56 nodes and 66 edges.
// From Abramson et al. (1996).
// Structure only — CPDs too large to hardcode.
func Hailfinder() *models.BayesianNetwork {
	return structureOnly(
		[]string{
			"N0_7muVerMo", "SubjVertMo", "QGVertMotion", "CombVerMo", "AreaMeso_ALS",
			"SatContMoist", "RaoContMoist", "CombMoisture", "AreaMoDryAir", "VISCloudCov",
			"IRCloudCover", "CombClouds", "CldShadeOth", "AMInstabMt", "InsInMt",
			"WndHodograph", "OutflowFrMt", "MorningBound", "Boundaries", "CldShadeConv",
			"CompPlFcst", "CapChange", "LoLevMoistAd", "InsChange", "MountainFcst",
			"Date", "Scenario", "ScenRelAMCIN", "MorningCIN", "AMCINInScen",
			"CapInScen", "ScenRelAMIns", "LIfr12ZDENSd", "AMDewptCalPl", "AMInsWliScen",
			"InsSclInScen", "ScenRel3_4", "LatestCIN", "LLIW", "CurPropConv",
			"ScnRelPlFcst", "PlainsFcst", "N34StarFcst", "R5Fcst", "Dewpoints",
			"LowLLapse", "MeanRH", "MidLLapse", "MvmtFeatures", "RHRatio",
			"SfcWndShfDis", "SynForcng", "TempDis", "WindAloft", "WindFieldMt",
			"WindFieldPln",
		},
		[][2]string{
			{"N0_7muVerMo", "CombVerMo"}, {"SubjVertMo", "CombVerMo"}, {"QGVertMotion", "CombVerMo"},
			{"CombVerMo", "AreaMeso_ALS"},
			{"SatContMoist", "CombMoisture"}, {"RaoContMoist", "CombMoisture"},
			{"AreaMeso_ALS", "AreaMoDryAir"}, {"CombMoisture", "AreaMoDryAir"},
			{"VISCloudCov", "CombClouds"}, {"IRCloudCover", "CombClouds"},
			{"AreaMoDryAir", "CldShadeOth"}, {"AreaMeso_ALS", "CldShadeOth"}, {"CombClouds", "CldShadeOth"},
			{"CldShadeOth", "InsInMt"}, {"AMInstabMt", "InsInMt"},
			{"InsInMt", "OutflowFrMt"}, {"WndHodograph", "OutflowFrMt"},
			{"OutflowFrMt", "Boundaries"}, {"WndHodograph", "Boundaries"}, {"MorningBound", "Boundaries"},
			{"InsInMt", "CldShadeConv"}, {"WndHodograph", "CldShadeConv"},
			{"Boundaries", "CompPlFcst"}, {"CldShadeConv", "CompPlFcst"}, {"AreaMeso_ALS", "CompPlFcst"}, {"CldShadeOth", "CompPlFcst"},
			{"CompPlFcst", "CapChange"},
			{"LoLevMoistAd", "InsChange"}, {"CompPlFcst", "InsChange"},
			{"InsInMt", "MountainFcst"},
			{"Date", "Scenario"},
			{"Scenario", "ScenRelAMCIN"},
			{"ScenRelAMCIN", "AMCINInScen"}, {"MorningCIN", "AMCINInScen"},
			{"AMCINInScen", "CapInScen"}, {"CapChange", "CapInScen"},
			{"Scenario", "ScenRelAMIns"},
			{"ScenRelAMIns", "AMInsWliScen"}, {"LIfr12ZDENSd", "AMInsWliScen"}, {"AMDewptCalPl", "AMInsWliScen"},
			{"AMInsWliScen", "InsSclInScen"}, {"InsChange", "InsSclInScen"},
			{"Scenario", "ScenRel3_4"},
			{"LatestCIN", "CurPropConv"}, {"LLIW", "CurPropConv"},
			{"Scenario", "ScnRelPlFcst"},
			{"CurPropConv", "PlainsFcst"}, {"InsSclInScen", "PlainsFcst"}, {"CapInScen", "PlainsFcst"}, {"ScnRelPlFcst", "PlainsFcst"},
			{"ScenRel3_4", "N34StarFcst"}, {"PlainsFcst", "N34StarFcst"},
			{"MountainFcst", "R5Fcst"}, {"N34StarFcst", "R5Fcst"},
			{"Scenario", "Dewpoints"},
			{"Scenario", "LowLLapse"},
			{"Scenario", "MeanRH"},
			{"Scenario", "MidLLapse"},
			{"Scenario", "MvmtFeatures"},
			{"Scenario", "RHRatio"},
			{"Scenario", "SfcWndShfDis"},
			{"Scenario", "SynForcng"},
			{"Scenario", "TempDis"},
			{"Scenario", "WindAloft"},
			{"Scenario", "WindFieldMt"},
			{"Scenario", "WindFieldPln"},
		},
	)
}

// Hepar2 returns the HEPAR II liver disease network with 70 nodes and 123 edges.
// From Onisko et al. (2001).
// Structure only — CPDs too large to hardcode.
func Hepar2() *models.BayesianNetwork {
	return structureOnly(
		[]string{
			"alcoholism", "vh_amn", "hepatotoxic", "THepatitis", "hospital",
			"surgery", "gallstones", "choledocholithotomy", "injections", "transfusion",
			"ChHepatitis", "sex", "age", "PBC", "fibrosis",
			"diabetes", "obesity", "Steatosis", "Cirrhosis", "Hyperbilirubinemia",
			"triglycerides", "RHepatitis", "fatigue", "bilirubin", "itching",
			"upper_pain", "fat", "pain_ruq", "pressure_ruq", "phosphatase",
			"skin", "ama", "le_cells", "joints", "pain",
			"proteins", "edema", "platelet", "inr", "bleeding",
			"flatulence", "alcohol", "encephalopathy", "urea", "ascites",
			"hepatomegaly", "hepatalgia", "density", "ESR", "alt",
			"ast", "amylase", "ggtp", "cholesterol", "hbsag",
			"hbsag_anti", "anorexia", "nausea", "spleen", "consciousness",
			"spiders", "jaundice", "albumin", "edge", "irregular_liver",
			"hbc_anti", "hcv_anti", "palms", "hbeag", "carcinoma",
		},
		[][2]string{
			{"hepatotoxic", "THepatitis"}, {"alcoholism", "THepatitis"},
			{"gallstones", "choledocholithotomy"},
			{"hospital", "injections"}, {"surgery", "injections"}, {"choledocholithotomy", "injections"},
			{"hospital", "transfusion"}, {"surgery", "transfusion"}, {"choledocholithotomy", "transfusion"},
			{"transfusion", "ChHepatitis"}, {"vh_amn", "ChHepatitis"}, {"injections", "ChHepatitis"},
			{"sex", "PBC"}, {"age", "PBC"},
			{"ChHepatitis", "fibrosis"},
			{"diabetes", "obesity"},
			{"obesity", "Steatosis"}, {"alcoholism", "Steatosis"},
			{"fibrosis", "Cirrhosis"}, {"Steatosis", "Cirrhosis"},
			{"age", "Hyperbilirubinemia"}, {"sex", "Hyperbilirubinemia"},
			{"Steatosis", "triglycerides"},
			{"hepatotoxic", "RHepatitis"},
			{"ChHepatitis", "fatigue"}, {"THepatitis", "fatigue"}, {"RHepatitis", "fatigue"},
			{"Hyperbilirubinemia", "bilirubin"}, {"PBC", "bilirubin"}, {"Cirrhosis", "bilirubin"}, {"gallstones", "bilirubin"}, {"ChHepatitis", "bilirubin"},
			{"bilirubin", "itching"},
			{"gallstones", "upper_pain"},
			{"gallstones", "fat"},
			{"Steatosis", "pain_ruq"}, {"Hyperbilirubinemia", "pain_ruq"},
			{"gallstones", "pressure_ruq"}, {"PBC", "pressure_ruq"}, {"ChHepatitis", "pressure_ruq"},
			{"RHepatitis", "phosphatase"}, {"THepatitis", "phosphatase"}, {"Cirrhosis", "phosphatase"}, {"ChHepatitis", "phosphatase"},
			{"bilirubin", "skin"},
			{"PBC", "ama"},
			{"PBC", "le_cells"},
			{"PBC", "joints"},
			{"PBC", "pain"}, {"joints", "pain"},
			{"Cirrhosis", "proteins"},
			{"Cirrhosis", "edema"},
			{"Cirrhosis", "platelet"}, {"PBC", "platelet"},
			{"ChHepatitis", "inr"}, {"Cirrhosis", "inr"}, {"THepatitis", "inr"}, {"Hyperbilirubinemia", "inr"},
			{"platelet", "bleeding"}, {"inr", "bleeding"},
			{"gallstones", "flatulence"},
			{"Cirrhosis", "alcohol"},
			{"Cirrhosis", "encephalopathy"}, {"PBC", "encephalopathy"},
			{"encephalopathy", "urea"},
			{"proteins", "ascites"},
			{"RHepatitis", "hepatomegaly"}, {"THepatitis", "hepatomegaly"}, {"Steatosis", "hepatomegaly"}, {"Hyperbilirubinemia", "hepatomegaly"},
			{"hepatomegaly", "hepatalgia"},
			{"encephalopathy", "density"},
			{"PBC", "ESR"}, {"ChHepatitis", "ESR"}, {"Steatosis", "ESR"}, {"Hyperbilirubinemia", "ESR"},
			{"ChHepatitis", "alt"}, {"RHepatitis", "alt"}, {"THepatitis", "alt"}, {"Steatosis", "alt"}, {"Cirrhosis", "alt"},
			{"ChHepatitis", "ast"}, {"RHepatitis", "ast"}, {"THepatitis", "ast"}, {"Steatosis", "ast"}, {"Cirrhosis", "ast"},
			{"gallstones", "amylase"},
			{"PBC", "ggtp"}, {"THepatitis", "ggtp"}, {"RHepatitis", "ggtp"}, {"Steatosis", "ggtp"}, {"ChHepatitis", "ggtp"}, {"Hyperbilirubinemia", "ggtp"},
			{"PBC", "cholesterol"}, {"Steatosis", "cholesterol"}, {"ChHepatitis", "cholesterol"},
			{"vh_amn", "hbsag"}, {"ChHepatitis", "hbsag"},
			{"vh_amn", "hbsag_anti"}, {"ChHepatitis", "hbsag_anti"}, {"hbsag", "hbsag_anti"},
			{"RHepatitis", "anorexia"}, {"THepatitis", "anorexia"},
			{"RHepatitis", "nausea"}, {"THepatitis", "nausea"},
			{"Cirrhosis", "spleen"}, {"RHepatitis", "spleen"}, {"THepatitis", "spleen"},
			{"encephalopathy", "consciousness"},
			{"Cirrhosis", "spiders"},
			{"bilirubin", "jaundice"},
			{"Cirrhosis", "albumin"},
			{"Cirrhosis", "edge"},
			{"Cirrhosis", "irregular_liver"},
			{"vh_amn", "hbc_anti"}, {"ChHepatitis", "hbc_anti"},
			{"vh_amn", "hcv_anti"}, {"ChHepatitis", "hcv_anti"},
			{"Cirrhosis", "palms"},
			{"vh_amn", "hbeag"}, {"ChHepatitis", "hbeag"},
			{"Cirrhosis", "carcinoma"}, {"PBC", "carcinoma"},
		},
	)
}

// Win95pts returns the Win95pts printer troubleshooting network with 76 nodes
// and 112 edges. From Heckerman et al. (1995).
// Structure only — CPDs too large to hardcode.
func Win95pts() *models.BayesianNetwork {
	return structureOnly(
		[]string{
			"AppOK", "DataFile", "AppData", "DskLocal", "PrtSpool",
			"PrtOn", "PrtPaper", "NetPrint", "PrtDriver", "PrtThread",
			"EMFOK", "GDIIN", "DrvSet", "DrvOK", "GDIOUT",
			"PrtSel", "PrtDataOut", "PrtPath", "NtwrkCnfg", "PTROFFLINE",
			"NetOK", "PrtCbl", "PrtPort", "CblPrtHrdwrOK", "LclOK",
			"DSApplctn", "PrtMpTPth", "DS_NTOK", "DS_LCLOK", "PC2PRT",
			"PrtMem", "PrtTimeOut", "FllCrrptdBffr", "TnrSpply", "PrtData",
			"Problem1", "AppDtGnTm", "PrntPrcssTm", "DeskPrntSpd", "PgOrnttnOK",
			"PrntngArOK", "ScrnFntNtPrntrFnt", "CmpltPgPrntd", "GrphcsRltdDrvrSttngs", "EPSGrphc",
			"NnPSGrphc", "PrtPScript", "PSGRAPHIC", "Problem4", "TrTypFnts",
			"FntInstlltn", "PrntrAccptsTrtyp", "TTOK", "NnTTOK", "Problem5",
			"LclGrbld", "NtGrbld", "GrbldOtpt", "HrglssDrtnAftrPrnt", "REPEAT",
			"AvlblVrtlMmry", "PSERRMEM", "TstpsTxt", "GrbldPS", "IncmpltPS",
			"PrtFile", "PrtIcon", "Problem6", "Problem3", "PrtQueue",
			"NtSpd", "Problem2", "PrtStatPaper", "PrtStatToner", "PrtStatMem",
			"PrtStatOff",
		},
		[][2]string{
			{"AppOK", "AppData"}, {"DataFile", "AppData"},
			{"AppData", "EMFOK"}, {"DskLocal", "EMFOK"}, {"PrtThread", "EMFOK"},
			{"AppData", "GDIIN"}, {"PrtSpool", "GDIIN"}, {"EMFOK", "GDIIN"},
			{"PrtDriver", "GDIOUT"}, {"GDIIN", "GDIOUT"}, {"DrvSet", "GDIOUT"}, {"DrvOK", "GDIOUT"},
			{"GDIOUT", "PrtDataOut"}, {"PrtSel", "PrtDataOut"},
			{"PrtPath", "NetOK"}, {"NtwrkCnfg", "NetOK"}, {"PTROFFLINE", "NetOK"},
			{"PrtCbl", "LclOK"}, {"PrtPort", "LclOK"}, {"CblPrtHrdwrOK", "LclOK"},
			{"AppData", "DS_NTOK"}, {"PrtPath", "DS_NTOK"}, {"PrtMpTPth", "DS_NTOK"}, {"NtwrkCnfg", "DS_NTOK"}, {"PTROFFLINE", "DS_NTOK"},
			{"AppData", "DS_LCLOK"}, {"PrtCbl", "DS_LCLOK"}, {"PrtPort", "DS_LCLOK"}, {"CblPrtHrdwrOK", "DS_LCLOK"},
			{"NetPrint", "PC2PRT"}, {"PrtDataOut", "PC2PRT"}, {"NetOK", "PC2PRT"}, {"LclOK", "PC2PRT"}, {"DSApplctn", "PC2PRT"}, {"DS_NTOK", "PC2PRT"}, {"DS_LCLOK", "PC2PRT"},
			{"PrtOn", "PrtData"}, {"PrtPaper", "PrtData"}, {"PC2PRT", "PrtData"}, {"PrtMem", "PrtData"}, {"PrtTimeOut", "PrtData"}, {"FllCrrptdBffr", "PrtData"}, {"TnrSpply", "PrtData"},
			{"PrtData", "Problem1"},
			{"PrtSpool", "AppDtGnTm"},
			{"PrtSpool", "PrntPrcssTm"},
			{"PrtMem", "DeskPrntSpd"}, {"AppDtGnTm", "DeskPrntSpd"}, {"PrntPrcssTm", "DeskPrntSpd"},
			{"PrtMem", "CmpltPgPrntd"}, {"PgOrnttnOK", "CmpltPgPrntd"}, {"PrntngArOK", "CmpltPgPrntd"},
			{"PrtMem", "NnPSGrphc"}, {"GrphcsRltdDrvrSttngs", "NnPSGrphc"}, {"EPSGrphc", "NnPSGrphc"},
			{"PrtMem", "PSGRAPHIC"}, {"GrphcsRltdDrvrSttngs", "PSGRAPHIC"}, {"EPSGrphc", "PSGRAPHIC"},
			{"NnPSGrphc", "Problem4"}, {"PrtPScript", "Problem4"}, {"PSGRAPHIC", "Problem4"},
			{"PrtMem", "TTOK"}, {"FntInstlltn", "TTOK"}, {"PrntrAccptsTrtyp", "TTOK"},
			{"PrtMem", "NnTTOK"}, {"ScrnFntNtPrntrFnt", "NnTTOK"}, {"FntInstlltn", "NnTTOK"},
			{"TrTypFnts", "Problem5"}, {"TTOK", "Problem5"}, {"NnTTOK", "Problem5"},
			{"AppData", "LclGrbld"}, {"PrtDriver", "LclGrbld"}, {"PrtMem", "LclGrbld"}, {"CblPrtHrdwrOK", "LclGrbld"},
			{"AppData", "NtGrbld"}, {"PrtDriver", "NtGrbld"}, {"PrtMem", "NtGrbld"}, {"NtwrkCnfg", "NtGrbld"},
			{"NetPrint", "GrbldOtpt"}, {"LclGrbld", "GrbldOtpt"}, {"NtGrbld", "GrbldOtpt"},
			{"AppDtGnTm", "HrglssDrtnAftrPrnt"},
			{"CblPrtHrdwrOK", "REPEAT"}, {"NtwrkCnfg", "REPEAT"},
			{"PrtPScript", "AvlblVrtlMmry"},
			{"PrtPScript", "PSERRMEM"}, {"AvlblVrtlMmry", "PSERRMEM"},
			{"PrtPScript", "TstpsTxt"}, {"AvlblVrtlMmry", "TstpsTxt"},
			{"GrbldOtpt", "GrbldPS"}, {"AvlblVrtlMmry", "GrbldPS"},
			{"CmpltPgPrntd", "IncmpltPS"}, {"AvlblVrtlMmry", "IncmpltPS"},
			{"PrtDataOut", "PrtFile"},
			{"NtwrkCnfg", "PrtIcon"}, {"PTROFFLINE", "PrtIcon"},
			{"GrbldOtpt", "Problem6"}, {"PrtPScript", "Problem6"}, {"GrbldPS", "Problem6"},
			{"CmpltPgPrntd", "Problem3"}, {"PrtPScript", "Problem3"}, {"IncmpltPS", "Problem3"},
			{"DeskPrntSpd", "NtSpd"}, {"NtwrkCnfg", "NtSpd"}, {"PrtQueue", "NtSpd"},
			{"NetPrint", "Problem2"}, {"DeskPrntSpd", "Problem2"}, {"NtSpd", "Problem2"},
			{"PrtPaper", "PrtStatPaper"},
			{"TnrSpply", "PrtStatToner"},
			{"PrtMem", "PrtStatMem"},
			{"PrtOn", "PrtStatOff"},
		},
	)
}

// Pathfinder returns the Pathfinder lymph-node disease network with 109 nodes
// and 195 edges. From Heckerman et al. (1992).
// Structure only — CPDs too large to hardcode.
func Pathfinder() *models.BayesianNetwork {
	return structureOnly(
		[]string{
			"Fault", "F1", "F97", "F2", "F78", "F3", "F4", "F5", "F53", "F6",
			"F7", "F56", "F8", "F9", "F10", "F55", "F52", "F11", "F12", "F13",
			"F14", "F15", "F16", "F17", "F18", "F19", "F41", "F44", "F20", "F90",
			"F21", "F22", "F23", "F24", "F25", "F26", "F27", "F28", "F92", "F29",
			"F98", "F30", "F31", "F32", "F33", "F34", "F35", "F36", "F37", "F84",
			"F96", "F38", "F39", "F40", "F42", "F43", "F45", "F46", "F47", "F85",
			"F48", "F49", "F50", "F51", "F83", "F54", "F57", "F58", "F59", "F60",
			"F61", "F62", "F63", "F64", "F65", "F66", "F67", "F68", "F69", "F72",
			"F86", "F70", "F71", "F73", "F74", "F75", "F76", "F77", "F79", "F80",
			"F81", "F82", "F87", "F88", "F89", "F91", "F93", "F94", "F95", "F99",
			"F100", "F105", "F101", "F102", "F103", "F104", "F106", "F107", "F108",
		},
		[][2]string{
			{"Fault", "F1"}, {"Fault", "F97"}, {"Fault", "F2"}, {"F97", "F2"},
			{"Fault", "F78"}, {"F97", "F78"},
			{"Fault", "F3"}, {"F97", "F3"}, {"F78", "F3"},
			{"Fault", "F4"}, {"F97", "F4"},
			{"Fault", "F5"}, {"F2", "F5"}, {"F97", "F5"},
			{"Fault", "F53"},
			{"Fault", "F6"}, {"F53", "F6"},
			{"Fault", "F7"},
			{"Fault", "F56"},
			{"Fault", "F8"}, {"F53", "F8"}, {"F56", "F8"},
			{"Fault", "F9"},
			{"Fault", "F10"},
			{"Fault", "F55"},
			{"Fault", "F52"}, {"F55", "F52"},
			{"Fault", "F11"}, {"F52", "F11"},
			{"Fault", "F12"},
			{"Fault", "F13"}, {"F8", "F13"},
			{"Fault", "F14"},
			{"Fault", "F15"},
			{"Fault", "F16"}, {"F15", "F16"},
			{"Fault", "F17"},
			{"Fault", "F18"}, {"F97", "F18"}, {"F17", "F18"},
			{"Fault", "F19"},
			{"Fault", "F41"}, {"Fault", "F44"},
			{"Fault", "F20"}, {"F41", "F20"}, {"F44", "F20"},
			{"Fault", "F90"},
			{"Fault", "F21"}, {"F90", "F21"},
			{"Fault", "F22"}, {"F21", "F22"},
			{"Fault", "F23"}, {"Fault", "F24"}, {"Fault", "F25"},
			{"Fault", "F26"}, {"Fault", "F27"}, {"Fault", "F28"},
			{"Fault", "F92"}, {"F92", "F29"},
			{"Fault", "F98"}, {"F97", "F98"},
			{"Fault", "F30"}, {"F98", "F30"},
			{"Fault", "F31"}, {"F20", "F31"}, {"F41", "F31"}, {"F44", "F31"},
			{"Fault", "F32"}, {"F56", "F32"},
			{"Fault", "F33"}, {"Fault", "F34"}, {"Fault", "F35"},
			{"Fault", "F36"}, {"Fault", "F37"},
			{"Fault", "F84"}, {"F20", "F84"}, {"F41", "F84"},
			{"Fault", "F96"}, {"F41", "F96"}, {"F84", "F96"},
			{"Fault", "F38"}, {"F41", "F38"}, {"F96", "F38"},
			{"Fault", "F39"}, {"F96", "F39"}, {"F84", "F39"}, {"F41", "F39"},
			{"Fault", "F40"}, {"F41", "F40"}, {"F96", "F40"},
			{"Fault", "F42"}, {"F97", "F42"},
			{"Fault", "F43"}, {"Fault", "F45"}, {"Fault", "F46"},
			{"Fault", "F47"}, {"F44", "F47"},
			{"Fault", "F85"}, {"F20", "F85"}, {"F44", "F85"},
			{"Fault", "F48"}, {"F44", "F48"}, {"F85", "F48"},
			{"Fault", "F49"}, {"F44", "F49"},
			{"Fault", "F50"}, {"F44", "F50"},
			{"Fault", "F51"},
			{"Fault", "F83"}, {"F41", "F83"}, {"F84", "F83"},
			{"Fault", "F54"}, {"F83", "F54"},
			{"Fault", "F57"}, {"Fault", "F58"}, {"Fault", "F59"},
			{"Fault", "F60"},
			{"Fault", "F61"}, {"F53", "F61"},
			{"Fault", "F62"}, {"F61", "F62"},
			{"Fault", "F63"}, {"Fault", "F64"}, {"Fault", "F65"},
			{"Fault", "F66"}, {"F61", "F66"},
			{"Fault", "F67"}, {"Fault", "F68"}, {"Fault", "F69"},
			{"Fault", "F72"},
			{"F20", "F86"}, {"F85", "F86"}, {"F84", "F86"}, {"F72", "F86"},
			{"Fault", "F70"}, {"F72", "F70"}, {"F86", "F70"},
			{"Fault", "F71"}, {"F72", "F71"},
			{"Fault", "F73"},
			{"Fault", "F74"}, {"F30", "F74"}, {"F98", "F74"},
			{"Fault", "F75"}, {"Fault", "F76"}, {"Fault", "F77"},
			{"Fault", "F79"}, {"Fault", "F80"},
			{"Fault", "F81"}, {"F72", "F81"}, {"F86", "F81"},
			{"Fault", "F82"}, {"F44", "F82"}, {"F85", "F82"},
			{"Fault", "F87"}, {"F31", "F87"}, {"F41", "F87"},
			{"Fault", "F88"}, {"F31", "F88"}, {"F44", "F88"},
			{"Fault", "F89"}, {"F88", "F89"}, {"F87", "F89"}, {"F31", "F89"}, {"F72", "F89"},
			{"Fault", "F91"},
			{"Fault", "F93"}, {"F83", "F93"}, {"F82", "F93"},
			{"Fault", "F94"}, {"F41", "F94"}, {"F96", "F94"},
			{"Fault", "F95"}, {"F72", "F95"},
			{"Fault", "F99"}, {"F98", "F99"},
			{"Fault", "F100"}, {"F97", "F100"}, {"F4", "F100"},
			{"Fault", "F105"}, {"F97", "F105"},
			{"Fault", "F101"}, {"F105", "F101"},
			{"Fault", "F102"}, {"F105", "F102"},
			{"Fault", "F103"}, {"F97", "F103"},
			{"Fault", "F104"}, {"F97", "F104"},
			{"Fault", "F106"}, {"F97", "F106"},
			{"Fault", "F107"}, {"F41", "F107"}, {"F44", "F107"},
			{"Fault", "F108"},
		},
	)
}

// Pigs returns the Pigs genetic pedigree network with 441 nodes and 592 edges.
// From Cowell et al. (1999).
// Structure only — CPDs too large to hardcode.
func Pigs() *models.BayesianNetwork {
	return structureOnly(
		pigsNodes,
		pigsEdges,
	)
}
