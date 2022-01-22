let spinnerLock = false;

function setTitle(doc) {
  const titleElement = document.getElementById("title");
  titleElement.innerHTML = `Editing ${doc.Name}`;
}

function registerTypeEvents(doc) {
  const textArea = document.getElementById("clip");
  textArea.oninput = async () => {
    setWordCount(textArea.value);
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

function setWordCount(text) {
  const len = text.length;
  let words = text.split(" ").length;
  if (len == 0) {
    words = 0;
  }

  document.getElementById("wordcount").innerText = `${words}`;
  document.getElementById("charcount").innerText = `${len}`;
}

async function refresh(id, timeout) {
  const doc = await getDocumentById(id, true);
  setText(doc);

  setTimeout(async () => {
    await refresh(id, timeout);
  }, timeout);
}

async function quitSave() {
  if (!spinnerLock) {
    spinnerLock = true;
    startSpinner("edit-spinner");
  }
  await requestSave();
  window.location.href = "/dash";
}

function copyText() {
  const textArea = document.getElementById('clip')
  textArea.select()
  textArea.setSelectionRange(0, 99999)
  document.execCommand('copy')
  textArea.blur()
}

async function setWriteMode(doc) {
  if (doc.ReadOnly) {
    console.log("This document is read-only.");
    const success = await probeWriteAccess(doc.Id);
    if (success) {
      console.log("[PROBE] Write allowed.");
      document.getElementById("readOnlyIndicator").innerText = ""
    } else {
      console.log("[PROBE] Write forbidden.");
      const clip = document.getElementById("clip");
      clip.disabled = true;
      document.getElementById("readOnlyIndicator").innerText = "read-only"
      return
    }
  }
  registerTypeEvents(doc);
}

window.onload = async () => {
  const textArea = document.getElementById("clip");
  textArea.value = "Loading Clipboard...";

  try {
    const url = window.location.href.split("/");
    const id = url[url.length - 1];

    const doc = await getDocumentById(id);
    setText(doc);

    setWriteMode(doc)

    setTitle(doc);

    const version = await getVersion();
    setVersion(version.Version, version.Production);

    setWordCount(doc.Content);

    if (doc.Refresh) {
      document.getElementById("refreshIndicator").innerText = "(live)";
      document.getElementById(
        "refreshIndicator"
      ).innerHTML = `<span>synchronized ${doc.RefreshInterval}ms  </span> <img id='syncSymbol' src='/static/media/sync.png'>`;
      console.log(
        `%cThis document supports live-editing @ ${doc.RefreshInterval / 1000}`,
        "color:cyan"
      );
      refresh(id, doc.RefreshInterval).then();
    } else {
      document.getElementById("refreshIndicator").innerText = "(static)";
      document.getElementById("refreshIndicator").innerHTML =
        "<span>not synchronized</span> <img id='syncSymbol' src='/static/media/syncoff.png'>";
      console.log(`%cThis document is static.`, "color:red");
    }
  } catch (err) {
    textArea.value = `An error occurred, unable to fetch content: ${err}`;
  }
};
