// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Component, OnDestroy, OnInit } from '@angular/core';
import { DataService } from "../services/data.service";
import { EventMqttService } from "../services/event-mqtt.service";


@Component({
  selector: 'app-inference-events',
  templateUrl: './inference-events.component.html',
  styleUrls: ['./inference-events.component.css']
})
export class InferenceEventsComponent implements OnInit, OnDestroy {
  constructor(public data: DataService, public eventService: EventMqttService) { }

  ngOnInit(): void {
  }

  ngOnDestroy(): void {
  }

  labelToIcon(label: string) {
    switch (label) {
      case 'vehicle':
        return 'directions_car'
      case 'bicycle':
        return 'pedal_bike'
      default:
        return label
    }
  }
}
