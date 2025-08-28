function main() {
  const websocket = new WebSocket("/ws"); // creating a websocket attempts to make a connection

  websocket.onopen = () => {
    console.log("connection opened");
    websocket.send("test");
  }

  websocket.onclose = () => {
    console.log("connection closed")
  }

  websocket.onmessage = e => {
    console.log(e.data);
  }

  setTimeout(() => {
    websocket.close();
  }, 1000)
}

document.addEventListener("DOMContentLoaded", main)
