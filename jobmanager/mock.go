package jobmanager

import (
	"github.com/usace/wat-go-sdk/pluginmanager"
)

func mockModelPayload(inputSource pluginmanager.ResourceInfo, outputDestination pluginmanager.ResourceInfo, eventParts string, plugin pluginmanager.Plugin) pluginmanager.Payload {
	mconfig := pluginmanager.ModelConfiguration{}
	inputs := make([]pluginmanager.LinkedDataDescription, 0)
	switch plugin.Name {
	case "hydrograph_scaler":
		mconfig.Name = "hydrographs"
		inputs = append(inputs, pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:      "Project File",
				Parameter: "Project Specification",
				Format:    ".json",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  inputSource.Fragment + "hsm.json",
			},
		})
		inputs = append(inputs, pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:      "Event Configuration",
				Parameter: "Event Configuration",
				Format:    ".json",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    outputDestination.Scheme,
				Authority: outputDestination.Authority,
				Fragment:  outputDestination.Fragment + "/hydrograph_scaler_Event Configuration.json",
			},
		})
		outputs := make([]pluginmanager.LinkedDataDescription, 1)
		outputs[0] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:      "muncie-White-RS-5696.24.csv",
				Parameter: "flow",
				Format:    "csv",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    outputDestination.Scheme,
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/muncie-White-RS-5696.24.csv",
			},
		}
		payload := pluginmanager.Payload{
			ModelConfiguration: mconfig,
			ModelLinks: pluginmanager.ModelLinks{
				Inputs:  inputs,
				Outputs: outputs,
			},
		}
		return payload
	case "ras-mutator":
		mconfig.Name = "Muncie"
		inputs = make([]pluginmanager.LinkedDataDescription, 2)
		inputs[0] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:   "self",
				Format: "object",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    "s3",
				Authority: inputSource.Authority,
				Fragment:  inputSource.Fragment + "models/Muncie/Muncie.p04.tmp.hdf", //this does not change
			},
		}
		inputs[1] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:   `/Event Conditions/Unsteady/Boundary Conditions/Flow Hydrographs/River: White  Reach: Muncie  RS: 15696.24`,
				Format: "object",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/muncie-White-RS-5696.24.csv",
			},
		}
		outputs := make([]pluginmanager.LinkedDataDescription, 1)
		outputs[0] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:   "self",
				Format: "object",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/Muncie.p04.tmp.hdf",
			},
		}
		payload := pluginmanager.Payload{
			ModelConfiguration: mconfig,
			ModelLinks: pluginmanager.ModelLinks{
				Inputs:  inputs,
				Outputs: outputs,
			},
		}
		return payload
	case "ras-unsteady":
		mconfig.Name = "Muncie"
		inputs = make([]pluginmanager.LinkedDataDescription, 5)
		inputs[0] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:   "Muncie.p04.tmp.hdf",
				Format: "object",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,        //this actually needs to change to output source authority.
				Fragment:  eventParts + "/Muncie.p04.tmp.hdf", //provided by the mutator - changes each event
			},
		}
		inputs[1] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:   "Muncie.b04",
				Format: "object",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    "s3",
				Authority: inputSource.Authority,
				Fragment:  inputSource.Fragment + "models/Muncie/Muncie.b04",
			},
		}
		inputs[2] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:   "Muncie.prj",
				Format: "object",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    "s3",
				Authority: inputSource.Authority,
				Fragment:  inputSource.Fragment + "models/Muncie/Muncie.prj",
			},
		}
		inputs[3] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:   "Muncie.x04",
				Format: "object",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    "s3",
				Authority: inputSource.Authority,
				Fragment:  inputSource.Fragment + "models/Muncie/Muncie.x04",
			},
		}
		inputs[4] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:   "Muncie.c04",
				Format: "object",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    "s3",
				Authority: inputSource.Authority,
				Fragment:  inputSource.Fragment + "models/Muncie/Muncie.c04",
			},
		}
		outputs := make([]pluginmanager.LinkedDataDescription, 3)
		outputs[0] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:   "Muncie.p04.hdf",
				Format: "object",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/Muncie.p04.hdf",
			},
		}
		outputs[1] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:   "Muncie.log",
				Format: "object",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/Muncie.log",
			},
		}
		outputs[2] = pluginmanager.LinkedDataDescription{
			DataDescription: pluginmanager.DataDescription{
				Name:   "Muncie.dss",
				Format: "object",
			},
			ResourceInfo: pluginmanager.ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/Muncie.dss",
			},
		}
		payload := pluginmanager.Payload{
			ModelConfiguration: mconfig,
			ModelLinks: pluginmanager.ModelLinks{
				Inputs:  inputs,
				Outputs: outputs,
			},
		}
		return payload
	}
	payload := pluginmanager.Payload{
		ModelConfiguration: mconfig,
		ModelLinks: pluginmanager.ModelLinks{
			Inputs: inputs,
		},
	}
	return payload
}
