package wat

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/usace/wat-go-sdk/plugin"
)

type linkedManifestStack []LinkedModelManifest

func (lms *linkedManifestStack) Push(lm LinkedModelManifest) {
	*lms = append(*lms, lm)
}

func (lms *linkedManifestStack) Pop() (LinkedModelManifest, error) {
	if len(*lms) == 0 {
		return LinkedModelManifest{}, errors.New("no more elements in the stack")
	}
	// identify the length of the stack
	id := len(*lms) - 1

	//find the last element
	lm := (*lms)[id]

	//remove the last element from the stack
	*lms = (*lms)[:id]

	//return the last element *pop*!
	return lm, nil
}

type provisionedResources struct { //this is very aws specific. can we make this provider agnostic? should it live here?
	LinkedManifestID      string
	ComputeEnvironmentARN *string
	JobDefinitionARN      *string
	JobARN                []*string
	QueueARN              *string
}

type DirectedAcyclicGraph struct {
	Models          []plugin.ModelIdentifier        `json:"models" yaml:"models"`
	LinkedManifests []LinkedModelManifest           `json:"linked_manifests" yaml:"linked_manifests"`
	Resources       map[string]provisionedResources `json:"provisioned_resources" yaml:"provisioned_resources"`
}

func (dag DirectedAcyclicGraph) TopologicallySort() ([]LinkedModelManifest, error) {
	//Kahn's Algorithm https://en.wikipedia.org/wiki/Topological_sorting
	S := linkedManifestStack{} //set of linked manifests with no upstream dependencies
	L := linkedManifestStack{}

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
		return S, errors.New("a DAG must contain at least one node with no dependencies satisfied by other linked manifests in the DAG")
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
				for _, visitedNode := range L {
					if m.ManifestID == visitedNode.ManifestID {
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
		return L, errors.New("the DAG contains a sub-cycle") //we could identify it by listing the elements in the dag that are not present in the Stack L
	}
	return L, nil
}

func (dag DirectedAcyclicGraph) Dependencies(lm LinkedModelManifest, eventIndex int) ([]*string, error) {
	dependencies := make([]*string, 0)
	for _, input := range lm.Inputs {
		for _, inputManifest := range dag.LinkedManifests {
			if lm.ManifestID == inputManifest.ManifestID {
				break //not dependent upon self.
			}
			if inputManifest.producesDependency(input) {
				resources, ok := dag.Resources[inputManifest.ManifestID]
				if ok {
					jobArn := resources.JobARN[eventIndex]
					dependencies = append(dependencies, jobArn)
				}
			}
		}
	}

	//deduplicate multiple arn references
	uniqueDependencies := make([]*string, 0)
	for _, s := range dependencies {
		contains := false
		for _, us := range uniqueDependencies {
			if s == us {
				contains = true
				break
			}
		}
		if !contains {
			uniqueDependencies = append(uniqueDependencies, s)
		}
	}
	return uniqueDependencies, nil
}

func (dag DirectedAcyclicGraph) GeneratePayload(lm LinkedModelManifest, eventIndex int, outputDestination plugin.ResourceInfo) (plugin.ModelPayload, error) {
	payload := plugin.ModelPayload{}
	payload.Model.Name = lm.Model.Name
	payload.Model.Alternative = lm.Model.Alternative
	payload.EventIndex = eventIndex
	payload.Id = uuid.NewSHA1(uuid.MustParse(lm.ManifestID), []byte(fmt.Sprintf("event%v", eventIndex))).String()
	//set inputs
	for _, input := range lm.Inputs {
		//try to link to other manifests first.
		resourcedInput, err := dag.linkToPluginOutput(input, eventIndex, outputDestination)
		if err != nil {
			//if links were not satisfied, link to model data defined in job manifest
			file, err := dag.linkToModelData(input, eventIndex, outputDestination)
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
	outputs, err := dag.setPayloadOutputDestinations(lm, eventIndex, outputDestination)
	if err != nil {
		return payload, errors.New("could not set outputs")
	}
	payload.Outputs = outputs
	return payload, nil
}

func (dag DirectedAcyclicGraph) linkToPluginOutput(linkedFile LinkedFileData, eventIndex int, outputDestination plugin.ResourceInfo) (plugin.ResourcedFileData, error) {
	resourcedInput := plugin.ResourcedFileData{
		FileName: linkedFile.FileName,
		ResourceInfo: plugin.ResourceInfo{
			Store: outputDestination.Store,
			Root:  outputDestination.Root,
		},
		InternalPaths: []plugin.ResourcedInternalPathData{},
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

func (dag DirectedAcyclicGraph) linkInternalPaths(linkedFile LinkedFileData, eventIndex int, outputDestination plugin.ResourceInfo) ([]plugin.ResourcedInternalPathData, error) {
	internalPaths := make([]plugin.ResourcedInternalPathData, len(linkedFile.InternalPaths))
	// currently not checking if a link is unsatisfied. it might be smart to error out if len(linkedFile.InternalPaths)!=numsuccessfullinks
	expectedLinks := len(linkedFile.InternalPaths)
	foundLinks := 0
	for idx, internalPath := range linkedFile.InternalPaths {
		internalPathId, outputFileName, ok := dag.producesInternalPath(internalPath)
		if ok {
			foundLinks++
			resourcedInput := plugin.ResourcedInternalPathData{
				PathName:     internalPath.PathName,
				FileName:     outputFileName,
				InternalPath: internalPathId,
				ResourceInfo: plugin.ResourceInfo{
					Store: outputDestination.Store,
					Root:  outputDestination.Root,
					Path:  fmt.Sprintf("%vevent_%v/%v", outputDestination.Path, eventIndex, outputFileName),
				},
			}
			internalPaths[idx] = resourcedInput
		}
	}
	if expectedLinks != foundLinks {
		return internalPaths, fmt.Errorf("expected %v links, found %v links", expectedLinks, foundLinks)
	}
	return internalPaths, nil
}

func (dag DirectedAcyclicGraph) linkToModelData(linkedFile LinkedFileData, eventIndex int, outputDestination plugin.ResourceInfo) (plugin.ResourcedFileData, error) {
	returnFile := plugin.ResourcedFileData{}
	file, ok := dag.producesModelFile(linkedFile)
	if ok {
		//check if there are internal file paths
		returnFile.Id = file.Id
		returnFile.FileName = file.FileName
		returnFile.ResourceInfo = file.ResourceInfo

		plugin.Log(plugin.Message{
			Message: fmt.Sprintf("\t\t%v | internal paths = %v\n", linkedFile.SourceDataId, len(linkedFile.InternalPaths)),
			Level:   plugin.INFO,
			Sender:  "DAG Linking",
		})

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

func (dag DirectedAcyclicGraph) setPayloadOutputDestinations(linkedManifest LinkedModelManifest, eventIndex int, outputDestination plugin.ResourceInfo) ([]plugin.ResourcedFileData, error) {
	outputs := make([]plugin.ResourcedFileData, len(linkedManifest.Outputs))
	for idx, output := range linkedManifest.Outputs {
		resourcedOutput := plugin.ResourcedFileData{
			FileName: output.FileName,
			ResourceInfo: plugin.ResourceInfo{
				Store: outputDestination.Store,
				Root:  outputDestination.Root,
				Path:  fmt.Sprintf("%vevent_%v/%v", outputDestination.Path, eventIndex, output.FileName),
			},
			InternalPaths: []plugin.ResourcedInternalPathData{},
		}
		outputs[idx] = resourcedOutput
	}
	return outputs, nil
}

func (dag DirectedAcyclicGraph) producesInternalPath(internalPath LinkedInternalPathData) (string, string, bool) {
	for _, lm := range dag.LinkedManifests {
		ip, fn, ok := lm.producesInternalPath(internalPath)
		if ok {
			return ip, fn, ok
		}
	}
	return "", "", false
}

func (dag DirectedAcyclicGraph) producesFile(linkedFile LinkedFileData) (plugin.FileData, bool) {
	for _, lm := range dag.LinkedManifests {
		f, ok := lm.producesFile(linkedFile.SourceDataId)
		if ok {
			return f, ok
		}
	}
	return plugin.FileData{}, false
}

func (dag DirectedAcyclicGraph) producesModelFile(linkedFile LinkedFileData) (plugin.ResourcedFileData, bool) {
	for _, model := range dag.Models {
		for _, modelFile := range model.Files {
			if modelFile.Id == linkedFile.SourceDataId {
				return modelFile, true
			}
		}
	}
	return plugin.ResourcedFileData{}, false
}
