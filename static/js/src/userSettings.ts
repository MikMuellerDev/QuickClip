interface User {
  Name: string;
  Password: string;
  Permissions: string[];
  WriteAllowed: string[];
}

async function getUsers(): Promise<User[]> {
  const url = "/api/users";
  const res = await fetch(url);
  const data = await res.json();
  //console.log(data);

  return data;
}

async function createUser(obj: User) {
  const url = `/api/user`;

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

async function modifyUser(username: string, obj: User) {
  const url = `/api/user/${username}`;

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

async function deleteUser(username: string) {
  const url = `/api/user/${username}`;

  const res = await fetch(url, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const data = await res.json();
  //console.log(data);
}


async function createNewUser() {
  await newUserPopup(
    {
      Name: "Username",
      Password: "",
      Permissions: [],
      WriteAllowed: [],
    },
    true
  );
  showPopup("Username");
}

function addUsers(users: User[]) {
  for (let user of users) {
    const parentDiv = document.getElementById(
      "doc-selector-div"
    ) as HTMLDivElement;

    const nodeItem = document.createElement("div");
    nodeItem.className = "mode-item station threeDp";

    const nodeItemLeft = document.createElement("div");
    nodeItemLeft.className = "node-item-left";

    const nodeItemLeftPicture = document.createElement("div");
    nodeItemLeftPicture.className = "node-item-picture";
    nodeItemLeftPicture.style.backgroundImage = `url(/static/media/user.png)`;

    const nodeItemLabels = document.createElement("div");
    nodeItemLabels.className = "node-item-labels";

    const nodeItemTitle = document.createElement("h2");
    nodeItemTitle.innerText = user.Name;

    parentDiv.appendChild(nodeItem);
    nodeItem.appendChild(nodeItemLeft);
    nodeItemLeft.appendChild(nodeItemLeftPicture);
    nodeItemLeft.appendChild(nodeItemLabels);
    nodeItemLabels.appendChild(nodeItemTitle);

    setTimeout(function () {
      nodeItem.style.opacity = "1";
    }, 50);

    nodeItem.onclick = async () => {
      await newUserPopup(user, false)
      showPopup(user.Name)
    }
  }
}

function removeUsers() {
  for (let i = 0; i < 2; i++) {
    const parentDiv = document.getElementById(
      "doc-selector-div"
    ) as HTMLDivElement;
    const children = parentDiv.children;
    for (let child of children) {
      child.remove();
    }
  }
}

window.onload = async () => {
  const userStatus = await getUserStatus();
  const version = await getVersion();
  setVersion(version.Version, version.Production);
  setButtonText(userStatus.LoggedIn);

  const users = await getUsers();
  console.log(users);
  addUsers(users);
};
