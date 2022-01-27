function createInput(label: string, placeholder: string) {
  const main = document.createElement("input");
  main.className = "text-input";
  main.type = "text";
  main.value = label;
  main.placeholder = placeholder;
  return main;
}

function createButton(label: string) {
  const button = document.createElement("button");
  const buttonText = document.createElement("a");
  buttonText.innerText = label;
  button.className = "button";
  button.appendChild(buttonText);
  return button;
}

function createSpan(label: string, className: string) {
  const span = document.createElement("span");
  span.className = className;
  span.innerText = label;
  return span;
}

async function deleteHandler(id: string) {
  console.log(`DELETE: ${id}`);
  await deleteClip(id);
}

async function modifyHandler(obj: Clip) {
  console.log(`MODIFY: ${obj.Name}`);
  await modifyClip(obj);
}

function newPopup(doc: Clip, create: boolean) {
  let descriptionTextPrev: string;

  const main = document.createElement("div");
  main.className = "edit-overlay";
  main.id = doc.Id;

  const title = document.createElement("h2");
  title.innerText = doc.Name;
  title.className = "popup-doc-title";

  const idTitle = document.createElement("h3");
  idTitle.innerText = `ID: ${doc.Id}`;
  idTitle.className = "popup-doc-id";

  const description = createSpan("edit-overlay-description", "");
  description.innerText = doc.Description;
  description.className = "popup-doc-description";

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

  const refreshIntervalInput = createInput(
    `${doc.RefreshInterval}`,
    "Refresh Interval"
  );

  const idInput = createInput(doc.Id, "unique ID");

  const descriptionInput = createInput(doc.Description, "Description");
  const nameInput = createInput(doc.Name, "Name");
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
  const descriptionSpan = createSpan("Description:", "");
  const refreshSpan = createSpan("Refresh Interval:", "");
  const idSpan = createSpan("unique ID:", "");

  labelDiv.appendChild(nameSpan);
  labelDiv.appendChild(descriptionSpan);
  labelDiv.appendChild(refreshSpan);

  nameInput.oninput = () => {
    title.innerText = nameInput.value;
  };

  descriptionInput.oninput = () => {
    console.log(descriptionInput.value.length);
    if (descriptionInput.value.length === 0) {
      description.innerText = doc.Name;
    } else if (descriptionInput.value.length < 34) {
      description.innerText = descriptionInput.value;
      descriptionTextPrev = descriptionInput.value;
    } else {
      descriptionInput.value = descriptionTextPrev;
    }
  };

  deleteButton.onclick = async () => {
    const confirm = prompt("Enter `yes` to confirm deletion", "");
    if (confirm == "yes") {
      await deleteHandler(doc.Id);
      removeDocuments();
      const documents = await getDocuments();
      if (documents !== null) addDocuments(documents, true);
    }
    setTimeout(() => {
      main.remove();
    }, 600);
    hidePopup(doc.Id);
  };

  saveButton.onclick = async () => {
    if (!create) {
      await modifyHandler({
        Id: doc.Id,
        Name: nameInput.value,
        Description: descriptionInput.value,
        RefreshInterval: parseInt(refreshIntervalInput.value),
        Refresh: doc.Refresh,
        ReadOnly: doc.ReadOnly,
        Restricted: doc.Restricted,
        Content: "",
      });
    } else {
      const tempDocs = await getDocuments();

      if (tempDocs != null) {
        for (let iterDoc of tempDocs) {
          if (iterDoc.Id == idInput.value) {
            saveButton.style.backgroundColor = "var(--clr-error)";
            alert("This ID already exists.");
            return;
          }
        }
      }

      await createClip({
        Id: idInput.value.replaceAll(" ", "-"),
        Name: nameInput.value,
        Description: descriptionInput.value,
        RefreshInterval: parseInt(refreshIntervalInput.value),
        Refresh: doc.Refresh,
        ReadOnly: doc.ReadOnly,
        Restricted: doc.Restricted,
        Content: "",
      });
    }

    hidePopup(doc.Id);
    setTimeout(() => {
      main.remove();
    }, 600);

    removeDocuments();
    addDocuments(await getDocuments(), true);
  };

  cancelButton.onclick = () => {
    hidePopup(doc.Id);
    setTimeout(() => {
      main.remove();
    }, 600);
  };

  main.appendChild(title);
  main.appendChild(description);
  main.appendChild(idTitle);
  inputDiv.appendChild(nameInput);
  inputDiv.appendChild(descriptionInput);
  inputDiv.appendChild(refreshIntervalInput);

  if (!create) {
    buttonsDiv.appendChild(deleteButton);
  } else {
    inputDiv.appendChild(idInput);
    labelDiv.appendChild(idSpan);
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

  // const docTraits = ["Restricted", "ReadOnly", "Refresh"];
  const docTraits: ("Restricted" | "ReadOnly" | "Refresh")[] = [
    "Restricted",
    "ReadOnly",
    "Refresh",
  ];
  for (let trait of docTraits) {
    const traitSwitchArr = createSlider();
    const traitSwitch = traitSwitchArr[0];
    const traitSwitchListener = traitSwitchArr[1];

    traitSwitchListener.checked = doc[trait];

    traitSwitchListener.onchange = () => {
      doc[trait] = traitSwitchListener.checked;
      console.log(`Trait ${trait} changed to ${doc[trait]}`);
    };

    const switchContainer = document.createElement("div");
    switchContainer.className = "inner-slider-div";

    const switchContainerSpan = document.createElement("span");
    switchContainerSpan.innerText = trait;

    switchContainer.appendChild(traitSwitch);
    switchContainer.appendChild(switchContainerSpan);
    sliderDiv.appendChild(switchContainer);
  }

  document.getElementsByTagName("body")[0].appendChild(main);

  return main;
}

function showPopup(id: string) {
  const allDocuments = document.getElementById(
    "doc-selector-div"
  ) as HTMLDivElement;
  const main = document.getElementById(id) as HTMLElement;
  main.style.opacity = "0";
  main.style.display = "flex";
  setTimeout(() => {
    main.style.opacity = "100%";
    allDocuments.style.filter = "brightness(50%)";
    main.style.zIndex = "1000";
  }, 100);
}

function hidePopup(id: string) {
  const allDocuments = document.getElementById(
    "doc-selector-div"
  ) as HTMLDivElement;
  allDocuments.style.zIndex = "-2000";
  const main = document.getElementById(id) as HTMLElement;

  main.style.opacity = "0%";
  main.style.transform = "translate(-50%, -50rem)";

  setTimeout(() => {
    main.style.display = "none";
    allDocuments.style.filter = "brightness(100%)";
    main.style.zIndex = "-1000";
    allDocuments.style.zIndex = "0";
  }, 400);
}

function createSlider(): [HTMLLabelElement, HTMLInputElement] {
  let switchL;
  let switchE;

  switchE = document.createElement("input");

  switchE.type = "checkbox";
  switchE.checked = false;

  const switchS = document.createElement("span");
  switchS.className = "slider round sixDp";

  switchL = document.createElement("label");
  switchL.className = "switch";
  switchL.appendChild(switchE);
  switchL.appendChild(switchS);

  return [switchL, switchE];
}
