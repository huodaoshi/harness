import { apiUrl, getApiBaseUrl } from "@/lib/harness/api-base";
import styles from "./page.module.scss";

/** 构建时不预渲染（避免无 backend 时 fetch 失败） */
export const dynamic = "force-dynamic";

async function probeConfig(): Promise<string> {
  try {
    const res = await fetch(apiUrl("/api/config"), {
      cache: "no-store",
    });
    if (!res.ok) return `HTTP ${res.status}`;
    const data = await res.json();
    return `needCode=${String(data.needCode)} hideUserApiKey=${String(data.hideUserApiKey)}`;
  } catch (e) {
    return e instanceof Error ? e.message : "fetch failed";
  }
}

export default async function Home() {
  const apiBase = getApiBaseUrl() || "(同域 rewrites → backend)";
  const configStatus = await probeConfig();

  return (
    <main className={styles.main}>
      <h1>Harness Chat</h1>
      <p className={styles.lead}>
        脚手架已就绪。请按 <code>MIGRATION.md</code> 从{" "}
        <code>NextChat/</code> 迁移业务代码。
      </p>

      <section className={styles.card}>
        <h2>连接</h2>
        <dl>
          <dt>API 基址</dt>
          <dd>
            <code>{apiBase}</code>
          </dd>
          <dt>GET /api/config</dt>
          <dd>
            <code>{configStatus}</code>
          </dd>
        </dl>
        <p className={styles.hint}>
          若 config 失败，请先启动 <code>backend</code>（:8080）并实现{" "}
          <code>/api/config</code> 路由。
        </p>
      </section>

      <section className={styles.card}>
        <h2>第一期目标</h2>
        <ul>
          <li>聊天主路径 + Mask</li>
          <li>Go：<code>/api/bytedance</code>、<code>/api/openai</code></li>
        </ul>
      </section>
    </main>
  );
}
