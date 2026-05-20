// @ts-nocheck — Harness Scaffold：MCP 客户端未启用
import type { McpRequestMessage, ServerConfig } from "./types";

export async function createClient(_id: string, _config: ServerConfig) {
  throw new Error("MCP is disabled in Harness scaffold");
}

export async function removeClient(_client: unknown) {}

export async function listTools(_client: unknown) {
  return { tools: [] };
}

export async function executeRequest(
  _client: unknown,
  _request: McpRequestMessage,
) {
  throw new Error("MCP is disabled in Harness scaffold");
}
