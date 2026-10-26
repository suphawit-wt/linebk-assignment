import type { FC } from 'react';

import styles from './SplashScreen.module.scss';

export const SplashScreen: FC = () => {
  return (
    <div className={styles.wrap}>
      <div className={styles.splash}>
        <div className={styles.loader} />
      </div>
    </div>
  );
};
