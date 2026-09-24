// Shared button styles so every pill on the site looks the same.
// Tailwind scans .ts files, so these full class strings are picked up at build time.

const BASE =
  'inline-flex items-center justify-center gap-2 rounded-full text-sm font-medium transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 disabled:opacity-60';

/** Solid dark pill (cream in dark mode): the one primary action on a surface. */
export const PILL_SOLID = `${BASE} bg-gray-900 dark:bg-gray-100 text-gray-50 dark:text-gray-900 hover:bg-gray-700 dark:hover:bg-white`;

/** Hairline outline pill: secondary actions. */
export const PILL_OUTLINE = `${BASE} border border-gray-300 dark:border-gray-800 text-gray-700 dark:text-gray-300 hover:border-gray-400 dark:hover:border-gray-600 hover:text-gray-900 dark:hover:text-gray-100`;
