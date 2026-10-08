// @vitest-environment jsdom
import { beforeAll, describe, expect, it } from "vitest";
import * as R from "../src/render.js";

beforeAll(() => {
  globalThis.requestAnimationFrame ??= (cb) =>
    setTimeout(() => cb(performance.now()), 0);
});

const result = (over = {}) => ({
  url: "https://evil.example/login",
  verdict: "Risky",
  score: 12,
  label: "High Risk",
  quip: "Likely unsafe. Better not to click. 🚫",
  bad: ["Very low traffic volume."],
  good: [],
  shortLink: null,
  incomplete: false,
  incompleteChecks: [],
  ...over,
});

describe("verdictCard", () => {
  it("shows the verdict, its badge and quip", () => {
    const card = R.verdictCard(result());
    expect(card.dataset.verdict).toBe("Risky");
    expect(card.textContent).toContain("Risky");
    expect(card.textContent).toContain("High Risk");
    expect(card.querySelector(".ring").getAttribute("aria-label")).toBe(
      "Trust score 12 out of 100",
    );
  });
});

describe("flags", () => {
  it("never parses text from the network as HTML", () => {
    const evil = '<img src=x onerror="alert(1)">';
    const section = R.flags(result({ bad: [evil] }));
    expect(section.querySelector("img")).toBeNull();
    expect(section.textContent).toContain(evil);
  });
  it("shows green flags for a safe result, and how many more there are", () => {
    const section = R.flags(
      result({ verdict: "Safe", bad: [], good: ["a", "b", "c", "d", "e"] }),
      3,
    );
    expect(section.querySelector("h3").textContent).toBe("Green flags");
    expect(section.querySelectorAll("li")).toHaveLength(3);
    expect(section.textContent).toContain("+2 more in the full report");
  });
  it("shows nothing when there are no flags", () => {
    expect(R.flags(result({ bad: [], good: [] }))).toBeNull();
  });
});

describe("notes", () => {
  it("uses the website's short-link wording", () => {
    const [note] = R.notes(
      result({
        shortLink: {
          url: "https://goo.su/TIcsAp",
          chain: ["https://goo.su/TIcsAp", "https://www.roblox.com.bn/x"],
          target: "https://www.roblox.com.bn/x",
          resolved: true,
        },
      }),
    );
    expect(note.textContent).toBe(
      "This is a short link on goo.su. It leads to www.roblox.com.bn, and this result is for that page.",
    );
  });
  it("says which checks didn't finish", () => {
    const [note] = R.notes(
      result({ incomplete: true, incompleteChecks: ["content_check"] }),
    );
    expect(note.textContent).toContain("the page itself");
  });
});

describe("wordmark", () => {
  it("takes the verdict's state on its dot", () => {
    expect(R.wordmark("Risky").querySelector(".dot").dataset.state).toBe(
      "Risky",
    );
  });
});

describe("redirect note", () => {
  it("says where a well-known site's redirect led", () => {
    const [note] = R.notes(
      result({
        redirectedFrom: {
          url: "https://sba.yandex.ru/redirect?url=x",
          chain: [],
          target: "https://nus.nbt.mybluehost.me/x",
        },
      }),
    );
    expect(note.textContent).toContain(
      "sends visitors on to nus.nbt.mybluehost.me",
    );
  });
});
