package wat

type ComputeResourceRequirements struct {
	LinkedManifestID   string `json:"linked_manifest_id" yaml:"linked_manifest_id"`
	ComputeEnvironment string `json:"compute_environment" yaml:"compute_environment"` //is this provided as JSON?
	//QUEUE string `json:"job_queue" yaml:"job_queue"`
	//JobDefinition string `json:"job_definition" yaml:"job_definition"`
}
