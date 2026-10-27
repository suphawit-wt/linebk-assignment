import type { FC } from 'react';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';

import { PinPage } from 'pages/pin';
import { SplashPage } from 'pages/splash';

import { AuthBootstrap } from './AuthBootstrap';
import { PublicRoute } from './PublicRoute';

export const AppRouterProvider: FC = () => {
  return (
    <BrowserRouter>
      <AuthBootstrap />
      <Routes>
        <Route path="/" element={<SplashPage />} />

        <Route element={<PublicRoute />}>
          <Route path="/pin" element={<PinPage />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
};
