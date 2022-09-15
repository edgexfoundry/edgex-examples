// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

// This file can be replaced during build by using the `fileReplacements` array.
// `ng build --prod` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

export const environment = {
  production: false,
  appServiceBaseUrl: 'http://localhost:59750/api/v2',
  defaultPipelineId: 'object_detection/person_vehicle_bike',

  mqtt: {
    host: 'localhost',
    port: 59001,
    path: '/mqtt',
    topic: 'incoming/data/edge-video-analytics/inference-event',
  }
};

/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
// import 'zone.js/dist/zone-error';  // Included with Angular CLI.
