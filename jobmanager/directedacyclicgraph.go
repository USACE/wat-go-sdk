package jobmanager

import (
	"errors"

	"github.com/usace/wat-go-sdk/plugindatamodel"
)

type DirectedAcyclicGraph struct {
	Models          []plugindatamodel.ModelIdentifier `json:"models" yaml:"models"`
	LinkedManifests []LinkedModelManifest             `json:"linked_manifests" yaml:"linked_manifests"`
	Resources       map[string]ProvisionedResources   `json:"provisioned_resources" yaml:"provisioned_resources"`
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
