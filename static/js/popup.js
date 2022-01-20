function createInput(label) {
  const main = document.createElement("input");
  main.className = "text-input";
  main.type = "text";
  main.placeholder = label;
  return main;
}

function createButton(label) {
  const button = document.createElement("button");
  const buttonText = document.createElement("a");
  buttonText.innerText = label;
  button.className = "button";
  button.appendChild(buttonText);
  return button;
}

function newPopup(id) {
  const main = document.createElement("div");
  main.className = "edit-overlay";
  main.id = id;

  const title = document.createElement("h2");

  const description = document.createElement("span");
  description.className = "edit-overlay-description";

  const mainDiv = document.createElement("div");
  mainDiv.className = "edit-overlay-main-div";

  const inputDiv = document.createElement("div");
  inputDiv.className = "edit-overlay-inputs";

  const buttonsDiv = document.createElement("div");
  buttonsDiv.className = "edit-overlay-buttons";

  const refreshIntervalInput = createInput("Refresh Interval");
  const descriptionInput = createInput("Description");
  const nameInput = createInput("Name");
  const deleteButton = createButton("delete");
  const saveButton = createButton("save");
  const cancelButton = createButton("cancel");

  description.innerText = "TEST";
  title.innerText = "TEST";

  main.appendChild(title);
  main.appendChild(description);
  inputDiv.appendChild(nameInput);
  inputDiv.appendChild(descriptionInput);
  inputDiv.appendChild(refreshIntervalInput);

  buttonsDiv.appendChild(deleteButton);
  buttonsDiv.appendChild(cancelButton);
  buttonsDiv.appendChild(saveButton);

  mainDiv.appendChild(inputDiv);
  mainDiv.appendChild(buttonsDiv);
  main.appendChild(mainDiv);
  document.getElementsByTagName("body")[0].appendChild(main);

  return main;
}

function showPopup(id) {
  const main = document.getElementById(id);
  main.style.zIndex = "1000";
}

function hidePopup(id) {
  const main = document.getElementById(id);
  main.style.zIndex = "0";
}
