package pluginmanager

//DataDescription describes what a model needs input for or what kind of output it can provide
type DataDescription struct {
	Name      string `json:"name" yaml:"name"`
	Parameter string `json:"parameter" yaml:"parameter"`
	Format    string `json:"format" yaml:"format"`
}

//ResourceInfo describes the elements to help find data
type ResourceInfo struct {
	//Might be S3, Redis, SQS, Azure, Local
	Scheme string `json:"scheme" yaml:"scheme"` //http or https for example
	//Bucket Name, service address, queue, bucket, root directory
	Authority string `json:"authority" yaml:"authority"` // //minio:9001 for example
	//path from bucket, key, omit for sqs?, path from bucket, path from root
	Path     string `json:"path,omitempty" yaml:"path,omitempty"`   //omit empty default value "/"
	Query    string `json:"query,omitempty" yaml:"query,omitempty"` //omit empty
	Fragment string `json:"fragment,omitempty" yaml:"fragment,omitempty"`
	//https://pkg.go.dev/go.lsp.dev/uri  consider this.
	/*
			    foo://example.com:8042/over/there?name=ferret#nose
		         \_/   \______________/\_________/ \_________/ \__/
		          |           |            |            |        |
		       scheme     authority       path        query   fragment
		          |   _____________________|__
		         / \ /                        \
		         urn:example:animal:ferret:nose
	*/
}

//LinkedDataDescription combines a DataDescription and a ResourceInfo object together to define how to access an input or an output.
type LinkedDataDescription struct {
	DataDescription string `json:"description" yaml:"description"`
	ResourceInfo    `json:"resource_info" yaml:"resource_info"`
}
type ModelLinks struct {
	LinkedInputs     []LinkedDataDescription `json:"linked_inputs" yaml:"linked_inputs"`
	NecessaryOutputs []LinkedDataDescription `json:"required_outputs" yaml:"required_outputs"`
}
