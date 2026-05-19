/**
 * Parse SSE events from a fetch Response body (POST stream).
 * @param {Response} response
 * @param {(event: { event: string, data: string }) => void | Promise<void>} onEvent
 */
export async function consumeSSE(response, onEvent) {
  if (!response.ok) {
    let detail = "";
    try {
      detail = await response.text();
    } catch {
      /* ignore */
    }
    throw new Error(detail || `HTTP ${response.status}`);
  }

  const reader = response.body?.getReader();
  if (!reader) {
    throw new Error("浏览器不支持流式响应");
  }

  const decoder = new TextDecoder();
  let buffer = "";

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    buffer += decoder.decode(value, { stream: true });
    const parts = buffer.split("\n\n");
    buffer = parts.pop() ?? "";

    for (const block of parts) {
      const evt = parseSSEBlock(block);
      if (evt) {
        await onEvent(evt);
      }
    }
  }

  if (buffer.trim()) {
    const evt = parseSSEBlock(buffer);
    if (evt) {
      await onEvent(evt);
    }
  }
}

function parseSSEBlock(block) {
  const lines = block.split("\n");
  let event = "message";
  const dataLines = [];

  for (const line of lines) {
    if (line.startsWith("event:")) {
      event = line.slice(6).trim();
    } else if (line.startsWith("data:")) {
      dataLines.push(line.slice(5).trimStart());
    }
  }

  if (dataLines.length === 0) return null;
  return { event, data: dataLines.join("\n") };
}
