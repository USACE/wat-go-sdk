package plugindatamodel

//LinkedModelManifest
type LinkedModelManifest struct {
	Plugin          `json:"plugin" yaml:"plugin"`
	ModelIdentifier `json:"model_identifier" yaml:"model_identifier"`
	Inputs          []LinkedData `json:"inputs" yaml:"inputs"`
	Outputs         []Data       `json:"outputs" yaml:"outputs"`
}

//can use this struct to create a payload for a plugin
