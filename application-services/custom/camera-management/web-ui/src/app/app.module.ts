// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { NgModule } from '@angular/core';

import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';

import { DragDropModule } from "@angular/cdk/drag-drop";

import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import { MatNativeDateModule } from '@angular/material/core';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatRadioModule } from '@angular/material/radio';
import { MatSelectModule } from '@angular/material/select';
import { MatSliderModule } from '@angular/material/slider';
import { MatSortModule } from "@angular/material/sort";
import { MatTableModule } from '@angular/material/table';
import { MatTabsModule } from '@angular/material/tabs';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatSnackBarModule } from '@angular/material/snack-bar';

import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { MatTooltipModule } from "@angular/material/tooltip";
import { PtzComponent } from './ptz/ptz.component';
import { ApiLogComponent } from './api-log/api-log.component';
import { InferenceEventsComponent } from './inference-events/inference-events.component';
import { CameraSelectorComponent } from './camera-selector/camera-selector.component';
import { LogRequestInterceptor } from "./interceptors/log-request-interceptor.service";
import { MatExpansionModule } from "@angular/material/expansion";
import { IMqttServiceOptions, MqttModule } from "ngx-mqtt";
import { environment } from "../environments/environment";
import { AllPipelinesComponent } from './all-pipelines/all-pipelines.component';

const MQTT_SERVICE_OPTIONS: IMqttServiceOptions = {
  hostname: window.location.hostname,
  port: environment.mqtt.port,
  path: environment.mqtt.path
};

@NgModule({
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    DragDropModule,
    MatButtonModule,
    MatCardModule,
    MatCheckboxModule,
    MatChipsModule,
    MatDatepickerModule,
    MatDialogModule,
    MatFormFieldModule,
    MatGridListModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatNativeDateModule,
    MatPaginatorModule,
    MatProgressSpinnerModule,
    MatRadioModule,
    MatSelectModule,
    MatSliderModule,
    MatSortModule,
    MatTableModule,
    MatTabsModule,
    MatToolbarModule,
    MatSnackBarModule,
    AppRoutingModule,
    MatTooltipModule,
    MatExpansionModule,
    MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
  ],
  declarations: [
    AppComponent,
    HomeComponent,
    PtzComponent,
    ApiLogComponent,
    InferenceEventsComponent,
    CameraSelectorComponent,
    AllPipelinesComponent,
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: LogRequestInterceptor,
      multi: true
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
