import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { useAppDispatch } from 'app/store';

import { AuthService, setAccessToken } from 'entities/auth';

import {
  TooManyRequestsError,
  UnauthorizedError,
} from 'shared/errors/api-error';

import styles from './PinPad.module.scss';

export const PinPad: React.FC = () => {
  const authService = new AuthService();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();

  const [pinDigits, setPinDigits] = useState<string[]>([]);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [attemptsLeft, setAttemptsLeft] = useState<number>(5);

  const handleKeyPress = async (val: string) => {
    if (pinDigits.length < 6) {
      const nextDigits = [...pinDigits, val];
      setPinDigits(nextDigits);

      if (nextDigits.length === 6) {
        try {
          const passcode = nextDigits.join('');
          const accessToken = await authService.signInWithPasscode(
            '000027b0e9a211ef95a30242ac180002',
            passcode,
          );

          dispatch(setAccessToken(accessToken));
          navigate('/main');
        } catch (err: unknown) {
          setPinDigits([]);

          if (err instanceof UnauthorizedError) {
            const nextAttemptCount = Math.max(0, attemptsLeft - 1);
            setAttemptsLeft(nextAttemptCount);

            setErrorMsg(
              `Invalid PIN Code.<br>You have ${nextAttemptCount} attempt${nextAttemptCount === 1 ? '' : 's'} left.`,
            );
          } else if (err instanceof TooManyRequestsError) {
            setErrorMsg('Too many sign-in attempts, please try again later.');
          } else {
            setErrorMsg('An unexpected error occurred. Please try again.');
          }
        }
      }
    }
  };

  const handleDelete = () => {
    setPinDigits((prev) => prev.slice(0, -1));
  };

  return (
    <div className={styles.pin}>
      <div className={styles.pinTop}>
        <span className={styles.pinPhoto}>
          <img
            src="https://dummyimage.com/200x200/999/fff"
            alt="User Profile"
          />
        </span>
        <h1 className={styles.pinName}>Clare</h1>
        <p
          className={styles.pinDsc}
          style={{ visibility: errorMsg ? 'visible' : 'hidden' }}
          dangerouslySetInnerHTML={{ __html: errorMsg || '&nbsp;' }}
        />
        <div className={styles.pinDots}>
          {[0, 1, 2, 3, 4, 5].map((idx) => (
            <span
              key={idx}
              className={`${styles.pinDot} ${idx < pinDigits.length ? styles.isFilled : ''}`}
            />
          ))}
        </div>
      </div>

      <div className={styles.pinBtm}>
        <a
          href="#signin_password"
          className={styles.pinLogin}
          onClick={(e) => {
            e.preventDefault();
            navigate('/signin-password');
          }}
        >
          Login with ID / Password
        </a>
        <span className={styles.pinKb}>Powered by TestLab</span>
        <div className={styles.pinKeys}>
          {['1', '2', '3', '4', '5', '6', '7', '8', '9'].map((digit) => (
            <button
              key={digit}
              type="button"
              className={styles.pinKey}
              onClick={() => handleKeyPress(digit)}
            >
              {digit}
            </button>
          ))}
          <span className={`${styles.pinKey} ${styles.pinKeySpace}`} />
          <button
            type="button"
            className={styles.pinKey}
            onClick={() => handleKeyPress('0')}
          >
            0
          </button>
          <button
            type="button"
            className={styles.pinKey}
            onClick={handleDelete}
          >
            ⌫
          </button>
        </div>
      </div>
    </div>
  );
};
