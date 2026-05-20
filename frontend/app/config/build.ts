import { DEFAULT_INPUT_TEMPLATE } from "../constant";

export const getBuildConfig = () => {
  const buildMode = process.env.BUILD_MODE ?? "standalone";
  const isApp = !!process.env.BUILD_APP;

  return {
    version: process.env.npm_package_version ?? "0.1.0",
    commitDate: "unknown",
    commitHash: "unknown",
    buildMode,
    isApp,
    template: process.env.DEFAULT_INPUT_TEMPLATE ?? DEFAULT_INPUT_TEMPLATE,
  };
};

export type BuildConfig = ReturnType<typeof getBuildConfig>;
