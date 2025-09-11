async function handleForm(form) {
  const formData = new FormData(form);
  const url = form.getAttribute("action");
  const method = form.getAttribute("method");
  const data = {
    username: formData.get("username"),
    password: formData.get("password")
  }

  const confirmPassword = formData.get("confirm-password");

  if (confirmPassword && confirmPassword != data["password"]) {
    alert("passwords are not the same");
    return
  }

  const options = createOptions(method, data);
  try {
    const response = await fetch(url, options);
    const data = await response.text();
    if (!response.ok) {
      throw new Error(data);
    }
    console.log(data);
    location.href = "/";
  } catch (error) {
    alert(error);
  }
}

function createOptions(method, data) {
  const options = {
    method: method,
    body: JSON.stringify(data),
    headers: {
      "Content-Type": "application/json",
    }
  }

  return options
}

function main() {
  const form = document.querySelector("form");

  form.addEventListener("submit", async e => {
    e.preventDefault();

    handleForm(e.target);
    return;
  })
}

document.addEventListener("DOMContentLoaded", main);
