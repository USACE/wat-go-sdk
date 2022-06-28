package jobmanager

import (
	"errors"
	"fmt"

	"github.com/USACE/filestore"
	"github.com/usace/wat-go-sdk/plugindatamodel"
	"gopkg.in/yaml.v3"
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
func (job Job) ProvisionResources() error {

	//depends on cloud-resources//
	return errors.New("resources!!!")
}

//destructresources
func (job Job) DestructResources() error {

	//depends on cloud-resources//
	return errors.New("ka-blewy!!!")
}

//GeneratePayloads
func (job Job) GeneratePayloads(fs filestore.FileStore) error {
	//generate payloads for all events up front?
	for i := job.EventStartIndex; i < job.EventStartIndex; i++ {
		//write out payloads to filestore. How do i get a handle on filestore from here?
		outputDestinationPath := fmt.Sprintf("%v/%v%v/", job.OutputDestination.Path, "event_", i)
		for _, n := range job.LinkedManifests {
			fmt.Println(n.ImageAndTag, outputDestinationPath)
			payload := plugindatamodel.ModelPayload{}
			bytes, err := yaml.Marshal(payload)
			if err != nil {
				panic(err)
			}
			//put payload in s3
			path := outputDestinationPath + n.Plugin.Name + "_payload.yml"
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
func (job Job) ComputeEvent(eventNumber int) error {
	for _, n := range job.LinkedManifests {
		submitTask(job.resources, n)
	}
	return errors.New(fmt.Sprintf("computing event %v", eventNumber))
}

func submitTask(resources []ProvisionedResources, manifest plugindatamodel.LinkedModelManifest) error {
	//depends on cloud-resources//
	//submit to batch.
	return errors.New("task for " + manifest.Plugin.Name)
}
func (job Job) generatePayload(lm plugindatamodel.LinkedModelManifest, eventindex int) (plugindatamodel.ModelPayload, error) {
	payload := plugindatamodel.ModelPayload{}
	payload.EventIndex = eventindex
	payload.Id = 1 //"make a uuid"
	for _, input := range lm.Inputs {
		foundMatch := false
		for _, provisionedresource := range job.resources {
			for _, output := range provisionedresource.LinkedManifest.Outputs {
				if input.Id == output.Id {
					//yay we found a match
					resourcedInput := plugindatamodel.ResourcedFileData{
						Id:       "", //make a uuid
						FileName: input.FileName,
						ResourceInfo: plugindatamodel.ResourceInfo{
							Store: job.OutputDestination.Store,
							Root:  job.OutputDestination.Root,
							Path:  job.OutputDestination.Path + output.FileName,
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
					if file.Id == input.Id {
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
			Id:       "", //make a uuid
			FileName: output.FileName,
			ResourceInfo: plugindatamodel.ResourceInfo{
				Store: job.OutputDestination.Store,
				Root:  job.OutputDestination.Root,
				Path:  job.OutputDestination.Path + output.FileName,
			},
			InternalPaths: []plugindatamodel.ResourcedInternalPathData{},
		}
		payload.Outputs = append(payload.Outputs, resourcedOutput)
	}
	return payload, nil
}
func (job Job) findDependencies(lm plugindatamodel.LinkedModelManifest) ([]*string, error) {
	dependencies := make([]*string, 0)
	for _, input := range lm.Inputs {
		foundMatch := false
		for _, provisionedresource := range job.resources {
			for _, outputs := range provisionedresource.LinkedManifest.Outputs {
				if input.Id == outputs.Id {
					//yay we found a match
					dependencies = append(dependencies, provisionedresource.JobARN)
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
	return dependencies, nil
}
