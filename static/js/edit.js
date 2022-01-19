let spinnerLock = false;

function setTitle(doc) {
  const titleElement = document.getElementById("title");
  titleElement.innerHTML = `Editing ${doc.Name}`;
}

function registerTypeEvents(doc) {
  const textArea = document.getElementById("clip");
  textArea.oninput = async () => {
    if (!spinnerLock) {
      spinnerLock = true;
      startSpinner("edit-spinner");
      await textInput(doc.Id, textArea.value);
      stopSpinner("edit-spinner");
      setTimeout(() => {
        spinnerLock = false;
      }, 1000);
    } else {
      await textInput(doc.Id, textArea.value);
    }
  };
}

function setText(doc) {
  const textArea = document.getElementById("clip");
  textArea.value = doc.Content;
}

async function mainloop(id) {
  const doc = await getDocumentById(id, true);
  setText(doc);

  setTimeout(async () => {
    await mainloop(id);
  }, 100);
}

async function quitSave() {
  if (!spinnerLock) {
    spinnerLock = true;
    startSpinner("edit-spinner");
  }
  await requestSave();
  window.location.href = "/dash";
}

window.onload = async () => {
  const textArea = document.getElementById("clip");
  textArea.value = "Loading Clipboard...";

  try {
    const url = window.location.href.split("/");
    const id = url[url.length - 1];

    const doc = await getDocumentById(id);
    setText(doc);

    registerTypeEvents(doc);
    setTitle(doc);

    const version = await getVersion();
    setVersion(version.Version, version.Production);

    mainloop(id).then();
  } catch (err) {
    textArea.value = `An error occurred, unable to fetch content: ${err}`;
  }
};
