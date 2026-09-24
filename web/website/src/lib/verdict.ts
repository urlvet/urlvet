// Verdict wording shared by the results page and the social preview image,
// so both always say the same thing.

export type Verdict = 'Safe' | 'Suspicious' | 'Risky';

export const VERDICT_COPY: Record<Verdict, { label: string; quip: string }> = {
  Safe: {
    label: 'Trusted',
    quip: 'Looks legit. Should be fine to click. 👌',
  },
  Suspicious: {
    label: 'Be Cautious',
    quip: "Hmm, something's off. Be careful. 🤔",
  },
  Risky: {
    label: 'High Risk',
    quip: 'Likely unsafe. Better not to click. 🚫',
  },
};

export function verdictCopy(verdict: string | undefined) {
  return VERDICT_COPY[(verdict ?? '') as Verdict] ?? VERDICT_COPY.Suspicious;
}
