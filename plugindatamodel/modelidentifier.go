package plugindatamodel

// ModelConfiguration
type ModelIdentifier struct {
	Name        string `json:"name" yaml:"name"`
	Alternative string `json:"alternative,omitempty" yaml:"alternative,omitempty"`
}
