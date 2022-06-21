package plugindatamodel

//LinkedData
type LinkedData struct {
	//Data
	Data `json:"data" yaml:"data"`
	//Provider a provider is a specific output data element from a manifest.
	SourceData Data          `json:"source_data" yaml:"source_data"`
	Provider   ModelManifest `json:"provider" yaml:"provider"` //could just be modelidentifier and plugin i believe.
}
