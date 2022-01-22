function addDocuments(documents) {
  for (let doc of documents) {
    const parentDiv = document.getElementById("doc-selector-div");

    const nodeItem = document.createElement("div");
    nodeItem.className = "mode-item station threeDp";
    nodeItem.id = doc.Id;

    const nodeItemLeft = document.createElement("div");
    nodeItemLeft.className = "node-item-left";

    const nodeItemLeftPicture = document.createElement("div");
    nodeItemLeftPicture.className = "node-item-picture";
    nodeItemLeftPicture.id = `${doc.Id}-picture`;
    nodeItemLeftPicture.style.backgroundImage = `url(/static/media/pen.png)`;

    const nodeItemLabels = document.createElement("div");
    nodeItemLabels.className = "node-item-labels";

    const nodeItemTitle = document.createElement("h2");
    nodeItemTitle.innerText = doc.Name;
    const nodeItemSubTitle = document.createElement("h3");
    nodeItemSubTitle.innerText = doc.Description;

    parentDiv.appendChild(nodeItem);
    nodeItem.appendChild(nodeItemLeft);
    nodeItemLeft.appendChild(nodeItemLeftPicture);
    nodeItemLeft.appendChild(nodeItemLabels);
    nodeItemLabels.appendChild(nodeItemTitle);
    nodeItemLabels.appendChild(nodeItemSubTitle);

    const spinner = document.createElement("div");
    spinner.className = "spinner";
    spinner.id = `${doc.Id}-spinner`;

    nodeItem.appendChild(spinner);

    setTimeout(function () {
      nodeItem.style.opacity = "1";
    }, 50);

    nodeItem.addEventListener("click", async () => {
      window.location.href = `/edit/${doc.Id}`;
    });
  }
}

window.onload = async function () {
  const userStatus = await getUserStatus();
  const documents = await getDocuments();
  const version = await getVersion();
  setButtonText(userStatus.LoggedIn);
  addDocuments(documents);
  setVersion(version.Version, version.Production);
};
