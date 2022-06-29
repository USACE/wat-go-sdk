package jobmanager_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/usace/wat-go-sdk/jobmanager"
	"gopkg.in/yaml.v3"
)

func TestReadJobManifest(t *testing.T) {
	path := "../exampledata/wat-job.yaml"
	jobManifest := jobmanager.JobManifest{}
	readobject(t, path, &jobManifest)
}
func TestComputePayloads(t *testing.T) {
	//read a jobmanifest into memory
	path := "../exampledata/wat-job.yaml"
	jobManifest := jobmanager.JobManifest{}
	readobject(t, path, &jobManifest)
	//construct a job manager
	jobManager, err := jobmanager.Init(jobManifest)
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
