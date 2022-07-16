package wat

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/usace/wat-go-sdk/plugin"
)

type Provider string

const (
	BATCH Provider = "AWS Batch"
	MOCK  Provider = "Mock"
)

type ComputeResourceRequirements struct {
	LinkedManifestID   string `json:"linked_manifest_id" yaml:"linked_manifest_id"`
	ComputeEnvironment string `json:"compute_environment" yaml:"compute_environment"` //is this provided as JSON?
	JobDefinition      string `json:"job_definition" yaml:"job_definition"`
	Queue              string `json:"job_queue" yaml:"job_queue"`
}
type CloudProvider interface {
	//initialize it with some sort of configuration?
	ProvisionResources(jobManager *JobManager) error
	TearDownResources(job Job) error
	ProcessTask(job *Job, eventIndex int, payloadPath string, linkedManifest LinkedModelManifest) error
}

type BatchCloudProvider struct {
	BatchSession *batch.Batch
}

func (b BatchCloudProvider) ProvisionResources(jobManager *JobManager) error {
	resources := make(map[string]provisionedResources, len(jobManager.job.Dag.LinkedManifests))
	for _, lm := range jobManager.job.Dag.LinkedManifests {
		computeResourceRequirements, err := jobManager.LinkedManifestComputeResources(lm.ManifestID)
		if err != nil {
			plugin.Log(plugin.Message{
				Message: err.Error(),
				Level:   plugin.ERROR,
				Sender:  jobManager.job.Id,
			})
			return err
		}
		//get and store the arn's
		computeEnvOutput, err := newComputeEnvironment(b.BatchSession, computeResourceRequirements.ComputeEnvironment)
		if err != nil {
			plugin.Log(plugin.Message{
				Message: err.Error(),
				Level:   plugin.ERROR,
				Sender:  jobManager.job.Id,
			})
			return err
		}
		computeEnvironmentArn := computeEnvOutput.ComputeEnvironmentArn
		queueOutput, err := newQueue(b.BatchSession, computeResourceRequirements.Queue, computeEnvironmentArn)
		if err != nil {
			plugin.Log(plugin.Message{
				Message: err.Error(),
				Level:   plugin.ERROR,
				Sender:  jobManager.job.Id,
			})
			return err
		}
		queueArn := queueOutput.JobQueueArn
		jobOutput, err := newJobDefinition(b.BatchSession, computeResourceRequirements.JobDefinition)
		if err != nil {
			plugin.Log(plugin.Message{
				Message: err.Error(),
				Level:   plugin.ERROR,
				Sender:  jobManager.job.Id,
			})
			return err
		}
		jobDefinitionArn := jobOutput.JobDefinitionArn
		lmResource := provisionedResources{
			LinkedManifestID:      lm.ManifestID,
			ComputeEnvironmentARN: computeEnvironmentArn,
			JobDefinitionARN:      jobDefinitionArn,
			JobARN:                []*string{},
			QueueARN:              queueArn,
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
func (b BatchCloudProvider) TearDownResources(job Job) error {

	for _, resources := range job.Dag.Resources {
		//kill all active jobs?
		for _, jobArn := range resources.JobARN {
			err := cancelJob(b.BatchSession, jobArn)
			plugin.Log(plugin.Message{
				Message: err.Error(),
				Level:   plugin.ERROR,
				Sender:  job.Id,
			})
		}
		_, err := deleteJobDefinition(b.BatchSession, resources.JobDefinitionARN)
		if err != nil {
			plugin.Log(plugin.Message{
				Message: err.Error(),
				Level:   plugin.ERROR,
				Sender:  job.Id,
			})
		}
		_, err = deleteQueue(b.BatchSession, resources.QueueARN)
		if err != nil {
			plugin.Log(plugin.Message{
				Message: err.Error(),
				Level:   plugin.ERROR,
				Sender:  job.Id,
			})
		}
		_, err = deleteComputeEnvironment(b.BatchSession, resources.ComputeEnvironmentARN)
		if err != nil {
			plugin.Log(plugin.Message{
				Message: err.Error(),
				Level:   plugin.ERROR,
				Sender:  job.Id,
			})
		}
	}
	return nil
}
func (b BatchCloudProvider) ProcessTask(job *Job, eventIndex int, payloadPath string, linkedManifest LinkedModelManifest) error {
	batchJobArn := "Placeholder for Batch response"
	//set job arn
	resources, ok := job.Dag.Resources[linkedManifest.ManifestID]
	dependsOn, err := job.Dag.Dependencies(linkedManifest, eventIndex)
	if err != nil {
		return err
	}
	jobType := batch.ArrayJobDependencyNToN
	dependencies := make([]*batch.JobDependency, len(dependsOn))
	for idx, s := range dependsOn {
		dependency := batch.JobDependency{
			JobId: s,
			Type:  &jobType,
		}
		dependencies[idx] = &dependency
	}
	commandvalues := linkedManifest.Plugin.Command
	command := make([]*string, len(commandvalues)+1) //is this always true?
	for idx, c := range commandvalues {
		command[idx] = &c
	}
	command[len(commandvalues)] = &payloadPath
	jobInput := batch.SubmitJobInput{
		DependsOn: dependencies,
		ContainerOverrides: &batch.ContainerOverrides{
			Command: command,
		},
		JobDefinition: resources.JobDefinitionARN, //need to verify this.
		JobName:       &payloadPath,               //this is unique but may be more than 128 characters.
		JobQueue:      resources.QueueARN,
		Parameters:    nil, //parameters?
		RetryStrategy: nil,
		Timeout:       nil,
	}
	output, err := b.BatchSession.SubmitJob(&jobInput)
	if err != nil {
		return err
	}
	batchJobArn = *output.JobId
	if ok {
		resources.JobARN = append(resources.JobARN, &batchJobArn)
		job.Dag.Resources[linkedManifest.ManifestID] = resources
	} else {
		return errors.New("task for " + linkedManifest.Plugin.Name)
	}
	return nil
}
func cancelJob(b *batch.Batch, jobId *string) error {
	reason := "mind your buisness."
	//terminate fails running jobs also.
	cancelInput := batch.TerminateJobInput{
		JobId:  jobId,
		Reason: &reason,
	}
	_, err := b.TerminateJob(&cancelInput)
	if err != nil {
		return err
	}
	return err
}
func InitalizeSession(config Config) (CloudProvider, error) {
	//check the config to see if it should be batch or some other provider?
	switch config.CloudProvider {
	case BATCH:
		provider := BatchCloudProvider{}
		var batchClient *batch.Batch
		awsconfig, err := config.PrimaryConfig()
		if err != nil {
			return provider, err
		}
		creds := credentials.NewStaticCredentials(
			awsconfig.AWS_ACCESS_KEY_ID,
			awsconfig.AWS_SECRET_ACCESS_KEY,
			"",
		)
		cfg := aws.NewConfig().WithRegion(awsconfig.AWS_REGION).WithCredentials(creds)
		s, err := session.NewSession(cfg)
		if err != nil {
			return provider, err
		}
		batchClient = batch.New(s)
		provider.BatchSession = batchClient
		return provider, nil
	case MOCK:
		return MockProvider{config: config}, nil
	default:
		return nil, errors.New("cloud provider unknown")
	}

}

const (
	COMPUTE_ENV    = "computeEnvironmentFile"
	JOB_DEFINITION = "jobDefinitionFile"
	JOB_QUEUE      = "jobQueueFile"
	NEW_JOB        = "newJob"
)

// Batch Service Wrappers
// Takes inputs or outputs from batch.*
func directiveFromJson(bs interface{}, path string, v any) error {
	var payloadFile string
	switch bs {
	case COMPUTE_ENV:
		payloadFile = path

	case JOB_DEFINITION:
		payloadFile = path

	case JOB_QUEUE:
		payloadFile = path

	case NEW_JOB:
		payloadFile = path

	default:
		return errors.New("unrecognized service")
	}

	instructions, err := ioutil.ReadFile(payloadFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(instructions, &v)
	if err != nil {
		return err
	}
	return nil
}

func newComputeEnvironment(bc *batch.Batch, computeEnvironmentPath string) (output *batch.CreateComputeEnvironmentOutput, err error) {
	var computeEnvironment batch.CreateComputeEnvironmentInput
	err = directiveFromJson(COMPUTE_ENV, computeEnvironmentPath, &computeEnvironment)
	if err != nil {
		return output, err
	}

	output, err = bc.CreateComputeEnvironment(&computeEnvironment)
	if err != nil {
		return output, err
	}

	return output, nil
}

func deleteComputeEnvironment(bc *batch.Batch, computeEnvironmentArn *string) (output *batch.DeleteComputeEnvironmentOutput, err error) {
	updateComputeEnvironmentData := batch.UpdateComputeEnvironmentInput{ComputeEnvironment: computeEnvironmentArn,
		State: aws.String("DISABLED")}

	_, err = bc.UpdateComputeEnvironment(&updateComputeEnvironmentData)
	if err != nil {
		return output, err
	}

	// Wait for AWS to update resources
	time.Sleep(90 * time.Second)

	deleteComputeEnvironmentData := batch.DeleteComputeEnvironmentInput{ComputeEnvironment: computeEnvironmentArn}

	output, err = bc.DeleteComputeEnvironment(&deleteComputeEnvironmentData)
	if err != nil {
		return output, err
	}

	return output, err
}

func newJobDefinition(bc *batch.Batch, path string) (output *batch.RegisterJobDefinitionOutput, err error) {
	var jobDefinition batch.RegisterJobDefinitionInput
	err = directiveFromJson(JOB_DEFINITION, path, &jobDefinition)
	if err != nil {
		return output, err
	}

	output, err = bc.RegisterJobDefinition(&jobDefinition)
	if err != nil {
		return output, err
	}

	// write to output file
	return output, err
}

func deleteJobDefinition(bc *batch.Batch, jobDefinitionArn *string) (output *batch.DeregisterJobDefinitionOutput, err error) {
	jobDefinitionDataInput := batch.DeregisterJobDefinitionInput{JobDefinition: jobDefinitionArn}
	_, err = bc.DeregisterJobDefinition(&jobDefinitionDataInput)

	if err != nil {
		return output, err
	}
	return output, err
}

func newQueue(bc *batch.Batch, path string, computeEnvironment *string) (output *batch.CreateJobQueueOutput, err error) {
	var jobQueue batch.CreateJobQueueInput
	err = directiveFromJson(JOB_QUEUE, path, &jobQueue)
	if err != nil {
		return output, err
	}

	// TODO: Think through the jobQueue.ComputeEnvironmentOrder list
	if *computeEnvironment != "" {
		jobQueue.ComputeEnvironmentOrder[0].ComputeEnvironment = computeEnvironment
	}

	output, err = bc.CreateJobQueue(&jobQueue)
	if err != nil {
		return output, err
	}

	return output, err
}

func deleteQueue(bc *batch.Batch, jobQueueArn *string) (output *batch.DeleteJobQueueOutput, err error) {

	updateQueueData := batch.UpdateJobQueueInput{JobQueue: jobQueueArn,
		State: aws.String("DISABLED")}

	updatedJobQueueData, err := bc.UpdateJobQueue(&updateQueueData)
	if err != nil {
		return output, err
	}

	// Wait for AWS to update resources
	time.Sleep(30 * time.Second)

	jobQueueData := batch.DeleteJobQueueInput{JobQueue: updatedJobQueueData.JobQueueName}
	_, err = bc.DeleteJobQueue(&jobQueueData)
	if err != nil {
		return output, err
	}

	return output, err
}
