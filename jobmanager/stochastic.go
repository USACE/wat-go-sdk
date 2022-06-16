package jobmanager

/*
import (
	"encoding/json"
	"fmt"
	"math/rand"
)

func (sj StochasticJob) GeneratePayloads(config config.WatConfig, fs filestore.FileStore, awsBatch *batch.Batch) error {
	//provision resources
	resources, err := sj.ProvisionResources(awsBatch)
	//create random seed generators.
	eventrg := rand.New(rand.NewSource(sj.InitialEventSeed))             //Natural Variability
	realizationrg := rand.New(rand.NewSource(sj.InitialRealizationSeed)) //KnowledgeUncertianty
	if err != nil {
		return err
	}
	nodes := sj.Dag.Nodes
	realizationRandomGeneratorByPlugin := make([]*rand.Rand, len(nodes))
	eventRandomGeneratorByPlugin := make([]*rand.Rand, len(nodes))
	for idx := range nodes {
		realizationSeeder := realizationrg.Int63()
		eventSeeder := eventrg.Int63()
		realizationRandomGeneratorByPlugin[idx] = rand.New(rand.NewSource(realizationSeeder))
		eventRandomGeneratorByPlugin[idx] = rand.New(rand.NewSource(eventSeeder))
	}
	for i := 0; i < sj.TotalRealizations; i++ { //knowledge uncertainty loop
		realizationIndexedSeeds := make([]IndexedSeed, len(nodes))
		for idx := range nodes {
			realizationSeed := realizationRandomGeneratorByPlugin[idx].Int63()
			realizationIndexedSeeds[idx] = IndexedSeed{Index: i, Seed: realizationSeed}
		}
		for j := 0; j < sj.EventsPerRealization; j++ { //natural variability loop
			//ultimately need to send messages for each task in the event (defined by the dag)
			//event randoms will spawn in unpredictable ways if we dont pre spawn them.
			pluginEventIndexedSeeds := make([]IndexedSeed, len(nodes))
			for idx := range nodes {
				pluginEventSeed := realizationRandomGeneratorByPlugin[idx].Int63()
				pluginEventIndexedSeeds[idx] = IndexedSeed{Index: j, Seed: pluginEventSeed}
			}
			go sj.ProcessDAG(config, i, j, realizationIndexedSeeds, pluginEventIndexedSeeds, fs, awsBatch, resources)
		}
	}
	fmt.Println("complete")
	return nil
}

func (sj StochasticJob) ProcessDAG(config config.WatConfig, realization int, event int, realizationIndexedSeeds []IndexedSeed, eventIndexedSeedsByPlugin []IndexedSeed, fs filestore.FileStore, awsBatch *batch.Batch, resources []ProvisionedResources) {
	outputDestinationPath := fmt.Sprintf("%v%v%v/%v%v", sj.Outputdestination.Fragment, "realization_", realization, "event_", event)
	for idx, n := range sj.Dag.Nodes {
		fmt.Println(n.ImageAndTag, outputDestinationPath)
		ec := EventConfiguration{
			OutputDestination: ResourceInfo{
				Scheme:    sj.Outputdestination.Scheme,
				Authority: sj.Outputdestination.Authority,
				Fragment:  outputDestinationPath,
			},
			Realization:     realizationIndexedSeeds[idx],
			Event:           eventIndexedSeedsByPlugin[idx],
			EventTimeWindow: sj.TimeWindow,
		}
		//write event configuration to s3.
		ecbytes, err := json.Marshal(ec)
		if err != nil {
			panic(err)
		}
		path := outputDestinationPath + "/" + n.Plugin.Name + "_Event Configuration.json"
		fmt.Println("putting object in fs:", path)
		_, err = fs.PutObject(path, ecbytes)
		if err != nil {
			fmt.Println("failure to push event configuration to filestore:", err)
			panic(err)
		}
		payload := Mock2DModelPayload(sj.Inputsource, ec.OutputDestination, outputDestinationPath, n.Plugin)
		bytes, err := yaml.Marshal(payload)
		if err != nil {
			panic(err)
		}
		//put payload in s3
		path = outputDestinationPath + "/" + n.Plugin.Name + "_payload.yml"
		fmt.Println("putting object in fs:", path)
		_, err = fs.PutObject(path, bytes)
		if err != nil {
			fmt.Println("failure to push payload to filestore:", err)
			panic(err)
		}
		//submit job to batch.
		//if n.Plugin.Name == "hydrograph_scaler" {
		/*s, err := utils.StartContainer(n.Plugin, path, config.EnvironmentVariables())
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Print(s)
*/
//}
