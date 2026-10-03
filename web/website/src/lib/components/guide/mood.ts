/** What the helper character is currently expressing. */
export type Mood =
  | 'idle'
  | 'wave'
  | 'thinking'
  | 'happy'
  | 'worried'
  | 'point'
  /** Checking a link: holds a magnifying glass. */
  | 'scanning'
  /** A Suspicious verdict: one raised eyebrow and a side-eye. */
  | 'suspicious'
  /** A Risky verdict: worried, shaking, holding up a shield. */
  | 'alarmed'
  /** A quick "ooh" when something lands, like a pasted link or a poke. */
  | 'ooh';
