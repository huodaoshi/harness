/** 与 backend GET /api/config 响应对齐（NextChat DangerConfig） */
export type DangerConfig = {
  needCode: boolean;
  hideUserApiKey: boolean;
  disableGPT4: boolean;
  hideBalanceQuery: boolean;
  disableFastLink: boolean;
  customModels: string;
  defaultModel: string;
  visionModels: string;
};
