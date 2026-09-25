// Guided tour steps. Each points at an element tagged data-guide="<target>";
// steps whose target isn't on the page right now are skipped.

/** landing: before a scan. results: once a result is on screen. */
export type TourStep = { target: string; title: string; text: string; on: 'landing' | 'results' };

export const TOUR_STEPS: TourStep[] = [
  {
    target: 'search',
    on: 'landing',
    title: 'Paste the link here',
    text: 'Full link, shortened link, or just a domain like example.com. I handle the formatting.',
  },
  {
    target: 'paste',
    on: 'landing',
    title: 'Or let me paste it',
    text: 'One click pastes whatever is on your clipboard. Then hit Scan.',
  },
  {
    target: 'examples',
    on: 'landing',
    title: 'Not sure what to try?',
    text: 'These run real scans. The red ones are lookalikes spelled with Cyrillic letters.',
  },
  {
    target: 'verdict',
    on: 'results',
    title: 'The verdict',
    text: 'Safe, Suspicious or Risky, with a trust score from 0 to 100 and a live screenshot of the page.',
  },
  {
    target: 'flags',
    on: 'results',
    title: 'Why I said that',
    text: 'Red flags pushed the score down, green flags pushed it up. Every reason is listed.',
  },
  {
    target: 'sections',
    on: 'results',
    title: 'The full breakdown',
    text: 'Each section has a one-glance status on the right. Open any of them for the details.',
  },
  {
    target: 'report',
    on: 'results',
    title: 'Think I got it wrong?',
    text: 'Report it. A human reads every report, and it helps me get better.',
  },
  {
    target: 'share',
    on: 'results',
    title: 'Share the result',
    text: 'Send it to whoever sent you the link. They will see the verdict without scanning again.',
  },
];

/** Steps for the current page (landing or results) whose target is on screen. */
export function availableSteps(): TourStep[] {
  const has = (target: string) => !!document.querySelector(`[data-guide="${target}"]`);
  const page = has('verdict') ? 'results' : 'landing';
  return TOUR_STEPS.filter((s) => s.on === page && has(s.target));
}
