import { consumeSSE } from "./sse.js";
import { userId } from "./user.js";
import { initProfileScreen } from "./profile.js";

const STORAGE_SESSION = "fwa_session_id";

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
const chatModeLabel = $("chat-mode-label");
const sessionHint = $("session-hint");
const messagesEl = $("messages");
const crisisBanner = $("crisis-banner");
const crisisBody = $("crisis-body");
const errorBar = $("error-bar");
const chatForm = $("chat-form");
const messageInput = $("message-input");
const btnSend = $("btn-send");

/** @type {"distress"|"normal"} */
let mode = "normal";
let sessionId = sessionStorage.getItem(STORAGE_SESSION) ?? "";
let streaming = false;
let inCrisis = false;

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

function enterChat(selectedMode) {
  mode = selectedMode;
  chatModeLabel.textContent =
    mode === "distress" ? "洪峰陪伴" : "轻松聊聊";
  sessionHint.textContent = sessionId
    ? `会话 ${sessionId.slice(0, 8)}…`
    : "";
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
  messageInput.disabled = active || inCrisis;
  btnSend.disabled = active || inCrisis;
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

async function streamMessage(text) {
  clearError();
  setStreaming(true);

  appendBubble("user", text);
  const assistantEl = appendBubble("assistant");
  assistantEl.classList.add("bubble--typing");

  const payload = {
    session_id: sessionId || undefined,
    user_id: userId(),
    mode,
    message: text,
  };

  try {
    const res = await fetch("/v1/sessions/stream", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

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
  if (streaming || inCrisis) return;

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
