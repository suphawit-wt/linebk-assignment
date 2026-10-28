import React, { useEffect, useMemo, useState } from 'react';
import { Helmet } from 'react-helmet-async';

import { Header } from 'widgets/header';
import { SplashScreen } from 'widgets/splash/ui/SplashScreen';

import { AccountModel } from 'entities/account/model/types';
import {
  BankService,
  BannerResponse,
  DebitCardResponse,
  RecentTransactionResponse,
  UserInfoResponse,
} from 'entities/bank';

import { axiosClient } from 'shared/api';

import styles from './MainPage.module.scss';

interface DashboardState {
  userInfo: UserInfoResponse | null;
  accounts: AccountModel[];
  transactions: RecentTransactionResponse[];
  debitCards: DebitCardResponse[];
  banners: BannerResponse[];
  isLoading: boolean;
  error: string | null;
}

export const MainPage: React.FC = () => {
  const bankService = useMemo(() => new BankService(axiosClient), []);
  const [state, setState] = useState<DashboardState>({
    userInfo: null,
    accounts: [],
    transactions: [],
    debitCards: [],
    banners: [],
    isLoading: true,
    error: null,
  });

  useEffect(() => {
    let isMounted = true;

    const fetchDashboardData = async () => {
      try {
        const [userInfo, accounts, transactions, debitCards, banners] =
          await Promise.all([
            bankService.getUserInfo(),
            bankService.getAccounts(),
            bankService.getRecentTransactions(),
            bankService.getDebitCards(),
            bankService.getBanners(),
          ]);

        if (isMounted) {
          setState({
            userInfo,
            accounts: Array.isArray(accounts) ? accounts : [],
            transactions: Array.isArray(transactions) ? transactions : [],
            debitCards: Array.isArray(debitCards) ? debitCards : [],
            banners: Array.isArray(banners) ? banners : [],
            isLoading: false,
            error: null,
          });
        }
      } catch (err) {
        if (isMounted) {
          console.error('Failed to load dashboard data from API specs', err);
          setState((prev) => ({
            ...prev,
            isLoading: false,
            error: 'Failed to load banking dashboard. Please try again.',
          }));
        }
      }
    };

    fetchDashboardData();

    return () => {
      isMounted = false;
    };
  }, [bankService]);

  if (state.isLoading) {
    return <SplashScreen />;
  }

  if (state.error) {
    return (
      <div className={styles.wrap}>
        <Header />
        <div style={{ padding: '40px', textAlign: 'center', color: 'red' }}>
          {state.error}
        </div>
      </div>
    );
  }

  return (
    <>
      <Helmet>
        <title>Main Page | LINE BK</title>
      </Helmet>

      <div className={styles.wrap}>
        <Header />
      </div>
    </>
  );
};
