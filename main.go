package main

import (
	"fmt"

	"github.com/usace/wat-go-sdk/jobmanager"
	"github.com/usace/wat-go-sdk/pluginmanager"
)

func main() {
	fmt.Print("wat-go-sdk")
	//load s3
	//load batch
	outputDestination := pluginmanager.ResourceInfo{
		Scheme:    "s3",
		Authority: "cloud-wat-dev",
		Path:      "/runs",
	}
	inputSource := pluginmanager.ResourceInfo{
		Scheme:    "s3",
		Authority: "cloud-wat-dev",
		Path:      "/data",
	}
	job := jobmanager.Job{
		EventCount:           20,
		DirectedAcyclicGraph: jobmanager.DirectedAcyclicGraph{},
		OutputDestination:    outputDestination,
		InputSource:          inputSource,
	}
	jobmanager := jobmanager.Init(job, nil, nil)
	err := jobmanager.ProcessJob()
	fmt.Println(err)
}

func mockDag() jobmanager.DirectedAcyclicGraph {
	manifests := make([]pluginmanager.ModelManifest, 3)
	manifests[0] = pluginmanager.ModelManifest{
		Plugin: pluginmanager.Plugin{
			Name:        "hydrograph_scaler",
			ImageAndTag: "williamlehman/hydrographscaler:v0.0.11",
			Command:     []string{"./main", "-payload"},
		},
	}
	manifests[1] = pluginmanager.ModelManifest{
		Plugin: pluginmanager.Plugin{
			Name:        "ras-mutator",
			ImageAndTag: "lawlerseth/ras-mutator:v0.1.1",
			Command:     []string{"./h5rasedit", "wat", "-m", "host.docker.internal:9000", "-f"},
		},
	}
	manifests[2] = pluginmanager.ModelManifest{
		Plugin: pluginmanager.Plugin{
			Name:        "ras-unsteady",
			ImageAndTag: "lawlerseth/ras-unsteady:v0.0.2",
			Command:     []string{"./watrun", "-m", "host.docker.internal:9000", "-f"},
		},
	}
	return jobmanager.DirectedAcyclicGraph{
		Nodes: manifests,
	}
}
