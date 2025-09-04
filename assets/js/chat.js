let connection;
let recipientId;

function handleMessageSend(event) {
  event.preventDefault();
  const formData = new FormData(this);
  const message = formData.get("message");
  if (message == "") {
    alert("message cannot be empty!")
    return
  }

  connection.send(JSON.stringify({
    recipient_id: recipientId,
    sender_id: 1,
    content: message,
    time: new Date(),
  }))
}

function main() {
  const messageForm = document.querySelector("#message-box");
  const messageBox = document.querySelector("#chatbox .main");
  const id = location.href.split("/").at(-1);
  recipientId = parseInt(id);
  const websocket = new WebSocket(`/chat/${id}`);
  connection = websocket;

  websocket.onopen = () => {
    console.log("working")
  }

  websocket.onclose = e => {
    console.log(e.code);
    console.log(e);
  }

  websocket.onmessage = e => {
    const message = JSON.parse(e.data)
    let sent = '';
    if (message.Message) {
      sent = 'right';
    }
    messageBox.insertAdjacentHTML("beforeend", `<p class="message ${sent}">${message.Message ? message.Message : message.content}</p>`);
  }

  messageForm.addEventListener("submit", handleMessageSend)
}

document.addEventListener("DOMContentLoaded", main);
