import { combineReducers } from '@reduxjs/toolkit';

import { authReducer } from 'entities/auth/model/authSlice';

export const rootReducer = combineReducers({
  auth: authReducer,
});
