const websocket: WebSocket = new WebSocket("/ws");
websocket.onopen = () => {
  console.log("OPEN");
};

websocket.onmessage = (_: MessageEvent) => {
  console.log("ws:// file updated");
  window.location.reload();
};

websocket.onclose = () => {
  console.log("CLOSED");
};

const timer = document.getElementById("timer")!;
const releaseDate = document.getElementById("release-date")!;
const options: Intl.DateTimeFormatOptions = { month: "long" };

const target = new Date(timer.dataset["releasetime"]!);
const months = [
  "Января",
  "Февраля",
  "Марта",
  "Апреля",
  "Мая",
  "Июня",
  "Июля",
  "Августа",
  "Сентября",
  "Октября",
  "Ноября",
  "Декабря",
];
releaseDate.innerHTML = target.getDate() + " " + months[target.getMonth()];

updateTimer();
var interval = setInterval(updateTimer, 1000);

function updateTimer() {
  const now = new Date();

  const distance = target.getTime() - now.getTime();
  if (distance < 0) {
    clearInterval(interval);
    return;
  }

  var days = Math.floor(distance / (1000 * 60 * 60 * 24));
  var hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
  var minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
  var seconds = Math.floor((distance % (1000 * 60)) / 1000);

  [days, hours, minutes, seconds].forEach((v, i) => {
    timer.children[2 * i]!.firstElementChild!.innerHTML = v.toString();
  });
}
