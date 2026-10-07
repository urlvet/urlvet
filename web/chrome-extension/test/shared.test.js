import { describe, expect, it } from "vitest";
import {
  extractLinks,
  formatUrl,
  hostIn,
  hostOf,
  isCheckable,
  reportUrl,
} from "../src/shared.js";

describe("hostIn", () => {
  it("matches a host or its subdomains, ignoring www.", () => {
    expect(hostIn("example.com", ["example.com"])).toBe(true);
    expect(hostIn("www.example.com", ["example.com"])).toBe(true);
    expect(hostIn("a.b.example.com", ["example.com"])).toBe(true);
    expect(hostIn("badexample.com", ["example.com"])).toBe(false);
  });
});

describe("links from a pasted message", () => {
  it("finds them the way the website does", () => {
    const msg =
      "Your parcel is held. Pay the fee at http://dhl-redelivery.example/pay. Or call us.";
    expect(extractLinks(msg).map(formatUrl)).toEqual([
      "http://dhl-redelivery.example/pay",
    ]);
  });
});

describe("helpers", () => {
  it("only checks web pages", () => {
    expect(isCheckable("https://example.com")).toBe(true);
    expect(isCheckable("chrome://extensions")).toBe(false);
    expect(isCheckable(undefined)).toBe(false);
  });
  it("links to the full report", () => {
    expect(reportUrl("https://url.vet/", "https://a.example/?x=1")).toBe(
      "https://url.vet/?q=https%3A%2F%2Fa.example%2F%3Fx%3D1",
    );
  });
  it("reads hosts", () => {
    expect(hostOf("https://WWW.Example.com/x")).toBe("www.example.com");
    expect(hostOf("nope")).toBe("");
  });
});
