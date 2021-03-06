package wat_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/usace/wat-go-sdk/plugin"
	"github.com/usace/wat-go-sdk/wat"
	"gopkg.in/yaml.v3"
)

func TestReadJobManifest(t *testing.T) {
	path := "../exampledata/wat-job.yaml"
	jobManifest := wat.JobManifest{}
	readObject(t, path, &jobManifest)
}

func TestReadLinkedManifest(t *testing.T) {
	path := "../exampledata/ras_mutator_linked_manifest.yaml"
	linkedManifest := wat.LinkedModelManifest{}
	readObject(t, path, &linkedManifest)
}

func TestComputePayloads(t *testing.T) {
	cfg, err := wat.InitConfig("../exampledata/watconfig.json")
	if err != nil {
		t.Fail()
	}
	//read a jobmanifest into memory
	plugin.SetLogLevel(plugin.DEBUG)
	path := "../exampledata/wat-job.yaml"
	jobManifest := wat.JobManifest{}
	readObject(t, path, &jobManifest)

	//construct a job manager
	jobManager, err := wat.Init(jobManifest, cfg)
	if err != nil {
		plugin.Log(plugin.Message{
			Message: err.Error(),
			Level:   plugin.ERROR,
		})
		t.Fail()
	}

	// validate -
	err = jobManager.Validate()
	if err != nil {
		plugin.Log(plugin.Message{
			Message: err.Error(),
			Level:   plugin.ERROR,
		})
		t.Fail()
	}

	//compute...
	err = jobManager.ProcessJob()
	if err != nil {
		plugin.Log(plugin.Message{
			Message: err.Error(),
			Level:   plugin.ERROR,
		})
		t.Fail()
	}

	pathOutput := "../exampledata/runs/event_0/ras-mutator_payload.yml"
	outputFile, err := os.Open(pathOutput)
	if err != nil {
		plugin.Log(plugin.Message{
			Message: err.Error(),
			Level:   plugin.ERROR,
		})
		t.Fail()
	}

	outputBytes, err := ioutil.ReadAll(outputFile)
	if err != nil {
		plugin.Log(plugin.Message{
			Message: err.Error(),
			Level:   plugin.ERROR,
		})
		t.Fail()
	}

	pathComparison := "../exampledata/ras-mutator_payload.yml"
	comparisonFile, err := os.Open(pathComparison)
	if err != nil {
		plugin.Log(plugin.Message{
			Message: err.Error(),
			Level:   plugin.ERROR,
		})
		t.Fail()
	}

	comparisonBytes, err := ioutil.ReadAll(comparisonFile)
	if err != nil {
		plugin.Log(plugin.Message{
			Message: err.Error(),
			Level:   plugin.ERROR,
		})
		t.Fail()
	}

	if !reflect.DeepEqual(outputBytes, comparisonBytes) {
		fmt.Println("outputBytes: ", string(outputBytes))
		fmt.Println("comparisonBytes: ", string(comparisonBytes))
		t.Fail()
	}
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
		//log.Println(err)
		plugin.Log(plugin.Message{
			Message: err.Error(),
			Level:   plugin.ERROR,
		})
		t.Fail()
	} else {
		plugin.Log(plugin.Message{
			Message: string(b),
			Level:   plugin.INFO,
		})
		newTestLine := "\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n"
		plugin.Log(plugin.Message{
			Message: newTestLine,
			Level:   plugin.INFO,
		})

		b2, err := yaml.Marshal(object)
		if err != nil {
			plugin.Log(plugin.Message{
				Message: err.Error(),
				Level:   plugin.ERROR,
			})
			t.Fail()
		}

		plugin.Log(plugin.Message{
			Message: string(b2),
			Level:   plugin.INFO,
		})

	}

}
