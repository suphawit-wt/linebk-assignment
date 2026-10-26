import type { FC } from 'react';
import { useSelector } from 'react-redux';
import { Navigate, Outlet } from 'react-router-dom';

import { selectIsAuthenticated } from 'entities/auth';

export const ProtectedRoute: FC = () => {
  const isAuthenticated = useSelector(selectIsAuthenticated);

  if (!isAuthenticated) {
    return <Navigate to="/pin" replace />;
  }
  return <Outlet />;
};
