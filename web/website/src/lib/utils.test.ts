import { describe, expect, it } from 'vitest';
import { extractLinks } from './utils';

describe('extractLinks', () => {
  it('finds the link in a forwarded message', () => {
    const msg =
      '🎁 Congratulations! You won a prize. Claim it now at https://prize-claim.top/win?id=42 before it expires!!';
    expect(extractLinks(msg)).toEqual(['https://prize-claim.top/win?id=42']);
  });

  it('strips sentence punctuation after a link', () => {
    expect(extractLinks('Check this (https://example.com/deal).')).toEqual([
      'https://example.com/deal',
    ]);
    expect(extractLinks('go to example.com, now')).toEqual(['example.com']);
  });

  it('returns every distinct link, in order', () => {
    const msg =
      'Deal: https://shop.example/a and also www.other.example/b, again https://shop.example/a';
    expect(extractLinks(msg)).toEqual(['https://shop.example/a', 'www.other.example/b']);
  });

  it('prefers explicit links over bare word.word matches', () => {
    expect(extractLinks('Mr.Smith says open https://example.com')).toEqual(['https://example.com']);
  });

  it('skips email addresses', () => {
    expect(extractLinks('write to support@bank.example or visit bank.example/help')).toEqual([
      'bank.example/help',
    ]);
  });

  it('handles lookalike letters from other alphabets', () => {
    expect(extractLinks('login at pаypal.com/verify')).toEqual(['pаypal.com/verify']);
  });

  it('returns nothing when there is no link', () => {
    expect(extractLinks('hello, how are you?')).toEqual([]);
  });
});
