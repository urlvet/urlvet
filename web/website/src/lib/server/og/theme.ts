// The site's palette and type, for drawing outside the browser.

export const OG = {
  width: 1200,
  height: 630,
  bg: '#0f0d0b',
  card: '#1a1714',
  line: '#2a2622',
  text: '#f3ede2',
  muted: '#a69c8c',
  faint: '#7c7366',
  accent: '#8eb0ff',
  serif: 'Instrument Serif',
  sans: 'Geist',
  mono: 'Geist Mono',
};

export const VERDICT_TONE: Record<string, { color: string; glow: string }> = {
  Safe: { color: '#34d399', glow: '16, 185, 129' },
  Suspicious: { color: '#facc15', glow: '234, 179, 8' },
  Risky: { color: '#f87171', glow: '239, 68, 68' },
};
