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

function setAccess(allow){
    if(allow){
        document.querySelector("#submit").enabled  = allow;
    document.querySelector("#chatbox").querySelector("textarea").enabled  = allow;
    }else{
        document.querySelector("#submit").disabled  = allow;
        document.querySelector("#chatbox").querySelector("textarea").disabled  = allow;
    }
    
    document.querySelector("#submit").style.backgroundColor=allow?"green":"gray";
    document.querySelector("#submit").style.cursor=allow?"pointer":"auto";
}

if (!window["WebSocket"]) {

    // alert("Error: Your browser does not support websockets.")
    addMesage("Error: Your browser does not support websockets.");
    setAccess(false); 

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
    socket.onclose = (e) => {
        //alert("Connection has been closed.");
        console.log("connection has been closed",e);
        setAccess(false); 
        
    }
    socket.onmessage = (e) => {
        console.log(`onmessage data is ${e.data}`);
        
        var msg = eval("("+e.data+")"); //eveal function convert JSON message
        //var messsage = `<strong> ${msg.Name} </strong>: <span>${msg.Message}</span>`;
        //addMesage(messsage);
        
        addMesage(msg.Name,msg.Message,msg.Picture)
    }
    socket.onopen = () =>{
        setAccess(true); 
        console.log("socket connection opened");
        
    }
    socket.onerror =  (e)=>{

        console.log("connection error",e);

    }
}