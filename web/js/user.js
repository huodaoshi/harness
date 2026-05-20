const STORAGE_ANON = "fwa_anon_id";
/** @deprecated 迁移自旧版 X-User-Id */
const STORAGE_ANON_LEGACY = "fwa_user_id";
const STORAGE_TOKEN = "fwa_access_token";

/** 浏览器内稳定的游客 UUID（每浏览器隔离）。 */
export function anonId() {
  let id = localStorage.getItem(STORAGE_ANON);
  if (!id) {
    id = localStorage.getItem(STORAGE_ANON_LEGACY);
    if (id) {
      localStorage.setItem(STORAGE_ANON, id);
    } else {
      id = crypto.randomUUID();
      localStorage.setItem(STORAGE_ANON, id);
    }
  }
  return id;
}

/** @deprecated 使用 anonId() */
export function userId() {
  return anonId();
}

export function accessToken() {
  return localStorage.getItem(STORAGE_TOKEN) || "";
}

export function setAccessToken(token) {
  if (token) localStorage.setItem(STORAGE_TOKEN, token);
  else localStorage.removeItem(STORAGE_TOKEN);
}

/** Wellness API 请求头：Bearer 优先，否则 X-Anon-ID。 */
export function apiHeaders(json = true) {
  const h = {};
  const token = accessToken();
  if (token) {
    h.Authorization = `Bearer ${token}`;
  } else {
    h["X-Anon-ID"] = anonId();
  }
  if (json) h["Content-Type"] = "application/json";
  return h;
}

export function profileURL() {
  return "/v1/profile";
}

/** 解析 401/429 等业务错误文案。 */
export async function apiErrorMessage(res, fallback = "") {
  let body = {};
  try {
    body = await res.json();
  } catch {
    /* ignore */
  }
  if (body.message) return body.message;
  if (res.status === 401) {
    if (body.code === 4014) return "匿名标识无效，请刷新页面后重试";
    return "未登录或登录已过期，请重新打开页面";
  }
  if (res.status === 429) return "请求过于频繁，请稍后再试";
  return fallback || `HTTP ${res.status}`;
}
