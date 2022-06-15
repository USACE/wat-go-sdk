package pluginmanager

//ModelConfiguration is a model name and an optional model alternative
type ModelConfiguration struct {
	Name        string `json:"model_name" yaml:"model_name"`
	Alternative string `json:"model_alternative,omitempty" yaml:"model_alternative,omitempty"`
}
