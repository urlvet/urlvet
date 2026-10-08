import { describe, expect, it } from 'vitest';
import { incompleteNote } from './incomplete';

describe('incompleteNote', () => {
  it('says nothing when every check ran', () => {
    expect(incompleteNote(undefined, false)).toBeNull();
    expect(incompleteNote([], false)).toBeNull();
  });

  it('names a PhishTank miss without doubting the result', () => {
    expect(incompleteNote(['phishtank_check'], false)).toBe(
      "Couldn't check the phishing database (PhishTank) this time. Everything else ran normally."
    );
  });

  it('lists what the site failed to answer, and suggests a retry', () => {
    expect(incompleteNote(['content_check', 'tls_combined_check'], true)).toBe(
      "Couldn't check the page itself and the security certificate this time, so this result may be missing signals. Try again in a minute."
    );
  });

  it('merges checks that share a description', () => {
    expect(incompleteNote(['dns_validity_check', 'ip_resolution', 'domain_rank'], true)).toContain(
      'DNS records and how popular the site is'
    );
  });

  it('names the Google checks', () => {
    expect(incompleteNote(['safe_browsing_check'], false)).toBe(
      "Couldn't check Google Safe Browsing this time. Everything else ran normally."
    );
  });
});
