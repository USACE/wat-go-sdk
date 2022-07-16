package wat

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/usace/wat-go-sdk/plugin"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/net/context"
)

type MockProvider struct {
}

func (m MockProvider) ProvisionResources(jobManager *JobManager) error {
	resources := make(map[string]provisionedResources, len(jobManager.job.Dag.LinkedManifests))
	for _, lm := range jobManager.job.Dag.LinkedManifests {
		computeEnvironmentArn := lm.ManifestID
		queueArn := lm.ManifestID
		jobDefinitionArn := lm.ManifestID
		lmResource := provisionedResources{
			LinkedManifestID:      lm.ManifestID,
			ComputeEnvironmentARN: &computeEnvironmentArn,
			JobDefinitionARN:      &jobDefinitionArn,
			JobARN:                []*string{},
			QueueARN:              &queueArn,
		}
		resources[lm.ManifestID] = lmResource
	}
	jobManager.job.Dag.Resources = resources
	plugin.Log(plugin.Message{
		Message: "provisioned resources",
		Level:   plugin.INFO,
		Sender:  jobManager.job.Id,
	})
	return nil
}
func (m MockProvider) TearDownResources(job Job) error {
	//remove all remaining containers.
	plugin.Log(plugin.Message{
		Message: "Kablooie",
		Level:   plugin.INFO,
		Sender:  job.Id,
	})
	return nil
}
func (m MockProvider) ProcessTask(job *Job, eventIndex int, payloadPath string, linkedManifest LinkedModelManifest) error {
	plugin.Log(plugin.Message{
		Message: "Processing Task",
		Level:   plugin.INFO,
		Sender:  job.Id,
	})
	resources, ok := job.Dag.Resources[linkedManifest.ManifestID]
	batchJobArn := payloadPath
	if ok {
		resources.JobARN = append(resources.JobARN, &batchJobArn)
		job.Dag.Resources[linkedManifest.ManifestID] = resources
		env := make([]string, 0)
		_, err := startContainer(linkedManifest.ImageAndTag, batchJobArn, env)
		if err != nil {
			plugin.Log(plugin.Message{
				Message: err.Error(),
				Level:   plugin.ERROR,
				Sender:  job.Id,
			})
		}
	} else {
		return errors.New("task for " + linkedManifest.Plugin.Name)
	}
	return nil
}
func startContainer(imageWithTag string, payloadPath string, environmentVariables []string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}
	//cli.NegotiateAPIVersion(ctx)
	reader, err := cli.ImagePull(ctx, imageWithTag, types.ImagePullOptions{})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	io.Copy(os.Stdout, reader)
	var chc *container.HostConfig
	var nnc *network.NetworkingConfig
	var vp *v1.Platform

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        imageWithTag,
		Cmd:          []string{"./main", "-payload=" + payloadPath},
		Tty:          true,
		AttachStdout: true,
		Env:          environmentVariables,
	}, chc, nnc, vp, "")
	if err != nil {
		return "", err
	}
	//retrieve container messages and parrot to lambda standard out.
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, Follow: true})
	if err != nil {
		return "", err
	}
	//defer out.Close()
	io.Copy(os.Stdout, out)
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	statuschn, errchn := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errchn:
		if err != nil {
			log.Fatal(err)
		}
	case status := <-statuschn:
		log.Printf("status.StatusCode: %#+v\n", status.StatusCode)
		//cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	}

	return resp.ID, err
}
