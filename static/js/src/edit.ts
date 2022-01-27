let spinnerLock = false;

function setTitle(doc: Clip) {
  const titleElement = document.getElementById("title") as HTMLHeadingElement ;
  titleElement.innerHTML = `Editing ${doc.Name}`;
}

function registerTypeEvents(doc: Clip) {
  const textArea = document.getElementById("clip") as HTMLInputElement;
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

function setText(doc: Clip) {
  const textArea = document.getElementById("clip") as HTMLTextAreaElement;
  textArea.value = doc.Content;
}

function setWordCount(text: string) {
  const len = text.length;
  let words = text.split(" ").length;
  if (len == 0) {
    words = 0;
  }

  (document.getElementById("wordcount") as HTMLSpanElement).innerText = `${words}`;
  (document.getElementById("charcount") as HTMLSpanElement).innerText = `${len}`;
}

async function refresh(id: string, timeout: number) {
  const doc = await getDocumentById(id, true);
  setText(doc);
  setWordCount(doc.Content)
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
  const textArea = document.getElementById("clip") as HTMLTextAreaElement ;
  textArea.select();
  textArea.setSelectionRange(0, 99999);
  document.execCommand("copy");
  textArea.blur();
}

async function setWriteMode(doc: Clip) {
  if (doc.ReadOnly) {
    console.log("This document is read-only.");
    const success = await probeWriteAccess(doc.Id);
    if (success) {
      console.log("[PROBE] Write allowed.");
      (document.getElementById("readOnlyIndicator") as HTMLSpanElement).innerText = "";
    } else {
      console.log("[PROBE] Write forbidden.");
      const clip = document.getElementById("clip") as HTMLTextAreaElement;
      clip.disabled = true;
      (document.getElementById("readOnlyIndicator") as HTMLSpanElement ).innerText = "read-only";
      if (doc.Refresh) {
        (document.getElementById("refreshIndicator") as HTMLSpanElement ).innerText = "(live)";
        (document.getElementById(
          "refreshIndicator"
        ) as HTMLSpanElement ).innerHTML = `<span>synchronized ${doc.RefreshInterval}ms  </span> <img id='syncSymbol' src='/static/media/sync.png'>`;
        console.log(
          `%cThis document supports live-editing @ ${
            doc.RefreshInterval / 1000
          }`,
          "color:cyan"
        );

        refresh(doc.Id, doc.RefreshInterval).then();
      } else {
        (document.getElementById("refreshIndicator") as HTMLSpanElement).innerText = "(static)";
        (document.getElementById("refreshIndicator") as HTMLSpanElement).innerHTML =
          "<span>not synchronized</span> <img id='syncSymbol' src='/static/media/syncoff.png'>";
        console.log(`%cThis document is static.`, "color:red");
      }
      return;
    }
  }
  registerTypeEvents(doc);
}

window.onload = async () => {
  const textArea = document.getElementById("clip") as HTMLTextAreaElement;
  textArea.value = "Loading Clipboard...";

  try {
    const url = window.location.href.split("/");
    const id = url[url.length - 1];

    const doc = await getDocumentById(id, false);
    setText(doc);

    setWriteMode(doc);

    setTitle(doc);

    const version = await getVersion();
    setVersion(version.Version, version.Production);

    setWordCount(doc.Content);
  } catch (err) {
    textArea.value = `An error occurred, unable to fetch content: ${err}`;
  }
};
