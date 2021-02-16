export default class PubSub {
  constructor() {
    const { protocol, host } = document.location;
    this.ws = new WebSocket(protocol.replace("http", "ws") + host);
    this.listeners = new Set();

    this.ws.addEventListener("message", ({ data }) => {
      try {
        const json = JSON.parse(data);
        for (const listener of this.listeners) {
          if (listener.key === json.key) {
            listener.callback(json.value);
          }
        }
      } catch (e) {}
    });
  }

  on(key, callback) {
    this.listeners.add({
      key,
      callback,
    });
  }

  emit(key, value) {
    this.ws.send(JSON.stringify({ key, value }));
  }
}
