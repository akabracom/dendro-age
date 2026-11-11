/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface DendroScalePattern {
  id?: number;
  title?: string;
  short_description?: string;
  full_description?: string;
  status?: string;
  image_url?: string;
  video_url?: string;
  event_year?: number | null;
  ring_count?: number | null;
  like_count?: number;
  liked?: number;
  mine?: number;
  creator_id?: number;
}

export interface PublishDendroScalePattern {
  short_description: string;
  full_description?: string;
  event_year: number;
  ring_count: number;
}

export interface LikeDendroScalePattern {
  like: 0 | 1;
}

export type GetDendroScalePatternsData = DendroScalePattern[];

export interface CreateDendroScalePatternPayload {
  title: string;
  /** @format binary */
  image?: File;
  /** @format binary */
  video?: File;
}

export type CreateDendroScalePatternData = DendroScalePattern;

export type GetDendroScalePatternFeedData = DendroScalePattern;

export type GetDendroScalePatternDraftData = DendroScalePattern;

export type PublishDendroScalePatternData = DendroScalePattern;

export type DeleteDendroScalePatternData = any;

export type LikeDendroScalePatternData = DendroScalePattern;
