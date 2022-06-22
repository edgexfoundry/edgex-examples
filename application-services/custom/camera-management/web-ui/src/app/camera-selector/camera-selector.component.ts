// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Component, OnDestroy, OnInit } from '@angular/core';
import { DataService } from "../services/data.service";
import { CameraApiService } from "../services/camera-api.service";
import { EventMqttService } from "../services/event-mqtt.service";

@Component({
  selector: 'app-camera-selector',
  templateUrl: './camera-selector.component.html',
  styleUrls: ['./camera-selector.component.css']
})
export class CameraSelectorComponent implements OnInit, OnDestroy {

  constructor(public data: DataService, public api: CameraApiService, public eventService: EventMqttService) {
  }

  ngOnInit(): void {
    this.api.updateCameraList();
    this.api.updatePipelinesList();
  }

  ngOnDestroy(): void {
  }

  cameraSelectionChanged(value) {
    this.api.updateProfiles(value);
    this.api.refreshPipelineStatus(value, true);
  }

  profileSelectionChanged(value) {
    this.api.updatePresets(this.data.selectedCamera, value);
  }

  startPipeline() {
    let tokens = this.data.selectedPipeline.split('/')
    this.api.startPipeline(this.data.selectedCamera, this.data.selectedProfile, tokens[0], tokens[1]);
  }

  shouldDisablePipeline() {
    return this.data.selectedCamera == undefined
      || this.data.selectedProfile == undefined
      || (this.data.pipelineMap[this.data.selectedCamera]
          && this.data.pipelineMap[this.data.selectedCamera].status.state == 'RUNNING')
  }
}
