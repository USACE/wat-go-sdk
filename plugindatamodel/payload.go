package plugindatamodel

// Payload
type Payload struct {
	ModelIdentifier   `json:"model_configuration" yaml:"model_configuration"`
	Inputs            []ResourcedData `json:"inputs" yaml:"inputs"`
	OutputDestination ResourceInfo    `json:"output_destination" yaml:"output_destination"`
	Outputs           []FileData      `json:"outputs" yaml:"outputs"`
}
