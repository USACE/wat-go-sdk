package pluginmanager

//ModelManifest is defined by a plugin, model configuration, compute resources and inputs and ouptuts necessary for a model, is recognizable by a Model Library MCAT
type ModelManifest struct {
	Plugin             `json:"plugin" yaml:"plugin"`
	ModelConfiguration `json:"model_configuration" yaml:"model_configuration"`
	//ModelComputeResources `json:"model_compute_resources" yaml:"model_compute_resources"`
	Inputs  []DataDescription `json:"inputs" yaml:"inputs"`
	Outputs []DataDescription `json:"outputs" yaml:"outputs"`
}
