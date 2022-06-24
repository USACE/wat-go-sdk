package plugindatamodel

//LinkedData
type ResourcedFileData struct {
	//FileName describes the name of the file that needs to be input or output.
	FileName string `json:"filename" yaml:"filename"`
	//Provider a provider is a specific output data element from a manifest.
	ResourceInfo  `json:"resource_info" yaml:"resource_info"`
	InternalPaths []ResourcedInternalPathData `json:"internal_paths,omitempty" yaml:"internal_paths,omitempty"`
}

type ResourcedInternalPathData struct {
	//PathName describes the internal path location to the data needed or produced.
	PathName     string `json:"pathname" yaml:"pathname"`
	FileName     string `json:"filename,omitempty" yaml:"filename,omitempty"`
	InternalPath string `json:"internal_path,omitempty" yaml:"internal_path,omitempty"`
	//Provider a provider is a specific output data element from a manifest.
	ResourceInfo `json:"resource_info" yaml:"resource_info"`
}
