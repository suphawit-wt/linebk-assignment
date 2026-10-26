export const ApiStatus = {
  Success: 'SUCCESS',
  InvalidInput: 'INVALID_INPUT',
  BadRequest: 'BAD_REQUEST',
  Unauthorized: 'UNAUTHORIZED',
  TokenNotMatched: 'TOKEN_NOT_MATCHED',
  ExpiredOTP: 'EXPIRED_OTP',
  Forbidden: 'FORBIDDEN',
  NotFound: 'NOT_FOUND',
  Conflict: 'CONFLICT',
  TooManyRequests: 'TOO_MANY_REQUESTS',
  Unexpected: 'UNEXPECTED',
} as const;

export type ApiStatus = (typeof ApiStatus)[keyof typeof ApiStatus];
