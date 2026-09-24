// Mentions of url.vet shown on the landing page.

export type ShoutoutType = 'tweet' | 'paper' | 'email' | 'package' | 'newsletter' | 'linkedin';

export type Shoutout = {
  name: string;
  handle: string;
  date: string;
  text: string;
  url: string;
  type: ShoutoutType;
  flag: string;
};

export const SHOUTOUTS: Shoutout[] = [
  {
    type: 'paper',
    name: 'Academic Research',
    handle: 'ijesr.org',
    date: 'Apr 26, 2025',
    text: 'An earlier version of this project was referenced and reproduced (👀) in an academic research paper.',
    url: 'https://www.ijesr.org/index.php/ijesr/article/view/377',
    flag: '🇮🇳',
  },
  {
    type: 'package',
    name: 'Dika Ardianta',
    handle: '@DikaArdnt',
    date: 'Apr 28, 2026',
    text: 'I ported the detection engine to PHP and published it as an unofficial community Composer package on packagist.org',
    url: 'https://packagist.org/packages/safesurf/safesurf',
    flag: '🇮🇩',
  },
  {
    type: 'tweet',
    name: 'Tom Dörr',
    handle: '@tom_doerr',
    date: 'Jan 16, 2026',
    text: 'Engine for phishing detection with a web UI and browser extension',
    url: 'https://twitter.com/tom_doerr/status/2012050578721915177',
    flag: '🇩🇪',
  },
  {
    type: 'newsletter',
    name: 'OSINTech',
    handle: 'osintech.substack.com',
    date: 'May 21, 2026',
    text: 'URLvet. Open-source phishing detection engine — get a trust score, a fully explainable verdict, and a shareable security report with live page preview, all in real time.',
    url: 'https://osintech.substack.com/p/osintechs-timeline-163-21052026?open=false#%C2%A7osint-tools-services-and-investigations',
    flag: '🇰🇿',
  },
  {
    type: 'tweet',
    name: 'NeoTeo.com',
    handle: '@NeoteoCom',
    date: 'Jan 16, 2026',
    text: 'SafeSurf, an open-source phishing detection engine with a web UI and browser extension. Easy to integrate to protect users and reduce fraud.',
    url: 'https://twitter.com/NeoteoCom/status/2012102861807554843',
    flag: '🇪🇸',
  },
  {
    type: 'tweet',
    name: 'Bryan',
    handle: '@so_sthbryan',
    date: 'May 16, 2026',
    text: 'Explainable phishing detection that scans URLs in real time. SafeSurf catches phishing sites before you interact with them.',
    url: 'https://x.com/so_sthbryan/status/2055390764339974377',
    flag: '🇺🇸',
  },
  {
    type: 'linkedin',
    name: 'Şevket Can Özöğretmen',
    handle: '@ozogretmen',
    date: 'May 14, 2026',
    text: 'Unsure where a link you click online will take you? url.vet is a powerful phishing detection engine that scans, scores, and transparently explains why a link is safe or not in real time.',
    url: 'https://www.linkedin.com/posts/ozogretmen_cybersecurity-opensource-phishingdetection-ugcPost-7460836352545341440-B4yF/',
    flag: '🇹🇷',
  },
];
