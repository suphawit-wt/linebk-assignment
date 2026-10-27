import type { FC } from 'react';
import { Helmet } from 'react-helmet-async';

import { PinPad } from 'features/auth/ui/pin-pad/PinPad';

import styles from './PinPage.module.scss';

export const PinPage: FC = () => {
  return (
    <>
      <Helmet>
        <title>Sign In Page | LINE BK</title>
      </Helmet>

      <div className={styles.wrap}>
        <main className={styles.container}>
          <PinPad />
        </main>
      </div>
    </>
  );
};
