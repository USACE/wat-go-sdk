package plugindatamodel_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/usace/wat-go-sdk/plugindatamodel"
	"gopkg.in/yaml.v3"
)

func TestReadPayload(t *testing.T) {
	path := "../exampledata/example_payload.yaml"
	payload := plugindatamodel.Payload{}
	readobject(t, path, payload)
}
func TestReadLinkedManifest(t *testing.T) {
	path := "../exampledata/example_linked_manifest.yaml"
	linkedmanifest := plugindatamodel.LinkedModelManifest{}
	readobject(t, path, linkedmanifest)
}
func TestReadManifest(t *testing.T) {
	path := "../exampledata/example_manifest.yaml"
	manifest := plugindatamodel.ModelManifest{}
	readobject(t, path, manifest)
}
func readobject(t *testing.T, path string, object interface{}) {
	file, err := os.Open(path)
	if err != nil {
		t.Fail()
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fail()
	}
	err = yaml.Unmarshal(b, &object)
	if err != nil {
		log.Println(err)
		t.Fail()
	} else {
		log.Println(string(b))
		log.Println(object)
	}

}

func TestWritePayload(t *testing.T) {
	inputs := make([]plugindatamodel.ResourcedData, 2)
	inputs[0] = plugindatamodel.ResourcedData{
		Data: plugindatamodel.Data{
			FileName: "Muncie.p04.tmp.hdf",
		},
		ResourceInfo: plugindatamodel.ResourceInfo{
			Scheme:    "s3",
			Authority: "cloud-wat-dev",
			Path:      "/data/models/Muncie",
		},
	}
	inputs[1] = plugindatamodel.ResourcedData{
		Data: plugindatamodel.Data{
			FileName: "/Event Conditions/ White  Reach: Muncie  RS: 15696.24",
		},
		ResourceInfo: plugindatamodel.ResourceInfo{
			Scheme:    "s3",
			Authority: "cloud-wat-dev",
			Path:      "/runs/realization_1/event_1/Muncie_RS_15696.24.csv",
		},
	}
	outputs := make([]plugindatamodel.Data, 2)
	outputs[0] = plugindatamodel.Data{
		FileName: "Muncie.p04.tmp.hdf",
	}
	outputs[1] = plugindatamodel.Data{
		FileName: "Muncie.log",
	}
	payload := plugindatamodel.Payload{
		ModelIdentifier: plugindatamodel.ModelIdentifier{
			Name:        "Muncie",
			Alternative: ".p04",
		},
		Inputs: inputs,
		OutputDestination: plugindatamodel.ResourceInfo{
			Scheme:    "s3",
			Authority: "cloud-wat-dev",
			Path:      "/runs/realization_1/event_1",
		},
		Outputs: outputs,
	}
	b, err := yaml.Marshal(payload)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	log.Println(string(b))
	log.Println(payload)
}
