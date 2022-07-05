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
	readobject(t, path, &jobManifest)
}
func TestReadLinkedManifest(t *testing.T) {
	path := "../exampledata/ras_mutator_linked_manifest.yaml"
	linkedmanifest := wat.LinkedModelManifest{}
	readobject(t, path, &linkedmanifest)
}
func TestComputePayloads(t *testing.T) {
	//read a jobmanifest into memory
	path := "../exampledata/wat-job.yaml"
	jobManifest := wat.JobManifest{}
	readobject(t, path, &jobManifest)
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
	pathoutput := "../exampledata/runs/event_0/ras-mutator_payload.yml"
	outputfile, err := os.Open(pathoutput)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	outputbytes, err := ioutil.ReadAll(outputfile)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	pathcomparison := "../exampledata/ras-mutator_payload.yml"
	comparisonfile, err := os.Open(pathcomparison)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	comparisonbytes, err := ioutil.ReadAll(comparisonfile)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if !reflect.DeepEqual(outputbytes, comparisonbytes) {
		fmt.Println(string(outputbytes))
		fmt.Println(string(comparisonbytes))
		t.Fail()
	}
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
