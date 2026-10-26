import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios';

import { store } from 'app/store';

import { clearAuth, setAccessToken } from 'entities/auth';
import { AuthService } from 'entities/auth/api/authService';

import { env } from '../constants/env';

export const axiosClient = axios.create({
  baseURL: env.REACT_APP_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 5000,
});

axiosClient.interceptors.request.use(
  (config) => {
    const state = store.getState();
    const accessToken = state.auth.accessToken;

    if (accessToken !== '') {
      config.headers.Authorization = `Bearer ${accessToken}`;
    }
    return config;
  },
  () => Promise.reject(new Error('Something went wrong!')),
);

axiosClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const authService = new AuthService();
        const newAccessToken: string = await authService.refreshToken();
        store.dispatch(setAccessToken(newAccessToken));
        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        return axiosClient(originalRequest);
      } catch (refreshError: unknown) {
        store.dispatch(clearAuth());

        if (typeof window !== 'undefined') {
          window.location.reload();
        }

        return Promise.reject(new Error(String(refreshError)));
      }
    }

    return Promise.reject(error);
  },
);
