import type { FC } from 'react';

import { DebitCardResponse } from 'entities/bank';

import styles from './DebitCardSwipe.module.scss';

export interface DebitCardSwipeProps {
  cards: DebitCardResponse[];
}

export const DebitCardSwipe: FC<DebitCardSwipeProps> = ({ cards }) => {
  return (
    <div
      className={`${styles.debitSwipeWrap} ${styles.mainLoading} ${styles.mainLoadingOrder6}`}
    >
      <div className={styles.debitSwipeInner}>
        <div className={styles.debitSwipeLst}>
          {cards.map((card, idx) => {
            const isInProgress = card.status === 'in-progress';
            const isWhiteCard = card.color === '#ffffff';

            return (
              <a
                key={idx}
                href="#card"
                className={styles.debitSwipeItem}
                style={{
                  backgroundColor: card.color,
                  color: isWhiteCard ? '#97999e' : '#fff',
                  ...(card.borderColor
                    ? ({
                        '--border-color': card.borderColor,
                      } as React.CSSProperties)
                    : {}),
                }}
              >
                <strong className={styles.debitSwipeName}>{card.name}</strong>

                {isInProgress ? (
                  <span className={styles.debitSwipeEtc}>In progress</span>
                ) : (
                  <span
                    className={`${styles.debitSwipeEtc} ${styles.debitSwipeEtcActive}`}
                  >
                    <span className={styles.debitSwipeEtcNum}>
                      {card.number || '9440 78•••• •••• 3115'}
                    </span>
                  </span>
                )}

                <span
                  className={styles.debitSwipeIssue}
                  style={{ color: isWhiteCard ? '#d3d3d2' : 'inherit' }}
                >
                  Issued by {card.issuer}
                </span>
              </a>
            );
          })}

          <a
            href="#all"
            className={`${styles.debitSwipeItem} ${styles.debitSwipeItemAll}`}
          >
            See all
          </a>
        </div>
      </div>
    </div>
  );
};
