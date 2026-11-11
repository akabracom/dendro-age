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

import type {
  CreateDendroScalePatternData,
  CreateDendroScalePatternPayload,
  DeleteDendroScalePatternData,
  GetDendroScalePatternDraftData,
  GetDendroScalePatternFeedData,
  GetDendroScalePatternsData,
  LikeDendroScalePattern,
  LikeDendroScalePatternData,
  PublishDendroScalePattern,
  PublishDendroScalePatternData,
} from "./data-contracts.ts";
import { HttpClient, type RequestParams } from "./http-client.ts";

export class DendroScalePatterns<SecurityDataType = unknown> {
  http: HttpClient<SecurityDataType>;

  constructor(http: HttpClient<SecurityDataType>) {
    this.http = http;
  }

  /**
   * No description
   *
   * @tags dendro_scale_patterns
   * @name GetDendroScalePatterns
   * @summary Список опубликованных паттернов
   * @request GET:/api/dendro_scale_patterns
   */
  getDendroScalePatterns = (
    query?: {
      date_from?: string;
      date_to?: string;
      title?: string;
    },
    params: RequestParams = {},
  ) =>
    this.http.request<GetDendroScalePatternsData, any>({
      path: `/api/dendro_scale_patterns`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
  /**
   * No description
   *
   * @tags dendro_scale_patterns
   * @name CreateDendroScalePattern
   * @summary Черновик с фото и видео
   * @request POST:/api/dendro_scale_patterns
   * @secure
   */
  createDendroScalePattern = (
    data: CreateDendroScalePatternPayload,
    params: RequestParams = {},
  ) =>
    this.http.request<CreateDendroScalePatternData, any>({
      path: `/api/dendro_scale_patterns`,
      method: "POST",
      body: data,
      secure: true,
      type: "multipart/form-data",
      format: "json",
      ...params,
    });
  /**
   * No description
   *
   * @tags dendro_scale_patterns
   * @name GetDendroScalePatternFeed
   * @summary Лента без идентификатора
   * @request GET:/api/dendro_scale_patterns/feed
   */
  getDendroScalePatternFeed = (params: RequestParams = {}) =>
    this.http.request<GetDendroScalePatternFeedData, any>({
      path: `/api/dendro_scale_patterns/feed`,
      method: "GET",
      format: "json",
      ...params,
    });
  /**
   * No description
   *
   * @tags dendro_scale_patterns
   * @name GetDendroScalePatternDraft
   * @summary Черновик текущего создателя
   * @request GET:/api/dendro_scale_patterns/draft
   * @secure
   */
  getDendroScalePatternDraft = (params: RequestParams = {}) =>
    this.http.request<GetDendroScalePatternDraftData, any>({
      path: `/api/dendro_scale_patterns/draft`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * No description
   *
   * @tags dendro_scale_patterns
   * @name PublishDendroScalePattern
   * @summary Публикация черновика
   * @request PUT:/api/dendro_scale_patterns/publish
   * @secure
   */
  publishDendroScalePattern = (
    data: PublishDendroScalePattern,
    params: RequestParams = {},
  ) =>
    this.http.request<PublishDendroScalePatternData, any>({
      path: `/api/dendro_scale_patterns/publish`,
      method: "PUT",
      body: data,
      secure: true,
      type: "application/json",
      format: "json",
      ...params,
    });
  /**
   * No description
   *
   * @tags dendro_scale_patterns
   * @name DeleteDendroScalePattern
   * @summary Логическое удаление
   * @request DELETE:/api/dendro_scale_patterns/{id}
   * @secure
   */
  deleteDendroScalePattern = (id: number, params: RequestParams = {}) =>
    this.http.request<DeleteDendroScalePatternData, any>({
      path: `/api/dendro_scale_patterns/${id}`,
      method: "DELETE",
      secure: true,
      ...params,
    });
  /**
   * No description
   *
   * @tags dendro_scale_patterns
   * @name LikeDendroScalePattern
   * @summary Лайк 1 или снятие 0
   * @request POST:/api/dendro_scale_patterns/{id}/like
   * @secure
   */
  likeDendroScalePattern = (
    id: number,
    data: LikeDendroScalePattern,
    params: RequestParams = {},
  ) =>
    this.http.request<LikeDendroScalePatternData, any>({
      path: `/api/dendro_scale_patterns/${id}/like`,
      method: "POST",
      body: data,
      secure: true,
      type: "application/json",
      format: "json",
      ...params,
    });
}
