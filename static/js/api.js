async function setCurrentMode(id) {
  const res = await fetch(`/api/mode/${id}`, { method: "post" });
  return await res.json();
}

async function getDocuments() {
  const url = "/api/clips";
  const res = await fetch(url);
  return (await res.json())["Clips"];
}

async function getDocumentById(id, useKeepalive) {
  let url = `/api/clip/${id}`;
  if (useKeepalive) {
    url = `/api/clip/refresh/${id}`;
  }
  const res = await fetch(url);
  return await res.json();
}

async function getUserStatus() {
  let url = `/api/user`;
  const res = await fetch(url);
  return await res.json();
}

async function textInput(id, content) {
  const url = `/api/clips/edit/${id}`;
  const cont = { Content: content };
  const res = await fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(cont),
  });

  const data = await res.json();
  console.log(data);
}

async function modifyClip(obj) {
  const url = `/api/clips/update`;

  const res = await fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(obj),
  });

  const data = await res.json();
  console.log(data);
}

async function deleteClip(id) {
  const url = `/api/clips/update/${id}`;
  const res = await fetch(url, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const data = await res.json();
  console.log(data);
}

async function requestSave() {
  const url = "/api/save";
  const res = await fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const data = await res.json();
  console.log(data);
}
