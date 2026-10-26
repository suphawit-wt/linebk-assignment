import { createSlice, type PayloadAction } from '@reduxjs/toolkit';

interface AuthState {
  accessToken: string;
  isAuthenticated: boolean;
}

const initialState: AuthState = {
  accessToken: '',
  isAuthenticated: false,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setAccessToken(state, action: PayloadAction<string>) {
      state.accessToken = action.payload;
      state.isAuthenticated = true;
    },
    clearAuth(state) {
      state.accessToken = '';
      state.isAuthenticated = false;
    },
  },
});

export const { setAccessToken, clearAuth } = authSlice.actions;
export const authReducer = authSlice.reducer;
