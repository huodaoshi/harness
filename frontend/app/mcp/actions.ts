"use server";

import {
  DEFAULT_MCP_CONFIG,
  McpConfigData,
  McpRequestMessage,
  ServerConfig,
  ServerStatusResponse,
} from "./types";

/** Harness Scaffold：MCP 未启用，Server Actions 为占位实现。 */

export async function getClientsStatus(): Promise<
  Record<string, ServerStatusResponse>
> {
  return {};
}

export async function getClientTools(_clientId: string) {
  return [];
}

export async function getAvailableClientsCount() {
  return 0;
}

export async function getAllTools() {
  return [];
}

export async function initializeMcpSystem() {}

export async function addMcpServer(
  _clientId: string,
  _config: ServerConfig,
): Promise<McpConfigData> {
  return DEFAULT_MCP_CONFIG;
}

export async function pauseMcpServer(_clientId: string): Promise<McpConfigData> {
  return DEFAULT_MCP_CONFIG;
}

export async function resumeMcpServer(_clientId: string): Promise<McpConfigData> {
  return DEFAULT_MCP_CONFIG;
}

export async function removeMcpServer(_clientId: string): Promise<McpConfigData> {
  return DEFAULT_MCP_CONFIG;
}

export async function restartAllClients(): Promise<McpConfigData> {
  return DEFAULT_MCP_CONFIG;
}

export async function executeMcpAction(
  _clientId: string,
  _request: McpRequestMessage,
) {
  throw new Error("MCP is disabled in Harness scaffold");
}

export async function getMcpConfigFromFile(): Promise<McpConfigData> {
  return DEFAULT_MCP_CONFIG;
}

export async function isMcpEnabled() {
  return false;
}
