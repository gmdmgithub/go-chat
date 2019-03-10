function addMesage(msg){
    var node = document.createElement("LI"); // Create a <li> node
    var textnode = document.createTextNode(msg); // Create a text node
    node.appendChild(textnode); // Append the text to <li>
    document.getElementById("messages").appendChild(node); // Append <li> to <ul> with id="myList"
}

if (!window["WebSocket"]) {

    // alert("Error: Your browser does not support websockets.")
    addMesage("Error: Your browser does not support websockets.");
    document.querySelector("#submit").disabled  = true;
    document.querySelector("#chatbox").querySelector("textarea").disabled  = true;
    document.querySelector("#submit").style.backgroundColor="gray";
    document.querySelector("#submit").style.cursor="auto"; 

} else {
    var socket = null;
    var messages = document.querySelectorAll("#messages");

    document.querySelector("#chatbox").addEventListener('submit', (e) => {
        e.preventDefault();
        var msgBox = document.querySelector("#chatbox").querySelector("textarea");

        console.log(`Message is sent ... value: ${msgBox.value}`);

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


    socket = new WebSocket(`ws://${socketAddress}/room`);
    socket.onclose = () => {
        alert("Connection has been closed.");
    }
    socket.onmessage = (e) => {
        addMesage(e.data);
    }
}