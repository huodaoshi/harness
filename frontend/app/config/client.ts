import { BuildConfig, getBuildConfig } from "./build";

export function getClientConfig() {
  if (typeof document !== "undefined") {
    return JSON.parse(queryMeta("config") || "{}") as BuildConfig;
  }

  if (typeof process !== "undefined") {
    return getBuildConfig();
  }
}

function queryMeta(key: string, defaultValue?: string): string {
  if (typeof document === "undefined") {
    return defaultValue ?? "";
  }
  const meta = document.head.querySelector(
    `meta[name='${key}']`,
  ) as HTMLMetaElement;
  return meta?.content ?? defaultValue ?? "";
}
