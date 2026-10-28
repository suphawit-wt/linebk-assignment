import axios, { type AxiosInstance } from 'axios';

import { AccountModel } from 'entities/account/model/types';

import { ApiStatus } from 'shared/constants';
import {
  NotFoundError,
  UnauthorizedError,
  UnexpectedError,
  ValidationError,
} from 'shared/errors/api-error';
import type { ApiResponse, ApiResponseWithData } from 'shared/types/api';

import type {
  AccountResponse,
  BannerResponse,
  DebitCardResponse,
  RecentTransactionResponse,
  UserInfoResponse,
} from './bankService.types';

interface BankServiceInterface {
  getUserInfo: () => Promise<UserInfoResponse>;
  getAccounts: () => Promise<AccountModel[]>;
  getRecentTransactions: () => Promise<RecentTransactionResponse[]>;
  getDebitCards: () => Promise<DebitCardResponse[]>;
  getBanners: () => Promise<BannerResponse[]>;
}

export class BankService implements BankServiceInterface {
  constructor(private readonly axiosClient: AxiosInstance) {}

  async getUserInfo(): Promise<UserInfoResponse> {
    try {
      const { data } =
        await this.axiosClient.get<ApiResponseWithData<UserInfoResponse>>(
          `v1/user`,
        );
      return data.data;
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const { status, message, errors } = error.response.data as ApiResponse;

        if (status === ApiStatus.InvalidInput && errors) {
          throw new ValidationError(message, errors);
        }

        switch (status) {
          case ApiStatus.Unauthorized:
            throw new UnauthorizedError(message);
          case ApiStatus.NotFound:
            throw new NotFoundError(message);
          default:
            throw new UnexpectedError();
        }
      }

      throw new UnexpectedError();
    }
  }

  async getAccounts(): Promise<AccountModel[]> {
    try {
      const { data } =
        await this.axiosClient.get<ApiResponseWithData<AccountResponse[]>>(
          `v1/accounts`,
        );

      return data.data.map((account) => ({
        ...account,
        amount: parseFloat(account.amount) || 0,
      }));
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const { status, message, errors } = error.response.data as ApiResponse;

        if (status === ApiStatus.InvalidInput && errors) {
          throw new ValidationError(message, errors);
        }

        switch (status) {
          case ApiStatus.Unauthorized:
            throw new UnauthorizedError(message);
          case ApiStatus.NotFound:
            throw new NotFoundError(message);
          default:
            throw new UnexpectedError();
        }
      }

      throw new UnexpectedError();
    }
  }

  async getRecentTransactions(): Promise<RecentTransactionResponse[]> {
    try {
      const { data } = await this.axiosClient.get<
        ApiResponseWithData<RecentTransactionResponse[]>
      >(`v1/transactions/recent`);
      return data.data;
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const { status, message, errors } = error.response.data as ApiResponse;

        if (status === ApiStatus.InvalidInput && errors) {
          throw new ValidationError(message, errors);
        }

        switch (status) {
          case ApiStatus.Unauthorized:
            throw new UnauthorizedError(message);
          case ApiStatus.NotFound:
            throw new NotFoundError(message);
          default:
            throw new UnexpectedError();
        }
      }

      throw new UnexpectedError();
    }
  }

  async getDebitCards(): Promise<DebitCardResponse[]> {
    try {
      const { data } =
        await this.axiosClient.get<ApiResponseWithData<DebitCardResponse[]>>(
          `v1/debit-cards`,
        );
      return data.data;
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const { status, message, errors } = error.response.data as ApiResponse;

        if (status === ApiStatus.InvalidInput && errors) {
          throw new ValidationError(message, errors);
        }

        switch (status) {
          case ApiStatus.Unauthorized:
            throw new UnauthorizedError(message);
          case ApiStatus.NotFound:
            throw new NotFoundError(message);
          default:
            throw new UnexpectedError();
        }
      }

      throw new UnexpectedError();
    }
  }

  async getBanners(): Promise<BannerResponse[]> {
    try {
      const { data } =
        await this.axiosClient.get<ApiResponseWithData<BannerResponse[]>>(
          `v1/banners`,
        );
      return data.data;
    } catch (error: unknown) {
      if (axios.isAxiosError(error) && error.response) {
        const { status, message, errors } = error.response.data as ApiResponse;

        if (status === ApiStatus.InvalidInput && errors) {
          throw new ValidationError(message, errors);
        }

        switch (status) {
          case ApiStatus.Unauthorized:
            throw new UnauthorizedError(message);
          case ApiStatus.NotFound:
            throw new NotFoundError(message);
          default:
            throw new UnexpectedError();
        }
      }

      throw new UnexpectedError();
    }
  }
}
