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
