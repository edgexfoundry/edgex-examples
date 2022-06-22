// Copyright (C) 2022 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

export interface GetProfilesResponse {
  Profiles?: (ProfilesEntity)[] | null;
}
export interface ProfilesEntity {
  AudioEncoderConfiguration: AudioEncoderConfiguration;
  AudioSourceConfiguration: AudioSourceConfiguration;
  Extension?: null;
  Fixed: boolean;
  MetadataConfiguration?: null;
  Name: string;
  PTZConfiguration: PTZConfiguration;
  Token: string;
  VideoAnalyticsConfiguration: VideoAnalyticsConfiguration;
  VideoEncoderConfiguration: VideoEncoderConfiguration;
  VideoSourceConfiguration: VideoSourceConfiguration;
}
export interface AudioEncoderConfiguration {
  Bitrate: number;
  Encoding: string;
  Multicast: Multicast;
  Name: string;
  SampleRate: number;
  SessionTimeout: string;
  Token: string;
  UseCount: number;
}
export interface Multicast {
  Address: Address;
  AutoStart: boolean;
  Port: number;
  TTL: number;
}
export interface Address {
  IPv4Address: string;
  Type: string;
}
export interface AudioSourceConfiguration {
  Name: string;
  SourceToken: string;
  Token: string;
  UseCount: number;
}
export interface PTZConfiguration {
  DefaultAbsolutePantTiltPositionSpace: string;
  DefaultContinuousPanTiltVelocitySpace: string;
  DefaultPTZSpeed: DefaultPTZSpeed;
  DefaultPTZTimeout: string;
  DefaultRelativePanTiltTranslationSpace: string;
  PanTiltLimits: PanTiltLimits;
  Token: string;
}
export interface DefaultPTZSpeed {
}
export interface PanTiltLimits {
  Range: Range;
}
export interface Range {
  URI: string;
  XRange: XRangeOrYRange;
  YRange: XRangeOrYRange;
}
export interface XRangeOrYRange {
  Max: number;
  Min: number;
}
export interface VideoAnalyticsConfiguration {
  AnalyticsEngineConfiguration: AnalyticsEngineConfiguration;
  Name: string;
  RuleEngineConfiguration: RuleEngineConfiguration;
  Token: string;
  UseCount: number;
}
export interface AnalyticsEngineConfiguration {
  AnalyticsModule?: (AnalyticsModuleEntity)[] | null;
}
export interface AnalyticsModuleEntity {
  Name: string;
  Parameters: Parameters;
  Type: string;
}
export interface Parameters {
  ElementItem?: (ElementItemEntity)[] | null;
  SimpleItem?: (SimpleItemEntity)[] | null;
}
export interface ElementItemEntity {
  Name: string;
}
export interface SimpleItemEntity {
  Name: string;
  Value: string;
}
export interface RuleEngineConfiguration {
  Rule: Rule;
}
export interface Rule {
  Name: string;
  Parameters: Parameters1;
  Type: string;
}
export interface Parameters1 {
  SimpleItem?: (SimpleItemEntity)[] | null;
}
export interface VideoEncoderConfiguration {
  Encoding: string;
  H264: H264;
  Multicast: Multicast;
  Name: string;
  Quality: number;
  RateControl: RateControl;
  Resolution: Resolution;
  SessionTimeout: string;
  Token: string;
  UseCount: number;
}
export interface H264 {
  GovLength: number;
  H264Profile: string;
}
export interface RateControl {
  BitrateLimit: number;
  EncodingInterval: number;
  FrameRateLimit: number;
}
export interface Resolution {
  Height: number;
  Width: number;
}
export interface VideoSourceConfiguration {
  Bounds: Bounds;
  Extension?: null;
  Name: string;
  SourceToken: string;
  Token: string;
  UseCount: number;
  ViewMode: string;
}
export interface Bounds {
  Height: number;
  Width: number;
  X: number;
  Y: number;
}

export interface PanTilt {
  Space: string;
  X: number;
  Y: number;
}

export interface PTZPosition {
  PanTilt: PanTilt;
}
export interface Preset {
  Name: string;
  PTZPosition: PTZPosition;
  Token: string;
}
export interface GetPresetsResponse {
  Preset: Preset[];
}

export interface AutoEvent {
  interval: string;
  onChange: boolean;
  sourceName: string;
}

export interface Other {
  Address: string;
  Port: string;
  Protocol: string;
}

export interface Onvif {
  Address: string;
  AuthMode: string;
  Port: string;
  SecretPath: string;
}

export interface Protocols {
  other: Other;
  Onvif: Onvif;
}

export interface Device {
  created: any;
  modified: any;
  id: string;
  name: string;
  description: string;
  adminState: string;
  operatingState: string;
  labels: string[];
  serviceName: string;
  profileName: string;
  autoEvents: AutoEvent[];
  protocols: Protocols;
}
