package plugindatamodel

//ModelManifest
type ModelManifest struct {
	ManifestID      int `json:"manifest_id" yaml:"manifest_id"`
	Plugin          `json:"plugin" yaml:"plugin"`
	ModelIdentifier `json:"model_identifier" yaml:"model_identifier"`
	Inputs          []FileData `json:"inputs" yaml:"inputs"`
	Outputs         []FileData `json:"outputs" yaml:"outputs"`
}

//can use this info to allow a user to link model manifests together.
type FileData struct {
	//Id is an internal element generated to identify any data element.
	Id int `json:"id,omitempty" yaml:"id,omitempty"`
	//FileName describes the name of the file that needs to be input or output.
	FileName string `json:"filename" yaml:"filename"`
	//InternalPaths (optional) describe the specific information in the file (e.g. /a/b/c/d/e/f for dss)
	InternalPaths []InternalPathData `json:"internal_paths,omitempty" yaml:"internal_paths,omitempty"`
}
type InternalPathData struct {
	//Id is an internal element generated to identify any data element.
	Id int `json:"id,omitempty" yaml:"id,omitempty"`
	//PathName describes the internal path location to the data needed or produced.
	PathName string `json:"pathname" yaml:"pathname"`
	//Type (optional) describes the type of information at the path (e.g. flow time-series)
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

//acceptable formats? format options?
//optional/required
