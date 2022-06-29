// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

import { HttpHeaders } from "@angular/common/http";

export const ApiLogIgnoreHeader: string = 'X-CameraApp-Ignore';
export const JsonHeaders = {
  headers: new HttpHeaders({
    'Content-Type': 'application/json; charset=utf-8',
  }),
};
