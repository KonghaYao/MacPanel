package response

type HomebrewStatus struct {
	IsExist      bool   `json:"isExist"`
	Version      string `json:"version"`
	Prefix       string `json:"prefix"`
	FormulaCount int    `json:"formulaCount"`
	CaskCount    int    `json:"caskCount"`
	BrewPath     string `json:"brewPath"`
	CellarPath   string `json:"cellarPath"`
}

type HomebrewPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"`
}

type HomebrewSearchResult struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type HomebrewMirrorConfig struct {
	BottleDomain  string `json:"bottleDomain"`
	APIDomain     string `json:"apiDomain"`
	BrewGitRemote string `json:"brewGitRemote"`
	CoreGitRemote string `json:"coreGitRemote"`
	CaskGitRemote string `json:"caskGitRemote"`
}

type HomebrewDoctorResult struct {
	Output string `json:"output"`
}
