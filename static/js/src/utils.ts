interface Clip {
  Name: string;
  Id: string;
  Content: string;
  Description: string;
  Restricted: boolean;
  Refresh: boolean;
  RefreshInterval: number;
  ReadOnly: boolean;
}

function setButtonText(loggedIn: boolean) {
  const button = document.getElementById("logout-button") as HTMLButtonElement;
  const buttonLabel = document.getElementById(
    "logout-button-text"
  ) as HTMLSpanElement;

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

function createNavbarButton(
  imageSrc: string,
  label: string,
  onclickCallback: Function
) {
  const main = document.createElement("div");
  main.className = "navbar-button";
  const image = document.createElement("img");
  image.src = imageSrc;

  const span = document.createElement("span") as HTMLSpanElement;
  span.innerText = label;
  main.appendChild(image);
  main.appendChild(span);
  main.onclick = () => {
    onclickCallback();
  };
  return main;
}

function addPasswordChangeButton() {
  const navbar = document.getElementById("navbar") as HTMLElement;
  const changePasswordButton = createNavbarButton(
    "/static/media/user.png",
    "user",
    showPasswordChangePopup
  );
  navbar.appendChild(changePasswordButton);
}

function showPasswordChangePopup() {
  let passwordShown = false;
  const main = document.createElement("div");
  main.className = "edit-overlay";
  main.id = "password-change-dialog";

  const name = document.createElement("h2");
  name.innerText = "Change Password";
  name.className = "popup-doc-title";

  const mainDiv = document.createElement("div");
  mainDiv.className = "edit-overlay-main-div";

  const inputDiv = document.createElement("div");
  inputDiv.className = "edit-overlay-inputs";

  const labelDiv = document.createElement("div");
  labelDiv.className = "edit-overlay-labels edit-overlay-labels-alt";

  const buttonsDiv = document.createElement("div");
  buttonsDiv.className = "edit-overlay-buttons";

  const passwordInput = createInput("", "Password");
  passwordInput.type = "password";

  passwordInput.onclick = () => {
    if (passwordShown) {
      passwordInput.type = "password";
    } else {
      passwordInput.type = "text";
    }
    passwordShown = !passwordShown;
  };

  const saveButton = createButton("update");
  const cancelButton = createButton("cancel");

  const passwordSpan = createSpan("Password:", "");

  labelDiv.appendChild(passwordSpan);

  saveButton.onclick = async () => {
    await changePassword(passwordInput.value);
    hidePopup("password-change-dialog");
    setTimeout(() => {
      main.remove();
    }, 600);
  };

  cancelButton.onclick = () => {
    hidePopup("password-change-dialog");
    setTimeout(() => {
      main.remove();
    }, 600);
  };

  main.appendChild(name);
  inputDiv.appendChild(passwordInput);
  buttonsDiv.appendChild(cancelButton);
  buttonsDiv.appendChild(saveButton);
  mainDiv.appendChild(labelDiv);
  mainDiv.appendChild(inputDiv);
  mainDiv.appendChild(buttonsDiv);
  main.appendChild(mainDiv);
  document.getElementsByTagName("body")[0].appendChild(main);
  showPopup("password-change-dialog");
}
