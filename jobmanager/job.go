package jobmanager

import (
	"errors"
	"fmt"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/usace/wat-go-sdk/pluginmanager"
	"gopkg.in/yaml.v3"
)

//JobManager
type JobManager struct {
	job           Job
	store         filestore.FileStore
	captainCrunch *batch.Batch
}

func Init(job Job, fs filestore.FileStore, batchClient *batch.Batch) JobManager {
	return JobManager{
		job:           job,
		store:         fs,
		captainCrunch: batchClient,
	}
}
func (jm JobManager) ProcessJob() error {
	resources, err := jm.job.ProvisionResources()
	fmt.Println(err)
	//add in defer and recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", r)
			fmt.Println("Tearing Down Resources")
			err = jm.job.DestructResources(resources)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	err = jm.job.GeneratePayloads(jm.store)
	fmt.Println(err)
	//create error channel.
	//create waitgroups to throttle compute resources?
	for i := 0; i < jm.job.EventCount; i++ {
		go func(index int) {
			err = jm.job.ComputeEvent(index, resources)
			fmt.Println(err)
		}(i)

	}
	//need a wait group or a buffer channel to stall the destruction until we finish the jobs
	err = jm.job.DestructResources(resources)
	fmt.Println(err)
	return errors.New("Job Processed!")
}

//Job
type Job struct {
	EventCount int
	DirectedAcyclicGraph
	OutputDestination pluginmanager.ResourceInfo
	InputSource       pluginmanager.ResourceInfo
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

//GeneratePayloads
func (job Job) GeneratePayloads(fs filestore.FileStore) error {
	//generate payloads for all events up front?
	for i := 0; i < job.EventCount; i++ {
		//write out payloads to filestore. How do i get a handle on filestore from here?
		outputDestinationPath := fmt.Sprintf("%v/%v%v", job.OutputDestination.Fragment, "event_", i)
		for _, n := range job.DirectedAcyclicGraph.Nodes {
			fmt.Println(n.ImageAndTag, outputDestinationPath)
			payload := outputDestinationPath + n.Plugin.Name + "payload.yml"
			bytes, err := yaml.Marshal(payload)
			if err != nil {
				panic(err)
			}
			//put payload in s3
			path := outputDestinationPath + "/" + n.Plugin.Name + "_payload.yml"
			fmt.Println("putting object in fs:", path)
			_, err = fs.PutObject(path, bytes)
			if err != nil {
				fmt.Println("failure to push payload to filestore:", err)
				panic(err)
			}
		}
	}
	return errors.New("payloads!!!")
}
func (job Job) ComputeEvent(eventNumber int, resources []ProvisionedResources) error {
	for _, n := range job.DirectedAcyclicGraph.Nodes {
		submitTask(resources, n)
	}
	return errors.New(fmt.Sprintf("computing event %v", eventNumber))
}

func submitTask(resources []ProvisionedResources, manifest pluginmanager.ModelManifest) error {
	//depends on cloud-resources//
	//submit to batch.
	return errors.New("task for " + manifest.Plugin.Name)
}
