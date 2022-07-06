package wat_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

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
	path := "../exampledata/wat-job.yaml"
	jobManifest := wat.JobManifest{}
	readObject(t, path, &jobManifest)

	//construct a job manager
	jobManager, err := wat.Init(jobManifest)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// validate -
	err = jobManager.Validate()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	//compute...
	err = jobManager.ProcessJob()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	pathOutput := "../exampledata/runs/event_0/ras-mutator_payload.yml"
	outputFile, err := os.Open(pathOutput)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	outputBytes, err := ioutil.ReadAll(outputFile)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	pathComparison := "../exampledata/ras-mutator_payload.yml"
	comparisonFile, err := os.Open(pathComparison)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	comparisonBytes, err := ioutil.ReadAll(comparisonFile)
	if err != nil {
		fmt.Println(err)
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
