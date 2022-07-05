package wat

import (
	"errors"
	"fmt"
	"os"

	"github.com/usace/wat-go-sdk/plugin"
	"gopkg.in/yaml.v3"
)

//Job
type Job struct {
	Id                string               `json:"job_identifier" yaml:"job_identifier"`
	EventStartIndex   int                  `json:"event_start_index" yaml:"event_start_index"`
	EventEndIndex     int                  `json:"event_end_index" yaml:"event_end_index"`
	Dag               DirectedAcyclicGraph `json:"directed_acyclic_graph" yaml:"directed_acyclic_graph"`
	OutputDestination plugin.ResourceInfo  `json:"output_destination" yaml:"output_destination"`
}
type PayloadProcessor func(payload plugin.ModelPayload, job Job, eventIndex int, modelManifest LinkedModelManifest) error

//ProvisionResources
func (job *Job) ProvisionResources() error {
	//make sure job arn list is provisioned for the total number of events to be computed.
	//depends on cloud-resources//
	resources := make(map[string]ProvisionedResources, len(job.Dag.LinkedManifests))
	for _, lm := range job.Dag.LinkedManifests {
		qarn := lm.ManifestID                  //@TODO: provisioned with batch
		computeEnviornmentArn := lm.ManifestID //@TODO: provisioned with batch
		lmResource := ProvisionedResources{
			LinkedManifestID:      lm.ManifestID,
			ComputeEnvironmentARN: &computeEnviornmentArn,
			JobARN:                []*string{},
			QueueARN:              &qarn,
		}
		resources[lm.ManifestID] = lmResource
	}
	job.Dag.Resources = resources
	return nil
}

//DestructResources
func (job Job) DestructResources() error {

	//depends on cloud-resources//
	fmt.Println("ka-blewy!!!")
	return nil
}
func (job Job) eventLevelOutputDirectory(eventIndex int) string {
	return fmt.Sprintf("%vevent_%v/", job.OutputDestination.Path, eventIndex)
}
func (job Job) generatePayloadPath(eventIndex int, manifest LinkedModelManifest) string {
	return fmt.Sprintf("%v%v_payload.yml", job.eventLevelOutputDirectory(eventIndex), manifest.Plugin.Name)
}
func (job Job) ValidateLinkages() error {
	return job.payloadLooper(payloadValidator)
}
func payloadValidator(payload plugin.ModelPayload, job Job, eventIndex int, modelManifest LinkedModelManifest) error {
	return nil
}
func (job Job) payloadLooper(processor PayloadProcessor) error {
	for eventIndex := job.EventStartIndex; eventIndex < job.EventEndIndex; eventIndex++ {
		//write out payloads to filestore. How do i get a handle on filestore from here?
		outputDestinationPath := job.eventLevelOutputDirectory(eventIndex)
		for _, n := range job.Dag.LinkedManifests {
			fmt.Println(n.ImageAndTag, outputDestinationPath)
			payload, err := job.Dag.GeneratePayload(n, eventIndex, job.OutputDestination)
			if err != nil {
				return err
			}
			err = processor(payload, job, eventIndex, n)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func payloadWriter(payload plugin.ModelPayload, job Job, eventIndex int, modelManifest LinkedModelManifest) error {
	bytes, err := yaml.Marshal(payload)
	if err != nil {
		return err
	}
	fmt.Println("")
	fmt.Println(string(bytes))
	//put payload in s3
	outputDestinationPath := job.eventLevelOutputDirectory(eventIndex)
	path := job.generatePayloadPath(eventIndex, modelManifest)
	fmt.Println("putting object in fs:", path)
	//_, err = fs.PutObject(path, bytes) //@TODO: replace with FileStore.
	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(outputDestinationPath, 0644)
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		fmt.Println("failure to push payload to filestore:", err)
		return err
	}
	return nil
}

//GeneratePayloads
func (job Job) GeneratePayloads() error {
	err := job.payloadLooper(payloadWriter)
	if err != nil {
		return err
	}
	fmt.Println("payloads!!!")
	return nil
}

//ComputeEvent
func (job Job) ComputeEvent(eventIndex int) error {
	for _, n := range job.Dag.LinkedManifests {
		job.submitTask(n, eventIndex)
	}
	fmt.Printf("computing event %v\n", eventIndex)
	return nil
}
func (job *Job) submitTask(manifest LinkedModelManifest, eventIndex int) error {
	//depends on cloud-resources//
	offset := eventIndex - job.EventStartIndex
	dependencies, err := job.Dag.Dependencies(manifest, offset)
	if err != nil {
		return err
	} else {
		fmt.Println(dependencies)
	}
	payloadPath := job.generatePayloadPath(eventIndex, manifest)
	fmt.Println(payloadPath)
	//submit to batch.
	batchjobarn := "batch arn returned." //@TODO: replace with call to batch
	//set job arn
	resources, ok := job.Dag.Resources[manifest.ManifestID]
	if ok {
		resources.JobARN = append(resources.JobARN, &batchjobarn)
		job.Dag.Resources[manifest.ManifestID] = resources
	} else {
		return errors.New("task for " + manifest.Plugin.Name)
	}
	return nil
}
