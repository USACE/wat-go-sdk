package plugindatamodel

// Plugin
type Plugin struct {
	Name        string   `json:"name" yaml:"name"`
	ImageAndTag string   `json:"image_and_tag" yaml:"image_and_tag"`
	Command     []string `json:"command" yaml:"command"`
}
