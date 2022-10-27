let socket = new WebSocket("ws://localhost:8099/socket");
let p = document.getElementById("progress");

function log(v) {
    p.innerHTML += v + "\n";
    p.scrollTop = p.scrollHeight;
}

socket.onopen = function(e) {
  log("[open] Connection established");
  log("Sending to server");
  socket.send("My name is John");
};

socket.onmessage = function(event) {
  log("[message] Data received from server: " + event.data);
};

socket.onclose = function(event) {
  if (event.wasClean) {
    log("[close] Connection closed cleanly, code=${event.code} reason=${event.reason}");
  } else {
    // e.g. server process killed or network down
    // event.code is usually 1006 in this case
    log("[close] Connection died");
  }
};

socket.onerror = function(error) {
  log("[error]");
};
