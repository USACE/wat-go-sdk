package wat_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	watplugin "github.com/usace/wat-go-sdk/plugin"
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
	//read a jobmanifest into memory
	watplugin.Logger.SetLogLevel(watplugin.ERROR)
	path := "../exampledata/wat-job.yaml"
	jobManifest := wat.JobManifest{}
	readObject(t, path, &jobManifest)

	//construct a job manager
	jobManager, err := wat.Init(jobManifest)
	if err != nil {
		watplugin.Logger.Log(watplugin.Log{
			Message: err.Error(),
			Level:   watplugin.ERROR,
		})
		t.Fail()
	}

	// validate -
	err = jobManager.Validate()
	if err != nil {
		watplugin.Logger.Log(watplugin.Log{
			Message: err.Error(),
			Level:   watplugin.ERROR,
		})
		t.Fail()
	}

	//compute...
	err = jobManager.ProcessJob()
	if err != nil {
		watplugin.Logger.Log(watplugin.Log{
			Message: err.Error(),
			Level:   watplugin.ERROR,
		})
		t.Fail()
	}

	pathOutput := "../exampledata/runs/event_0/ras-mutator_payload.yml"
	outputFile, err := os.Open(pathOutput)
	if err != nil {
		watplugin.Logger.Log(watplugin.Log{
			Message: err.Error(),
			Level:   watplugin.ERROR,
		})
		t.Fail()
	}

	outputBytes, err := ioutil.ReadAll(outputFile)
	if err != nil {
		watplugin.Logger.Log(watplugin.Log{
			Message: err.Error(),
			Level:   watplugin.ERROR,
		})
		t.Fail()
	}

	pathComparison := "../exampledata/ras-mutator_payload.yml"
	comparisonFile, err := os.Open(pathComparison)
	if err != nil {
		watplugin.Logger.Log(watplugin.Log{
			Message: err.Error(),
			Level:   watplugin.ERROR,
		})
		t.Fail()
	}

	comparisonBytes, err := ioutil.ReadAll(comparisonFile)
	if err != nil {
		watplugin.Logger.Log(watplugin.Log{
			Message: err.Error(),
			Level:   watplugin.ERROR,
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
		watplugin.Logger.Log(watplugin.Log{
			Message: err.Error(),
			Level:   watplugin.ERROR,
		})
		t.Fail()
	} else {
		watplugin.Logger.Log(watplugin.Log{
			Message: string(b),
			Level:   watplugin.INFO,
		})
		newTestLine := "\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n"
		watplugin.Logger.Log(watplugin.Log{
			Message: newTestLine,
			Level:   watplugin.INFO,
		})

		b2, err := yaml.Marshal(object)
		if err != nil {
			watplugin.Logger.Log(watplugin.Log{
				Message: err.Error(),
				Level:   watplugin.ERROR,
			})
			t.Fail()
		}

		watplugin.Logger.Log(watplugin.Log{
			Message: string(b2),
			Level:   watplugin.INFO,
		})

	}

}
