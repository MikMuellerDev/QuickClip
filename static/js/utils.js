function setButtonText(loggedIn) {
  const button = document.getElementById("loggout-button");
  if (loggedIn) {
    button.innerText = "logout";
    button.onclick = () => {
      window.location.href = "/logout";
    };
  } else {
    button.innerText = "login";
    button.onclick = () => {
      window.location.href = "/login";
    };
  }
}
