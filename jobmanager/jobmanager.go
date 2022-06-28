package jobmanager

import (
	"errors"
	"fmt"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/service/batch"
)

//JobManager
type JobManager struct {
	job           Job
	store         filestore.FileStore
	captainCrunch *batch.Batch
}

func Init(job Job, fs filestore.FileStore, batchClient *batch.Batch) JobManager {
	return JobManager{
		job:           job,
		store:         fs,
		captainCrunch: batchClient,
	}
}
func (jm JobManager) ProcessJob() error {
	err := jm.job.ProvisionResources()
	fmt.Println(err)
	//add in defer and recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", r)
			fmt.Println("Tearing Down Resources")
			err = jm.job.DestructResources()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	err = jm.job.GeneratePayloads(jm.store)
	fmt.Println(err)
	//create error channel.
	//create waitgroups to throttle compute resources?
	for i := jm.job.EventStartIndex; i < jm.job.EventEndIndex; i++ {
		go func(index int) {
			err = jm.job.ComputeEvent(index)
			fmt.Println(err)
		}(i)

	}
	//need a wait group or a buffer channel to stall the destruction until we finish the jobs
	err = jm.job.DestructResources()
	fmt.Println(err)
	return errors.New("Job Processed!")
}
