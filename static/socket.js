function addMesage(name, msg,url){
    var node = document.createElement("LI"); // Create a <li> node
    //var textnode = document.createTextNode(msg); // Create a text node
    var imgNode = document.createElement('img')
    var strongMsg = document.createElement('strong') //create stronge node
    var  spanMsg = document.createElement('span'); // Create span node
    var textnode = document.createTextNode(":\u00A0");
    strongMsg.innerHTML = name;
    spanMsg.innerHTML = msg;
    if(url){
        imgNode.src=url;
        imgNode.className="picture";
        node.appendChild(imgNode);
        node.appendChild(textnode);
    }
    node.appendChild(strongMsg)
    node.appendChild(textnode);
    node.appendChild(spanMsg); // Append the text to <li>
    document.getElementById("messages").appendChild(node); // Append <li> to <ul> with id="myList"
}

function disableDialog(){
    document.querySelector("#submit").disabled  = true;
    document.querySelector("#chatbox").querySelector("textarea").disabled  = true;
    document.querySelector("#submit").style.backgroundColor="gray";
    document.querySelector("#submit").style.cursor="auto";
}

if (!window["WebSocket"]) {

    // alert("Error: Your browser does not support websockets.")
    addMesage("Error: Your browser does not support websockets.");
    disableDialog(); 

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
        //socket.send(msgBox.value);
        socket.send(JSON.stringify({"Message": msgBox.value}));
        msgBox.value = "";
        return false;
    });


    socket = new WebSocket(`ws://${socketAddress}/room`);
    socket.onclose = () => {
        alert("Connection has been closed.");
    }
    socket.onmessage = (e) => {
        console.log(`onmessage data is ${e.data}`);
        
        var msg = eval("("+e.data+")");
        //var messsage = `<strong> ${msg.Name} </strong>: <span>${msg.Message}</span>`;
        //addMesage(messsage);
        
        addMesage(msg.Name,msg.Message,msg.Picture)
    }
}