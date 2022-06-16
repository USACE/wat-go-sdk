package jobmanager

import (
	"errors"
	"fmt"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/usace/wat-go-sdk/pluginmanager"
)

//JobManager
type JobManager struct {
	Job
	store         filestore.FileStore
	captainCrunch *batch.Batch
}

func Init(job Job, fs filestore.FileStore, batchClient *batch.Batch) JobManager {
	return JobManager{
		Job:           job,
		store:         fs,
		captainCrunch: batchClient,
	}
}
func (jm JobManager) ProcessJob() error {
	resources, err := jm.ProvisionResources()
	fmt.Println(err)
	//add in defer and recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", r)
			jm.DestructResources(resources)
		}
	}()
	err = jm.GeneratePayloads()
	fmt.Println(err)
	for i := 0; i < jm.Job.TaskCount; i++ {
		err = jm.ComputeEvent(i, resources)
		fmt.Println(err)
	}
	err = jm.DestructResources(resources)
	fmt.Println(err)
	return errors.New("Job Processed!")
}

//Job
type Job struct {
	TaskCount int
	DirectedAcyclicGraph
}

type DirectedAcyclicGraph struct {
	Nodes []pluginmanager.ModelManifest
	Links []pluginmanager.ModelLinks
}
type ProvisionedResources struct {
	Plugin                pluginmanager.Plugin
	ComputeEnvironmentARN *string
	JobARN                *string
	QueueARN              *string
}

//provisionresources
func (job Job) ProvisionResources() ([]ProvisionedResources, error) {

	//depends on cloud-resources//
	return nil, errors.New("resources!!!")
}

//provisionresources
func (job Job) DestructResources([]ProvisionedResources) error {

	//depends on cloud-resources//
	return errors.New("ka-blewy!!!")
}

//does this thing need to "run" or "compute"
func (job Job) GeneratePayloads() error {
	for i := 0; i < job.TaskCount; i++ {
		//write out payloads to filestore. How do i get a handle on filestore from here?
	}
	return errors.New("payloads!!!")
}
func (job Job) ComputeEvent(eventNumber int, resources []ProvisionedResources) error {
	for _, n := range job.DirectedAcyclicGraph.Nodes {
		submitTask(resources, n)
	}
	return errors.New(fmt.Sprintf("computing event %v", eventNumber))
}

//does this thing need to "run" or "compute"
func submitTask(resources []ProvisionedResources, manifest pluginmanager.ModelManifest) error {
	//depends on cloud-resources//
	//submit to batch.
	return errors.New("task for " + manifest.Plugin.Name)
}
