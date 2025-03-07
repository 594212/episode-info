const websocket: WebSocket = new WebSocket("/ws");
websocket.onopen = () => {
  console.log("OPEN");
};

websocket.onmessage = (_: MessageEvent) => {
  console.log("ws:// file updated");
  window.location.reload();
};

websocket.onclose = () => {
  console.log("CLOSED");
};
