package plugin

//Plugin a plugin is defined by its name an image and tag and the command necessary to call it to run a payload.
type Plugin struct {
	Name        string   `json:"name" yaml:"name"`
	ImageAndTag string   `json:"image_and_tag" yaml:"image_and_tag"`
	Command     []string `json:"command" yaml:"command"`
}

//ModelIdentifier a model identifier describes a model at a high level, the name, a given alternative (simulation, plan, etc.) and the associated files that define the geometry (optional)
type ModelIdentifier struct {
	Name        string              `json:"name" yaml:"name"`
	Alternative string              `json:"alternative,omitempty" yaml:"alternative,omitempty"`
	Files       []ResourcedFileData `json:"files,omitempty" yaml:"files,omitempty"`
}
type Store string

const (
	S3    Store = "S3"
	LOCAL Store = "Local"
	//others?
)

//ResourceInfo defines the elements to resolve to a file object, what store it uses, the root (or bucket) and then the path to the asset.
type ResourceInfo struct {
	Store Store  `json:"store" yaml:"store"`                   // s3, azure, local, queue?
	Root  string `json:"root" yaml:"root"`                     // bucket, rootdir, queue?
	Path  string `json:"path,omitempty" yaml:"path,omitempty"` // path to object
}

/* up to date as of 6/21/2022
The main elements of concern in a plugin are the following elements:
1. A Manifest File
2. A Linked Manifest file
3. A Payload File

## Manifest
	Responsible for describing the inputs and outputs specific to a model and a plugin
	- describes inputs and outputs
	- is generated by an MCAT based on a model for a specific plugin.
	- is unique to a model and plugin combination.
## Linked Manifest File
	Specific to a model and plugin, describes the specific inputs (and what manifest will be generating the input) and outputs for a model based on the manifest definition
	- describes selected inputs and what manifest will be generating the input and outputs
	- is defined by a user to reflect sources for specific inputs
	- changes based on a model, and the overal wat job it is designed to be linked in.
## Payload file
	Specific to a plugin, model, and event
	- describes where inputs can be found and where to put outputs
	- is generated by wat from a set of linked manifests (effectively a DAG)
	- changes event by event for a given model
*/
