import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { SplashScreen } from 'widgets/splash/ui/SplashScreen';

export const SplashPage: React.FC = () => {
  const navigate = useNavigate();

  useEffect(() => {
    const timer = setTimeout(() => {
      navigate('/pin');
    }, 300);
    return () => clearTimeout(timer);
  }, [navigate]);

  return <SplashScreen />;
};
