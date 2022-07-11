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
	Id                      string                        `json:"job_identifier" yaml:"job_identifier"`
	EventStartIndex         int                           `json:"event_start_index" yaml:"event_start_index"`
	EventEndIndex           int                           `json:"event_end_index" yaml:"event_end_index"`
	Models                  []plugin.ModelIdentifier      `json:"models" yaml:"models"`
	LinkedManifestResources []plugin.ResourceInfo         `json:"linked_manifests" yaml:"linked_manifests"`
	ComputeResources        []ComputeResourceRequirements `json:"resource_requirements" yaml:"resource_requirements"`
	OutputDestination       plugin.ResourceInfo           `json:"output_destination" yaml:"output_destination"`
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

func (jm JobManager) ProcessJob() error {
	err := jm.job.ProvisionResources()
	if err != nil {
		return err
	}

	// add in defer and recover
	defer func() {
		if r := recover(); r != nil {
			plugin.Log(plugin.Message{
				Message: fmt.Sprintf("Recovered %v\nTearing Down Resources", r),
				Level:   plugin.ERROR,
				Sender:  jm.job.Id,
			})
			err = jm.job.DestructResources()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	err = jm.job.GeneratePayloads() //jm.store
	if err != nil {
		fmt.Println(err)
		return err
	}

	//create error channel.
	for i := jm.job.EventStartIndex; i < jm.job.EventEndIndex; i++ {
		go func(index int) {
			err = jm.job.ComputeEvent(index)
			if err != nil {
				plugin.Log(plugin.Message{
					Message: fmt.Sprintf("%v\n", err),
					Level:   plugin.ERROR,
					Sender:  jm.job.Id,
				})
			}
			plugin.Log(plugin.Message{
				Message: fmt.Sprintf("Computed %v\n", index),
				Level:   plugin.INFO,
				Sender:  jm.job.Id,
			})
		}(i)

	}
	//need a wait group or a buffer channel to stall the destruction until we finish the jobs
	err = jm.job.DestructResources()
	if err != nil {
		plugin.Log(plugin.Message{
			Message: fmt.Sprintf("%v\n", err),
			Level:   plugin.ERROR,
			Sender:  jm.job.Id,
		})
		return err
	}
	plugin.Log(plugin.Message{
		Message: fmt.Sprint("\nJob Processed!\n\n"),
		Level:   plugin.INFO,
		Sender:  jm.job.Id,
	})
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
	plugin.Log(plugin.Message{
		Message: "Placeholder: Deallocate / Deregister / Destroy resources",
		Level:   plugin.INFO,
		Sender:  job.Id,
	})
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
			plugin.Log(plugin.Message{
				Message: fmt.Sprint(fmt.Sprintln("\n", n.ImageAndTag), fmt.Sprintln("\t", outputDestinationPath)),
				Level:   plugin.INFO,
				Sender:  job.Id,
			})
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
	// put payload in s3
	path := job.generatePayloadPath(eventIndex, modelManifest)

	outputResourceInfo := plugin.ResourceInfo{
		Store: job.OutputDestination.Store,
		Root:  job.OutputDestination.Root,
		Path:  path,
	}
	plugin.UpLoadFile(outputResourceInfo, bytes)

	if err != nil {
		plugin.Log(plugin.Message{
			Message: fmt.Sprintf("failure to push payload to filestore: %v\n", err),
			Level:   plugin.ERROR,
			Sender:  job.Id,
		})

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
	// fmt.Println("Payloads Generated")
	return nil
}

//ComputeEvent
func (job Job) ComputeEvent(eventIndex int) error {
	for _, n := range job.Dag.LinkedManifests {
		job.submitTask(n, eventIndex)
	}
	plugin.Log(plugin.Message{
		Message: fmt.Sprintf("computing event %v\n", eventIndex),
		Level:   plugin.INFO,
		Sender:  job.Id,
	})
	return nil
}

func (job *Job) submitTask(manifest LinkedModelManifest, eventIndex int) error {

	//depends on cloud-resources//
	offset := eventIndex - job.EventStartIndex
	dependencies, err := job.Dag.Dependencies(manifest, offset)
	if err != nil {
		return err
	} else {
		plugin.Log(plugin.Message{
			Message: fmt.Sprint(dependencies),
			Level:   plugin.INFO,
			Sender:  job.Id,
		})
	}

	payloadPath := job.generatePayloadPath(eventIndex, manifest)
	plugin.Log(plugin.Message{
		Message: payloadPath,
		Level:   plugin.INFO,
		Sender:  job.Id,
	})
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
