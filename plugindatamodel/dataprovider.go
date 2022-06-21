package plugindatamodel

//Provider
type Provider struct {
	//ModelIdentifiier
	ModelIdentifier `json:"model_identifier" yaml:"model_identifier"`
	//Data
	Data `json:"data" yaml:"data"`
}
