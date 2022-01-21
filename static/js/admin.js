// window.onload = function () {
//   const pop = newPopup("myPop");
//   showPopup("myPop");
// };

window.onload = async function () {
  console.log("IT WORKS");
  const documents = await getDocuments();
  const pop = newPopup(documents[0]);
  showPopup(documents[0].Id);
};
