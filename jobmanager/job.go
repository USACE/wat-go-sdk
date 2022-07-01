package jobmanager

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/usace/wat-go-sdk/plugindatamodel"
	"gopkg.in/yaml.v3"
)

//Job
type Job struct {
	Id                string                            `json:"job_identifier" yaml:"job_identifier"`
	EventStartIndex   int                               `json:"event_start_index" yaml:"event_start_index"`
	EventEndIndex     int                               `json:"event_end_index" yaml:"event_end_index"`
	Models            []plugindatamodel.ModelIdentifier `json:"models" yaml:"models"`
	LinkedManifests   []LinkedModelManifest             `json:"linked_manifests" yaml:"linked_manifests"`
	OutputDestination plugindatamodel.ResourceInfo      `json:"output_destination" yaml:"output_destination"`
	resources         []ProvisionedResources
}

//ProvisionResources
func (job *Job) ProvisionResources() error {
	//make sure job arn list is provisioned for the total number of events to be computed.
	//depends on cloud-resources//
	job.resources = make([]ProvisionedResources, len(job.LinkedManifests))
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
		for _, n := range job.LinkedManifests {
			fmt.Println(n.ImageAndTag, outputDestinationPath)
			payload, err := job.generatePayload(n, eventIndex)
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
	for _, n := range job.LinkedManifests {
		job.submitTask(n, eventIndex)
	}
	fmt.Println(fmt.Sprintf("computing event %v", eventIndex))
	return nil
}
func (job Job) submitTask(manifest LinkedModelManifest, eventIndex int) error {
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
		if resource.LinkedManifestID == manifest.ManifestID {
			offsetIndex := eventIndex - job.EventStartIndex //incase we start at a non zero index..
			resource.JobARN[offsetIndex] = &batchjobarn
			break
		}
	}
	return errors.New("task for " + manifest.Plugin.Name)
}
func (job Job) generatePayload(lm LinkedModelManifest, eventindex int) (plugindatamodel.ModelPayload, error) {
	payload := plugindatamodel.ModelPayload{}
	payload.EventIndex = eventindex
	payload.Id = uuid.NewSHA1(uuid.MustParse(lm.ManifestID), []byte(fmt.Sprintf("event%v", eventindex))).String()
	//set inputs
	for _, input := range lm.Inputs {
		//try to link to other manifests first.
		resourcedInput, err := job.linkToPluginOutput(input, eventindex)
		if err != nil {
			//if links were not satisfied, link to model data defined in job manifest
			file, err := job.linkToModelData(input, eventindex)
			if err != nil {
				//if link not found, fail out.
				return payload, err
			}
			payload.Inputs = append(payload.Inputs, file)
		} else {
			payload.Inputs = append(payload.Inputs, resourcedInput)
		}
	}
	//set output destinations
	outputs, err := job.setPayloadOutputDestinations(lm, eventindex)
	if err != nil {
		return payload, errors.New("could not set outputs")
	}
	payload.Outputs = outputs
	return payload, nil
}
func (job Job) linkToPluginOutput(linkedFile LinkedFileData, eventIndex int) (plugindatamodel.ResourcedFileData, error) {
	resourcedInput := plugindatamodel.ResourcedFileData{
		FileName: linkedFile.FileName,
		ResourceInfo: plugindatamodel.ResourceInfo{
			Store: job.OutputDestination.Store,
			Root:  job.OutputDestination.Root,
		},
		InternalPaths: []plugindatamodel.ResourcedInternalPathData{},
	}
	for _, linkedManifest := range job.LinkedManifests {
		for _, output := range linkedManifest.Outputs {
			if linkedFile.SourceDataId == output.Id {
				//yay we found a match
				resourcedInput.Path = fmt.Sprintf("%vevent_%v/%v", job.OutputDestination.Path, eventIndex, output.FileName)
				//check if there are internal file paths
				if len(linkedFile.InternalPaths) > 0 {
					//link internal paths if there are any.
					resourcedInternalPaths, err := job.linkInternalPaths(linkedFile, eventIndex)
					if err != nil {
						return resourcedInput, err
					}
					resourcedInput.InternalPaths = resourcedInternalPaths
				}
				return resourcedInput, nil
			}
		}
	}
	return resourcedInput, errors.New("no link found")
}
func (job Job) linkInternalPaths(linkedFile LinkedFileData, eventIndex int) ([]plugindatamodel.ResourcedInternalPathData, error) {
	internalPaths := make([]plugindatamodel.ResourcedInternalPathData, len(linkedFile.InternalPaths))
	//currently not checking if a link is unsatisfied. it might be smart to error out if len(linkedFile.InternalPaths)!=numsuccessfullinks
	for idx, internalPath := range linkedFile.InternalPaths {
		for _, linkedManifest := range job.LinkedManifests {
			for _, output := range linkedManifest.Outputs {
				if internalPath.SourceFileID == output.Id {
					//yay we found a match
					ip := ""
					if len(output.InternalPaths) > 0 {
						for _, internalpath := range output.InternalPaths {
							if internalpath.Id == internalPath.SourcePathID {
								ip = internalpath.PathName
							}
						}
					}
					resourcedInput := plugindatamodel.ResourcedInternalPathData{
						PathName:     internalPath.PathName,
						FileName:     output.FileName,
						InternalPath: ip,
						ResourceInfo: plugindatamodel.ResourceInfo{
							Store: job.OutputDestination.Store,
							Root:  job.OutputDestination.Root,
							Path:  fmt.Sprintf("%vevent_%v/%v", job.OutputDestination.Path, eventIndex, output.FileName),
						},
					}
					internalPaths[idx] = resourcedInput
				}
			}
		}
	}
	return internalPaths, nil
}
func (job Job) linkToModelData(linkedFile LinkedFileData, eventIndex int) (plugindatamodel.ResourcedFileData, error) {
	returnFile := plugindatamodel.ResourcedFileData{}
	for _, model := range job.Models {
		for _, file := range model.Files {
			if file.Id == linkedFile.SourceDataId {
				//check if there are internal file paths
				returnFile.Id = file.Id
				returnFile.FileName = file.FileName
				returnFile.ResourceInfo = file.ResourceInfo
				fmt.Printf("there are %v internal paths on input %v\n", len(linkedFile.InternalPaths), linkedFile.SourceDataId)
				if len(linkedFile.InternalPaths) > 0 {
					resourcedInternalPaths, err := job.linkInternalPaths(linkedFile, eventIndex)
					if err != nil {
						return returnFile, err
					}
					returnFile.InternalPaths = resourcedInternalPaths
				}
				return returnFile, nil
			}
		}
	}
	return returnFile, errors.New("could not find a match")
}

