let ws;

function Connect() {
    ws = new WebSocket("ws://" + window.location.host + "/ws");
    console.log("Websocket - status: " + ws.readyState);

    ws.onopen = function () {
        console.log("Connection opened..." + this.readyState);
    };

    ws.onmessage = function (msg) {
        UpdateStatus(msg.data);
    };

    ws.onerror = function (msg) {
        console.log("Error occured sending..." + msg.data);
    };

    ws.onclose = function (msg) {
        alert("Disconnected - status " + this.readyState);
    };
}