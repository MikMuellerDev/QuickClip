interface Clip {
  Name: string;
  Id: string;
  Content: string;
  Description: string;
  Restricted: boolean;
  Refresh: boolean;
  RefreshInterval: number;
  ReadOnly: boolean;
  // [key: string]: any;
}

// type Clip struct {
// 	Name            string
// 	Id              string
// 	Content         string
// 	Description     string
// 	Restricted      bool
// 	Refresh         bool
// 	RefreshInterval int
// 	ReadOnly        bool
// }

function setButtonText(loggedIn: boolean) {
  const button = document.getElementById("logout-button") as HTMLButtonElement;
  const buttonLabel = document.getElementById(
    "logout-button-text"
  ) as HTMLSpanElement;

  if (loggedIn) {
    buttonLabel.innerText = "logout";
    button.onclick = () => {
      window.location.href = "/logout";
    };
  } else {
    buttonLabel.innerText = "login";
    button.onclick = () => {
      window.location.href = "/login";
    };
  }
}
