interface UserStatus {
  User: string;
  LoggedIn: boolean;
}

async function getDocuments(): Promise<Clip[]> {
  const url = "/api/clips";
  const res = await fetch(url);
  return (await res.json())["Clips"];
}

async function probeWriteAccess(id: string): Promise<boolean> {
  const url = `/api/clip/probe/${id}`;
  const res = await fetch(url);
  return (await res.json())["Success"];
}

async function getDocumentById(
  id: string,
  useKeepalive: boolean
): Promise<Clip> {
  let url = `/api/clip/${id}`;
  if (useKeepalive) {
    url = `/api/clip/refresh/${id}`;
  }
  const res = await fetch(url);
  return await res.json();
}

async function getUserStatus(): Promise<UserStatus> {
  let url = `/api/user`;
  const res = await fetch(url);
  return await res.json();
}

async function textInput(id: string, content: string) {
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
  console.log(`Write successful: ${data.Success}`);
}

async function modifyClip(obj: Clip) {
  const url = `/api/clips/update`;

  const res = await fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(obj),
  });

  const data = await res.json();
  //console.log(data);
}

async function createClip(obj: Clip) {
  const url = `/api/clips/add`;

  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(obj),
  });

  const data = await res.json();
  //console.log(data);
}

async function deleteClip(id: string) {
  const url = `/api/clips/update/${id}`;
  const res = await fetch(url, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const data = await res.json();
  //console.log(data);
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
  //console.log(data);
}
