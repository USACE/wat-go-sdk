package wat

type ProvisionedResources struct {
	LinkedManifestID      string //plugindatamodel.LinkedModelManifest
	ComputeEnvironmentARN *string
	JobARN                []*string
	QueueARN              *string
}
