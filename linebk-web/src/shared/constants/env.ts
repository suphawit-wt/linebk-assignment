import { createEnv } from '@t3-oss/env-core';
import { z } from 'zod';

export const env = createEnv({
  clientPrefix: 'REACT_APP_',
  client: {
    REACT_APP_API_BASE_URL: z.url(),
  },
  runtimeEnv: {
    REACT_APP_API_BASE_URL: process.env.REACT_APP_API_BASE_URL,
  },
  emptyStringAsUndefined: true,
});
