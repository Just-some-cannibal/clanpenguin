window.onload = function() {
    let socket = new WebSocket("ws://clanpenguin.play.ai/ws");
    let chat = document.getElementById("chat");
    let submit = document.getElementById("submit");
    let text = document.getElementById("text");
  
    function send () {
        if(text.value !== "") {
          socket.send(JSON.stringify({data: text.value}));
        }
        text.value = "";
    }
  
		socket.onmessage = function(e) {
			chat.innerHTML += 
			    `<span class="message"> ${e.data} <span>`;
      chat.scrollTop = chat.scrollHeight;
		}
    
		submit.onclick = function() {
      send();
		}
    
    text.addEventListener("keyup", function(e) {
        if(e.keyCode == 13) {
            send();
        }
    });
}