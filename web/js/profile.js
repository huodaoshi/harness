import { apiErrorMessage, apiHeaders, profileURL } from "./user.js";

const peopleList = () => document.getElementById("people-list");
const selfNote = () => document.getElementById("profile-self-note");
const currentIssue = () => document.getElementById("profile-current-issue");
const profileStatus = () => document.getElementById("profile-status");
const profileForm = () => document.getElementById("profile-form");

let returnTo = "welcome";

function personRow(person = { label: "", relation: "", note: "" }) {
  const row = document.createElement("div");
  row.className = "person-row";
  row.innerHTML = `
    <input type="text" class="person-label" placeholder="称呼" value="${escapeAttr(person.label)}" />
    <input type="text" class="person-relation" placeholder="关系" value="${escapeAttr(person.relation)}" />
    <input type="text" class="person-note" placeholder="备注" value="${escapeAttr(person.note)}" />
    <button type="button" class="btn-icon btn-remove-person" aria-label="删除">×</button>
  `;
  row.querySelector(".btn-remove-person").addEventListener("click", () => {
    row.remove();
  });
  return row;
}

function escapeAttr(s) {
  return String(s ?? "")
    .replace(/&/g, "&amp;")
    .replace(/"/g, "&quot;")
    .replace(/</g, "&lt;");
}

function collectPeople() {
  return [...peopleList().querySelectorAll(".person-row")].map((row) => ({
    label: row.querySelector(".person-label")?.value.trim() ?? "",
    relation: row.querySelector(".person-relation")?.value.trim() ?? "",
    note: row.querySelector(".person-note")?.value.trim() ?? "",
  })).filter((p) => p.label || p.relation || p.note);
}

function setStatus(msg, isError = false) {
  const el = profileStatus();
  el.textContent = msg;
  el.hidden = !msg;
  el.classList.toggle("error-bar", isError);
  el.classList.toggle("profile-status--ok", !isError && !!msg);
}

export async function loadProfileForm() {
  setStatus("加载中…");
  try {
    const res = await fetch(profileURL(), { headers: apiHeaders(false) });
    if (res.status === 401 || res.status === 429) {
      throw new Error(await apiErrorMessage(res));
    }
    if (!res.ok) {
      throw new Error(await res.text());
    }
    const data = await res.json();
    selfNote().value = data.self?.note ?? "";
    currentIssue().value = data.current_issue ?? "";
    peopleList().replaceChildren();
    for (const p of data.people ?? []) {
      peopleList().appendChild(personRow(p));
    }
    setStatus("");
  } catch (err) {
    setStatus(err instanceof Error ? err.message : "加载失败", true);
  }
}

export async function saveProfile() {
  const body = {
    self: { note: selfNote().value.trim() },
    people: collectPeople(),
    current_issue: currentIssue().value.trim(),
  };
  setStatus("保存中…");
  try {
    const res = await fetch(profileURL(), {
      method: "PUT",
      headers: apiHeaders(),
      body: JSON.stringify(body),
    });
    if (res.status === 401 || res.status === 429) {
      throw new Error(await apiErrorMessage(res));
    }
    if (!res.ok) {
      throw new Error(await res.text());
    }
    setStatus("已保存");
    return true;
  } catch (err) {
    setStatus(err instanceof Error ? err.message : "保存失败", true);
    return false;
  }
}

export function initProfileScreen({ onNavigate }) {
  document.getElementById("btn-add-person")?.addEventListener("click", () => {
    peopleList().appendChild(personRow());
  });

  document.getElementById("btn-profile-back")?.addEventListener("click", () => {
    onNavigate(returnTo);
  });

  profileForm()?.addEventListener("submit", async (e) => {
    e.preventDefault();
    const ok = await saveProfile();
    if (ok) {
      setTimeout(() => onNavigate(returnTo), 400);
    }
  });

  return {
    open(fromScreen) {
      returnTo = fromScreen;
      loadProfileForm();
      onNavigate("profile");
    },
  };
}
