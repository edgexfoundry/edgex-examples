// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { environment } from '../../environments/environment';
import { Device, GetPresetsResponse, GetProfilesResponse } from "./camera-api.types";
import { DataService } from "./data.service";
import { Pipeline, PipelineInfoStatus, PipelineStatus, StartPipelineRequest } from "./pipeline-api.types";
import { ApiLogIgnoreHeader, JsonHeaders } from "../constants";

@Injectable({
  providedIn: 'root',
})
export class CameraApiService {
  camerasUrl = `${environment.appServiceBaseUrl}/cameras`;

  constructor(private httpClient: HttpClient, private data: DataService) {
  }

  makeCameraUrl(cameraName: string, path: string) {
    return `${this.camerasUrl}/${cameraName}${path}`
  }

  makePipelineUrl(cameraName: string, actionPath: string) {
    return this.makeCameraUrl(cameraName, `/pipeline${actionPath}`);
  }

  makeProfileUrl(cameraName: string, profileToken: string, path: string) {
    return this.makeCameraUrl(cameraName, `/profiles/${profileToken}`) + path
  }

  makePresetUrl(cameraName: string, profileToken: string, presetToken: string) {
    return this.makeProfileUrl(cameraName, profileToken, `/presets/`) + presetToken
  }

  updateCameraList() {
    this.data.cameras = undefined;
    this.httpClient
      .get<Device[]>(this.camerasUrl)
      .subscribe({
        next: data => {
          this.data.cameras = data;
          if (data.length > 0) {
            this.data.selectedCamera = data[0].name;
            this.updateProfiles(this.data.selectedCamera);
            this.refreshPipelineStatus(this.data.selectedCamera, true);
          } else {
            this.data.selectedCamera = undefined;
          }
        }, error: _ => {
          this.data.cameras = undefined;
        }
      });
  }

  updateProfiles(cameraName: string) {
    this.data.selectedProfile = undefined;
    this.data.profiles = undefined;

    this.httpClient.get<GetProfilesResponse>(
      this.makeCameraUrl(cameraName, '/profiles'))
      .subscribe({
        next: data => {
          this.data.profiles = data.Profiles;
          this.data.selectedProfile = data.Profiles.length > 0 ? data.Profiles[0].Token : undefined;
          this.updatePresets(cameraName, this.data.selectedProfile);
        }, error: _ => {
          this.data.selectedProfile = undefined;
          this.data.profiles = undefined;
        }
      });
  }

  updatePresets(cameraName: string, profileToken: string) {
    this.data.presets = undefined;
    this.httpClient.get<GetPresetsResponse>(
      this.makeProfileUrl(cameraName, profileToken, '/presets'))
      .subscribe({
        next: data => {
          this.data.presets = data.Preset;
        }, error: _ => {
          this.data.presets = undefined;
        }
      });
  }

  gotoPreset(cameraName: string, profileToken: string, presetToken: string) {
    return this.httpClient.post<any>(
      this.makePresetUrl(cameraName, profileToken, presetToken), '')
      .subscribe();
  }

  startPipeline(cameraName: string, profileToken: string, name: string, version: string) {
    let url = this.makePipelineUrl(cameraName, '/start');
    let req = new StartPipelineRequest(profileToken, name, version);
    this.httpClient.post<any>(url, JSON.stringify(req), JsonHeaders).subscribe(_ => {
      // todo: do at an interval?
      this.refreshPipelineStatus(this.data.selectedCamera, true);
    });
  }

  stopPipeline(cameraName: string) {
    let url = this.makePipelineUrl(cameraName, '/stop');
    this.httpClient.post<any>(url, '').subscribe(_ => {
      this.refreshPipelineStatus(this.data.selectedCamera, true);
    });
  }

  refreshPipelineStatus(cameraName: string, doNotLog: boolean) {
    let url = this.makePipelineUrl(cameraName, '/status');
    let headers = {};
    if (doNotLog) {
      // we use this as a hack to tell the http-interceptor that we do not want to
      // log these requests. this way the log is not cluttered by tons of
      // pipeline status requests
      headers[ApiLogIgnoreHeader] = 'true';
    }
    this.httpClient.get<PipelineStatus>(url, {
      headers: new HttpHeaders(headers)
    }).subscribe({
      next: status => {
        this.data.pipelineStatus = status;
      }, error: _ => {
        this.data.pipelineStatus = undefined;
      }
    });
  }

  refreshAllPipelineStatuses() {
    let url = environment.appServiceBaseUrl + "/pipelines/status/all";
    let headers = {};
    headers[ApiLogIgnoreHeader] = 'true';
    this.httpClient.get<Map<string, PipelineInfoStatus>>(url, {
      headers: new HttpHeaders(headers)
    }).subscribe({
      next: statuses => {
        if (statuses == null) {
          this.data.pipelineMap.clear();
          return;
        }

        // delete any keys that are no longer valid
        for (const key of this.data.pipelineMap.keys()) {
          if (!statuses.hasOwnProperty(key)) {
            console.log('deleting key:', key);
            this.data.pipelineMap.delete(key);
          }
        }
        // add/update existing and new keys
        for (const [key, value] of Object.entries(statuses)) {
          let status = this.data.pipelineMap.get(key)
          if (status === undefined) {
            // new value, so just add it
            this.data.pipelineMap.set(key, <PipelineInfoStatus>value);
          } else {
            // existing value, so update it
            status.status = value.status
            status.info = value.info
          }
        }
      }, error: _ => {
        this.data.pipelineMap.clear();
      }
    });
  }

  updatePipelinesList() {
    this.data.pipelines = undefined;
    this.httpClient.get<Pipeline[]>(environment.appServiceBaseUrl + "/pipelines")
      .subscribe({
        next: data => {
          this.data.pipelines = data.filter(p => {
            if (this.shouldShowPipeline(p)) {
              if (`${p.name}/${p.version}` === environment.defaultPipelineId) {
                this.data.selectedPipeline = environment.defaultPipelineId;
              }
              return true;
            }
            return false;
          });
        }, error: _ => {
          this.data.pipelines = undefined;
        }
      });
  }

  // filter out non-working or non-interesting pipelines
  shouldShowPipeline(p: Pipeline): boolean {
    return p.name !== 'audio_detection'
      && p.name !== 'video_decode'
      && p.version !== 'app_src_dst'
      && p.version !== 'object_zone_count'
      && p.version !== 'object_line_crossing';
  }
}
