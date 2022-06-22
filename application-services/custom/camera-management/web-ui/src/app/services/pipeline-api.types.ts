// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

export interface PipelineStatus {
  avg_fps: number;
  elapsed_time: number;
  id: number;
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
  info: PipelineInfo
  status: PipelineStatus
}

export interface Pipeline {
  description: string;
  name: string;
  type: string;
  version: string;
}

export class StartPipelineRequest {
  profile_token: string
  pipeline_name: string
  pipeline_version: string

  constructor(profile_token: string, pipeline_name: string, pipeline_version: string) {
    this.profile_token = profile_token;
    this.pipeline_name = pipeline_name;
    this.pipeline_version = pipeline_version;
  }
}
