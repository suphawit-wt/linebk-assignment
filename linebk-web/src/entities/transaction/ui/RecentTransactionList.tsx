import type { FC } from 'react';

import { RecentTransactionResponse } from 'entities/bank';

import styles from './RecentTransactionList.module.scss';

export interface RecentTransactionListProps {
  transactions: RecentTransactionResponse[];
  className?: string;
}

export const RecentTransactionList: FC<RecentTransactionListProps> = ({
  transactions,
  className = '',
}) => {
  return (
    <div
      className={`${styles.rctlyWrap} ${styles.mainLoading} ${styles.mainLoadingOrder5} ${className}`}
    >
      <ul className={styles.rctlyLst}>
        {transactions.map((tx, idx) => (
          <li key={idx} className={styles.rctlyItem}>
            <a href="#recent" className={styles.rctlyLink}>
              <span
                className={`${styles.rctlyThumb} ${tx.isBank ? styles.isBank : ''}`}
              >
                <img src={tx.image} alt={tx.name} />
              </span>
              <span className={styles.rctlyName}>{tx.name}</span>
            </a>
          </li>
        ))}
      </ul>
    </div>
  );
};
