const STORAGE_USER = "fwa_user_id";

/** Stable anonymous user id for MVP (isolated per browser). */
export function userId() {
  let id = localStorage.getItem(STORAGE_USER);
  if (!id) {
    id = crypto.randomUUID();
    localStorage.setItem(STORAGE_USER, id);
  }
  return id;
}

export function apiHeaders(json = true) {
  const h = { "X-User-Id": userId() };
  if (json) h["Content-Type"] = "application/json";
  return h;
}

export function profileURL() {
  return `/v1/profile?user_id=${encodeURIComponent(userId())}`;
}
