window.onload = async function () {
  const userStatus = await getUserStatus();
  const documents = await getDocuments();
  const version = await getVersion();
  setVersion(version.Version, version.Production);
  setButtonText(userStatus.LoggedIn);
  if (documents !== null) {
    addDocuments(documents, true);
  }
};

function createNewDocument() {
  const id = `clip-${Math.random()
  .toString(36)
  .replace(/[^a-z]+/g, "")
  .substr(0, 10)}`
  newPopup(
    {
      Id: id,
      Name: "new Document",
      Description: "describe this document",
      RefreshInterval: 500,
      Refresh: false,
      ReadOnly: false,
      Restricted: true,
      Content: ""
    },
    true
  );
  showPopup(id);
}
