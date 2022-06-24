package plugindatamodel

//ResourceInfo
type ResourceInfo struct {
	Store string `json:"store" yaml:"store"`                   // s3, azure, local, queue?
	Root  string `json:"root" yaml:"root"`                     // bucket, rootdir, queue?
	Path  string `json:"path,omitempty" yaml:"path,omitempty"` // path to object
}
