async function newUserPopup(user: User, create: boolean) {
  const main = document.createElement("div");
  let passwordShown = false;
  main.className = "edit-overlay";
  console.log(user.Name);

  main.id = user.Name;

  const name = document.createElement("h2");
  name.innerText = user.Name;
  name.className = "popup-doc-title";

  const mainDiv = document.createElement("div");
  mainDiv.className = "edit-overlay-main-div";

  const sliderDiv = document.createElement("div");
  sliderDiv.className = "edit-overlay-slider-div";

  const inputDiv = document.createElement("div");
  inputDiv.className = "edit-overlay-inputs";

  const labelDiv = document.createElement("div");
  labelDiv.className = "edit-overlay-labels";

  const buttonsDiv = document.createElement("div");
  buttonsDiv.className = "edit-overlay-buttons";

  const nameInput = createInput(user.Name, "Username");
  if (!create) {
    nameInput.disabled = true;
  }
  const passwordInput = createInput(user.Password, "Password");
  passwordInput.type = "password";

  passwordInput.onclick = () => {
    if (passwordShown) {
      passwordInput.type = "password";
    } else {
      passwordInput.type = "text";
    }
    passwordShown = !passwordShown;
  };

  const deleteButton = createButton("delete");
  deleteButton.style.backgroundColor = "var(--clr-error)";

  let saveButton: HTMLButtonElement;

  if (create) {
    saveButton = createButton("create");
  } else {
    saveButton = createButton("save");
  }

  const cancelButton = createButton("cancel");

  const nameSpan = createSpan("Name:", "");
  const passwordSpan = createSpan("Password:", "");

  labelDiv.appendChild(nameSpan);
  labelDiv.appendChild(passwordSpan);

  nameInput.oninput = () => {
    name.innerText = nameInput.value;
  };

  deleteButton.onclick = async () => {
    const confirm = prompt("Enter `yes` to confirm deletion", "");
    if (confirm == "yes") {
      await deleteHandler(user.Name);
      removeUsers();
      const users = await getUsers();
      if (users !== null) addUsers(users);
    }
    setTimeout(() => {
      main.remove();
    }, 600);
    hidePopup(user.Name);
  };

  saveButton.onclick = async () => {
    const tempUsers = await getUsers();
    if (tempUsers != null && user.Name !== nameInput.value) {
      for (let iterUser of tempUsers) {
        if (iterUser.Name == nameInput.value) {
          saveButton.style.backgroundColor = "var(--clr-error)";
          setTimeout(() => {
            saveButton.style.backgroundColor = "var(--clr-purple)";
          }, 2000);
          alert("This user already exists.");
          return;
        }
      }
    }
    if (!create) {
      modifyUser(user.Name, {
        Name: nameInput.value,
        Password: passwordInput.value,
        Permissions: user.Permissions,
        WriteAllowed: user.WriteAllowed,
      });
    } else {
      createUser({
        Name: nameInput.value,
        Password: passwordInput.value,
        Permissions: user.Permissions,
        WriteAllowed: user.WriteAllowed,
      });
    }

    hidePopup(user.Name);
    setTimeout(() => {
      main.remove();
    }, 600);

    removeUsers();
    addUsers(await getUsers());
  };

  cancelButton.onclick = () => {
    hidePopup(user.Name);
    setTimeout(() => {
      main.remove();
    }, 600);
  };

  main.appendChild(name);
  inputDiv.appendChild(nameInput);
  inputDiv.appendChild(passwordInput);

  if (!create) {
    buttonsDiv.appendChild(deleteButton);
  }

  buttonsDiv.appendChild(cancelButton);
  buttonsDiv.appendChild(saveButton);

  mainDiv.appendChild(labelDiv);
  mainDiv.appendChild(inputDiv);
  mainDiv.appendChild(buttonsDiv);
  main.appendChild(mainDiv);
  main.appendChild(sliderDiv);
  main.style.zIndex = "-1000";
  main.style.display = "none";

  const clipboards = await getDocuments();
  for (let clip of clipboards) {
    const traitSwitchArr = createSlider();
    const traitSwitch = traitSwitchArr[0];
    const traitSwitchListener = traitSwitchArr[1];

    traitSwitchListener.checked = user.Permissions.includes(clip.Name);

    traitSwitchListener.onchange = () => {
      if (traitSwitchListener.checked) {
        if (!user.Permissions.includes(clip.Name)) {
          user.Permissions.push(clip.Name);
        }
      } else {
        let permissionsTemp: string[] = [];
        for (let item of user.Permissions) {
          if (item !== clip.Name) {
            permissionsTemp.push(clip.Name);
          }
        }
        user.Permissions = permissionsTemp;
      }
      console.log(user.Permissions);
    };

    const switchContainer = document.createElement("div");
    switchContainer.className = "inner-slider-div";

    const switchContainerSpan = document.createElement("span");
    switchContainerSpan.innerText = clip.Name;

    switchContainer.appendChild(traitSwitch);
    switchContainer.appendChild(switchContainerSpan);
    sliderDiv.appendChild(switchContainer);
  }

  for (let clip of clipboards) {
    const traitSwitchArr = createSlider();
    const traitSwitch = traitSwitchArr[0];
    const traitSwitchListener = traitSwitchArr[1];

    traitSwitchListener.checked = user.WriteAllowed.includes(clip.Name);

    traitSwitchListener.onchange = () => {
      if (traitSwitchListener.checked) {
        if (!user.WriteAllowed.includes(clip.Name)) {
          user.WriteAllowed.push(clip.Name);
        }
      } else {
        let permissionsTemp: string[] = [];
        for (let item of user.WriteAllowed) {
          if (item !== clip.Name) {
            permissionsTemp.push(clip.Name);
          }
        }
        user.WriteAllowed = permissionsTemp;
      }
      console.log(user.WriteAllowed);
    };

    const switchContainer = document.createElement("div");
    switchContainer.className = "inner-slider-div";

    const switchContainerSpan = document.createElement("span");
    switchContainerSpan.innerText = clip.Name;

    switchContainer.appendChild(traitSwitch);
    switchContainer.appendChild(switchContainerSpan);
    sliderDiv.appendChild(switchContainer);
  }

  document.getElementsByTagName("body")[0].appendChild(main);

  return main;
}
