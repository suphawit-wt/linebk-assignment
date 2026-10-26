import axios from 'axios';

import { ApiStatus, env } from 'shared/constants';
import {
  BadRequestError,
  TooManyRequestsError,
  UnauthorizedError,
  UnexpectedError,
  ValidationError,
} from 'shared/errors/api-error';
import type { ApiResponse, ApiResponseWithData } from 'shared/types/api';

import type { TokenResponse } from './authService.types';

const API_PREFIX = 'auth';

const apiClient = axios.create({
  baseURL: env.REACT_APP_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 2000,
  withCredentials: true,
});

interface AuthServiceInterface {
  signInWithPasscode: (id: string, passcode: string) => Promise<string>;
  signInWithPassword: (id: string, password: string) => Promise<string>;
  refreshToken: () => Promise<string>;
}

export class AuthService implements AuthServiceInterface {
  async signInWithPasscode(id: string, passcode: string): Promise<string> {
    try {
      const { data } = await apiClient.post<ApiResponseWithData<TokenResponse>>(
        `v1/${API_PREFIX}/signin/passcode`,
        { id, passcode },
        { withCredentials: true },
      );
      return data.data.accessToken;
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const { status, message, errors } = error.response.data as ApiResponse;

        if (status === ApiStatus.InvalidInput && errors) {
          throw new ValidationError(message, errors);
        }

        switch (status) {
          case ApiStatus.BadRequest:
            throw new BadRequestError(message);
          case ApiStatus.TooManyRequests:
            throw new TooManyRequestsError(message);
          case ApiStatus.Unauthorized:
            throw new UnauthorizedError(message);
          default:
            throw new UnexpectedError();
        }
      }

      throw new UnexpectedError();
    }
  }

  async signInWithPassword(id: string, password: string): Promise<string> {
    try {
      const { data } = await apiClient.post<ApiResponseWithData<TokenResponse>>(
        `v1/${API_PREFIX}/signin/password`,
        { id, password },
        { withCredentials: true },
      );
      return data.data.accessToken;
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const { status, message, errors } = error.response.data as ApiResponse;

        if (status === ApiStatus.InvalidInput && errors) {
          throw new ValidationError(message, errors);
        }

        switch (status) {
          case ApiStatus.BadRequest:
            throw new BadRequestError(message);
          case ApiStatus.TooManyRequests:
            throw new TooManyRequestsError(message);
          case ApiStatus.Unauthorized:
            throw new UnauthorizedError(message);
          default:
            throw new UnexpectedError();
        }
      }

      throw new UnexpectedError();
    }
  }

  async refreshToken(): Promise<string> {
    try {
      const { data } = await apiClient.post<ApiResponseWithData<TokenResponse>>(
        `v1/${API_PREFIX}/refresh`,
        {},
        { withCredentials: true },
      );
      return data.data.accessToken;
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const { status, message } = error.response.data as ApiResponse;

        switch (status) {
          case ApiStatus.BadRequest:
            throw new BadRequestError(message);
          case ApiStatus.Unauthorized:
            throw new UnauthorizedError(message);
          default:
            throw new UnexpectedError();
        }
      }

      throw new UnexpectedError();
    }
  }
}
