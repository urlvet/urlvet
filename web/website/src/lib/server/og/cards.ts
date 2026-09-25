// Social preview cards (1200x630), drawn with satori in the site's editorial style.
import { verdictCopy } from '../../verdict';
import { el, h, type Node } from './h';
import { OG, VERDICT_TONE } from './theme';

const R = 36;
const CIRC = 2 * Math.PI * R;

/** "● url.vet" mark used in the top-left of every card. */
function brand(): Node {
  return h(
    'div',
    {
      alignItems: 'center',
      gap: 14,
      fontFamily: OG.sans,
      fontSize: 30,
      fontWeight: 500,
      color: OG.text,
    },
    h('div', { width: 12, height: 12, borderRadius: 999, background: OG.accent }),
    'url.vet'
  );
}

/** Hand-drawn underline, like the one under "url.vet" in the site tagline. */
function swash(width: number): Node {
  return el(
    'svg',
    {
      width,
      height: 12,
      viewBox: '0 0 100 8',
      preserveAspectRatio: 'none',
      style: { position: 'absolute', left: 0, bottom: -6 },
    },
    el('path', {
      d: 'M1 5.5C20 2.5 45 1.8 99 4.2',
      fill: 'none',
      stroke: OG.accent,
      strokeWidth: 2.2,
      strokeLinecap: 'round',
    })
  );
}

function glow(rgb: string, opacity: number): Node {
  return h('div', {
    position: 'absolute',
    top: -260,
    left: 200,
    width: 800,
    height: 560,
    backgroundImage: `radial-gradient(ellipse at center, rgba(${rgb}, ${opacity}), rgba(${rgb}, 0) 70%)`,
  });
}

function footer(): Node {
  return h(
    'div',
    {
      justifyContent: 'space-between',
      alignItems: 'center',
      borderTop: `1px solid ${OG.line}`,
      paddingTop: 28,
    },
    h(
      'div',
      { fontFamily: OG.serif, fontStyle: 'italic', fontSize: 32, color: OG.muted },
      'sketchy link? just url.vet it'
    ),
    h('div', { fontFamily: OG.mono, fontSize: 20, color: OG.faint, letterSpacing: 1 }, 'url.vet')
  );
}

function domainSize(domain: string): number {
  if (domain.length <= 16) return 72;
  if (domain.length <= 24) return 60;
  if (domain.length <= 34) return 48;
  return 40;
}

function scoreRing(score: number, color: string): Node {
  const offset = CIRC - (Math.max(0, Math.min(100, score)) / 100) * CIRC;
  return h(
    'div',
    {
      position: 'relative',
      width: 230,
      height: 230,
      alignItems: 'center',
      justifyContent: 'center',
    },
    el(
      'svg',
      {
        width: 230,
        height: 230,
        viewBox: '0 0 88 88',
        style: { position: 'absolute', top: 0, left: 0 },
      },
      el('circle', { cx: 44, cy: 44, r: R, fill: 'none', stroke: OG.line, strokeWidth: 4 }),
      el('circle', {
        cx: 44,
        cy: 44,
        r: R,
        fill: 'none',
        stroke: color,
        strokeWidth: 4,
        strokeLinecap: 'round',
        strokeDasharray: `${CIRC}`,
        strokeDashoffset: `${offset}`,
        transform: 'rotate(-90 44 44)',
      })
    ),
    h(
      'div',
      { flexDirection: 'column', alignItems: 'center' },
      h(
        'div',
        { fontFamily: OG.sans, fontSize: 72, fontWeight: 500, color, lineHeight: 1 },
        String(score)
      ),
      h(
        'div',
        { fontFamily: OG.mono, fontSize: 18, color: OG.faint, marginTop: 8, letterSpacing: 2 },
        'TRUST SCORE'
      )
    )
  );
}

