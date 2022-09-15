// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

export const environment = {
  production: true,
  appServiceBaseUrl: '/api/v2',
  defaultPipelineId: 'object_detection/person_vehicle_bike',

  mqtt: {
    host: null,
    port: 59001,
    path: '/mqtt',
    topic: 'incoming/data/edge-video-analytics/inference-event',
  }
};
