import { harnessAuthHeaders } from "./auth";
import { getChatMode } from "./chat-mode";

/** 合并 LLM / API 请求头；关怀模式时附加 X-Harness-Mode 与鉴权头 */
export function withHarnessHeaders(
  init?: HeadersInit,
): Record<string, string> {
  const out: Record<string, string> = {};
  if (init instanceof Headers) {
    init.forEach((v, k) => {
      out[k] = v;
    });
  } else if (Array.isArray(init)) {
    for (const [k, v] of init) out[k] = v;
  } else if (init) {
    Object.assign(out, init);
  }
  if (getChatMode() === "wellness") {
    out["X-Harness-Mode"] = "wellness";
    Object.assign(out, harnessAuthHeaders());
  }
  return out;
}
