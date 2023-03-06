// Copyright (C) 2022-2023 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

export interface PipelineStatus {
  avg_fps: number;
  elapsed_time: number;
  id: string;
  start_time: number;
  state: string;
}

export interface PipelineInfo {
  // id is the instance_id assigned by the pipeline server once it is started
  id: string;
  // name is the first part of the pipeline's full name. In the case of 'object_detection/person_vehicle_bike'
  // the name is 'object_detection'
  name: string;
  // version is the second part of the pipeline's full name. In the case of 'object_detection/person_vehicle_bike'
  // the version is 'person_vehicle_bike'
  version: string;
  // profile is the ProfileToken for the specific stream
  profile: string;
}

export interface PipelineInfoStatus {
  camera: string;
  info: PipelineInfo;
  status: PipelineStatus;
}

export interface Pipeline {
  description: string;
  name: string;
  type: string;
  version: string;
}

export interface OnvifConfig {
  profile_token: string;
}

export interface USBConfig {
  InputFps?: string;
  InputImageSize?: string;
  InputPixelFormat?: string;
  OutputFrames?: string;
  OutputFps?: string;
  OutputImageSize?: string;
  OutputAspect?: string;
  OutputVideoCodec?: string;
  OutputVideoQuality?: string;
}

export class StartPipelineRequest {
  onvif?: OnvifConfig;
  usb?: USBConfig;
  pipeline_name: string;
  pipeline_version: string;

  static forUSB(pipeline_name: string, pipeline_version: string, usb: USBConfig): StartPipelineRequest {
    const req = new StartPipelineRequest();
    req.usb = usb;
    req.pipeline_name = pipeline_name;
    req.pipeline_version = pipeline_version;
    return req;
  }

  static forOnvif(pipeline_name: string, pipeline_version: string, onvif: OnvifConfig): StartPipelineRequest {
    const req = new StartPipelineRequest();
    req.onvif = onvif;
    req.pipeline_name = pipeline_name;
    req.pipeline_version = pipeline_version;
    return req;
  }
}
