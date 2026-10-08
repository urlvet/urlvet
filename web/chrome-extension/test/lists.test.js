import { describe, expect, it } from "vitest";
import { isKnownSite } from "../src/lists.js";

describe("isKnownSite", () => {
  const sites = new Set([
    "google.com",
    "github.com",
    "paypal.com",
    "amazonaws.com",
  ]);
  const userContent = new Set([
    "docs.google.com",
    "sites.google.com",
    "github.io",
    "amazonaws.com",
  ]);

  it.each([
    ["google.com", true],
    ["www.google.com", true],
    ["mail.google.com", true],
    ["www.paypal.com", true],
    ["docs.google.com", false], // anyone can publish a Google Doc
    ["sites.google.com", false],
    ["evil.github.io", false],
    ["bucket.s3.amazonaws.com", false],
    ["paypal.com.evil.example", false], // the brand at the front isn't the site
    ["google.com.example", false],
    ["example.com", false],
  ])("%s → %s", (host, known) => {
    expect(isKnownSite(host, sites, userContent)).toBe(known);
  });
});
