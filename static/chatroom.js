window.onload = function () {
    let socket = new WebSocket("ws://" + location.host + "/ws")
    let chat = document.getElementById("chat");
    let submit = document.getElementById("submit");
    let text = document.getElementById("text");

    function addMessage(user, text) {
        let message = document.createElement("span");
        let textNode = document.createTextNode(text);
        let userNode = document.createTextNode(user);
        message.className = "message";
        message.append(textNode);
        chat.append(message);
    }

    function send() {
        if (text.value !== "") {
            socket.send(JSON.stringify({ data: {user: "", text: text.value}, protocol: "broadcast" }));
        }
        text.value = "";
    }

    socket.onmessage = function (e) {
        let response = JSON.parse(e.data);
        if (response.protocol == "err") {
            console.error(response.data);
        } else if (response.protocol == "broadcast") {
            addMessage(response.data.user, response.data.text);
            chat.scrollTop = chat.scrollHeight;
        } else if (response.protocol == "get") {
            for (let message of response.data) {
                addMessage(message.user, message.text);
            }
            chat.scrollTop = chat.scrollHeight;
        }
    }

    socket.onopen = function () {
        socket.send(JSON.stringify({ data: "", protocol: "get" }))
    }

    submit.onclick = function () {
        send();
    }

    text.addEventListener("keyup", function (e) {
        if (e.keyCode == 13) {
            send();
        }
    });
}