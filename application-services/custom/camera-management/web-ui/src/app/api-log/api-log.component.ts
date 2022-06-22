// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Component, OnDestroy, OnInit } from '@angular/core';
import { DataService } from "../services/data.service";

@Component({
  selector: 'app-api-log',
  templateUrl: './api-log.component.html',
  styleUrls: ['./api-log.component.css']
})
export class ApiLogComponent implements OnInit, OnDestroy {
  constructor(public data: DataService) { }

  ngOnInit(): void {
  }
  ngOnDestroy(): void {
  }

}
