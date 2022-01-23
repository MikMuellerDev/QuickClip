window.onload = async function () {
  const userStatus = await getUserStatus();
  const documents = await getDocuments();
  const version = await getVersion();
  setButtonText(userStatus.LoggedIn);
  addDocuments(documents);
  setVersion(version.Version, version.Production);
};
