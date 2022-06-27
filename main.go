package main

import (
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"

	"gopkg.in/yaml.v3"
)

func main() {
	fmt.Println("wat-go-sdk")

	manifestFile := "/app/exampledata/manifest-1.yaml"
	linkedManifest := "/app/exampledata/linked-manifest-1.yaml"

	var manifest Manifest

	data, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		fmt.Println("Error", err)
	}

	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		fmt.Println("Error", err)
	}

	// Add Unique ID's to create links in manifest
	manifest.AddLinks()

	linkedData, err := yaml.Marshal(manifest)
	if err != nil {
		fmt.Println("Error", err)
	}

	err = ioutil.WriteFile(linkedManifest, linkedData, 0644)
	if err != nil {
		fmt.Println("Error", err)
	}
}

type Manifest struct {
	ManifestID      string          `yaml:"manifest_id"`
	Plugin          Plugin          `yaml:"plugin"`
	ModelIdentifier ModelIdentifier `yaml:"model_identifier"`
	Inputs          []Input         `yaml:"inputs"`
	Outputs         []Output        `yaml:"outputs"`
}

func (m *Manifest) AddLinks() {
	m.ManifestID = uuid.New().String()
	for i := range m.Inputs {
		m.Inputs[i].NewID()
		for j := range m.Inputs[i].InternalPaths {
			m.Inputs[i].InternalPaths[j].NewID()
		}
	}

	for i := range m.Outputs {
		m.Outputs[i].NewID()
	}
}

type Plugin struct {
	Name        string   `yaml:"name"`
	ImageAndTag string   `yaml:"image_and_tag"`
	Command     []string `yaml:"command"`
}
type ModelIdentifier struct {
	Name        string `yaml:"name"`
	Alternative string `yaml:"alternative"`
}
type InternalPath struct {
	PathID   string `yaml:"path_id,omitempty"`
	Pathname string `yaml:"pathname"`
}
type Input struct {
	InputID       string         `yaml:"input_id,omitempty"`
	Filename      string         `yaml:"filename"`
	InternalPaths []InternalPath `yaml:"internal_paths"`
}
type Output struct {
	OutputID string `yaml:"output_id,omitempty"`
	Filename string `yaml:"filename"`
}

type DataIdentifier interface {
	ReadID() string
}

func (m Manifest) ReadID() string {
	return m.ManifestID
}

func (i Input) ReadID() string {
	return i.InputID
}

func (internalPath InternalPath) ReadID() string {
	return internalPath.PathID
}

func (output Output) ReadID() string {
	return output.OutputID
}

func (m *Manifest) NewID() {
	m.ManifestID = uuid.New().String()
}

func (input *Input) NewID() {
	input.InputID = uuid.New().String()
}

func (internalPath *InternalPath) NewID() {
	internalPath.PathID = uuid.New().String()
}

func (output *Output) NewID() {
	output.OutputID = uuid.New().String()
}
