import type { FC } from 'react';
import { useSelector } from 'react-redux';
import { Navigate, Outlet } from 'react-router-dom';

import { selectIsAuthenticated } from 'entities/auth';

export const PublicRoute: FC = () => {
  const isAuthenticated = useSelector(selectIsAuthenticated);

  if (isAuthenticated) {
    return <Navigate to="/main" replace />;
  }
  return <Outlet />;
};
