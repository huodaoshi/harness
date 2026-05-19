/** 通用聊天（/api/* + 本地历史） */
export type ChatMode = "general" | "wellness";

const STORAGE_KEY = "harness-chat-mode";

export function getChatMode(): ChatMode {
  if (typeof window === "undefined") return "general";
  const v = localStorage.getItem(STORAGE_KEY);
  return v === "wellness" ? "wellness" : "general";
}

/** 第一期仅埋点；第二期由设置页调用 */
export function setChatMode(mode: ChatMode): void {
  if (typeof window === "undefined") return;
  localStorage.setItem(STORAGE_KEY, mode);
}
