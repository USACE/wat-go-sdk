package wat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

type CloudProvider interface {
	//initialize it with some sort of configuration?
	ProvisionResources(job *Job) error
	TearDownResources(job *Job) error
	ProcessTask(job *Job, eventIndex int, payloadPath string, linkedManifest LinkedModelManifest) error
}
type BatchCloudProvider struct {
	BatchSession *batch.Batch
}

func (b BatchCloudProvider) ProvisionResources(job *Job) error {
	return nil
}
func (b BatchCloudProvider) TearDownResources(job *Job) error {
	return nil
}
func (b BatchCloudProvider) ProcessTask(job *Job, eventIndex int, payloadPath string, linkedManifest LinkedModelManifest) error {
	return nil
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
func (bp AWSBatchPayload) DirectiveFromJson(bs interface{}, v any) error {
	var payloadFile string

	switch bs {
	case COMPUTE_ENV:
		payloadFile = bp.ComputeEnvironmentFile

	case JOB_DEFINITION:
		payloadFile = bp.JobDefinitionFile

	case JOB_QUEUE:
		payloadFile = bp.JobQueueFile

	case NEW_JOB:
		payloadFile = bp.NewJob

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

func (bp AWSBatchPayload) NewComputeEnvironment(bc *batch.Batch) (output *batch.CreateComputeEnvironmentOutput, err error) {
	var computeEnvironment batch.CreateComputeEnvironmentInput
	err = bp.DirectiveFromJson(COMPUTE_ENV, &computeEnvironment)
	if err != nil {
		return output, err
	}

	output, err = bc.CreateComputeEnvironment(&computeEnvironment)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (bp AWSBatchPayload) DeleteComputeEnvironment(bc *batch.Batch) (output *batch.DeleteComputeEnvironmentOutput, err error) {
	var computeEnvironment batch.CreateComputeEnvironmentOutput
	err = bp.DirectiveFromJson(COMPUTE_ENV, &computeEnvironment)
	if err != nil {
		return output, err
	}

	updateComputeEnvironmentData := batch.UpdateComputeEnvironmentInput{ComputeEnvironment: computeEnvironment.ComputeEnvironmentName,
		State: aws.String("DISABLED")}

	_, err = bc.UpdateComputeEnvironment(&updateComputeEnvironmentData)
	if err != nil {
		return output, err
	}

	// Wait for AWS to update resources
	time.Sleep(90 * time.Second)

	deleteComputeEnvironmentData := batch.DeleteComputeEnvironmentInput{ComputeEnvironment: computeEnvironment.ComputeEnvironmentName}

	output, err = bc.DeleteComputeEnvironment(&deleteComputeEnvironmentData)
	if err != nil {
		return output, err
	}

	return output, err
}

func (bp AWSBatchPayload) NewJobDefinition(bc *batch.Batch) (output *batch.RegisterJobDefinitionOutput, err error) {
	var jobDefinition batch.RegisterJobDefinitionInput
	err = bp.DirectiveFromJson(JOB_DEFINITION, &jobDefinition)
	if err != nil {
		fmt.Println("Error", err)
	}

	output, err = bc.RegisterJobDefinition(&jobDefinition)
	if err != nil {
		return output, err
	}

	// write to output file
	return output, err
}

func (bp AWSBatchPayload) DeleteJobDefinition(bc *batch.Batch) (output *batch.DeregisterJobDefinitionOutput, err error) {
	var jobDefinition batch.RegisterJobDefinitionOutput
	err = bp.DirectiveFromJson(JOB_DEFINITION, &jobDefinition)
	if err != nil {
		fmt.Println("Error", err)
	}

	jobDefinitionDataInput := batch.DeregisterJobDefinitionInput{JobDefinition: jobDefinition.JobDefinitionArn}

	_, err = bc.DeregisterJobDefinition(&jobDefinitionDataInput)

	if err != nil {
		return output, err
	}
	return output, err
}

func (bp AWSBatchPayload) NewQueue(bc *batch.Batch, computeEnvironment *string) (output *batch.CreateJobQueueOutput, err error) {
	var jobQueue batch.CreateJobQueueInput
	err = bp.DirectiveFromJson(JOB_QUEUE, &jobQueue)
	if err != nil {
		fmt.Println("Error", err)
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

func (bp AWSBatchPayload) DeleteQueue(bc *batch.Batch) (output *batch.DeleteJobQueueOutput, err error) {
	var jobQueue batch.CreateJobQueueInput
	err = bp.DirectiveFromJson(JOB_QUEUE, &jobQueue)
	if err != nil {
		fmt.Println("Error", err)
	}

	updateQueueData := batch.UpdateJobQueueInput{JobQueue: jobQueue.JobQueueName,
		State: aws.String("DISABLED")}

	updatedJobQueueData, err := bc.UpdateJobQueue(&updateQueueData)
	if err != nil {
		fmt.Println("Error....", err)
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

//@TODO this looks wrong - this looks like create jobqueue not submit job.
func (bp AWSBatchPayload) SubmitJob(bc *batch.Batch, computeEnvironment *string) (output *batch.CreateJobQueueOutput, err error) {
	var jobQueue batch.CreateJobQueueInput
	err = bp.DirectiveFromJson(JOB_QUEUE, &jobQueue)
	if err != nil {
		fmt.Println("Error", err)
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
