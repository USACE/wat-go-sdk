package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/usace/wat-go-sdk/plugin"
	"github.com/usace/wat-go-sdk/wat"
	"gopkg.in/yaml.v3"
)

func main() {
	fmt.Println("wat-go-sdk testing")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		_, cPresent := params["compute"]
		if cPresent {
			cfg, err := wat.InitConfig("./exampledata/watconfig.json")
			if err != nil {
				log.Fatal("error reading config")
			}
			//read a jobmanifest into memory
			plugin.SetLogLevel(plugin.DEBUG)
			path := "../exampledata/wat-job.yaml"
			jobManifest := wat.JobManifest{}
			readObject(path, &jobManifest)

			//construct a job manager
			jobManager, err := wat.Init(jobManifest, cfg)
			if err != nil {
				plugin.Log(plugin.Message{
					Message: err.Error(),
					Level:   plugin.ERROR,
				})
				log.Fatal("errors with the job manager")
			}

			// validate -
			err = jobManager.Validate()
			if err != nil {
				plugin.Log(plugin.Message{
					Message: err.Error(),
					Level:   plugin.ERROR,
				})
				log.Fatal("errors with the job.")
			}

			//compute...
			err = jobManager.ProcessJob()
		} else {
			http.Error(w, "404 not found.", http.StatusNotFound)
		}

	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func readObject(path string, object interface{}) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = yaml.Unmarshal(b, object)

	if err != nil {
		//log.Println(err)
		plugin.Log(plugin.Message{
			Message: err.Error(),
			Level:   plugin.ERROR,
		})
		log.Fatal(err.Error())
	} else {
		plugin.Log(plugin.Message{
			Message: string(b),
			Level:   plugin.INFO,
		})
	}
}
