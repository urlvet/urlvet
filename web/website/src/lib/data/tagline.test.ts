import { describe, expect, it } from 'vitest';
import { TAGLINE_LINKS, TAGLINE_MAX_LENGTH } from './tagline';

describe('TAGLINE_LINKS', () => {
  it('keeps every link short enough to read on a phone', () => {
    const tooLong = TAGLINE_LINKS.filter((l) => l.length > TAGLINE_MAX_LENGTH);
    expect(tooLong).toEqual([]);
  });
});
