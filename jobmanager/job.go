package jobmanager

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/usace/wat-go-sdk/plugindatamodel"
)

//Job
type Job struct {
	Id                string                                `json:"job_identifier" yaml:"job_identifier"`
	EventStartIndex   int                                   `json:"event_start_index" yaml:"event_start_index"`
	EventEndIndex     int                                   `json:"event_end_index" yaml:"event_end_index"`
	Models            []plugindatamodel.ModelIdentifier     `json:"models" yaml:"models"`
	LinkedManifests   []plugindatamodel.LinkedModelManifest `json:"linked_manifests" yaml:"linked_manifests"`
	OutputDestination plugindatamodel.ResourceInfo          `json:"output_destination" yaml:"output_destination"`
	resources         []ProvisionedResources
}

//provisionresources
func (job *Job) ProvisionResources() error {
	//make sure job arn list is provisioned for the total number of events to be computed.
	//depends on cloud-resources//
	job.resources = make([]ProvisionedResources, len(job.LinkedManifests))
	return errors.New("resources!!!")
}

//destructresources
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
		for _, n := range job.LinkedManifests {
			fmt.Println(n.ImageAndTag, outputDestinationPath)
			payload, err := job.generatePayload(n, eventIndex)
			if err != nil {
				return err
			}
			fmt.Println(payload)
			/*bytes, err := yaml.Marshal(payload)
			if err != nil {
				panic(err)
			}*/
			//put payload in s3
			path := outputDestinationPath + n.Plugin.Name + "_payload.yml"
			fmt.Println("putting object in fs:", path)
			//_, err = fs.PutObject(path, bytes)
			if err != nil {
				fmt.Println("failure to push payload to filestore:", err)
				return err
			}
		}
	}
	fmt.Println("payloads!!!")
	return nil
}
func (job Job) ComputeEvent(eventIndex int) error {
	for _, n := range job.LinkedManifests {
		job.submitTask(n, eventIndex)
	}
	return errors.New(fmt.Sprintf("computing event %v", eventIndex))
}

func (job Job) submitTask(manifest plugindatamodel.LinkedModelManifest, eventIndex int) error {
	//depends on cloud-resources//
	dependencies, err := job.findDependencies(manifest, eventIndex)
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
	for _, resource := range job.resources {
		if resource.LinkedManifest.ManifestID == manifest.ManifestID {
			offsetIndex := eventIndex - job.EventStartIndex //incase we start at a non zero index..
			resource.JobARN[offsetIndex] = &batchjobarn
			break
		}
	}
	return errors.New("task for " + manifest.Plugin.Name)
}
func (job Job) generatePayload(lm plugindatamodel.LinkedModelManifest, eventindex int) (plugindatamodel.ModelPayload, error) {
	payload := plugindatamodel.ModelPayload{}
	payload.EventIndex = eventindex
	payload.Id = uuid.New().String()
	for _, input := range lm.Inputs {
		foundMatch := false
		for _, linkedManifest := range job.LinkedManifests {
			for _, output := range linkedManifest.Outputs {
				if input.SourceDataId == output.Id {
					//yay we found a match
					resourcedInput := plugindatamodel.ResourcedFileData{
						//Id:       uuid.New().String(),
						FileName: input.FileName,
						ResourceInfo: plugindatamodel.ResourceInfo{
							Store: job.OutputDestination.Store,
							Root:  job.OutputDestination.Root,
							Path:  fmt.Sprintf("%vevent_%v/%v", job.OutputDestination.Path, eventindex, output.FileName),
						},
						InternalPaths: []plugindatamodel.ResourcedInternalPathData{},
					}
					//check if there are internal file paths
					if len(input.InternalPaths) > 0 {
						panic("oh no... do something fancy?")
					}
					payload.Inputs = append(payload.Inputs, resourcedInput)
					foundMatch = true
					break
				}
			}
			if foundMatch {
				break
			}
		}
		if !foundMatch {
			//this will trigger on all wat job model files.
			for _, model := range job.Models {
				for _, file := range model.Files {
					if file.Id == input.SourceDataId {
						payload.Inputs = append(payload.Inputs, file)
						foundMatch = true
						break
					}
				}
				if foundMatch {
					break
				}
			}
			if !foundMatch {
				return payload, errors.New("failed to find a match to an input dependency")
			}
		}
	}
	//@TODO set output destinations!!
	for _, output := range lm.Outputs {
		resourcedOutput := plugindatamodel.ResourcedFileData{
			//Id:       uuid.New().String(),
			FileName: output.FileName,
			ResourceInfo: plugindatamodel.ResourceInfo{
				Store: job.OutputDestination.Store,
				Root:  job.OutputDestination.Root,
				Path:  fmt.Sprintf("%vevent_%v/%v", job.OutputDestination.Path, eventindex, output.FileName),
			},
			InternalPaths: []plugindatamodel.ResourcedInternalPathData{},
		}
		payload.Outputs = append(payload.Outputs, resourcedOutput)
	}
	return payload, nil
}
func (job Job) findDependencies(lm plugindatamodel.LinkedModelManifest, eventindex int) ([]*string, error) {
	dependencies := make([]*string, 0)
	offsetIndex := eventindex - job.EventStartIndex
	for _, input := range lm.Inputs {
		foundMatch := false
		for _, provisionedresource := range job.resources {
			for _, outputs := range provisionedresource.LinkedManifest.Outputs {
				if input.Id == outputs.Id {
					//yay we found a match
					dependencies = append(dependencies, provisionedresource.JobARN[offsetIndex])
					foundMatch = true
					break
				}
			}
			if foundMatch {
				break
			}
		}
		if !foundMatch {
			//this will trigger on all wat job model files.
			for _, model := range job.Models {
				for _, file := range model.Files {
					if file.Id == input.Id {
						//no dependency to add here. but we did find the match.
						foundMatch = true
						break
					}
				}
				if foundMatch {
					break
				}
			}
			if !foundMatch {
				return dependencies, errors.New("failed to find a match to an input dependency")
			}
		}
	}
	//deduplicate multiple arn references

	return dependencies, nil
}
