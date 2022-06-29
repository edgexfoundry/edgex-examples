// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Component, OnInit } from '@angular/core';
import { DataService } from "../services/data.service";
import { CameraApiService } from "../services/camera-api.service";
import { ChangeDetectorRef } from '@angular/core';

@Component({
  selector: 'app-all-pipelines',
  templateUrl: './all-pipelines.component.html',
  styleUrls: ['./all-pipelines.component.css']
})
export class AllPipelinesComponent implements OnInit {

  constructor(public data: DataService, public api: CameraApiService, private ref: ChangeDetectorRef) {
  }

  ngOnInit(): void {
  }

  ngOnDestroy(): void {
  }
}
