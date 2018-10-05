let ws;

function Connect() {
    ws = new WebSocket("ws://" + window.location.host + "/ws");
    console.log("Websocket chat - status: " + ws.readyState);

    ws.onopen = function() {
        console.log("Connection chat opened..." + this.readyState);
    };

    ws.onmessage = function(msg) {
        UpdateStatus(msg.data);
    };

    ws.onerror = function(msg) {
        console.log("Error chat occured sending..." + msg.data);
    };

    ws.onclose = function(msg) {
        console.log("Disconnected chat - status " + this.readyState);
    };
}

function UpdateStatus(jsonMessage) {
    console.log(jsonMessage)
}