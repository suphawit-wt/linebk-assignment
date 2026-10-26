import type { HttpStatusCode } from 'axios';

import type { ApiStatus } from 'shared/constants';

interface ErrResponse {
  field: string;
  message: string;
}

interface PaginationResponse {
  page: number;
  perPage: number;
  totalCount: number;
  hasNext: boolean;
}

export interface ApiResponse {
  status: ApiStatus;
  statusCode: HttpStatusCode;
  message: string;
  errors?: ErrResponse[];
}

export interface ApiResponseWithData<T> {
  status: ApiStatus;
  statusCode: HttpStatusCode;
  message: string;
  data: T;
}

export interface ApiResponseWithDataPagination<T> {
  status: ApiStatus;
  statusCode: HttpStatusCode;
  message: string;
  data: T;
  pagination: PaginationResponse;
}
