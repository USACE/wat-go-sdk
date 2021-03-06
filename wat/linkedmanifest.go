package wat

import "github.com/usace/wat-go-sdk/plugin"

// LinkedModelManifest represents a model manifest that has been linked (i.e. has dependencies)
type LinkedModelManifest struct {
	ManifestID    string `json:"linked_manifest_id" yaml:"linked_manifest_id"`
	plugin.Plugin `json:"plugin" yaml:"plugin"`
	plugin.Model  `json:"model" yaml:"model"`
	Inputs        []LinkedFileData  `json:"inputs" yaml:"inputs"`
	Outputs       []plugin.FileData `json:"outputs" yaml:"outputs"`
}

// LinkedFileData
type LinkedFileData struct {
	// Id is an internal element generated to identify any data element.
	Id string `json:"id,omitempty" yaml:"id,omitempty"`
	// FileName describes the name of the file that needs to be input or output.
	FileName string `json:"filename" yaml:"filename"`
	// Provider a provider is a specific output data element from a manifest.
	SourceDataId  string                   `json:"source_data_identifier" yaml:"source_data_identifier"`
	InternalPaths []LinkedInternalPathData `json:"internal_paths,omitempty" yaml:"internal_paths,omitempty"`
}

func (lf LinkedFileData) HasInternalPaths() bool {
	return len(lf.InternalPaths) > 0
}

// LinkedInternalPathData
type LinkedInternalPathData struct {
	// Id is an internal element generated to identify any data element.
	Id string `json:"id,omitempty" yaml:"id,omitempty"`
	// PathName describes the internal path location to the data needed or produced.
	PathName     string `json:"pathname" yaml:"pathname"`
	SourcePathID string `json:"source_path_identifier,omitempty" yaml:"source_path_identifier,omitempty"`
	SourceFileID string `json:"source_file_identifier" yaml:"source_file_identifier"`
}

func (lm LinkedModelManifest) producesFile(fileId string) (plugin.FileData, bool) {
	for _, output := range lm.Outputs {
		if fileId == output.Id {
			return output, true
		}
	}
	return plugin.FileData{}, false
}

func (lf LinkedModelManifest) producesInternalPath(internalPath LinkedInternalPathData) (string, string, bool) {
	output, ok := lf.producesFile(internalPath.SourceFileID)
	if ok {
		if len(output.InternalPaths) > 0 {
			for _, ip := range output.InternalPaths {
				if internalPath.SourcePathID == ip.Id {
					return ip.PathName, output.FileName, true
				}
			}
		}
		return "", output.FileName, true
	}
	return "", "", false
}

func (lm LinkedModelManifest) producesDependency(linkedFile LinkedFileData) bool {
	for _, output := range lm.Outputs {
		if linkedFile.SourceDataId == output.Id {
			return true
		}
		if linkedFile.HasInternalPaths() {
			for _, internalPath := range linkedFile.InternalPaths {
				if internalPath.SourceFileID == output.Id {
					return true
				}
			}
		}
	}
	return false
}
