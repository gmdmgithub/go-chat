var socket = null;
var messages = document.querySelectorAll("#messages");


document.querySelector("#chatbox").addEventListener('submit', (e) => {
    e.preventDefault();
    var msgBox = document.querySelector("#chatbox").querySelector("textarea");

    console.log(`What is goin on? ${msgBox.value}`);
    
    if (!msgBox.value) { 
        alert("Error: Message is empty");
        return false;
    }
    if (!socket) {
        alert("Error: There is no socket connection.");
        return false;
    }
    socket.send(msgBox.value);
    msgBox.value = "";
    return false;
});

if (!window["WebSocket"]) {
    alert("Error: Your browser does not support websockets.")
} else {
    socket = new WebSocket("ws://localhost:8081/room");
    socket.onclose = () => {
        alert("Connection has been closed.");
    }
    socket.onmessage = (e) => {
        var node = document.createElement("LI"); // Create a <li> node
        var textnode = document.createTextNode(e.data); // Create a text node
        node.appendChild(textnode); // Append the text to <li>
        document.getElementById("messages").appendChild(node); // Append <li> to <ul> with id="myList"
    }
}