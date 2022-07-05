package wat

import (
	"fmt"
)

//JobManager
type JobManager struct {
	job Job
	//store         filestore.FileStore
	//captainCrunch *batch.Batch
}

func Init(jobManifest JobManifest) (JobManager, error) { //, fs filestore.FileStore, batchClient *batch.Batch) JobManager {
	jobManager := JobManager{}
	job, err := jobManifest.ConvertToJob()
	if err != nil {
		return jobManager, err
	}
	orderedManifests, err := job.Dag.TopologicallySort()
	if err != nil {
		return jobManager, err
	}
	job.Dag.LinkedManifests = orderedManifests //*/
	jobManager.job = job
	return jobManager, nil
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
	err = jm.job.GeneratePayloads() //jm.store
	fmt.Println(err)
	//create error channel.
	for i := jm.job.EventStartIndex; i < jm.job.EventEndIndex; i++ {
		go func(index int) {
			err = jm.job.ComputeEvent(index)
			fmt.Println(err)
		}(i)

	}
	//need a wait group or a buffer channel to stall the destruction until we finish the jobs
	err = jm.job.DestructResources()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Job Processed!")
	return nil
}
func (jm JobManager) Validate() error {
	err := jm.job.ValidateLinkages() //evaluate if this can be trimmed down to "validateLinkages"
	if err != nil {
		return err
	}
	_, err = jm.job.Dag.TopologicallySort()
	if err != nil {
		return err
	}

	return nil
}