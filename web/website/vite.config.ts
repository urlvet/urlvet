import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import adapterAuto from '@sveltejs/adapter-auto';
import adapterStatic from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit({
      preprocess: vitePreprocess(),
      // Docker builds set DOCKER_BUILD=1 → adapter-static (served by Caddy).
      // All other environments (Vercel, local dev) → adapter-auto.
      adapter:
        process.env.DOCKER_BUILD === '1' ? adapterStatic({ fallback: '200.html' }) : adapterAuto(),
    }),
  ],
  test: {
    environment: 'jsdom',
    setupFiles: ['src/test/setup.ts'],
    include: ['src/**/*.{test,spec}.{js,ts}'],
  },
});
