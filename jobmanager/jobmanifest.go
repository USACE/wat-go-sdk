package jobmanager

import (
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/usace/wat-go-sdk/plugindatamodel"
	"gopkg.in/yaml.v3"
)

//Job
type JobManifest struct {
	Id                      string                            `json:"job_identifier" yaml:"job_identifier"`
	EventStartIndex         int                               `json:"event_start_index" yaml:"event_start_index"`
	EventEndIndex           int                               `json:"event_end_index" yaml:"event_end_index"`
	Models                  []plugindatamodel.ModelIdentifier `json:"models" yaml:"models"`
	LinkedManifestResources []plugindatamodel.ResourceInfo    `json:"linked_manifests" yaml:"linked_manifests"`
	OutputDestination       plugindatamodel.ResourceInfo      `json:"output_destination" yaml:"output_destination"`
}

func (jm JobManifest) ConvertToJob() (Job, error) {
	linkedManifests := make([]plugindatamodel.LinkedModelManifest, len(jm.LinkedManifestResources))
	for idx, resourceInfo := range jm.LinkedManifestResources {
		lm := plugindatamodel.LinkedModelManifest{}
		file, err := os.Open(resourceInfo.Path) //replace with filestore? injected?
		if err != nil {
			panic(err)
		}
		defer file.Close()
		b, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(b, lm)
		if err != nil {
			panic(err)
		}
		linkedManifests[idx] = lm
	}
	job := Job{
		Id:                uuid.New().String(), //make a uuid
		EventStartIndex:   jm.EventStartIndex,
		EventEndIndex:     jm.EventEndIndex,
		Models:            jm.Models,
		LinkedManifests:   linkedManifests,
		OutputDestination: jm.OutputDestination,
	}
	return job, nil
}
