function newPopup(id) {
  //   <div id="edit-overlay">
  //     <h2 id="edit-overlay-title">loading..</h2>
  //     <span id="edit-overlay-description">Description</span>
  //   </div>;

  const main = document.createElement("div");
  main.className = "edit-overlay";
  main.id = id;
  const title = document.createElement("h2");

  const description = document.createElement("span");
  description.className = "edit-overlay-description";

  description.innerText = "TEST";
  title.innerText = "TEST";

  main.appendChild(title);
  main.appendChild(description);

  document.getElementsByTagName("body")[0].appendChild(main);

  //   return main;
}

function showPopup(id) {
  const main = document.getElementById(id);
  main.style.zIndex = "1000";
}

function hidePopup(id) {
  const main = document.getElementById(id);
  main.style.zIndex = "0";
}
