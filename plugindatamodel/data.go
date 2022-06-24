package plugindatamodel

type FileData struct {
	//FileName describes the name of the file that needs to be input or output.
	FileName string `json:"filename" yaml:"filename"`
	//InternalPaths (optional) describe the specific information in the file (e.g. /a/b/c/d/e/f for dss)
	InternalPaths string `json:"internal_paths,omitempty" yaml:"internal_paths,omitempty"`
}
type InternalPathData struct {
	//PathName describes the internal path location to the data needed or produced.
	PathName string `json:"pathname" yaml:"pathname"`
	//Type (optional) describes the type of information at the path (e.g. flow time-series)
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

//acceptable formats? format options?
//optional/required
