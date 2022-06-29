// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

export interface BoundingBox {
  x_max: number;
  x_min: number;
  y_max: number;
  y_min: number;
}

export interface Detection {
  bounding_box: BoundingBox;
  confidence: number;
  label: string;
  label_id: number;
}

export interface EventObject {
  detection: Detection;
  h: number;
  roi_type: string;
  w: number;
  x: number;
  y: number;
}

export interface Resolution {
  height: number;
  width: number;
}

export interface InferenceEvent {
  objects: EventObject[];
  resolution: Resolution;
  source: string;
  timestamp: number;
}
