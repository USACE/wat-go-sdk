package plugindatamodel

//ResourceInfo
type ResourceInfo struct {
	Scheme    string `json:"scheme" yaml:"scheme"`                         // s3, azure, local, queue?
	Authority string `json:"authority" yaml:"authority"`                   // bucket, rootdir, queue?
	Path      string `json:"path,omitempty" yaml:"path,omitempty"`         // path to object
	Fragment  string `json:"fragment,omitempty" yaml:"fragment,omitempty"` // hdf path, dss path, etc
}