func (job Job) setPayloadOutputDestinations(linkedManifest LinkedModelManifest, eventIndex int) ([]plugindatamodel.ResourcedFileData, error) {
	outputs := make([]plugindatamodel.ResourcedFileData, len(linkedManifest.Outputs))
	for idx, output := range linkedManifest.Outputs {
		resourcedOutput := plugindatamodel.ResourcedFileData{
			FileName: output.FileName,
			ResourceInfo: plugindatamodel.ResourceInfo{
				Store: job.OutputDestination.Store,
				Root:  job.OutputDestination.Root,
				Path:  fmt.Sprintf("%vevent_%v/%v", job.OutputDestination.Path, eventIndex, output.FileName),
			},
			InternalPaths: []plugindatamodel.ResourcedInternalPathData{},
		}
		outputs[idx] = resourcedOutput
	}
	return outputs, nil
}
func (job Job) findDependencies(lm LinkedModelManifest, eventindex int) ([]*string, error) {
	dependencies := make([]*string, 0)
	offsetIndex := eventindex - job.EventStartIndex
	for _, input := range lm.Inputs {
		foundMatch := false
		for _, provisionedresource := range job.resources {
			lm := LinkedModelManifest{} //get the right linkedmodel manifest based on the resource linkmanifest id.
			for _, outputs := range lm.Outputs {

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
