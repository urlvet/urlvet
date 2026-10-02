import { describe, expect, it } from 'vitest';
import { easterEgg, localAddressMessage } from './eastereggs';

describe('easterEgg', () => {
  it.each([
    ['url.vet', 'Checking me?'],
    ['https://www.url.vet/about', 'Checking me?'],
    ['https://github.com/urlvet/urlvet', 'where I live'],
    ['https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'Never gonna give you up'],
    ['https://youtu.be/dQw4w9WgXcQ', 'Never gonna give you up'],
    ['pаypal.com', 'Spot the difference'],
    ['notascam.lol', 'from our tagline'],
    ['paypal.com', 'The real one'],
    ['https://www.apple.com/iphone', 'Cyrillic twin'],
    ['example.com', 'most innocent'],
    ['en.wikipedia.org/wiki/Phishing', 'Citation needed'],
    ['1.1.1.1', 'phone book'],
    ['https://invoice.zip', 'isn’t a file'.replace('’', "'")],
  ])('%s', (url, expected) => {
    expect(easterEgg(url)).toContain(expected);
  });

  it('stays quiet for ordinary links', () => {
    expect(easterEgg('github.com')).toBeNull();
    expect(easterEgg('https://www.youtube.com/watch?v=abc')).toBeNull();
    expect(easterEgg('paypal.com.evil.top')).toBeNull();
  });
});

describe('localAddressMessage', () => {
  it.each(['localhost', 'http://localhost:3000', '127.0.0.1', 'http://[::1]/'])(
    '%s is your own computer',
    (url) => {
      expect(localAddressMessage(url)).toContain('your own computer');
    }
  );

  it.each(['192.168.1.1', '10.0.0.5', '172.20.1.1', 'printer.local'])(
    '%s is a private network',
    (url) => {
      expect(localAddressMessage(url)).toContain('private network');
    }
  );

  it('lets public addresses through', () => {
    expect(localAddressMessage('example.com')).toBeNull();
    expect(localAddressMessage('172.32.0.1')).toBeNull();
    expect(localAddressMessage('8.8.8.8')).toBeNull();
  });
});
