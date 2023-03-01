// Copyright (C) 2022-2023 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Component, OnDestroy, OnInit } from '@angular/core';
import { DataService } from '../services/data.service';
import { PtzService } from '../services/ptz.service';
import { CameraApiService } from '../services/camera-api.service';

@Component({
  selector: 'app-ptz',
  templateUrl: './ptz.component.html',
  styleUrls: ['./ptz.component.css'],
})
export class PtzComponent implements OnInit, OnDestroy {

  constructor(public data: DataService, public ptz: PtzService, public api: CameraApiService) {
  }

  ngOnInit(): void {
  }

  ngOnDestroy(): void {
  }

  isPtzDisabled(): boolean {
    return this.data.cameraFeatures === undefined || this.data.cameraFeatures.PTZ === false;
  }

  isZoomDisabled(): boolean {
    return this.data.cameraFeatures === undefined || this.data.cameraFeatures.Zoom === false;
  }
}
