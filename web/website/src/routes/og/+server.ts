import type { RequestHandler } from '@sveltejs/kit';
import satori from 'satori';
import sharp from 'sharp';
import { resultCard, siteCard } from '$lib/server/og/cards';
import { loadFonts } from '$lib/server/og/fonts';
import { OG } from '$lib/server/og/theme';

// Social preview image.
//   /og                              → site card
//   /og?domain=example.com&v=Safe&s=92 → scan result card (v and s optional)
export const GET: RequestHandler = async ({ url }) => {
  const domain = url.searchParams.get('domain')?.trim().slice(0, 253);
  const verdict = url.searchParams.get('v') ?? 'Suspicious';
  const s = url.searchParams.get('s');
  const score = s !== null && /^\d{1,3}$/.test(s) ? Number(s) : null;

  const card = domain ? resultCard(domain, verdict, score) : siteCard();
  const fonts = await loadFonts();
  // satori's element type is React-shaped; our plain objects match it structurally.
  const svg = await satori(card as Parameters<typeof satori>[0], {
    width: OG.width,
    height: OG.height,
    fonts,
  });
  const png = await sharp(Buffer.from(svg)).png().toBuffer();

  return new Response(new Uint8Array(png), {
    headers: {
      'Content-Type': 'image/png',
      'Cache-Control': 'public, max-age=3600',
    },
  });
};
