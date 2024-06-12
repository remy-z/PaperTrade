window.onload = function () {
    connectWebSocket()
}

function connectWebSocket() {
    if (window["WebSocket"]) {

        socket = new WebSocket("ws://" + document.location.host + "/ws")

        socket.onopen = function (event) {
            console.log("socket-opened")
        }

        socket.onclose = function (event) {
            console.log("socket-closed")
            console.log("You have been disconnected. Reload the page to reconnect")
        }

        socket.onmessage = function (event) {
            const eventData = JSON.parse(event.data);
            const eventObject = Object.assign(new WebSocketEvent, eventData);
            routeEvent(eventObject);
        }
    } else {
        alert("browser does not support websockets")
    }
}

function routeEvent(event) {
    if (event.type === undefined) {
        alert('no type field in the event');
    }

    switch (event.type) {
        case "recieve_err":
            const errorEvent = Object.assign(new RecieveErrorEvent, event.payload);
            console.log(errorEvent)
            break;
        case "recieve_price":
            const priceEvent = Object.assign(new RecievePriceEvent, event.payload)
            console.log(priceEvent)
            break;
        default:
            alert("unsupported message type");
            break;
    }
}


function sendEvent(eventName, payload) {
    const event = new WebSocketEvent(eventName, payload)
    console.log(event)
    socket.send(JSON.stringify(event))
}

function sendSub(symbol) {
    if (symbol != null) {
        let outgoingEvent = new SendSubEvent(symbol);
        sendEvent("send_sub", outgoingEvent)
    }
}

function sendUnsub(symbol) {
    if (symbol != null) {
        let outgoingEvent = new SendUnsubEvent(symbol);
        sendEvent("send_unsub", outgoingEvent)
    }
}


// TYPES

class WebSocketEvent {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

class RecieveErrorEvent {
    constructor(message) {
        this.message = message;
    }
}

class SendSubEvent {
    constructor(symbol) {
        this.symbol = symbol;
    }
}

class SendUnsubEvent {
    constructor(symbol) {
        this.symbol = symbol;
    }
}

class RecievePriceEvent {
    constructor(symbol, price) {
        this.symbol = symbol;
        this.price = price
    }
}