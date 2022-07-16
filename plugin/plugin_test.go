package plugin_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/usace/wat-go-sdk/plugin"
	"gopkg.in/yaml.v3"
)

func TestReadPayload(t *testing.T) {
	path := "../exampledata/ras-mutator_payload.yml"
	payload := plugin.ModelPayload{}
	readObject(t, path, &payload)
}

func TestReadManifest(t *testing.T) {
	path := "../exampledata/ras_mutator_manifest.yaml"
	manifest := plugin.ModelManifest{}
	readObject(t, path, &manifest)
}

func readObject(t *testing.T, path string, object interface{}) {
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
	inputs := make([]plugin.ResourcedFileData, 2)
	inputs[0] = plugin.ResourcedFileData{
		FileName: "Muncie.p04.tmp.hdf",
		ResourceInfo: plugin.ResourceInfo{
			Store: "s3",
			Root:  "cloud-wat-dev",
			Path:  "/data/models/Muncie",
		},
	}
	inputs[1] = plugin.ResourcedFileData{
		FileName: "/Event Conditions/ White  Reach: Muncie  RS: 15696.24",
		ResourceInfo: plugin.ResourceInfo{
			Store: "s3",
			Root:  "cloud-wat-dev",
			Path:  "/runs/realization_1/event_1/Muncie_RS_15696.24.csv",
		},
	}
	outputs := make([]plugin.ResourcedFileData, 2)
	outputs[0] = plugin.ResourcedFileData{
		FileName: "Muncie.p04.tmp.hdf",
		ResourceInfo: plugin.ResourceInfo{
			Store: "s3",
			Root:  "cloud-wat-dev",
			Path:  "/runs/realization_1/event_1",
		},
	}
	outputs[1] = plugin.ResourcedFileData{
		FileName: "Muncie.log",
		ResourceInfo: plugin.ResourceInfo{
			Store: "s3",
			Root:  "cloud-wat-dev",
			Path:  "/runs/realization_1/event_1",
		},
	}
	payload := plugin.ModelPayload{
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

//commenting this test out until i have a better setup with docker compose.
/*func TestLoadPayload(t *testing.T) {
	// this will relies on having s3mock and a docker compose running.
	plugin.InitConfigFromPath("../exampledata/config.json")
	p, err := plugin.LoadPayload("ras-mutator_payload.yml")
	if err != nil {
		plugin.Log(plugin.Message{
			Sender:  "Plugin Test Suite",
			Message: err.Error(),
			Level:   plugin.ERROR,
		})
		t.Fail()
	}
	fmt.Println(p)
}
func TestWriteLog(t *testing.T) {
	//log no progress or status.
	plugin.SetLogLevel(plugin.WARN)
	plugin.Log(plugin.Message{
		Message: "ERROR",
		Level:   plugin.ERROR,
	})
	plugin.Log(plugin.Message{
		Message: "WARNING",
		Level:   plugin.WARN,
		Sender:  "testing",
	})
	plugin.Log(plugin.Message{
		Message: "INFO",
		Level:   plugin.INFO,
	})
}
*/
