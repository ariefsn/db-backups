import type { Handle } from '@sveltejs/kit';
import { handleProxy } from 'sveltekit-proxy';

const apiPath = '/api';

export const handle: Handle = async ({ event, resolve }) => {
  if (event.url.pathname.startsWith(apiPath)) {
    return handleProxy({
      target: process.env.API_URL || 'http://localhost:8080',
      rewrite: (path) => path.replace(apiPath, ''),
      onResponse: ({ response, duration }) => {
        console.log(`[Proxy] ${event.request.method} ${event.url.pathname} -> ${response.status} (${duration.toFixed(2)}ms)`);
      },
      onError: ({ error, request }) => {
        console.error('[Proxy Error]', error, request.url);
      }
    })({ event, resolve });
  }

  return resolve(event);
};
