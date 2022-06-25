package plugindatamodel

//LinkedData
type LinkedFileData struct {
	//Id is an internal element generated to identify any data element.
	Id int `json:"id,omitempty" yaml:"id,omitempty"`
	//FileName describes the name of the file that needs to be input or output.
	FileName string `json:"filename" yaml:"filename"`
	//Provider a provider is a specific output data element from a manifest.
	SourceData          int                      `json:"source_data_identifier" yaml:"source_data_identifier"`
	ProducingManifestID int                      `json:"producing_manifest_identifier" yaml:"producing_manifest_identifier"` //could just be modelidentifier and plugin i believe.
	InternalPaths       []LinkedInternalPathData `json:"internal_paths,omitempty" yaml:"internal_paths,omitempty"`
}

//LinkedInternalPathData
type LinkedInternalPathData struct {
	//Id is an internal element generated to identify any data element.
	Id int `json:"id,omitempty" yaml:"id,omitempty"`
	//PathName describes the internal path location to the data needed or produced.
	PathName   string `json:"pathname" yaml:"pathname"`
	SourcePath string `json:"source_path_name,omitempty" yaml:"source_path_name,omitempty"`
	//Provider a provider is a specific output data element from a manifest.
	SourceData          int `json:"source_data_identifier" yaml:"source_data_identifier"`
	ProducingManifestID int `json:"producing_manifest_identifier" yaml:"producing_manifest_identifier"` //could just be modelidentifier and plugin i believe.
}
