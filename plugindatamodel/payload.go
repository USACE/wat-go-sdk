package plugindatamodel

// Payload
type Payload struct {
	Id int `json:"payload_id" yaml:"payload_id"`
	//ModelIdentifier   `json:"model_configuration" yaml:"model_configuration"`
	Inputs  []ResourcedFileData `json:"inputs" yaml:"inputs"`
	Outputs []ResourcedFileData `json:"outputs" yaml:"outputs"`
}
