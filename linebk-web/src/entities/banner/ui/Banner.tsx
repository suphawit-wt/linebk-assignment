import type { FC } from 'react';

import { BannerResponse } from 'entities/bank';

import styles from './Banner.module.scss';

export interface BannerProps {
  banner: BannerResponse;
  className?: string;
  onClick?: (banner: BannerResponse) => void;
}

export const Banner: FC<BannerProps> = ({
  banner,
  className = '',
  onClick,
}) => {
  const { title, description, image } = banner;

  return (
    <a
      href="#banner"
      className={`${styles.mainProd} ${className}`}
      onClick={(e) => {
        e.preventDefault();
        onClick?.(banner);
      }}
    >
      <span className={styles.mainProdCmsIco}>
        <img src={image} alt="" />
      </span>
      <strong className={styles.mainProdTit}>{title}</strong>
      <p className={styles.mainProdDsc}>{description}</p>
    </a>
  );
};
