// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { Component } from '@angular/core';
import { BreakpointObserver } from "@angular/cdk/layout";
import { DataService } from "./services/data.service";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  public title = 'Camera Management UI';
  public navbarCollapseShow = false;
  darkTheme: boolean;

  constructor(public data: DataService, private bo: BreakpointObserver) {
    bo.observe('(prefers-color-scheme: dark)').subscribe((state) => {
      this.useDarkTheme = state.matches;
    });
  }

  get useDarkTheme() {
    return this.darkTheme;
  }

  set useDarkTheme(v: boolean) {
    this.darkTheme = v;
    if (v) {
      document.body.classList.add('dark-theme');
    } else {
      document.body.classList.remove('dark-theme');
  }
  }
}
