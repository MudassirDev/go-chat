function handleRegisterForm() {

}

async function handleLoginForm(form) {
  const formData = new FormData(form);
  const data = {};
  let sendData = false;

  for (const key of formData.keys()) {
    sendData = true;
    data[key] = formData.get(key);
  }

  const url = form.getAttribute("action");
  const method = form.getAttribute("method");
  const options = createOptions(method, sendData ? data : null);

  try {
    const response = await fetch(url, options);
    const data = await response.text();
    if (!response.ok) {
      throw new Error(data);
    }
    console.log(data);
    location.href = "/";
  } catch (error) {
    console.log(error);
  }
}

function createOptions(method, data) {
  const options = {
    method: method,
  }

  if (data) {
    options.body = JSON.stringify(data);
  }

  if (method == "POST") {
    options.headers = {
      "Content-Type": "application/json",
    }
  }

  return options
}

function main() {
  const form = document.querySelector("form");
  const isRegister = location.href.split("/").at(-1) == "register";

  form.addEventListener("submit", async e => {
    e.preventDefault();

    if (isRegister) {
      handleRegisterForm()
      return;
    }

    handleLoginForm(e.target)
  })
}

document.addEventListener("DOMContentLoaded", main);
