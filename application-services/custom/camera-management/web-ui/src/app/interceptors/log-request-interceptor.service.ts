// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Injectable } from '@angular/core';
import {
  HttpErrorResponse,
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
  HttpResponse
} from '@angular/common/http';
import { Observable, tap } from "rxjs";
import { APILogItem, DataService } from "../services/data.service";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ApiLogIgnoreHeader } from "../constants";

@Injectable({
  providedIn: 'root',
})
export class LogRequestInterceptor implements HttpInterceptor {
  constructor(private data: DataService, private snackbar: MatSnackBar) {}

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {

    if (request.headers.has(ApiLogIgnoreHeader)) {
      return next.handle(request);
    }

    let item = new APILogItem(request);
    // this.data.apiLog.push(item);
    this.data.apiLog.unshift(item);
    while (this.data.apiLog.length > 5) {
      this.data.apiLog.pop();
    }
    return next.handle(request).pipe(tap(event => {
      if (event instanceof HttpResponse) {
        item.response = event as HttpResponse<any>;
      }
      return event;
    }, error => {
      if (error instanceof HttpErrorResponse) {
        item.response = error as HttpErrorResponse;
        this.snackbar.open(`${error.status} ${error.statusText}\n${error.error.substring(0, 60)}...`, '', {
          duration: 2500,
          panelClass: ['mat-toolbar', 'mat-warn', 'error-snackbar'],
        });
      }
    }));
  }
}
