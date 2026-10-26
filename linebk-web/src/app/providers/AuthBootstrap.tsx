import type { FC } from 'react';
import { useEffect, useMemo } from 'react';

import { useAppDispatch } from 'app/store';

import { clearAuth, setAccessToken } from 'entities/auth';
import { AuthService } from 'entities/auth/api/authService';

export const AuthBootstrap: FC = () => {
  const authService = useMemo(() => new AuthService(), []);
  const dispatch = useAppDispatch();

  useEffect(() => {
    const refreshToken = async () => {
      try {
        const token = await authService.refreshToken();
        dispatch(setAccessToken(token));
      } catch (err) {
        console.error('Token refresh error', err);
        dispatch(clearAuth());
      }
    };

    void refreshToken();
  }, [authService, dispatch]);

  return null;
};
