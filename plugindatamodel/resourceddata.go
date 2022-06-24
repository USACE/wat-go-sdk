package plugindatamodel

//LinkedData
type ResourcedData struct {
	//Data
	FileData `json:"data" yaml:"data"`
	//Provider a provider is a specific output data element from a manifest.
	ResourceInfo `json:"resource_info" yaml:"resource_info"`
}
