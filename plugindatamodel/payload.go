package plugindatamodel

// Payload
type Payload struct {
	Id int `json:"payload_id" yaml:"payload_id"`
	//ModelIdentifier   `json:"model_configuration" yaml:"model_configuration"`
	Inputs            []ResourcedFileData `json:"inputs" yaml:"inputs"`
	OutputDestination ResourceInfo        `json:"output_destination" yaml:"output_destination"`
	Outputs           []FileData          `json:"outputs" yaml:"outputs"`
}
