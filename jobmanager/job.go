package jobmanager

import (
	"errors"
	"fmt"
	"os"

	"github.com/usace/wat-go-sdk/plugindatamodel"
	"gopkg.in/yaml.v3"
)

//Job
type Job struct {
	Id                string                       `json:"job_identifier" yaml:"job_identifier"`
	EventStartIndex   int                          `json:"event_start_index" yaml:"event_start_index"`
	EventEndIndex     int                          `json:"event_end_index" yaml:"event_end_index"`
	Dag               DirectedAcyclicGraph         `json:"directed_acyclic_graph" yaml:"directed_acyclic_graph"`
	OutputDestination plugindatamodel.ResourceInfo `json:"output_destination" yaml:"output_destination"`
}

//ProvisionResources
func (job *Job) ProvisionResources() error {
	//make sure job arn list is provisioned for the total number of events to be computed.
	//depends on cloud-resources//
	job.Dag.Resources = make(map[string]ProvisionedResources, len(job.Dag.LinkedManifests))
	return errors.New("resources!!!")
}

//DestructResources
func (job Job) DestructResources() error {

	//depends on cloud-resources//
	return errors.New("ka-blewy!!!")
}

//GeneratePayloads
func (job Job) GeneratePayloads() error {
	//generate payloads for all events up front?
	for eventIndex := job.EventStartIndex; eventIndex < job.EventEndIndex; eventIndex++ {
		//write out payloads to filestore. How do i get a handle on filestore from here?
		outputDestinationPath := fmt.Sprintf("%vevent_%v/", job.OutputDestination.Path, eventIndex)
		for _, n := range job.Dag.LinkedManifests {
			fmt.Println(n.ImageAndTag, outputDestinationPath)
			payload, err := job.Dag.GeneratePayload(n, eventIndex, job.OutputDestination)
			if err != nil {
				return err
			}
			//fmt.Println(payload)
			bytes, err := yaml.Marshal(payload)
			if err != nil {
				return err
			}
			fmt.Println("")
			fmt.Println(string(bytes))
			//put payload in s3
			path := outputDestinationPath + n.Plugin.Name + "_payload.yml"
			fmt.Println("putting object in fs:", path)
			//_, err = fs.PutObject(path, bytes)
			if _, err = os.Stat(path); os.IsNotExist(err) {
				os.MkdirAll(outputDestinationPath, 0644)
			}
			err = os.WriteFile(path, bytes, 0644)
			if err != nil {
				fmt.Println("failure to push payload to filestore:", err)
				return err
			}
		}
	}
	fmt.Println("payloads!!!")
	return nil
}

//ComputeEvent
func (job Job) ComputeEvent(eventIndex int) error {
	for _, n := range job.Dag.LinkedManifests {
		job.submitTask(n, eventIndex)
	}
	fmt.Println(fmt.Sprintf("computing event %v", eventIndex))
	return nil
}
func (job Job) submitTask(manifest LinkedModelManifest, eventIndex int) error {
	//depends on cloud-resources//
	offset := eventIndex - job.EventStartIndex
	dependencies, err := job.Dag.Dependencies(manifest, offset)
	if err != nil {
		return err
	} else {
		fmt.Print(dependencies)
	}
	payloadPath := fmt.Sprintf("%vevent_%v/%v_payload.yml", job.OutputDestination.Path, eventIndex, manifest.Plugin.Name)
	fmt.Println(payloadPath)
	//submit to batch.
	batchjobarn := "batch arn returned."
	//set job arn
	resources, ok := job.Dag.Resources[manifest.ManifestID]
	if ok {
		resources.JobARN[offset] = &batchjobarn
		job.Dag.Resources[manifest.ManifestID] = resources
	} else {
		return errors.New("task for " + manifest.Plugin.Name)
	}
	return nil
}
