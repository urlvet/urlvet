// Small surprises for curious people: a line shown on the result for a few
// special links. Order matters; the first match wins.

import { TAGLINE_LINKS } from '../data/tagline';

type Egg = { test: (u: URL, host: string) => boolean; say: string };

const RICKROLL = 'dQw4w9WgXcQ';

const EGGS: Egg[] = [
  {
    test: (_, h) => h === 'url.vet',
    say: "Checking me? Bold. I'd do the same.",
  },
  {
    test: (u, h) => h === 'github.com' && u.pathname.toLowerCase().startsWith('/urlvet/urlvet'),
    say: "That's where I live. A star wouldn't hurt.",
  },
  {
    test: (_, h) => h === 'abhizaik.com',
    say: "That's the person who built me. Seems trustworthy, but I'm biased.",
  },
  {
    test: (u, h) =>
      ((h === 'youtube.com' || h === 'm.youtube.com') && u.searchParams.get('v') === RICKROLL) ||
      (h === 'youtu.be' && u.pathname === `/${RICKROLL}`),
    say: "Safe to click. Whether you'll want to is another matter. Never gonna give you up.",
  },
  {
    // Lookalikes spelled with Cyrillic letters, like the examples on the home page.
    test: (_, h) => /[Ѐ-ӿ]/.test(h) || h.startsWith('xn--'),
    say: "Spot the difference? There isn't one you can see, and that's the whole trick.",
  },
  {
    test: (_, h) => TAGLINE_LINKS.includes(h),
    say: "That one's from our tagline. We made it up, but someone will register it eventually.",
  },
  {
    test: (_, h) => h === 'paypal.com',
    say: "The real one. The fakes swap in a Cyrillic «а»; this one doesn't.",
  },
  {
    test: (_, h) => h === 'apple.com',
    say: 'The real one. The fakes spell it аррӏе, every letter swapped for a Cyrillic twin.',
  },
  {
    test: (_, h) => h === 'google.com',
    say: 'Ah yes, the website you use to search websites.',
  },
  {
    test: (_, h) => ['example.com', 'example.org', 'example.net'].includes(h),
    say: 'The most innocent website on the internet. It exists so other pages can point at it.',
  },
  {
    test: (_, h) => h === 'wikipedia.org' || h.endsWith('.wikipedia.org'),
    say: 'Citation needed? Not for this one.',
  },
  {
    test: (_, h) => h === 'neverssl.com',
    say: 'A site that skips HTTPS on purpose, so Wi-Fi sign-in pages can load. The one time “not secure” is the point.',
  },
  {
    test: (_, h) => ['1.1.1.1', '1.0.0.1', '8.8.8.8', '8.8.4.4', '9.9.9.9'].includes(h),
    say: "That's a DNS resolver: the internet's phone book, not really a website.",
  },
  {
    test: (_, h) => h.endsWith('.zip') || h.endsWith('.mov'),
    say: "A “.zip” that isn't a file. These are real web addresses now, and scammers noticed.",
  },
];

function parse(raw: string): { url: URL; host: string } | null {
  try {
    const url = new URL(/^https?:\/\//i.test(raw) ? raw : `http://${raw}`);
    const host = url.hostname.toLowerCase().replace(/^www\./, '');
    return { url, host };
  } catch {
    return null;
  }
}

/** A quip for a special link, or null. */
export function easterEgg(raw: string): string | null {
  const p = parse(raw);
  if (!p) return null;
  return EGGS.find((e) => e.test(p.url, p.host))?.say ?? null;
}

/** A friendly reason not to scan addresses that only exist on the visitor's side. */
export function localAddressMessage(raw: string): string | null {
  const p = parse(raw);
  if (!p) return null;
  const h = p.url.hostname.toLowerCase().replace(/^\[|\]$/g, '');
  if (
    h === 'localhost' ||
    h.endsWith('.localhost') ||
    h === '::1' ||
    /^127\./.test(h) ||
    h === '0.0.0.0'
  ) {
    return "That's your own computer. We can't see it from here, and neither can anyone else.";
  }
  if (
    /^10\./.test(h) ||
    /^192\.168\./.test(h) ||
    /^172\.(1[6-9]|2\d|3[01])\./.test(h) ||
    h.endsWith('.local')
  ) {
    return "That's an address on a private network, like a home router. It isn't reachable from the internet, so there's nothing for us to check.";
  }
  return null;
}
