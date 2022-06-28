package plugindatamodel

// Payload
type ModelPayload struct {
	Id         int                 `json:"payload_id" yaml:"payload_id"`
	EventIndex int                 `json:"event_index" yaml:"event_index"`
	Inputs     []ResourcedFileData `json:"inputs" yaml:"inputs"`
	Outputs    []ResourcedFileData `json:"outputs" yaml:"outputs"`
}

//ResourcedFileData
type ResourcedFileData struct {
	//Id is an internal element generated to identify any data element.
	Id int `json:"id,omitempty" yaml:"id,omitempty"`
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
