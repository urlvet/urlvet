// Questions and answers for the How it works page. Worded the way people search,
// since they also feed the FAQ rich result. Keep answers in step with what the
// scanner and the privacy page actually do.

export type Faq = {
  q: string;
  a: string;
  /** Optional short list shown under the answer. */
  points?: string[];
};

export const FAQ_GROUPS: { title: string; items: Faq[] }[] = [
  {
    title: 'Checking a link',
    items: [
      {
        q: 'How do I check if a link is safe?',
        a: 'Paste it into the box on the url.vet home page and press Scan. A few seconds later you get a verdict (Safe, Suspicious or Risky), a score from 0 to 100, and the reasons behind it. You never have to open the link yourself.',
      },
      {
        q: 'How do I check a link without clicking it?',
        a: 'On a phone, press and hold the link and choose Copy link. On a computer, right-click it and choose Copy link address. Then paste it into url.vet. The link is never opened on your device.',
      },
      {
        q: 'Can I check a link from WhatsApp, SMS, email or a QR code?',
        a: 'Yes. Copy the link and paste it here, or paste the whole message and url.vet will find the links in it. For a QR code, scan it with your camera and copy the link instead of opening it.',
      },
      {
        q: 'Is it safe to check a dangerous link here?',
        a: 'Yes. Our server visits the link, not your device. Nothing from the page runs in your browser; you only see a screenshot.',
      },
      {
        q: 'Will the website know I checked it?',
        a: 'It sees a visit from our server, not from you. Your IP address is never sent to it.',
      },
      {
        q: 'Can I use url.vet on my phone?',
        a: "Yes. It works in any phone browser, with nothing to install. If you like, add it to your home screen from your browser's menu. On Android you can then share a link from another app straight to url.vet.",
      },
    ],
  },
  {
    title: 'Understanding the result',
    items: [
      {
        q: 'What do Safe, Suspicious and Risky mean?',
        a: 'Every link gets a score from 0 to 100. Safe (65 and up) means no warning signs turned up. Suspicious (30 to 64) means there are some warning signs, or not enough evidence either way, so be careful with logins and payments. Risky (below 30) means strong signs of phishing or a scam, so stay away.',
      },
      {
        q: "Why does my result say some checks didn't finish?",
        a: "Some checks depend on other services, like the site itself, its domain records or the PhishTank database, and those can be slow or unavailable. The result tells you which checks were missed. If a missed check could change the verdict, the result is marked incomplete and isn't saved, so scanning again a minute later retries it.",
      },
      {
        q: 'Why did a site I trust come back Suspicious?',
        a: "New or little-known sites have fewer signals in their favour, so they can score lower. The flags show exactly why. If you think it's wrong, use the Report button on the result. A human reads every report.",
      },
      {
        q: 'Does Safe mean the site is honest?',
        a: 'It means no warning signs turned up. A real shop can still sell bad products, and a brand-new scam may not show any signs yet. Use url.vet as one layer of defense, not the only one.',
      },
      {
        q: 'How accurate is url.vet? Can it be wrong?',
        a: 'It can. url.vet runs a fixed set of checks and shows every reason, so you can judge the result yourself. A new, little-known site can look suspicious, and a scam that went live minutes ago may show no signs yet. If a result looks wrong, use the Report button.',
      },
    ],
  },
  {
    title: 'Staying safe',
    items: [
      {
        q: 'What is phishing?',
        a: 'Phishing is a scam where a fake message or website pretends to be someone you trust, like your bank, a delivery company or your employer, to steal your password, card number or money. It usually arrives as a link in an email, SMS or chat message.',
      },
      {
        q: 'How can I tell if a link is a phishing scam?',
        a: "Common signs: the address misspells a real brand (paypa1.com), swaps in lookalike letters from another alphabet, was registered only days ago, uses an unusual ending like .top or .zip, or the page asks for a password or card details on a site the brand doesn't own. url.vet checks for all of these and lists any it finds.",
      },
      {
        q: 'What is typosquatting?',
        a: "Typosquatting is registering an address that is one slip away from a real one, like paypa1.com or amaz0n.com, hoping you won't notice. Some fakes go further and swap in letters from other alphabets that look identical. url.vet compares every address against well-known sites to catch both.",
      },
      {
        q: 'Are shortened links like bit.ly or tinyurl safe?',
        a: "A shortener isn't dangerous by itself, but it hides where the link really goes, which is why scammers like them. Paste the short link into url.vet: it follows every redirect and shows the full path to the page you would land on.",
      },
      {
        q: 'Does the padlock or https mean a website is safe?',
        a: 'No. The padlock only means the connection is private. It says nothing about who runs the site. Most phishing sites have a padlock too, so look at the address itself, not the icon.',
      },
      {
        q: 'Can I get a virus just by clicking a link?',
        a: "On an up-to-date phone or browser, simply opening a page rarely infects you, but it can happen. The bigger risk is what comes next: typing your password or card details into a fake page, or downloading and opening a file. Keep your browser updated and check links you don't trust first.",
      },
      {
        q: 'I already clicked a suspicious link. What should I do?',
        a: "Don't enter anything on the page. If you already typed a password, change it on the real site and turn on two-step verification. If you entered card or bank details, call your bank right away. Then scan the link here to see what it was.",
      },
      {
        q: 'What are some simple tips to stay safe from scam links?',
        a: 'A few habits stop most scams:',
        points: [
          "Check before you click. If you weren't expecting a link, paste it into url.vet first.",
          'Read the address, not the page. Scam pages copy the look of real ones; the address gives them away.',
          'Be wary of urgency. "Your account will be blocked today" is how scammers stop you thinking.',
          "Never type a password, OTP or card number on a page you reached from a message. Open the app or type the site's address yourself.",
          'Turn on two-step verification for your email, bank and social accounts.',
          'Keep your phone and browser updated.',
          'If an offer sounds too good to be true, it is.',
        ],
      },
    ],
  },
  {
    title: 'About url.vet',
    items: [
      {
        q: 'Is url.vet free? Do I need an account?',
        a: "It's free, with no account, no signup and no ads.",
      },
      {
        q: 'Do you keep the links I check?',
        a: "Results are cached for up to 24 hours so a repeat check loads instantly, then deleted. We don't log who checked what, and there are no accounts. The privacy page has the details.",
      },
      {
        q: 'Is url.vet open source? Who runs it?',
        a: 'Yes. All of the code, including the checks and the scoring, is public on GitHub under the AGPL-3.0 license, so anyone can read it or run their own copy. One developer builds and runs it. The About page has the story.',
      },
      {
        q: 'Is url.vet an alternative to VirusTotal or URLScan?',
        a: "For checking links, yes. url.vet is free and open source, and it explains its verdict in plain words instead of listing raw data. It checks web addresses and pages only; it doesn't scan files or use antivirus engines.",
      },
      {
        q: 'What is the paperclip in the corner?',
        a: "That's Vetty, url.vet's helper, inspired by Clippy from old office software. Tap him and he will check a link for you, explain a result in simple words, or help you warn the person who sent it. He is not an AI chatbot, and he never collects your data, shows ads, or asks you to sign up. You can hide him from his menu. The About page has his story.",
      },
    ],
  },
];

/** A stable id for linking straight to a question, e.g. #faq-what-is-phishing. */
export function faqId(q: string): string {
  return `faq-${q
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')}`;
}

/** The full answer as one string, for the FAQ structured data. */
export function faqText(f: Faq): string {
  return f.points ? `${f.a} ${f.points.join(' ')}` : f.a;
}
