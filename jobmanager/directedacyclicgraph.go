package jobmanager

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/usace/wat-go-sdk/plugindatamodel"
)

type LinkedManifestStack []LinkedModelManifest

func (lms *LinkedManifestStack) Push(lm LinkedModelManifest) {
	*lms = append(*lms, lm)
}
func (lms *LinkedManifestStack) Pop() (LinkedModelManifest, error) {
	if len(*lms) == 0 {
		return LinkedModelManifest{}, errors.New("no more elements in the stack")
	}
	id := len(*lms) - 1
	lm := (*lms)[id]
	*lms = (*lms)[:id]
	return lm, nil
}

type DirectedAcyclicGraph struct {
	Models          []plugindatamodel.ModelIdentifier `json:"models" yaml:"models"`
	LinkedManifests []LinkedModelManifest             `json:"linked_manifests" yaml:"linked_manifests"`
	Resources       map[string]ProvisionedResources   `json:"provisioned_resources" yaml:"provisioned_resources"`
}

func (dag DirectedAcyclicGraph) TopologicallySort() ([]LinkedModelManifest, error) {
	//Kahn's Algorithm
	S := LinkedManifestStack{} //set of linked manifests with no upstream dependencies
	L := LinkedManifestStack{}
	for _, lm := range dag.LinkedManifests {
		noDependencies := true
		for _, input := range lm.Inputs {
			_, ok := dag.producesFile(input)
			if ok {
				noDependencies = false
			} else {
				if len(input.InternalPaths) > 0 {
					for _, ip := range input.InternalPaths {
						_, _, ipok := dag.producesInternalPath(ip)
						if ipok {
							noDependencies = false
						}
					}
				}
			}
		}
		if noDependencies {
			S.Push(lm)
		}
	}
	if len(S) == 0 {
		return S, errors.New("cannot have a DAG without at least one start node")
	}
	for len(S) > 0 {
		n, err := S.Pop()
		if err != nil {
			return S, err
		}
		L.Push(n)
		for _, m := range dag.LinkedManifests {
			noOtherDependencies := true
			for _, input := range m.Inputs {
				_, dagok := dag.producesFile(input)
				if dagok {
					inL := false
					//should i check for anything in L?
					//_, ok := n.producesFile(input.SourceDataId)
					for _, Ln := range L {
						_, ok := Ln.producesFile(input.SourceDataId)
						if ok {
							inL = true
						}
					}
					if !inL {
						noOtherDependencies = false
					}
				} else {
					if len(input.InternalPaths) > 0 {
						for _, ip := range input.InternalPaths {
							_, _, ipok := dag.producesInternalPath(ip)
							if ipok {
								inL := false
								//should i check for anything in L?
								for _, Ln := range L {
									_, _, ok := Ln.producesInternalPath(ip)
									if ok {
										inL = true
									}
								}
								if !inL {
									noOtherDependencies = false
								}
							}
						}
					}
				}
			}
			if noOtherDependencies {
				visited := false
				for _, vistedNode := range L {
					if m.ManifestID == vistedNode.ManifestID {
						visited = true
					}
				}
				for _, addedToS := range S {
					if m.ManifestID == addedToS.ManifestID {
						visited = true // added to s but not yet popped off the stack
					}
				}
				if !visited {
					S.Push(m)
				}
			}
		}
	}
	if len(L) != len(dag.LinkedManifests) {
		return L, errors.New("something went wrong in the sorting.")
	}
	return L, nil
}
func (dag DirectedAcyclicGraph) Dependencies(manifestUUID string, eventIndex int) ([]*string, error) {
	//get the dependencies for a given manifestUUID and eventIndex.
	//get the linked manifest for a given manifestUUID
	linkedManifest := LinkedModelManifest{}
	for _, lm := range dag.LinkedManifests {
		if lm.ManifestID == manifestUUID {
			linkedManifest = lm
			break
		}
	}
	return dag.findDependencies(linkedManifest, eventIndex)
}
func (dag DirectedAcyclicGraph) findDependencies(lm LinkedModelManifest, eventIndex int) ([]*string, error) {
	dependencies := make([]*string, 0)
	for _, input := range lm.Inputs {
		foundMatch := false
		for _, inputmanifest := range dag.LinkedManifests {
			if lm.ManifestID == inputmanifest.ManifestID {
				break
			}
			//lm := LinkedModelManifest{} //get the right linkedmodel manifest based on the resource linkmanifest id.
			for _, outputs := range inputmanifest.Outputs {

				if input.Id == outputs.Id {
					//yay we found a match
					resources, ok := dag.Resources[inputmanifest.ManifestID]
					if ok {
						dependencies = append(dependencies, resources.JobARN[eventIndex])
						foundMatch = true
						break
					} else {
						return dependencies, errors.New("resources not provisioned for this input")
					}
				}
			}
			if foundMatch {
				break
			}
		}
		if !foundMatch {
			//this will trigger on all wat job model files.
			for _, model := range dag.Models {
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

func (dag DirectedAcyclicGraph) GeneratePayload(lm LinkedModelManifest, eventindex int, outputDestination plugindatamodel.ResourceInfo) (plugindatamodel.ModelPayload, error) {
	payload := plugindatamodel.ModelPayload{}
	payload.EventIndex = eventindex
	payload.Id = uuid.NewSHA1(uuid.MustParse(lm.ManifestID), []byte(fmt.Sprintf("event%v", eventindex))).String()
	//set inputs
	for _, input := range lm.Inputs {
		//try to link to other manifests first.
		resourcedInput, err := dag.linkToPluginOutput(input, eventindex, outputDestination)
		if err != nil {
			//if links were not satisfied, link to model data defined in job manifest
			file, err := dag.linkToModelData(input, eventindex, outputDestination)
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
	outputs, err := dag.setPayloadOutputDestinations(lm, eventindex, outputDestination)
	if err != nil {
		return payload, errors.New("could not set outputs")
	}
	payload.Outputs = outputs
	return payload, nil
}
func (dag DirectedAcyclicGraph) linkToPluginOutput(linkedFile LinkedFileData, eventIndex int, outputDestination plugindatamodel.ResourceInfo) (plugindatamodel.ResourcedFileData, error) {
	resourcedInput := plugindatamodel.ResourcedFileData{
		FileName: linkedFile.FileName,
		ResourceInfo: plugindatamodel.ResourceInfo{
			Store: outputDestination.Store,
			Root:  outputDestination.Root,
		},
		InternalPaths: []plugindatamodel.ResourcedInternalPathData{},
	}
	output, ok := dag.producesFile(linkedFile)
	if ok {
		resourcedInput.Path = fmt.Sprintf("%vevent_%v/%v", outputDestination.Path, eventIndex, output.FileName)
		//check if there are internal file paths
		if linkedFile.HasInternalPaths() {
			//link internal paths if there are any.
			resourcedInternalPaths, err := dag.linkInternalPaths(linkedFile, eventIndex, outputDestination)
			if err != nil {
				return resourcedInput, err
			}
			resourcedInput.InternalPaths = resourcedInternalPaths
		}
		return resourcedInput, nil
	}
	return resourcedInput, errors.New("no link found")
}
func (dag DirectedAcyclicGraph) linkInternalPaths(linkedFile LinkedFileData, eventIndex int, outputDestination plugindatamodel.ResourceInfo) ([]plugindatamodel.ResourcedInternalPathData, error) {
	internalPaths := make([]plugindatamodel.ResourcedInternalPathData, len(linkedFile.InternalPaths))
	//currently not checking if a link is unsatisfied. it might be smart to error out if len(linkedFile.InternalPaths)!=numsuccessfullinks
	for idx, internalPath := range linkedFile.InternalPaths {
		internalpathid, outputFileName, ok := dag.producesInternalPath(internalPath)
		if ok {
			resourcedInput := plugindatamodel.ResourcedInternalPathData{
				PathName:     internalPath.PathName,
				FileName:     outputFileName,
				InternalPath: internalpathid,
				ResourceInfo: plugindatamodel.ResourceInfo{
					Store: outputDestination.Store,
					Root:  outputDestination.Root,
					Path:  fmt.Sprintf("%vevent_%v/%v", outputDestination.Path, eventIndex, outputFileName),
				},
			}
			internalPaths[idx] = resourcedInput
		}
	}
	return internalPaths, nil
}
func (dag DirectedAcyclicGraph) linkToModelData(linkedFile LinkedFileData, eventIndex int, outputDestination plugindatamodel.ResourceInfo) (plugindatamodel.ResourcedFileData, error) {
	returnFile := plugindatamodel.ResourcedFileData{}
	file, ok := dag.producesModelFile(linkedFile)
	if ok {
		//check if there are internal file paths
		returnFile.Id = file.Id
		returnFile.FileName = file.FileName
		returnFile.ResourceInfo = file.ResourceInfo
		fmt.Printf("there are %v internal paths on input %v\n", len(linkedFile.InternalPaths), linkedFile.SourceDataId)
		if len(linkedFile.InternalPaths) > 0 {
			resourcedInternalPaths, err := dag.linkInternalPaths(linkedFile, eventIndex, outputDestination)
			if err != nil {
				return returnFile, err
			}
			returnFile.InternalPaths = resourcedInternalPaths
		}
		return returnFile, nil
	}
	return returnFile, errors.New("could not find a match")
}
func (dag DirectedAcyclicGraph) setPayloadOutputDestinations(linkedManifest LinkedModelManifest, eventIndex int, outputDestination plugindatamodel.ResourceInfo) ([]plugindatamodel.ResourcedFileData, error) {
	outputs := make([]plugindatamodel.ResourcedFileData, len(linkedManifest.Outputs))
	for idx, output := range linkedManifest.Outputs {
		resourcedOutput := plugindatamodel.ResourcedFileData{
			FileName: output.FileName,
			ResourceInfo: plugindatamodel.ResourceInfo{
				Store: outputDestination.Store,
				Root:  outputDestination.Root,
				Path:  fmt.Sprintf("%vevent_%v/%v", outputDestination.Path, eventIndex, output.FileName),
			},
			InternalPaths: []plugindatamodel.ResourcedInternalPathData{},
		}
		outputs[idx] = resourcedOutput
	}
	return outputs, nil
}
func (dag DirectedAcyclicGraph) producesInternalPath(internalpath LinkedInternalPathData) (string, string, bool) {
	for _, lm := range dag.LinkedManifests {
		ip, fn, ok := lm.producesInternalPath(internalpath)
		if ok {
			return ip, fn, ok
		}
	}
	return "", "", false
}
func (dag DirectedAcyclicGraph) producesFile(linkedFile LinkedFileData) (plugindatamodel.FileData, bool) {
	for _, lm := range dag.LinkedManifests {
		f, ok := lm.producesFile(linkedFile.SourceDataId)
		if ok {
			return f, ok
		}
	}
	return plugindatamodel.FileData{}, false
}
func (dag DirectedAcyclicGraph) producesModelFile(linkedFile LinkedFileData) (plugindatamodel.ResourcedFileData, bool) {
	for _, model := range dag.Models {
		for _, modelFile := range model.Files {
			if modelFile.Id == linkedFile.SourceDataId {
				return modelFile, true
			}
		}
	}
	return plugindatamodel.ResourcedFileData{}, false
}
