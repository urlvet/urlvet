// Everything Vetty says, in one place. Tone: the casual url.vet voice,
// with a wink at the old office assistant.
import type { AnalyzeResult } from '../../types';

export const NAME = 'Vetty';

export const INTRO = [
  `Hi, I'm ${NAME}! I'm suspicious of everything. Professionally.`,
  'Got a link from someone? Let me squint at it before you click.',
];

/** Said when someone keeps tapping Vetty, one line per extra tap. */
export const POKES = [
  'Yes?',
  "I'm right here.",
  'Do you touch all your paperclips like this?',
  'Okay, that tickles.',
];

/** Heading of the "More" list in Vetty's menu. */
export const MORE_TITLE = "Here's everything else I can do.";

export const NUDGE = "It looks like you're about to check a link. I can help with that.";
export const NUDGE_CTA = 'See what I can do';

/** Greeting when reopening Vetty, by page. */
export const GREETINGS: Record<string, string[]> = {
  home: [
    "It looks like you're checking a link. Would you like help?",
    'Back again? Good. Suspicious links hate this one trick: checking them.',
    "Need a hand? I don't have many, but they're both yours.",
  ],
  result: [
    'Want me to walk you through this one?',
    'I read the fine print so you do not have to. Ask away.',
  ],
  scanning: ['Hold on, I am reading the fine print…'],
  error: ["That one didn't go through. Happens to the best of us."],
  'how-it-works': ['Reading the manual? I respect that.', 'Ah, a fellow fan of documentation.'],
  about: ["Getting to know us? I'm the friendly one.", 'Yes, I also work here.'],
  other: ['Lost? I know a good place to check links.'],
};

export const TIPS = [
  'Tip: press / anywhere to jump to the search bar.',
  'Tip: shortened links hide where they really go. Paste them here first.',
  'Tip: a page can look exactly like your bank and still not be your bank. Check the address.',
  'Tip: every result has a share link, so you can send it to whoever sent you the link.',
];

/** Score bands, matching the backend thresholds. */
export const SCORE_BANDS = [
  {
    verdict: 'Risky',
    from: 0,
    to: 29,
    dot: 'bg-red-500',
    text: 'text-red-600 dark:text-red-400',
    meaning: 'Strong signs of phishing or a scam. Stay away.',
  },
  {
    verdict: 'Suspicious',
    from: 30,
    to: 64,
    dot: 'bg-yellow-500',
    text: 'text-yellow-700 dark:text-yellow-400',
    meaning: 'Mixed signals. Be careful with logins and payments.',
  },
  {
    verdict: 'Safe',
    from: 65,
    to: 100,
    dot: 'bg-emerald-500',
    text: 'text-emerald-700 dark:text-emerald-400',
    meaning: 'Nothing worrying turned up.',
  },
];

export const SCORES_FOOTNOTE =
  'I can be wrong, so every result has a Report button. A human reads every report.';

export const SCAN_FAILED = [
  'Scans usually fail for one of three reasons:',
  '• a typo in the address,',
  '• the site is down,',
  '• or it blocks automated visitors (some banks do).',
  'Double-check the link and try again.',
];

export const HIDE_BYE = "Fine. I'll be down in the footer. I've been retired before.";

/** "Who is Vetty?": what he is, what he can do, and why he's here. */
export const MEET = {
  lead: `I'm ${NAME}, url.vet's helper.`,
  why: "Checking a link shouldn't need a security degree. I'm here so anyone can use url.vet: I explain results in plain words and show you where everything is.",
  canTitle: 'What I can do',
  can: [
    'Check a link for you, or a whole message with links in it',
    'Explain a result in plain words',
    'Tell you what the score means',
    'Walk you through the page',
    'Help you report a result that looks wrong',
  ],
  note: "I only know about links you check here, and I forget them when you leave. Hide me any time; I'll wait in the footer.",
};

/** A message to send back to whoever shared the link, worded for the verdict. */
export function warnMessage(verdict: string, domain: string, score: number, link: string): string {
  if (verdict === 'Risky')
    return `I checked this link on url.vet and it looks like a scam (${domain}, ${score}/100). Please don't open it or enter any details there.\n\nWhy: ${link}`;
  if (verdict === 'Suspicious')
    return `I checked this link on url.vet. ${domain} looks suspicious (${score}/100), so better not to enter any details there.\n\nWhy: ${link}`;
  return `I checked this link on url.vet and ${domain} looks safe (${score}/100).\n\nDetails: ${link}`;
}

/** Said once, unprompted, when a result comes back Risky. */
export const RISKY_ALERT = 'This one looks dangerous. Want me to warn whoever sent it?';
export const RISKY_ALERT_CTA = 'Warn them';

export const WARN_INTRO: Record<string, string> = {
  Risky: "Let's warn whoever sent it. Edit the message if you like:",
  Suspicious: 'Let them know before they open it. Edit the message if you like:',
  Safe: 'Let them know it checked out. Edit the message if you like:',
};

export const REAL_SITE = {
  brand: (d: string) => `This page is pretending to be ${d}.`,
  lookalike: (d: string) => `This address is made to look like ${d}.`,
  advice: (d: string) =>
    `If you meant to visit ${d}, don't use the link you were sent. Go to the real site instead, or check it first.`,
};

export const CLEAN_LINK = {
  lead: 'This link carries tracking tags. They tell the sender who clicked and where from.',
  done: 'Same page, without the tags:',
};

export const CHECK_FOR_ME =
  "Paste the link, or the whole message it came in. I'll put it in the search bar and check it.";

export const TRY_EXAMPLE =
  'Pick one. The green ones are legit links and red ones are fakes spelled with Cyrillic letters.';

/** Plain-language explanation of a scan result. */
export function explain(result: AnalyzeResult): {
  lead: string;
  lists: { title: string; items: string[] }[];
  advice: string;
} {
  const verdict = result.result?.verdict ?? 'Suspicious';
  const score = result.result?.final_score ?? 0;
  const bad = (result.result?.reasons?.bad_reasons ?? []).slice(0, 3);
  const good = (result.result?.reasons?.good_reasons ?? []).slice(0, 3);
  const domain = result.domain;

  if (verdict === 'Safe') {
    return {
      lead: `Good news: ${domain} looks legit (${score}/100).`,
      lists: [
        { title: "Why I'm relaxed", items: good },
        { title: 'Small things I noticed', items: bad },
      ].filter((l) => l.items.length),
      advice:
        "Still, if this link came in an unexpected message, it's safer to type the website's address yourself than to click it.",
    };
  }
  if (verdict === 'Risky') {
    return {
      lead: `I'd stay away from ${domain} (${score}/100).`,
      lists: [{ title: 'What worried me', items: bad }].filter((l) => l.items.length),
      advice:
        "Don't enter anything there. If you think I'm wrong, report it and a human will check.",
    };
  }
  return {
    lead: `I'm on the fence about ${domain} (${score}/100).`,
    lists: [
      { title: 'What worried me', items: bad },
      { title: 'In its favour', items: good },
    ].filter((l) => l.items.length),
    advice:
      "If this page asks you to log in or pay, don't do it here. Type the website's address into your browser yourself and log in from there.",
  };
}

export const pick = <T>(list: T[]): T => list[Math.floor(Math.random() * list.length)];
