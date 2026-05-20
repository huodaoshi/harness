const STORAGE_ANON = "harness_anon_id";
const STORAGE_TOKEN = "harness_access_token";

function storage(): Storage | null {
  if (typeof window === "undefined") return null;
  return window.localStorage;
}

/** 浏览器内稳定的游客 UUID。SSR 时返回空字符串。 */
export function getAnonId(): string {
  const s = storage();
  if (!s) return "";
  let id = s.getItem(STORAGE_ANON);
  if (!id) {
    id = crypto.randomUUID();
    s.setItem(STORAGE_ANON, id);
  }
  return id;
}

export function getAccessToken(): string | null {
  return storage()?.getItem(STORAGE_TOKEN) ?? null;
}

export function setAccessToken(token: string | null): void {
  const s = storage();
  if (!s) return;
  if (token) s.setItem(STORAGE_TOKEN, token);
  else s.removeItem(STORAGE_TOKEN);
}

/** Wellness API：Bearer 优先，否则 X-Anon-ID。 */
export function harnessAuthHeaders(): Record<string, string> {
  const token = getAccessToken();
  if (token) return { Authorization: `Bearer ${token}` };
  const anon = getAnonId();
  if (!anon) return {};
  return { "X-Anon-ID": anon };
}

export type WellnessApiErrorBody = {
  code?: number;
  message?: string;
};

/** 401/429 等 JSON 错误文案（关怀 API 约定）。 */
export function wellnessErrorMessage(
  status: number,
  body: WellnessApiErrorBody,
  fallback?: string,
): string {
  if (body.message) return body.message;
  if (status === 401) {
    if (body.code === 4014) return "匿名标识无效，请刷新页面后重试";
    return "未登录或登录已过期，请重新打开页面";
  }
  if (status === 429) return "请求过于频繁，请稍后再试";
  return fallback ?? `HTTP ${status}`;
}
