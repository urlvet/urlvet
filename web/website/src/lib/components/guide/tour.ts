// Guided tour steps. Each points at an element tagged data-guide="<target>";
// steps whose target isn't on the page right now are skipped.

/** landing: before a scan. results: once a result is on screen. */
export type TourStep = {
  target: string;
  title: string;
  text: string;
  on: 'landing' | 'results';
  /** Room around the highlight, in px (default 8); less where neighbours are close. */
  pad?: number;
};

export const TOUR_STEPS: TourStep[] = [
  {
    target: 'search',
    on: 'landing',
    title: 'Paste the link here',
    text: 'A full link, a short link, or just a name like example.com. You can paste a whole message too, and I’ll find the links in it.',
  },
  {
    target: 'paste',
    on: 'landing',
    pad: 2, // Scan sits right beside it
    title: 'Or paste in one click',
    text: 'This pastes whatever you last copied.',
  },
  {
    target: 'scan',
    on: 'landing',
    title: 'Then press Scan',
    text: 'Pressing Enter works too. The check takes a few seconds.',
  },
  {
    target: 'examples',
    on: 'landing',
    title: 'Just want to try it?',
    text: 'These check real sites. The red ones are fakes with lookalike letters.',
  },
  {
    target: 'verdict',
    on: 'results',
    title: 'The answer',
    text: 'Safe, Suspicious or Risky, with a score out of 100. Higher is safer.',
  },
  {
    target: 'screenshot',
    on: 'results',
    title: 'What the page looks like',
    text: 'A picture of the page, so you can see it without opening it yourself.',
  },
  {
    target: 'flags',
    on: 'results',
    title: 'Why',
    text: 'Red flags are what worried me, green flags what reassured me.',
  },
  {
    target: 'sections',
    on: 'results',
    title: 'All the details',
    text: 'Open any row to see exactly what I checked. The dot on the right is the short answer.',
  },
  {
    target: 'report',
    on: 'results',
    title: 'Think I got it wrong?',
    text: 'Tell us. A person reads every report.',
  },
  {
    target: 'share',
    on: 'results',
    title: 'Share the result',
    text: 'Send it to whoever sent you the link, so they can see it too.',
  },
  {
    target: 'learn-more',
    on: 'landing',
    title: 'Want to know more?',
    text: '“How it works” explains every check. “Privacy” says what happens to the links you check.',
  },
];

/** Steps for the current page (landing or results) whose target is on screen. */
export function availableSteps(): TourStep[] {
  const has = (target: string) => !!document.querySelector(`[data-guide="${target}"]`);
  const page = has('verdict') ? 'results' : 'landing';
  return TOUR_STEPS.filter((s) => s.on === page && has(s.target));
}