/** Card for a shared scan result. */
export function resultCard(domain: string, verdict: string, score: number | null): Node {
  const tone = VERDICT_TONE[verdict] ?? VERDICT_TONE.Suspicious;
  const copy = verdictCopy(verdict);
  const quip = copy.quip.replace(/ \S+$/u, ''); // satori can't draw the emoji

  return h(
    'div',
    {
      position: 'relative',
      width: OG.width,
      height: OG.height,
      background: OG.bg,
      flexDirection: 'column',
      justifyContent: 'space-between',
      padding: '56px 72px',
    },
    glow(tone.glow, 0.16),
    h(
      'div',
      { justifyContent: 'space-between', alignItems: 'center' },
      brand(),
      h(
        'div',
        { fontFamily: OG.mono, fontSize: 20, color: OG.faint, letterSpacing: 3 },
        'SCAN RESULT'
      )
    ),
    h(
      'div',
      { justifyContent: 'space-between', alignItems: 'center', gap: 40 },
      h(
        'div',
        { flexDirection: 'column', flex: 1, minWidth: 0 },
        h(
          'div',
          {
            fontFamily: OG.mono,
            fontSize: domainSize(domain) * 0.5,
            color: OG.muted,
            maxWidth: 760,
            overflow: 'hidden',
            whiteSpace: 'nowrap',
            textOverflow: 'ellipsis',
          },
          domain
        ),
        h(
          'div',
          { alignItems: 'center', gap: 24, marginTop: 6 },
          h(
            'div',
            {
              fontFamily: OG.serif,
              fontSize: verdict.length > 6 ? 136 : 168,
              color: OG.text,
              lineHeight: 1.05,
              letterSpacing: -2,
            },
            verdict
          ),
          h(
            'div',
            {
              fontFamily: OG.mono,
              fontSize: 20,
              letterSpacing: 2,
              color: tone.color,
              border: `1.5px solid ${tone.color}`,
              borderRadius: 999,
              padding: '8px 18px',
              whiteSpace: 'nowrap',
              flexShrink: 0,
              marginTop: 24,
            },
            copy.label.toUpperCase()
          )
        ),
        h('div', { fontFamily: OG.sans, fontSize: 32, color: OG.muted, marginTop: 4 }, quip)
      ),
      score === null ? null : scoreRing(score, tone.color)
    ),
    footer()
  );
}

/** Card for the site itself (landing page shares). */
export function siteCard(): Node {
  return h(
    'div',
    {
      position: 'relative',
      width: OG.width,
      height: OG.height,
      background: OG.bg,
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
    },
    glow('214, 170, 120', 0.14),
    h(
      'div',
      {
        fontFamily: OG.serif,
        fontSize: 210,
        lineHeight: 1,
        letterSpacing: -4,
        alignItems: 'baseline',
      },
      h('div', { color: OG.faint }, 'url'),
      h('div', { color: OG.accent }, '.'),
      h('div', { color: OG.text }, 'vet')
    ),
    h(
      'div',
      {
        fontFamily: OG.serif,
        fontStyle: 'italic',
        fontSize: 52,
        color: OG.text,
        marginTop: 28,
        alignItems: 'baseline',
        gap: 14,
      },
      h('div', {}, 'sketchy link? just'),
      h('div', { position: 'relative' }, 'url.vet', swash(124)),
      h('div', {}, 'it')
    ),
    h(
      'div',
      { fontFamily: OG.sans, fontSize: 28, color: OG.muted, marginTop: 30 },
      'Vet any URL for phishing, scams & suspicious redirects'
    ),
    h(
      'div',
      {
        gap: 36,
        marginTop: 56,
        fontFamily: OG.mono,
        fontSize: 20,
        color: OG.faint,
        letterSpacing: 1,
      },
      ...['OPEN SOURCE', 'NO SIGNUP', 'EXPLAINS EVERY VERDICT'].map((t) =>
        h(
          'div',
          { alignItems: 'center', gap: 12 },
          h('div', { width: 6, height: 6, borderRadius: 999, background: OG.accent }),
          t
        )
      )
    )
  );
}
