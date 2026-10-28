import type { FC } from 'react';

import styles from './Header.module.scss';

export const Header: FC = () => {
  return (
    <header className={styles.header}>
      <a href="#" className={`${styles.headerLft} ${styles.headerMenu}`}>
        <span className={`${styles.blind}`}>Menu</span>
      </a>
      <button
        type="button"
        className={`${styles.headerRgt} ${styles.headerCxl}`}
      >
        <span className={`${styles.blind}`}>Cancel</span>
      </button>
    </header>
  );
};
