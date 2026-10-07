import { describe, expect, it } from 'vitest';
import { redirectNote, shortLinkNote } from './shortlink';

describe('shortLinkNote', () => {
  it('says nothing for an ordinary link', () => {
    expect(shortLinkNote(undefined)).toBeNull();
  });

  it('names where a short link leads', () => {
    expect(
      shortLinkNote({
        url: 'https://goo.su/TIcsAp',
        chain: ['https://goo.su/TIcsAp', 'https://www.roblox.com.bn/users/1/profile'],
        target: 'https://www.roblox.com.bn/users/1/profile',
        resolved: true,
      })
    ).toBe(
      'This is a short link on goo.su. It leads to www.roblox.com.bn, and this result is for that page.'
    );
  });

  it('names short links passed through on the way', () => {
    expect(
      shortLinkNote({
        url: 'https://bit.ly/x',
        chain: ['https://bit.ly/x', 'https://tinyurl.com/y', 'https://evil.example/login'],
        target: 'https://evil.example/login',
        resolved: true,
      })
    ).toBe(
      'This is a short link on bit.ly. It leads to evil.example (through tinyurl.com), and this result is for that page.'
    );
  });

  it('says when the destination could not be found', () => {
    expect(
      shortLinkNote({ url: 'https://bom.so/abc', chain: ['https://bom.so/abc'], resolved: false })
    ).toBe(
      "This is a short link on bom.so. We couldn't see where it leads, so this result is for the short link itself and can't be better than Suspicious."
    );
  });
});

describe('redirectNote', () => {
  it('says nothing without a redirect', () => {
    expect(redirectNote(undefined)).toBeNull();
  });
  it('names the site that redirected and where to', () => {
    expect(
      redirectNote({
        url: 'https://sba.yandex.ru/redirect?url=x',
        chain: ['https://sba.yandex.ru/redirect?url=x', 'https://nus.nbt.mybluehost.me/x'],
        target: 'https://nus.nbt.mybluehost.me/x',
      })
    ).toBe(
      'This link on sba.yandex.ru sends visitors on to nus.nbt.mybluehost.me. This result is for that page, plus anything wrong with the link itself.'
    );
  });
});
