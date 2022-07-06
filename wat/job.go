package wat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/usace/wat-go-sdk/plugin"
	"gopkg.in/yaml.v3"
)

// JobManifest
type JobManifest struct {
	Id                      string                   `json:"job_identifier" yaml:"job_identifier"`
	EventStartIndex         int                      `json:"event_start_index" yaml:"event_start_index"`
	EventEndIndex           int                      `json:"event_end_index" yaml:"event_end_index"`
	Models                  []plugin.ModelIdentifier `json:"models" yaml:"models"`
	LinkedManifestResources []plugin.ResourceInfo    `json:"linked_manifests" yaml:"linked_manifests"`
	OutputDestination       plugin.ResourceInfo      `json:"output_destination" yaml:"output_destination"`
}

func (jm JobManifest) ConvertToJob() (Job, error) {
	job := Job{
		Id:                uuid.New().String(), //make a uuid version 4
		EventStartIndex:   jm.EventStartIndex,
		EventEndIndex:     jm.EventEndIndex,
		OutputDestination: jm.OutputDestination,
	}

	linkedManifests := make([]LinkedModelManifest, len(jm.LinkedManifestResources))

	for idx, resourceInfo := range jm.LinkedManifestResources {
		// fmt.Println(resourceInfo.Path)
		lm := LinkedModelManifest{}
		file, err := os.Open(resourceInfo.Path) //replace with filestore? injected?
		if err != nil {
			return job, err
		}

		defer file.Close()
		b, err := ioutil.ReadAll(file)
		if err != nil {
			return job, err
		}

		err = yaml.Unmarshal(b, &lm)
		if err != nil {
			return job, err
		}

		linkedManifests[idx] = lm
	}
	job.Dag = DirectedAcyclicGraph{
		Models:          jm.Models,
		LinkedManifests: linkedManifests,
		Resources:       map[string]provisionedResources{},
	}
	return job, nil
}

// JobManager
type JobManager struct {
	job Job
	//store         filestore.FileStore
	//captainCrunch *batch.Batch
}

func Init(jobManifest JobManifest) (JobManager, error) { //, fs filestore.FileStore, batchClient *batch.Batch) JobManager {
	jobManager := JobManager{}
	job, err := jobManifest.ConvertToJob()
	if err != nil {
		return jobManager, err
	}

	orderedManifests, err := job.Dag.TopologicallySort()
	if err != nil {
		return jobManager, err
	}

	job.Dag.LinkedManifests = orderedManifests
	jobManager.job = job

	return jobManager, nil
}

func (jm JobManager) ProcessJob(logLevel string) error {
	err := jm.job.ProvisionResources()
	if err != nil {
		return err
	}

	// add in defer and recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", r)
			fmt.Println("Tearing Down Resources")
			err = jm.job.DestructResources()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	err = jm.job.GeneratePayloads(logLevel) //jm.store
	if err != nil {
		fmt.Println(err)
		return err
	}

	//create error channel.
	for i := jm.job.EventStartIndex; i < jm.job.EventEndIndex; i++ {
		go func(index int) {
			err = jm.job.ComputeEvent(index, logLevel)
			fmt.Println(err)
		}(i)

	}
	//need a wait group or a buffer channel to stall the destruction until we finish the jobs
	err = jm.job.DestructResources()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("\nJob Processed!\n\n")
	return nil
}

func (jm JobManager) Validate() error {
	err := jm.job.ValidateLinkages() //evaluate if this can be trimmed down to "validateLinkages"
	if err != nil {
		return err
	}
	_, err = jm.job.Dag.TopologicallySort()
	if err != nil {
		return err
	}

	return nil
}

// Job
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
	resources := make(map[string]provisionedResources, len(job.Dag.LinkedManifests))
	for _, lm := range job.Dag.LinkedManifests {
		queueArn := lm.ManifestID              //@TODO: provisioned with batch
		computeEnvironmentArn := lm.ManifestID //@TODO: provisioned with batch
		lmResource := provisionedResources{
			LinkedManifestID:      lm.ManifestID,
			ComputeEnvironmentARN: &computeEnvironmentArn,
			JobARN:                []*string{},
			QueueARN:              &queueArn,
		}
		resources[lm.ManifestID] = lmResource
	}
	job.Dag.Resources = resources
	return nil
}

//DestructResources
func (job Job) DestructResources() error {

	//depends on cloud-resources//
	fmt.Println("\nPlaceholder: Deallocate / Deregister / Destroy resources")
	return nil
}

func (job Job) eventLevelOutputDirectory(eventIndex int) string {
	return fmt.Sprintf("%vevent_%v/", job.OutputDestination.Path, eventIndex)
}

func (job Job) generatePayloadPath(eventIndex int, manifest LinkedModelManifest) string {
	return fmt.Sprintf("%v%v_payload.yml", job.eventLevelOutputDirectory(eventIndex), manifest.Plugin.Name)
}

func (job Job) ValidateLinkages() error {
	return job.payloadLooper(payloadValidator, "Info")
}

func payloadValidator(payload plugin.ModelPayload, job Job, eventIndex int, modelManifest LinkedModelManifest) error {
	return nil
}

func (job Job) payloadLooper(processor PayloadProcessor, logLevel string) error {

	for eventIndex := job.EventStartIndex; eventIndex < job.EventEndIndex; eventIndex++ {
		//write out payloads to filestore. How do i get a handle on filestore from here?
		outputDestinationPath := job.eventLevelOutputDirectory(eventIndex)

		for _, n := range job.Dag.LinkedManifests {

			if logLevel == "Info" {
				fmt.Println("\n", n.ImageAndTag)
				fmt.Println("\t", outputDestinationPath)
			}

			payload, err := job.Dag.GeneratePayload(n, eventIndex, job.OutputDestination, logLevel)
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

	// fmt.Println("")
	// fmt.Println(string(bytes))

	// put payload in s3
	outputDestinationPath := job.eventLevelOutputDirectory(eventIndex)
	path := job.generatePayloadPath(eventIndex, modelManifest)

	// fmt.Println("putting object in fs:", path)
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
func (job Job) GeneratePayloads(logLevel string) error {
	err := job.payloadLooper(payloadWriter, logLevel)
	if err != nil {
		return err
	}
	// fmt.Println("Payloads Generated")
	return nil
}

//ComputeEvent
func (job Job) ComputeEvent(eventIndex int, logLevel string) error {
	for _, n := range job.Dag.LinkedManifests {
		job.submitTask(n, eventIndex, logLevel)
	}
	if logLevel == "Info" {
		fmt.Printf("computing event %v\n", eventIndex)
	}

	return nil
}

func (job *Job) submitTask(manifest LinkedModelManifest, eventIndex int, logLevel string) error {

	//depends on cloud-resources//
	offset := eventIndex - job.EventStartIndex
	dependencies, err := job.Dag.Dependencies(manifest, offset)
	if err != nil {
		return err
	} else {
		fmt.Println(dependencies)
	}

	payloadPath := job.generatePayloadPath(eventIndex, manifest)
	if logLevel == "Info" {
		fmt.Println(payloadPath)
	}

	//submit to batch.
	//@TODO: replace with call to batch
	batchJobArn := "Placeholder for Batch response"

	//set job arn
	resources, ok := job.Dag.Resources[manifest.ManifestID]

	if ok {
		resources.JobARN = append(resources.JobARN, &batchJobArn)
		job.Dag.Resources[manifest.ManifestID] = resources
	} else {
		return errors.New("task for " + manifest.Plugin.Name)
	}
	return nil
}
