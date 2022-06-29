// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Injectable } from '@angular/core';
import { HttpRequest, HttpResponseBase } from "@angular/common/http";
import { Device, Preset, ProfilesEntity } from "./camera-api.types";
import { Pipeline, PipelineInfoStatus, PipelineStatus } from "./pipeline-api.types";

/**
 * Represents a page that users can navigate to, either via
 * a URL route, the navbar, or from other component interactions.
 *
 * @interface
 */
export interface Page {
  page: 'home';
  caption: string;
}

export class APILogItem {
  request: HttpRequest<any>;
  response: HttpResponseBase;
  constructor(request: HttpRequest<any>) {
    this.request = request;
  }
}

@Injectable({
  providedIn: 'root',
})
export class DataService {
  public cameras: Device[];
  public selectedCamera: string;

  public profiles: ProfilesEntity[];
  public selectedProfile: string;

  public pipelineStatus: PipelineStatus;

  public pipelineMap: Map<string, PipelineInfoStatus>;

  public pipelines: Pipeline[];
  public selectedPipeline: string;

  public presets: Preset[];

  public apiLog: APILogItem[];

  // pages is a list of all tabs that are navigable by the user, accessible
  // via the routing module and also by clicking tabs
  public pages: Page[] = [
    { page: 'home', caption: 'Home' },
  ];

  // currentPage is set whenever the user taps/clicks on a tab or navigates
  // to a page via the routing module. Its value is important because it is
  // visually bound to highlighting the currently viewed tab
  public currentPage: 'home';

  constructor() {
    this.apiLog = new Array<APILogItem>();
    this.pipelineMap = new Map<string, PipelineInfoStatus>();
  }
}
