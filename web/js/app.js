import { consumeSSE } from "./sse.js";
import { apiErrorMessage, apiHeaders } from "./user.js";
import { initProfileScreen } from "./profile.js";

const STORAGE_SESSION = "fwa_session_id";
const STORAGE_SUMMARY_DISMISS = "fwa_summary_dismissed_";

const $ = (id) => document.getElementById(id);

const screenWelcome = $("screen-welcome");
const screenChat = $("screen-chat");
const screenProfile = $("screen-profile");
const disclaimerCheck = $("disclaimer-check");
const btnDistress = $("btn-distress");
const btnNormal = $("btn-normal");
const btnProfileWelcome = $("btn-profile-welcome");
const btnProfileChat = $("btn-profile-chat");
const btnBack = $("btn-back");
const btnEndSession = $("btn-end-session");
const chatModeLabel = $("chat-mode-label");
const sessionHint = $("session-hint");
const messagesEl = $("messages");
const crisisBanner = $("crisis-banner");
const crisisBody = $("crisis-body");
const summaryCard = $("summary-card");
const summaryList = $("summary-list");
const btnSummaryClose = $("btn-summary-close");
const errorBar = $("error-bar");
const chatForm = $("chat-form");
const messageInput = $("message-input");
const btnSend = $("btn-send");

/** @type {"distress"|"normal"} */
let mode = "normal";
let sessionId = sessionStorage.getItem(STORAGE_SESSION) ?? "";
let streaming = false;
let inCrisis = false;
let sessionEnded = false;

function setButtonsEnabled(enabled) {
  btnDistress.disabled = !enabled;
  btnNormal.disabled = !enabled;
  if (btnProfileWelcome) btnProfileWelcome.disabled = !enabled;
}

disclaimerCheck.addEventListener("change", () => {
  setButtonsEnabled(disclaimerCheck.checked);
});

function showScreen(name) {
  const welcome = name === "welcome";
  const chat = name === "chat";
  const profile = name === "profile";
  screenWelcome.classList.toggle("screen--active", welcome);
  screenWelcome.hidden = !welcome;
  screenChat.classList.toggle("screen--active", chat);
  screenChat.hidden = !chat;
  if (screenProfile) {
    screenProfile.classList.toggle("screen--active", profile);
    screenProfile.hidden = !profile;
  }
}

const profileUI = initProfileScreen({ onNavigate: showScreen });

function resetChatSession() {
  sessionId = "";
  sessionEnded = false;
  inCrisis = false;
  sessionStorage.removeItem(STORAGE_SESSION);
  messagesEl.replaceChildren();
  crisisBanner.hidden = true;
  summaryCard.hidden = true;
  screenChat.classList.remove("screen--crisis");
  sessionHint.textContent = "";
}

function enterChat(selectedMode) {
  resetChatSession();
  mode = selectedMode;
  chatModeLabel.textContent =
    mode === "distress" ? "洪峰陪伴" : "轻松聊聊";
  showScreen("chat");
  messageInput.focus();
}

btnDistress.addEventListener("click", () => enterChat("distress"));
btnNormal.addEventListener("click", () => enterChat("normal"));

btnProfileWelcome?.addEventListener("click", () => profileUI.open("welcome"));
btnProfileChat?.addEventListener("click", () => profileUI.open("chat"));

btnBack.addEventListener("click", () => {
  if (streaming) return;
  showScreen("welcome");
});

function showError(msg) {
  errorBar.textContent = msg;
  errorBar.hidden = false;
}

function clearError() {
  errorBar.hidden = true;
  errorBar.textContent = "";
}

function appendBubble(role, text = "") {
  const el = document.createElement("div");
  el.className = `bubble bubble--${role}`;
  el.textContent = text;
  messagesEl.appendChild(el);
  messagesEl.scrollTop = messagesEl.scrollHeight;
  return el;
}

function setStreaming(active) {
  streaming = active;
  messageInput.disabled = active || inCrisis || sessionEnded;
  btnSend.disabled = active || inCrisis || sessionEnded;
  if (btnEndSession) btnEndSession.disabled = active || inCrisis || !sessionId;
}

function showGateCard(body, title = "你此刻的安全最重要") {
  inCrisis = true;
  const heading = crisisBanner.querySelector("h2");
  if (heading) heading.textContent = title;
  crisisBody.textContent = body;
  crisisBanner.hidden = false;
  screenChat.classList.add("screen--crisis");
  messageInput.disabled = true;
  btnSend.disabled = true;
}

function summaryDismissKey(id) {
  return STORAGE_SUMMARY_DISMISS + id;
}

