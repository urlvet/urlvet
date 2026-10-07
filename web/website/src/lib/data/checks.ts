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
      risk("A short link we couldn't follow to where it leads"),
      risk('A high-risk ending, like .top or .zip'),
      risk('A very long link'),
      risk('More than two subdomains'),
      risk('Words like "login" or "verify" in the link'),
      trust('A verified ending, like .gov, .edu or .bank'),
      info(
        'Short links are followed, through any further short links, and the page they lead to is checked instead'
      ),
    ],
    against: '138 high-risk endings, 238 verified endings and 2,671 known link shorteners.',
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
      risk('A certificate issued for a different address (browsers warn about it)'),
      risk('A password form on a page without HTTPS'),
      info('Who issued the certificate, and when it expires'),
      info('Whether it is publicly logged'),
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
      risk('Registered in the last 90 days'),
      risk('Not among the top million sites, and new or of unknown age'),
      trust('Among the 10,000 most-visited sites'),
      trust('Among the 50,000 most-visited sites'),
      trust('Ranked, with lower traffic'),
      trust('Running for more than a year, more so after three and five'),
      trust('Signed DNS records (DNSSEC)'),
      info('How random the name looks'),
      info(
        "A site on a hosting service or site builder (github.io, vercel.app, weebly.com) is judged on its own, not on the host's reputation"
      ),
    ],
    against:
      'The top million sites for traffic rank, the 5,000 most-visited for lookalike spellings, and 514 brands for names inside an address.',
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
      risk('Cloudflare or the host has flagged the page or taken it down'),
      risk('The page sends you on by script, to a raw IP address or another site'),
      risk('A page hidden in a WordPress system folder, a sign of a hacked site'),
      risk('It downloads a program from an unknown site, or one any user could upload'),
      trust('It names a brand that does own the address'),
      info('Tracking pixels, and requests for personal details'),
    ],
    against:
      '514 brands across banking, payments, crypto, shopping, delivery, telecoms, government and more, and the addresses they really use. Lookalike letters are read as the plain ones they imitate.',
  },
  {
    id: 'threats',
    title: 'Has anyone already reported it?',
    label: 'Threat feeds',
    checks: 4,
    desc: 'Checks lists of links already reported as phishing or malware: PhishTank, URLhaus and Google Safe Browsing. The lists are kept on our server, and Google only ever sees short codes made from the link, never the link.',
    signals: [
      risk('On a list of confirmed phishing or malware links'),
      risk('Other pages on the same site are on those lists'),
      risk('Google lists it as dangerous'),
    ],
  },
];

export const TOTAL_CHECKS = CHECK_GROUPS.reduce((n, g) => n + g.checks, 0);
