// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ApiLogComponent } from './api-log.component';

describe('ApiLogComponent', () => {
  let component: ApiLogComponent;
  let fixture: ComponentFixture<ApiLogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ApiLogComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ApiLogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
