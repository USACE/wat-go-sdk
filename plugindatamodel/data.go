package plugindatamodel

type Data struct {
	//FileName describes the name of the file that needs to be input or output.
	FileName string `json:"filename" yaml:"filename"`
	//Path describes the specific information in the file (e.g. /a/b/c/d/e/f for dss)
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

/*
type DataContainer struct {
	Type        string `json:"type" yaml:"type"`
	DataElement `json:"data" yaml:"data"`
}
type DataElement interface {
	DataElementType() string
}

//FileData
type File struct {
	//FileName
	FileName string `json:"filename" yaml:"filename"`
}

func (fd File) DataElementType() string {
	return strings.ToLower(reflect.TypeOf(fd).Name())
}

//FileAndPathData
type FileAndPath struct {
	Type string `json:"type" yaml:"type"`
	//FileName
	FileName string `json:"filename" yaml:"filename"`
	//Path describes the specific information in the file (e.g. /a/b/c/d/e/f for dss)
	Path string `json:"path" yaml:"path"`
}
*/
//acceptable formats? format options?
//optional/required
