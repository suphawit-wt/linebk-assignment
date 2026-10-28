import type { FC } from 'react';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';

import { MainPage } from 'pages/main';
import { PinPage } from 'pages/pin';
import { SplashPage } from 'pages/splash';

import { AuthBootstrap } from './AuthBootstrap';
import { ProtectedRoute } from './ProtectedRoute';
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

        <Route element={<ProtectedRoute />}>
          <Route path="/main" element={<MainPage />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
};
