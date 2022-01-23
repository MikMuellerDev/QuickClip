function setButtonText(loggedIn) {
  const button = document.getElementById("logout-button");
  const buttonLabel = document.getElementById("logout-button-text");

  if (loggedIn) {
    buttonLabel.innerText = "logout";
    button.onclick = () => {
      window.location.href = "/logout";
    };
  } else {
    buttonLabel.innerText = "login";
    button.onclick = () => {
      window.location.href = "/login";
    };
  }
}
