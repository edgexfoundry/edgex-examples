// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PtzComponent } from './ptz.component';

describe('PtzComponent', () => {
  let component: PtzComponent;
  let fixture: ComponentFixture<PtzComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PtzComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PtzComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
