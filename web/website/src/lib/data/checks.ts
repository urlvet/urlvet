// The check groups on the How it works page: what each one looks for, and whether a
// finding counts towards risk or trust. The rules live in
// server/internal/analyzer/result.go; no point values here, since those get tuned.
// Result sections link here by id (#check-url, #check-network, ...).

export type Signal = {
  text: string;
  /** Which side of the score it counts towards; 'info' is shown but not scored. */
  kind: 'risk' | 'trust' | 'info';
};

export type CheckGroup = {
  id: string;
  /** The question this group answers, in plain words. */
  title: string;
  /** The technical name, for people who know it. */
  label: string;
  /** How many of the scan's checks belong to this group. */
  checks: number;
  desc: string;
  signals: Signal[];
  /** What the findings are compared against. */
  against?: string;
};

const risk = (text: string): Signal => ({ text, kind: 'risk' });
const trust = (text: string): Signal => ({ text, kind: 'trust' });
const info = (text: string): Signal => ({ text, kind: 'info' });

export const CHECK_GROUPS: CheckGroup[] = [
  {
    id: 'url',
    title: 'Does the address itself look wrong?',
    label: 'URL structure',
    checks: 8,
    desc: 'Reads the link before visiting anything. Many scams give themselves away in the address alone.',
    signals: [
      risk('A raw IP address instead of a name'),
      risk('Encoded lookalike letters (punycode)'),
      risk('Letters from another alphabet that look the same'),
      risk('An unregulated ending'),
      risk('A very deep path'),
      risk('A link shortener that hides the destination'),
      risk('A high-risk ending, like .top or .zip'),
      risk('A very long link'),
      risk('More than two subdomains'),
      risk('Words like "login" or "verify" in the link'),
      trust('A verified ending, like .gov, .edu or .bank'),
    ],
    against: '138 high-risk endings, 238 verified endings and 2,668 known link shorteners.',
  },
  {
    id: 'network',
    title: 'Where does the link really go?',
    label: 'Redirects & HTTP',
    checks: 1,
    desc: 'Makes one real request from our server and follows every redirect to the final page.',
    signals: [
      risk('It lands on a different site than the one in the link'),
      risk('More than three redirects that leave the site'),
      trust('The site insists on a secure connection (HSTS)'),
      info('The response code the site returns'),
    ],
  },
  {
    id: 'dns',
    title: 'Is the site set up like a real one?',
    label: 'DNS',
    checks: 2,
    desc: 'Looks up where the site lives. Throwaway scam sites often skip the basics.',
    signals: [
      risk('Missing or incomplete name-server setup'),
      risk('No mail server for the domain'),
      info('The addresses the site resolves to'),
    ],
  },
  {
    id: 'tls',
    title: 'Is the connection secure?',
    label: 'TLS / SSL',
    checks: 1,
    desc: "Reads the site's security certificate. A valid certificate proves the connection is private, not that the site is honest, so most of this is shown for you to read.",
    signals: [
      risk('A password form on a page without HTTPS'),
      info('Who issued the certificate, and when it expires'),
      info('Whether it matches the address and is publicly logged'),
    ],
  },
  {
    id: 'domain',
    title: 'How old and well-known is the site?',
    label: 'Domain intelligence',
    checks: 4,
    desc: 'Real sites build up a history. Scam sites are usually days old, unknown, and named after someone else.',
    signals: [
      risk('One or two letters away from a well-known site'),
      risk('Registered in the last 30 days'),
      risk("A brand's name inside someone else's address"),
      risk('Registered less than a year ago'),
      risk('Not among the top million sites'),
      trust('Among the 10,000 most-visited sites'),
      trust('Among the 50,000 most-visited sites'),
      trust('Ranked, with lower traffic'),
      trust('Registered more than five years ago'),
      trust('Signed DNS records (DNSSEC)'),
      info('How random the name looks'),
    ],
    against:
      'The top million sites for traffic rank, and the 5,000 most-visited for lookalike spellings.',
  },
  {
    id: 'content',
    title: 'What does the page ask you for?',
    label: 'Page content',
    checks: 1,
    desc: 'Fetches the page and reads its forms: what they ask for, and where that information would be sent.',
    signals: [
      risk("It claims to be a brand that doesn't own the address"),
      risk('A form sends a password, card or personal details to another site'),
      risk('A login form on a new or unknown site'),
      risk('A hidden frame loading another page'),
      risk('Card or payment fields'),
      risk('A form sends an email or username to another site'),
      trust('It names a brand that does own the address'),
      info('Tracking pixels, and requests for personal details'),
    ],
    against: '134 well-known brands and the addresses they really use.',
  },
  {
    id: 'threats',
    title: 'Has anyone already reported it?',
    label: 'Threat feeds',
    checks: 1,
    desc: 'Asks PhishTank, a public list of phishing links reported and reviewed by volunteers.',
    signals: [risk('Confirmed as phishing'), risk('Reported as phishing, not yet reviewed')],
  },
];

export const TOTAL_CHECKS = CHECK_GROUPS.reduce((n, g) => n + g.checks, 0);
