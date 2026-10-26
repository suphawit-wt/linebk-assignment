import { createSelector } from 'reselect';

import type { RootState } from 'app/store/store';

const selectAuth = (state: RootState) => state.auth;

export const selectAccessToken = createSelector(
  selectAuth,
  (auth) => auth.accessToken,
);

export const selectIsAuthenticated = createSelector(
  selectAuth,
  (auth) => auth.isAuthenticated,
);
