// Loads the site's fonts (Instrument Serif, Geist, Geist Mono) for satori.
// Without a browser User-Agent, the Google Fonts CSS API serves plain TTF files, which
// satori can read (it cannot read WOFF2). Fonts are fetched once and kept in memory.

type Weight = 400 | 500 | 600;
export type OgFont = {
  name: string;
  data: ArrayBuffer;
  weight: Weight;
  style: 'normal' | 'italic';
};

const CSS_URL =
  'https://fonts.googleapis.com/css2?family=Instrument+Serif:ital@0;1&family=Geist:wght@400;500&family=Geist+Mono:wght@400;500';

let cache: Promise<OgFont[]> | null = null;

export function loadFonts(): Promise<OgFont[]> {
  cache ??= fetchFonts().catch((err) => {
    cache = null; // retry on the next request
    throw err;
  });
  return cache;
}

async function fetchFonts(): Promise<OgFont[]> {
  const css = await (await fetch(CSS_URL)).text();
  const faces = [...css.matchAll(/@font-face\s*{([^}]*)}/g)].map((m) => m[1]);

  return Promise.all(
    faces.map(async (face) => {
      const name = /font-family:\s*'([^']+)'/.exec(face)?.[1] ?? 'Geist';
      const style = /font-style:\s*italic/.test(face) ? 'italic' : 'normal';
      const weight = Number(/font-weight:\s*(\d+)/.exec(face)?.[1] ?? 400) as Weight;
      const url = /url\(([^)]+)\)/.exec(face)?.[1];
      if (!url) throw new Error(`no font url for ${name}`);
      const data = await (await fetch(url)).arrayBuffer();
      return { name, data, weight, style } satisfies OgFont;
    })
  );
}
