package jobmanager

import "github.com/usace/wat-go-sdk/plugindatamodel"

type ProvisionedResources struct {
	LinkedManifest        plugindatamodel.LinkedModelManifest
	ComputeEnvironmentARN *string
	JobARN                []*string
	QueueARN              *string
}
