// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Component } from '@angular/core';
import { DataService } from '../services/data.service';
import { EventMqttService } from "../services/event-mqtt.service";
import { CameraApiService } from "../services/camera-api.service";

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent {
  private interval: any;

  constructor(public data: DataService, public api: CameraApiService, public eventService: EventMqttService) {
    // the following line allows the navbar to properly highlight which
    // page is active when refreshing the page (i.e. when visiting
    // site.com/<this_page>). Without this, it will default to 'home' as the
    // active page
    this.data.currentPage = 'home';
  }

  ngOnInit(): void {
    this.interval = setInterval(() => {
      this.api.refreshAllPipelineStatuses();
    }, 1000);
  }

  ngOnDestroy(): void {
    clearInterval(this.interval);
  }
}

