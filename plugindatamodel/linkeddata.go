package plugindatamodel

//LinkedData
type LinkedData struct {
	//Data
	FileData `json:"data" yaml:"data"`
	//Provider a provider is a specific output data element from a manifest.
	SourceData          FileData `json:"source_data" yaml:"source_data"`
	ProducingManifestID int      `json:"producing_manifest_identifier" yaml:"producing_manifest_identifier"` //could just be modelidentifier and plugin i believe.
}
