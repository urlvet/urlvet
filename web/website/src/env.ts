import { defineEnvVars } from '@sveltejs/kit/env';

export const variables = defineEnvVars({
  // Go API base URL. Read when the app starts, so one build works for any backend.
  PUBLIC_BASE_URL: {
    public: true,
    schema: (value) => value || 'http://localhost:8080/api/v1',
  },
});
