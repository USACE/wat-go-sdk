package plugindatamodel

//Data
type Data struct {
	//Name - usually describes a location, or a need e.g. (Muncie RS 15893.2 or Event Configuration)
	Name string `json:"name" yaml:"name"`
	//Parameter - convention is to name the parameter in lower case.
	Parameter string `json:"parameter" yaml:"parameter"` //file, flow, stage, timewindow, eventconfiguration - definately not an enum because it could be anything
}

//acceptable formats? format options?
//optional/required
