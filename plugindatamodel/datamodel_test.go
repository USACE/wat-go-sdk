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
	path := "../exampledata/payload_update.yaml"
	payload := plugindatamodel.ModelPayload{}
	readobject(t, path, &payload)
}

func TestReadManifest(t *testing.T) {
	path := "../exampledata/manifest_update.yaml"
	manifest := plugindatamodel.ModelManifest{}
	readobject(t, path, &manifest)
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
	err = yaml.Unmarshal(b, object)
	if err != nil {
		log.Println(err)
		t.Fail()
	} else {
		log.Println(string(b))
		log.Println()
		log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		log.Println()
		b2, err := yaml.Marshal(object)
		if err != nil {
			log.Println(err)
			t.Fail()
		}
		log.Println(string(b2))
	}

}

func TestWritePayload(t *testing.T) {
	inputs := make([]plugindatamodel.ResourcedFileData, 2)
	inputs[0] = plugindatamodel.ResourcedFileData{
		FileName: "Muncie.p04.tmp.hdf",
		ResourceInfo: plugindatamodel.ResourceInfo{
			Store: "s3",
			Root:  "cloud-wat-dev",
			Path:  "/data/models/Muncie",
		},
	}
	inputs[1] = plugindatamodel.ResourcedFileData{
		FileName: "/Event Conditions/ White  Reach: Muncie  RS: 15696.24",
		ResourceInfo: plugindatamodel.ResourceInfo{
			Store: "s3",
			Root:  "cloud-wat-dev",
			Path:  "/runs/realization_1/event_1/Muncie_RS_15696.24.csv",
		},
	}
	outputs := make([]plugindatamodel.ResourcedFileData, 2)
	outputs[0] = plugindatamodel.ResourcedFileData{
		FileName: "Muncie.p04.tmp.hdf",
		ResourceInfo: plugindatamodel.ResourceInfo{
			Store: "s3",
			Root:  "cloud-wat-dev",
			Path:  "/runs/realization_1/event_1",
		},
	}
	outputs[1] = plugindatamodel.ResourcedFileData{
		FileName: "Muncie.log",
		ResourceInfo: plugindatamodel.ResourceInfo{
			Store: "s3",
			Root:  "cloud-wat-dev",
			Path:  "/runs/realization_1/event_1",
		},
	}
	payload := plugindatamodel.ModelPayload{
		/*ModelIdentifier: plugindatamodel.ModelIdentifier{
			Name:        "Muncie",
			Alternative: ".p04",
		},*/
		Inputs:  inputs,
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
