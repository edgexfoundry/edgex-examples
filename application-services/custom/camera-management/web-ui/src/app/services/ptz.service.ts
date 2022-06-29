// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { DataService } from './data.service';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class PtzService {
  makePtzUrl(cameraName: string, profileToken: string, actionPath: string) {
    return `${environment.appServiceBaseUrl}/cameras/${cameraName}/profiles/${profileToken}/ptz${actionPath}`
  }

  constructor(private httpClient: HttpClient, public data: DataService) {
  }

  post(cameraName: string, profileToken: string, actionPath: string) {
    let url = this.makePtzUrl(cameraName, profileToken, actionPath);
    return this.httpClient.post<any>(url, '')
      .subscribe();
  }

  downLeft(cameraName: string, profileToken: string) {
    return this.post(cameraName, profileToken, '/down-left');
  }

  down(cameraName: string, profileToken: string) {
    return this.post(cameraName, profileToken, '/down');
  }

  downRight(cameraName: string, profileToken: string) {
    return this.post(cameraName, profileToken, '/down-right');
  }

  upLeft(cameraName: string, profileToken: string) {
    return this.post(cameraName, profileToken, '/up-left');
  }

  up(cameraName: string, profileToken: string) {
    return this.post(cameraName, profileToken, '/up');
  }

  upRight(cameraName: string, profileToken: string) {
    return this.post(cameraName, profileToken, '/up-right');
  }

  left(cameraName: string, profileToken: string) {
    return this.post(cameraName, profileToken, '/left');
  }

  right(cameraName: string, profileToken: string) {
    return this.post(cameraName, profileToken, '/right');
  }

  home(cameraName: string, profileToken: string) {
    return this.post(cameraName, profileToken, '/home');
  }

  zoomIn(cameraName: string, profileToken: string) {
    return this.post(cameraName, profileToken, '/zoom-in');
  }

  zoomOut(cameraName: string, profileToken: string) {
    return this.post(cameraName, profileToken, '/zoom-out');
  }
}