function showSummaryCard(summary3) {
  if (!sessionId || sessionStorage.getItem(summaryDismissKey(sessionId)) === "1") {
    return;
  }
  summaryList.replaceChildren();
  for (const line of summary3) {
    const li = document.createElement("li");
    li.textContent = line;
    summaryList.appendChild(li);
  }
  summaryCard.hidden = false;
}

function hideSummaryCard() {
  summaryCard.hidden = true;
  if (sessionId) {
    sessionStorage.setItem(summaryDismissKey(sessionId), "1");
  }
}

btnSummaryClose?.addEventListener("click", hideSummaryCard);

async function endSession() {
  if (!sessionId || sessionEnded || streaming) return;
  clearError();
  try {
    const res = await fetch("/v1/sessions/end", {
      method: "POST",
      headers: apiHeaders(),
      body: JSON.stringify({ session_id: sessionId }),
    });
    if (res.status === 401 || res.status === 429) {
      throw new Error(await apiErrorMessage(res));
    }
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error ?? err.message ?? `HTTP ${res.status}`);
    }
    const data = await res.json();
    sessionEnded = true;
    messageInput.disabled = true;
    btnSend.disabled = true;
    if (btnEndSession) btnEndSession.disabled = true;
    if (Array.isArray(data.summary3)) {
      showSummaryCard(data.summary3);
    }
  } catch (err) {
    showError(err instanceof Error ? err.message : "结束会话失败");
  }
}

btnEndSession?.addEventListener("click", endSession);

async function streamMessage(text) {
  if (sessionEnded) {
    showError("本场对话已结束，请返回首页开始新会话。");
    return;
  }
  clearError();
  setStreaming(true);

  appendBubble("user", text);
  const assistantEl = appendBubble("assistant");
  assistantEl.classList.add("bubble--typing");

  const payload = {
    session_id: sessionId || undefined,
    mode,
    message: text,
  };

  try {
    const res = await fetch("/v1/sessions/stream", {
      method: "POST",
      headers: apiHeaders(),
      body: JSON.stringify(payload),
    });

    if (res.status === 401 || res.status === 429) {
      showError(await apiErrorMessage(res));
      assistantEl.remove();
      return;
    }

    if (res.status === 409) {
      const err = await res.json().catch(() => ({}));
      showError(err.error ?? "会话已满，请结束本场对话。");
      assistantEl.remove();
      return;
    }

    await consumeSSE(res, async ({ event, data }) => {
      let parsed;
      try {
        parsed = JSON.parse(data);
      } catch {
        parsed = { raw: data };
      }

      switch (event) {
        case "token":
          assistantEl.textContent += parsed.text ?? "";
          messagesEl.scrollTop = messagesEl.scrollHeight;
          break;
        case "done":
          assistantEl.classList.remove("bubble--typing");
          if (parsed.session_id) {
            sessionId = parsed.session_id;
            sessionStorage.setItem(STORAGE_SESSION, sessionId);
            sessionHint.textContent = `会话 ${sessionId.slice(0, 8)}…`;
          }
          if (btnEndSession) btnEndSession.disabled = false;
          break;
        case "crisis":
          assistantEl.remove();
          showGateCard(
            parsed.body ?? "请尽快寻求身边或专业帮助。",
            "你此刻的安全最重要"
          );
          break;
        case "medical":
          assistantEl.remove();
          showGateCard(
            parsed.body ?? "我无法提供医学建议，请联系专业机构。",
            "关于医疗与用药"
          );
          break;
        case "error":
          assistantEl.classList.remove("bubble--typing");
          if (parsed.code === "content_blocked") {
            showError(parsed.message ?? "这条内容无法继续讨论");
          } else if (parsed.code === "session_message_cap") {
            showError(parsed.message ?? "本场对话已达到消息上限");
          } else {
            showError(parsed.message ?? parsed.code ?? "服务暂时不可用");
          }
          break;
        default:
          break;
      }
    });
  } catch (err) {
    assistantEl.classList.remove("bubble--typing");
    showError(err instanceof Error ? err.message : "网络错误");
  } finally {
    setStreaming(false);
  }
}

chatForm.addEventListener("submit", (e) => {
  e.preventDefault();
  if (streaming || inCrisis || sessionEnded) return;

  const text = messageInput.value.trim();
  if (!text) return;

  messageInput.value = "";
  messageInput.style.height = "auto";
  streamMessage(text);
});

messageInput.addEventListener("input", () => {
  messageInput.style.height = "auto";
  messageInput.style.height = `${Math.min(messageInput.scrollHeight, 120)}px`;
});

messageInput.addEventListener("keydown", (e) => {
  if (e.key === "Enter" && !e.shiftKey) {
    e.preventDefault();
    chatForm.requestSubmit();
  }
});
