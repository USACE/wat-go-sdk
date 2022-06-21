package plugindatamodel_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/usace/wat-go-sdk/plugindatamodel"
	"gopkg.in/yaml.v3"
)

func TestPostCompute(t *testing.T) {
	file, err := os.Open("../exampledata/example_payload.yaml")
	defer file.Close()
	if err != nil {
		t.Fail()
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fail()
	}
	payload := plugindatamodel.Payload{}
	err = yaml.Unmarshal(b, payload)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	log.Println(string(b))
	log.Println(payload)
}
