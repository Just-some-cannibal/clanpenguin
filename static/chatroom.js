window.onload = function () {
    let socket = new WebSocket("ws://" + location.host + "/ws")
    let chat = document.getElementById("chat");
    let submit = document.getElementById("submit");
    let text = document.getElementById("text");
    let overlay = document.getElementById("overlay");
    let overlayEnabled = false;
    let overlayOK = document.getElementById("overlay-ok");
    let name = document.getElementById("name");

    function showOverlay(bool) {
        overlayEnabled = bool;
        if (bool) {
            overlay.style.display = "flex";
        } else {
            overlay.style.display = "none";
        }
    }

    function addMessage(user, text) {
        let message = document.createElement("span");
        let textNode = document.createTextNode(text);
        let userText = document.createElement("span")
        let userNode = document.createTextNode(user + ": ");
        userText.append(userNode);
        message.className = "message";
        message.append(userNode);
        message.append(textNode);
        chat.append(message);
    }

    function send() {
        if (text.value !== "") {
            socket.send(JSON.stringify({ data: {user: name.value, text: text.value}, protocol: "broadcast" }));
        }
        text.value = "";
    }

    function handleError(error) {
        if (error == "spam") {
            text.disabled = true;
            setTimeout(function() {
                text.disabled = false;
                text.focus();
                showOverlay(false);
            }, 5000);
            showOverlay(true);
        }
        console.error(error);
    }

    overlayOK.onclick = function() {
        showOverlay(false);
    }

    socket.onmessage = function (e) {
        let response = JSON.parse(e.data);
        if (response.protocol == "err") {
            handleError(response.data);
        } else if (response.protocol == "broadcast") {
            addMessage(response.data.user, response.data.text);
            console.log(response.data.user);
            chat.scrollTop = chat.scrollHeight;
        } else if (response.protocol == "get") {
            for (let message of response.data) {
                console.log(message)
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
        if (e.keyCode == 13 && document) {
            send();
        }
    });
    document.addEventListener("keydown", function(e) {
        if ((e.keyCode === 27 || e.keyCode === 27) && overlayEnabled) {
            showOverlay(false);
        }
    });
}