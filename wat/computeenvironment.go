package wat

type ComputeResourceRequirements struct {
	LinkedManifestID   string `json:"linked_manifest_id" yaml:"linked_manifest_id"`
	ComputeEnvironment string `json:"compute_environment" yaml:"compute_environment"` //is this provided as JSON?
	Definition         string `json:"Definition"`
	Queue              string `json:"jobQueueFile"`
	//QUEUE string `json:"job_queue" yaml:"job_queue"`
	//JobDefinition string `json:"job_definition" yaml:"job_definition"`
}

//from seth's implementation
type AWSBatchPayload struct {
	ComputeEnvironmentFile string `json:"computeEnvironmentFile"`
	JobDefinitionFile      string `json:"jobDefinitionFile"`
	JobQueueFile           string `json:"jobQueueFile"`
	NewJob                 string `json:"newJob"` //not sure why we need this atm.
}
