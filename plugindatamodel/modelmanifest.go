package plugindatamodel

//ModelManifest
type ModelManifest struct {
	ManifestID      int `json:"manifest_id" yaml:"manifest_id"`
	Plugin          `json:"plugin" yaml:"plugin"`
	ModelIdentifier `json:"model_identifier" yaml:"model_identifier"`
	Inputs          []Data `json:"inputs" yaml:"inputs"`
	Outputs         []Data `json:"outputs" yaml:"outputs"`
}

//can use this info to allow a user to link model manifests together.
