package plugindatamodel

type Data struct {
	//FileName describes the name of the file that needs to be input or output.
	FileName string `json:"filename" yaml:"filename"`
	//Path describes the specific information in the file (e.g. /a/b/c/d/e/f for dss)
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

//acceptable formats? format options?
//optional/required
