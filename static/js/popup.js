function createInput(label, placeholder) {
  const main = document.createElement("input");
  main.className = "text-input";
  main.type = "text";
  main.value = label;
  main.placeholder = placeholder;
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

function createSpan(label, className) {
  const span = document.createElement("span");
  span.className = className;
  span.innerText = label;
  return span;
}

function deleteHandler(id) {
  console.log(`DELETE: ${id}`);
  deleteClip(id);
}

function modifyHandler(obj) {
  console.log(`MODIFY: ${obj.Name}`);
  modifyClip(obj);
}

function newPopup(doc) {
  const main = document.createElement("div");
  main.className = "edit-overlay";
  main.id = doc.id;

  const title = document.createElement("h2");

  // const description = document.createElement("span");
  // description.className = "edit-overlay-description";

  const description = createSpan("edit-overlay-description");

  const mainDiv = document.createElement("div");
  mainDiv.className = "edit-overlay-main-div";

  const inputDiv = document.createElement("div");
  inputDiv.className = "edit-overlay-inputs";

  const labelDiv = document.createElement("div");
  labelDiv.className = "edit-overlay-labels";

  const buttonsDiv = document.createElement("div");
  buttonsDiv.className = "edit-overlay-buttons";

  const refreshIntervalInput = createInput(
    doc.RefreshInterval,
    "Refresh Interval"
  );

  const descriptionInput = createInput(doc.Description, "Description");
  const nameInput = createInput(doc.Name, "Name");
  const deleteButton = createButton("delete");
  const saveButton = createButton("save");
  const cancelButton = createButton("cancel");

  const nameSpan = createSpan("Name", "");
  const descriptionSpan = createSpan("Description", "");
  const refreshSpan = createSpan("Refresh Interval", "");

  deleteButton.onclick = () => {
    deleteHandler(doc.Id);
  };

  saveButton.onclick = () => {
    modifyHandler({
      Id: doc.Id,
      Name: nameInput.value,
      Description: descriptionInput.value,
      RefreshInterval: parseInt(refreshIntervalInput.value),
      Refresh: true,
      Restricted: false,
    });
  };

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

  mainDiv.appendChild(labelDiv);
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
