class ChatBox {
  constructor() {
    // html elements
    this.messageInput = document.querySelector("#test-message");
    this.chatbox = document.querySelector("#chatbox .main");
    this.messageForm = document.querySelector("#message-box");
    this.recordBtn = document.querySelector("#record-audio");

    // initializing empty connection
    this.connection = null;

    // recipientID
    this.recipientId = null;

    // connection data
    this.messageTypes = {
      TEXT_TYPE: "TEXT",
      AUDIO_TYPE: "AUDIO",
    };
    Object.freeze(this.messageTypes);
  }

  async initialize(id) {
    // closing old connections to prevent leak
    if (this.connection != null) {
      this.connection.close();
      this.connection = null;
    }

    this.recipientId = id;
    await this.loadMessages();
    this.connection = new WebSocket(`/chat/${this.recipientId}`);

    // recorder
    this.mediaRecorder = null;
    this.audioChunks = [];

    this.connection.onmessage = e => {
      this.manageMessageHTML(e);
    };
  }

  // sending message to the websocket
  sendMessage(messageType, content) {
    const rawData = {
      recipient_id: this.recipientId,
      message_type: messageType,
      time: new Date(),
    };
    if (messageType == this.messageTypes.AUDIO_TYPE) {
      rawData["content_data"] = content;
    } else {
      rawData["content"] = content;
    }
    const data = JSON.stringify(rawData);

    if (!this.connection.readyState == WebSocket.OPEN) {
      // will do proper handling later
      alert("cannot send");
      return
    }
    this.connection.send(data)
  }

  // handle text message
  handleTextMessage() {
    const message = this.messageInput.value;
    if (message == "") {
      alert("cannot send empty message!");
      return
    }
    this.sendMessage(this.messageTypes.TEXT_TYPE, message);
    this.messageInput.value = "";
  }

  // handle audio message
  async handleAudioMessage() {
    const constraints = {
      audio: true,
    };
    try {
      const stream = await navigator.mediaDevices.getUserMedia(constraints);
      this.mediaRecorder = new MediaRecorder(stream);
    } catch (error) {
      console.log(error);
      return;
    }

    this.mediaRecorder.ondataavailable = e => {
      this.audioChunks.push(e.data);
    }

    this.mediaRecorder.onstop = async () => {
      const blob = new Blob(this.audioChunks);
      this.audioChunks = [];

      try {
        const arrayBuffer = await blob.arrayBuffer();
        const rawData = new Uint8Array(arrayBuffer);
        const data = Array.from(rawData);
        this.sendMessage(this.messageTypes.AUDIO_TYPE, data);
      } catch (error) {
        console.log(error);
      } finally {
        this.mediaRecorder = null;
      }
    }

    this.addRecordingPopup();
    this.mediaRecorder.start();
  }

  // adds a popup for recording
  addRecordingPopup() {
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

    stopButton.onclick = () => {
      console.log('Stop recording clicked');
      this.mediaRecorder.stop();
      popup.remove();
    };

    popup.appendChild(message);
    popup.appendChild(stopButton);

    document.body.appendChild(popup);
  }

  // handle dom manipulation
  manageMessageHTML(e) {
    const rawData = e.data;
    const data = JSON.parse(rawData);

    if (data.MessageType == this.messageTypes.TEXT_TYPE) {
      const message = this.createTextMessage(data);
      this.chatbox.append(message);
      return
    }
    const message = this.createAudioMessage(data);
    this.chatbox.append(message);
  }

  // create text message
  createTextMessage(data) {
    const message = document.createElement("div");
    message.classList.add("message");
    if (data.RecipientID == this.recipientId) {
      message.classList.add("right");
    }

    const p = document.createElement("p");
    p.innerText = data.Content;

    message.append(p);
    return message
  }

  // create audio message
  createAudioMessage(data) {
    const message = document.createElement("div");
    message.classList.add("message");
    if (data.RecipientID == this.recipientId) {
      message.classList.add("right");
    }

    const audio = document.createElement("audio");
    audio.src = `/${data.Content}`;
    audio.controls = true;

    message.append(audio);
    return message
  }

  // adding event listeners to respective elements
  addEventListeners() {
    this.messageForm.addEventListener("submit", e => {
      e.preventDefault();
      this.handleTextMessage();
    })

    this.recordBtn.addEventListener("click", () => {
      this.handleAudioMessage();
    });
  }

  async loadMessages() {
    try {
      const response = await fetch(`/users/${this.recipientId}`);
      if (!response.ok) {
        throw new Error("failed to load messages")
      }
      const data = await response.text();
      this.chatbox.innerHTML = data;
    } catch (error) {
      console.log(error);
    }
  }
}

function handleClickUser(target, chatbox) {
  document.querySelector("#chatbox").classList.add("active");
  chatbox.initialize(target.id);
}

function main() {
  const allUsers = document.querySelectorAll(".user");
  const chatbox = new ChatBox();

  chatbox.addEventListeners();

  allUsers.forEach(user => {
    user.addEventListener("click", e => {
      handleClickUser(e.target, chatbox);
    })
  })
}

document.addEventListener("DOMContentLoaded", main);
