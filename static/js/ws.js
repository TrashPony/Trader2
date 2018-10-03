var ws;

function Connect() {
    ws = new WebSocket("ws://" + window.location.host + "/ws");
    console.log("Websocket chat - status: " + ws.readyState);

    ws.onopen = function() {
        console.log("Connection chat opened..." + this.readyState);
    };

    ws.onmessage = function(msg) {
        console.log(msg);
        NewChatMessage(msg.data);
    };

    ws.onerror = function(msg) {
        console.log("Error chat occured sending..." + msg.data);
    };

    ws.onclose = function(msg) {
        console.log("Disconnected chat - status " + this.readyState);
    };
}

function Message() {
    var chatInput = document.getElementById("chatInput");
    var text = chatInput.value;
    if (text !== "") {
        chatInput.value = null;

        ws.send(JSON.stringify({
            event: "NewMessage",
            message: text
        }));
    }
}

function NewChatMessage(jsonMessage) {
    var event = JSON.parse(jsonMessage).event;

    if (event === "NewMessage") {

    }
}
