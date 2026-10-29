import type { FC } from 'react';

import { AccountModel } from '../model/types';

import styles from './Account.module.scss';

export interface AccountProps {
  account: AccountModel;
  className?: string;
}

const formatAccountTypeName = (typeSlug: string): string => {
  const dictionary: Record<string, string> = {
    'saving-account': 'Saving Account',
    'credit-loan': 'Credit Loan',
    'goal-saving-account': 'Travel New York',
  };

  return (
    dictionary[typeSlug] ||
    typeSlug
      .split('-')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')
  );
};

const formatCurrencySymbol = (currencyCode: string): string => {
  return currencyCode === 'THB' ? '฿' : currencyCode;
};

export const Account: FC<AccountProps> = ({ account, className = '' }) => {
  const {
    type,
    amount,
    currency,
    accountNumber,
    issuer,
    color,
    isMainAccount,
    progress,
    flags,
  } = account;

  const displayTitle = formatAccountTypeName(type);
  const displayCurrency = formatCurrencySymbol(currency || 'THB');
  const formattedAmount = amount.toLocaleString(undefined, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });

  return (
    <div
      className={`${styles.mainAcc} ${
        isMainAccount ? styles.mainAccLarge : styles.mainAccIsSmall
      } ${className}`}
      style={{ backgroundColor: color }}
    >
      <div className={styles.mainAccTop}>
        <h2 className={styles.mainAccName}>{displayTitle}</h2>
        <span className={styles.mainAccAmount}>
          {displayCurrency}
          {formattedAmount}
        </span>

        {flags &&
          flags.map((flag, idx) => (
            <span
              key={idx}
              className={`${styles.mainAccFlag} ${
                flag.toLowerCase() === 'disbursement'
                  ? styles.mainAccFlagLock
                  : ''
              }`}
            >
              {flag}
            </span>
          ))}

        <span
          className={`${styles.mainAccDetail} ${
            !isMainAccount ? styles.mainAccDetailNum : ''
          }`}
        >
          {isMainAccount
            ? `Smart account ${accountNumber}`
            : `${displayTitle} ${accountNumber}`}
        </span>
        {issuer && (
          <span className={styles.mainAccDetail}>Powered by {issuer}</span>
        )}
      </div>

      {isMainAccount && (
        <div className={styles.mainAccBottom}>
          <div className={styles.mainAccLinkBox}>
            <div className={styles.mainAccLinkItem}>
              <a
                href="#withdrawal"
                className={`${styles.mainAccLink} ${styles.mainAccLinkWithdrawal}`}
              >
                Withdrawal
              </a>
              <a
                href="#qr"
                className={`${styles.mainAccLink} ${styles.mainAccLinkQr}`}
              >
                QR scan
              </a>
              <a
                href="#addmoney"
                className={`${styles.mainAccLink} ${styles.mainAccLinkAddmoney}`}
              >
                Add money
              </a>
            </div>
          </div>
        </div>
      )}

      {progress !== undefined && (
        <div className={styles.mainAccCircle}>
          <svg width="100%" height="100%" viewBox="0 0 42 42">
            <circle
              cx="21"
              cy="21"
              r="15.91549430918954"
              fill="transparent"
              stroke="rgba(0,0,0,0.07)"
              strokeWidth="1.5"
            />
            <circle
              cx="21"
              cy="21"
              r="15.91549430918954"
              fill="transparent"
              stroke="#fff"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeDashoffset="25"
              style={{ strokeDasharray: `${progress} ${100 - progress}` }}
            />
          </svg>
          <div className={styles.mainAccNum}>
            <span className={styles.percent}>{progress}</span>
            <span className={styles.unit}>%</span>
          </div>
        </div>
      )}

      <button
        type="button"
        className={`${styles.mainAccMore} ${
          !isMainAccount ? styles.mainAccMoreSmall : ''
        }`}
      >
        <span className={styles.blind}>More Action</span>
      </button>
    </div>
  );
};
