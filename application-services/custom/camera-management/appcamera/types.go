//
// Copyright (C) 2022 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

// Note: The code in this file was created from actual JSON payloads, using 1 or more of the
//       many JSON -> Go struct converters available.

type StreamUriRequest struct {
	StreamSetup  StreamSetup `json:"StreamSetup"`
	ProfileToken string      `json:"ProfileToken"`
}
type Transport struct {
	Protocol string `json:"Protocol"`
}
type StreamSetup struct {
	Stream    string    `json:"Stream"`
	Transport Transport `json:"Transport"`
}

type ProfilesResponse struct {
	Profiles []struct {
		AudioEncoderConfiguration struct {
			Bitrate   int    `json:"Bitrate"`
			Encoding  string `json:"Encoding"`
			Multicast struct {
				Address struct {
					IPv4Address string `json:"IPv4Address"`
					Type        string `json:"Type"`
				} `json:"Address"`
				AutoStart bool `json:"AutoStart"`
				Port      int  `json:"Port"`
				TTL       int  `json:"TTL"`
			} `json:"Multicast"`
			Name           string `json:"Name"`
			SampleRate     int    `json:"SampleRate"`
			SessionTimeout string `json:"SessionTimeout"`
			Token          string `json:"Token"`
			UseCount       int    `json:"UseCount"`
		} `json:"AudioEncoderConfiguration"`
		AudioSourceConfiguration struct {
			Name        string `json:"Name"`
			SourceToken string `json:"SourceToken"`
			Token       string `json:"Token"`
			UseCount    int    `json:"UseCount"`
		} `json:"AudioSourceConfiguration"`
		Extension             interface{} `json:"Extension"`
		Fixed                 bool        `json:"Fixed"`
		MetadataConfiguration interface{} `json:"MetadataConfiguration"`
		Name                  string      `json:"Name"`
		PTZConfiguration      struct {
			DefaultAbsolutePantTiltPositionSpace  string `json:"DefaultAbsolutePantTiltPositionSpace"`
			DefaultContinuousPanTiltVelocitySpace string `json:"DefaultContinuousPanTiltVelocitySpace"`
			DefaultPTZSpeed                       struct {
			} `json:"DefaultPTZSpeed"`
			DefaultPTZTimeout                      string `json:"DefaultPTZTimeout"`
			DefaultRelativePanTiltTranslationSpace string `json:"DefaultRelativePanTiltTranslationSpace"`
			PanTiltLimits                          struct {
				Range struct {
					URI    string `json:"URI"`
					XRange struct {
						Max int `json:"Max"`
						Min int `json:"Min"`
					} `json:"XRange"`
					YRange struct {
						Max int `json:"Max"`
						Min int `json:"Min"`
					} `json:"YRange"`
				} `json:"Range"`
			} `json:"PanTiltLimits"`
			Token string `json:"Token"`
		} `json:"PTZConfiguration"`
		Token                       string `json:"Token"`
		VideoAnalyticsConfiguration struct {
			AnalyticsEngineConfiguration struct {
				AnalyticsModule []struct {
					Name       string `json:"Name"`
					Parameters struct {
						ElementItem []struct {
							Name string `json:"Name"`
						} `json:"ElementItem"`
						SimpleItem []struct {
							Name  string `json:"Name"`
							Value string `json:"Value"`
						} `json:"SimpleItem"`
					} `json:"Parameters"`
					Type string `json:"Type"`
				} `json:"AnalyticsModule"`
			} `json:"AnalyticsEngineConfiguration"`
			Name                    string `json:"Name"`
			RuleEngineConfiguration struct {
				Rule struct {
					Name       string `json:"Name"`
					Parameters struct {
						SimpleItem []struct {
							Name  string `json:"Name"`
							Value string `json:"Value"`
						} `json:"SimpleItem"`
					} `json:"Parameters"`
					Type string `json:"Type"`
				} `json:"Rule"`
			} `json:"RuleEngineConfiguration"`
			Token    string `json:"Token"`
			UseCount int    `json:"UseCount"`
		} `json:"VideoAnalyticsConfiguration"`
		VideoEncoderConfiguration struct {
			Encoding string `json:"Encoding"`
			H264     struct {
				GovLength   int    `json:"GovLength"`
				H264Profile string `json:"H264Profile"`
			} `json:"H264"`
			Multicast struct {
				Address struct {
					IPv4Address string `json:"IPv4Address"`
					Type        string `json:"Type"`
				} `json:"Address"`
				AutoStart bool `json:"AutoStart"`
				Port      int  `json:"Port"`
				TTL       int  `json:"TTL"`
			} `json:"Multicast"`
			Name        string `json:"Name"`
			Quality     int    `json:"Quality"`
			RateControl struct {
				BitrateLimit     int `json:"BitrateLimit"`
				EncodingInterval int `json:"EncodingInterval"`
				FrameRateLimit   int `json:"FrameRateLimit"`
			} `json:"RateControl"`
			Resolution struct {
				Height int `json:"Height"`
				Width  int `json:"Width"`
			} `json:"Resolution"`
			SessionTimeout string `json:"SessionTimeout"`
			Token          string `json:"Token"`
			UseCount       int    `json:"UseCount"`
		} `json:"VideoEncoderConfiguration"`
		VideoSourceConfiguration struct {
			Bounds struct {
				Height int `json:"Height"`
				Width  int `json:"Width"`
				X      int `json:"X"`
				Y      int `json:"Y"`
			} `json:"Bounds"`
			Extension   interface{} `json:"Extension"`
			Name        string      `json:"Name"`
			SourceToken string      `json:"SourceToken"`
			Token       string      `json:"Token"`
			UseCount    int         `json:"UseCount"`
			ViewMode    string      `json:"ViewMode"`
		} `json:"VideoSourceConfiguration"`
	} `json:"Profiles"`
}

type StreamUriResponse struct {
	MediaURI struct {
		InvalidAfterConnect bool   `json:"InvalidAfterConnect"`
		InvalidAfterReboot  bool   `json:"InvalidAfterReboot"`
		Timeout             string `json:"Timeout"`
		URI                 string `json:"Uri"`
	} `json:"MediaUri"`
}

type PipelineRequest struct {
	Source      Source      `json:"source"`
	Destination Destination `json:"destination"`
}
type Source struct {
	URI  string `json:"uri"`
	Type string `json:"type"`
}
type Metadata struct {
	Type  string `json:"type"`
	Host  string `json:"host"`
	Topic string `json:"topic"`
}
type Frame struct {
	Type string `json:"type"`
	Path string `json:"path"`
}
type Destination struct {
	Metadata Metadata `json:"metadata"`
	Frame    Frame    `json:"frame"`
}

type GetPresetsResponse struct {
	Preset []struct {
		Name        string `json:"Name"`
		PTZPosition struct {
			PanTilt struct {
				Space string `json:"Space"`
				X     int    `json:"X"`
				Y     int    `json:"Y"`
			} `json:"PanTilt"`
		} `json:"PTZPosition"`
		Token string `json:"Token"`
	} `json:"Preset"`
}

type StartPipelineRequest struct {
	ProfileToken    string `json:"profile_token"`
	PipelineName    string `json:"pipeline_name"`
	PipelineVersion string `json:"pipeline_version"`
}
