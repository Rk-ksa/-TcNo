import { offlineSafeImageSrc, withAssetCacheBust } from "../stores/offlineMode";

export interface AccountAvatarInput {
  /** Streamer mode: never show the real face, whatever else is true. */
  streamer: boolean;
  /** Machine-local salt from the backend; see generatedAvatar. */
  salt: string;
  platformKey: string;
  /** Login name or platform ID — the seed the generated avatar is built from. */
  accountKey: string;
  imageUrl?: string | null;
  /** An avatar that is being fetched right now, not one that does not exist. */
  pending: boolean;
  epoch: number;
  offline: boolean;
  /** The platform's static placeholder, still used while an avatar is in flight. */
  fallback: string;
}

/**
 * The single rule for what an account tile shows.
 *
 * Real profile images are shown when available. Otherwise use the platform
 * artwork passed as `fallback`, including streamer mode, so no generated colour
 * tile is mistaken for an account photo.
 */
export function accountAvatarSrc(input: AccountAvatarInput): string {
  if (input.streamer) {
    return input.fallback;
  }
  if (input.pending) {
    return input.fallback;
  }
  const real = offlineSafeImageSrc(
    input.offline,
    withAssetCacheBust(input.imageUrl, input.epoch),
    "",
  );
  if (real) {
    return real;
  }
  return input.fallback;
}
