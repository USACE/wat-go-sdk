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
	file, err := os.Open("../exampledata/example_payload.yaml")
	if err != nil {
		t.Fail()
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fail()
	}
	payload := plugindatamodel.Payload{}
	err = yaml.Unmarshal(b, &payload)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	log.Println(string(b))
	log.Println(payload)
}

func TestWritePayload(t *testing.T) {
	inputs := make([]plugindatamodel.ResourcedData, 2)
	inputs[0] = plugindatamodel.ResourcedData{
		Data: plugindatamodel.Data{
			Name:      "Muncie.p04.tmp.hdf",
			Parameter: "temporary hdf file",
		},
		ResourceInfo: plugindatamodel.ResourceInfo{
			Scheme:    "s3",
			Authority: "cloud-wat-dev",
			Path:      "/data/models/Muncie",
		},
	}
	inputs[1] = plugindatamodel.ResourcedData{
		Data: plugindatamodel.Data{
			Name:      "/Event Conditions/ White  Reach: Muncie  RS: 15696.24",
			Parameter: "flow time-series table",
		},
		ResourceInfo: plugindatamodel.ResourceInfo{
			Scheme:    "s3",
			Authority: "cloud-wat-dev",
			Path:      "/runs/realization_1/event_1/Muncie_RS_15696.24.csv",
		},
	}
	outputs := make([]plugindatamodel.Data, 2)
	outputs[0] = plugindatamodel.Data{
		Name:      "Muncie.p04.tmp.hdf",
		Parameter: "temporary hdf file",
	}
	outputs[1] = plugindatamodel.Data{
		Name:      "Muncie.log",
		Parameter: "log file",
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
