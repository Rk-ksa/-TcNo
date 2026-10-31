import { describe, expect, it } from "vitest";
import { accountAvatarSrc, type AccountAvatarInput } from "./accountAvatarSrc";

const base: AccountAvatarInput = {
  streamer: false,
  salt: "salt",
  platformKey: "Steam",
  accountKey: "76561198000000001",
  imageUrl: "https://avatars.example/real.jpg",
  pending: false,
  epoch: 0,
  offline: false,
  fallback: "/img/platform/Steam.svg",
};

describe("accountAvatarSrc", () => {
  it("shows the real avatar when there is nothing to hide", () => {
    expect(accountAvatarSrc(base)).toContain("https://avatars.example/real.jpg");
  });

  it("uses the platform icon instead of the real avatar in streamer mode", () => {
    const src = accountAvatarSrc({ ...base, streamer: true });
    expect(src).toBe("/img/platform/Steam.svg");
  });

  it("uses the platform icon when the account has no avatar", () => {
    expect(accountAvatarSrc({ ...base, imageUrl: "" })).toBe("/img/platform/Steam.svg");
  });

  it("leaves an in-flight avatar on the platform placeholder", () => {
    expect(accountAvatarSrc({ ...base, imageUrl: "", pending: true })).toBe("/img/platform/Steam.svg");
  });

  it("uses the platform icon when offline blocks a remote avatar", () => {
    expect(accountAvatarSrc({ ...base, offline: true })).toBe("/img/platform/Steam.svg");
  });

  it("uses the platform icon when account identity data is missing", () => {
    expect(accountAvatarSrc({ ...base, imageUrl: "", accountKey: "" })).toBe("/img/platform/Steam.svg");
  });

  it("keeps using the platform icon in streamer mode without an account key", () => {
    const src = accountAvatarSrc({ ...base, streamer: true, accountKey: "" });
    expect(src).toBe("/img/platform/Steam.svg");
  });
});
