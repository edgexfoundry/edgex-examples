// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Injectable } from '@angular/core';
import { IMqttMessage, MqttService } from "ngx-mqtt";
import { Subscription } from "rxjs";
import { InferenceEvent } from "./event-mqtt.types";
import { environment } from "../../environments/environment";

@Injectable({
  providedIn: "root",
  deps: [
    MqttService
  ]
})
export class EventMqttService {
  private subscription: Subscription;
  private paused: boolean = true;

  public eventCount: number = 0;
  public lastEvent: InferenceEvent;

  constructor(private _mqttService: MqttService) {
  }


  isPaused(): boolean {
    return this.paused;
  }

  pause() {
    if (!this.paused) {
      this.paused = true;
      this.unsubscribe();
    }
  }

  resume() {
    if (this.paused) {
      this.paused = false;
      this.subscribe();
    }
  }

  private unsubscribe() {
    if (this.subscription) {
      this.subscription.unsubscribe();
      this.subscription = undefined;
    }
  }

  private subscribe() {
    this.subscription = this._mqttService
      .observe(environment.mqtt.topic)
      .subscribe((data: IMqttMessage) => {
        this.lastEvent = JSON.parse(data.payload.toString()) as InferenceEvent;
        this.eventCount++;
      });
  }
}
