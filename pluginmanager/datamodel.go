package pluginmanager

// Payload
type Payload struct {
	ModelConfiguration `json:"model_configuration" yaml:"model_configuration"`
	ModelLinks         `json:"model_links" yaml:"model_links"`
}

// ModelLinks
type ModelLinks struct {
	Inputs  []LinkedDataDescription `json:"inputs" yaml:"inputs"`
	Outputs []LinkedDataDescription `json:"outputs" yaml:"outputs"`
}

//LinkedDataDescription
type LinkedDataDescription struct {
	DataDescription `json:"description" yaml:"description"`
	ResourceInfo    `json:"resource_info" yaml:"resource_info"`
}

//ModelManifest
type ModelManifest struct {
	Plugin             `json:"plugin" yaml:"plugin"`
	ModelConfiguration `json:"model_configuration" yaml:"model_configuration"`
	//ModelComputeResources `json:"model_compute_resources" yaml:"model_compute_resources"`
	Inputs  []DataDescription `json:"inputs" yaml:"inputs"`
	Outputs []DataDescription `json:"outputs" yaml:"outputs"`
}

// ModelConfiguration
type ModelConfiguration struct {
	Name        string `json:"name" yaml:"name"`
	Alternative string `json:"alternative,omitempty" yaml:"alternative,omitempty"`
}

//ResourceInfo
type ResourceInfo struct {
	Scheme    string `json:"scheme" yaml:"scheme"`                         // s3, azure, local, queue?
	Authority string `json:"authority" yaml:"authority"`                   // bucket, rootdir, queue?
	Path      string `json:"path,omitempty" yaml:"path,omitempty"`         // path to object
	Fragment  string `json:"fragment,omitempty" yaml:"fragment,omitempty"` // hdf path, dss path, etc
}

//DataDescription
type DataDescription struct {
	Name      string `json:"name" yaml:"name"`
	Parameter string `json:"parameter" yaml:"parameter"`
	Format    string `json:"format" yaml:"format"`
}

// Plugin
type Plugin struct {
	Name        string   `json:"name" yaml:"name"`
	ImageAndTag string   `json:"image_and_tag" yaml:"image_and_tag"`
	Command     []string `json:"command" yaml:"command"`
}
