let connection;
let recipientId;

function createPopup(mediaRecorder) {
  const popup = document.createElement('div');

  popup.style.position = 'absolute';
  popup.style.top = '50%';
  popup.style.left = '50%';
  popup.style.transform = 'translate(-50%, -50%)';
  popup.style.width = '50%';
  popup.style.height = '50%';
  popup.style.background = 'green';
  popup.style.border = '2px solid white';
  popup.style.display = 'flex';
  popup.style.flexDirection = 'column';
  popup.style.justifyContent = 'center';
  popup.style.alignItems = 'center';
  popup.style.zIndex = '1';

  const message = document.createElement('p');
  message.textContent = 'RECORDING YOUR MESSAGE!';

  const stopButton = document.createElement('button');
  stopButton.textContent = 'Click here to stop';

  stopButton.style.outline = 'none';
  stopButton.style.border = '1px solid white';
  stopButton.style.padding = '20px';
  stopButton.style.cursor = 'pointer';

  stopButton.onclick = function() {
    console.log('Stop recording clicked');
    mediaRecorder.stop();
    popup.remove();
  };

  popup.appendChild(message);
  popup.appendChild(stopButton);

  document.body.appendChild(popup);
}

async function handleRecording(stream) {
  const mediaRecorder = new MediaRecorder(stream);
  let chunks = [];

  createPopup(mediaRecorder);

  mediaRecorder.ondataavailable = e => {
    chunks.push(e.data);
  }

  mediaRecorder.onstop = async () => {
    console.log("recording stopped");
    console.log(chunks);

    const blob = new Blob(chunks, { type: mediaRecorder.mimeType });
    chunks = [];
    try {
      const data = await blob.arrayBuffer();
      const dataToSend = new Uint8Array(data);
      console.log(Array.from(dataToSend));
    } catch (error) {
      console.log(err);
      alert("failed to send the message!");
    }
  }

  mediaRecorder.start();
}

async function handleRecordAudio() {
  const constraints = {
    audio: true,
  }

  let timeout = setTimeout(() => {
    alert("please choose and option!");
    return
  }, 5000)

  try {
    const stream = await navigator.mediaDevices.getUserMedia(constraints);
    handleRecording(stream);
  } catch (error) {
    console.log(error);
    alert("please allow microphone permissions to send a voice message!")
  } finally {
    clearTimeout(timeout);
    timeout = null;
  }
}

function handleMessageSend(event) {
  event.preventDefault();
  const formData = new FormData(this);
  const message = formData.get("message");

  this.querySelector('[name="message"]').value = "";

  if (message == "") {
    alert("message cannot be empty!")
    return
  }

  if (connection.readyState == WebSocket.OPEN) {
    connection.send(JSON.stringify({
      recipient_id: recipientId,
      sender_id: 1,
      content: message,
      time: new Date(),
    }))
  }
}

function main() {
  const messageForm = document.querySelector("#message-box");
  const messageBox = document.querySelector("#chatbox .main");
  const recordAudioBtn = document.querySelector("#record-audio");
  const id = location.href.split("/").at(-1);
  recipientId = parseInt(id);
  const websocket = new WebSocket(`/chat/${id}`);
  connection = websocket;

  websocket.onopen = () => {
    console.log("working")
  }

  websocket.onclose = e => {
    console.log(e);
  }

  websocket.onmessage = e => {
    const message = JSON.parse(e.data)
    let sent = '';
    if (message.Content) {
      sent = 'right';
    }
    messageBox.insertAdjacentHTML("beforeend", `<p class="message ${sent}">${message.Content ? message.Content : message.content}</p>`);
  }

  messageForm.addEventListener("submit", handleMessageSend)
  recordAudioBtn.addEventListener("click", handleRecordAudio)
}

document.addEventListener("DOMContentLoaded", main);
