'use strict'

var usernamePage = document.querySelector('#username-page');
var chatPage = document.querySelector('#chat-page');
var usernameForm = document.querySelector('#usernameForm');
var messageForm = document.querySelector('#messageForm');
var messageInput = document.querySelector('#message');
var messageArea = document.querySelector('#messageArea');
var connectingElement = document.querySelector('.connecting');

var stompClient = null;
var username =null;

var colors = [
    '#2196F3', '#f32196', '#7FFFD4', '#808080',
    '#8B0000', '#E9967A', '#483D8B', '#E6E6FA'
];

function connect(event) {
    console.log("Connecting...");
    username = document.querySelector('#name').value.trim();
    if(username){
        usernamePage.classList.add('hidden');
        chatPage.classList.remove('hidden');

        var socket = new SockJS('/javatechie');
        stompClient = Stomp.over(socket);

        stompClient.connect({},onConnected, onError);
    }
    event.preventDefault();
}

function onConnected(){
    stompClient.subscribe('/topic/public', onMessageReceived);

    stompClient.send("/app/chat.register",
        {},
        JSON.stringify({sender:username, type: 'JOIN'})
    )
    connectingElement.classList.add('hidden');
}

function onError(){
    connectingElement.textContent='Could not connect to WebSocket server. Please refresh this page to try again!';
    connectingElement.style.color='red';
}


function send(event){
    var messageContent = messageInput.value.trim();

    if(messageContent && stompClient){
        var chatMessage = {
            sender:username,
            content: messageInput.value,
            type:'CHAT'
        };
        stompClient.send("/app/chat.send", {}, JSON.stringify(chatMessage));
        messageInput.value ='';
    }
    event.preventDefault();
}


function onMessageReceived(payload){
    var message = JSON.parse(payload.body);
    var messageElement = document.createElement('li');

    if(message.type === 'JOIN'){
        messageElement.classList.add('event-message');
        message.content = message.sender + ' joined!';
    }else if (message.type === 'LEAVE'){
        messageElement.classList.add('event-message');
        message.content = message.sender + 'left!';
    }else{
        messageElement.classList.add('chat-message');
        var avatarElement = document.createElement('i');
        var avatarText = document.createTextNode(message.sender[0]);
        avatarElement.appendChild(avatarText);
        avatarElement.style['background-color'] = getAvatarColor(message.sender);

        messageElement.appendChild(avatarElement);

        var usernameElement = document.createElement('span');
        var usernameText = document.createTextNode(message.sender);
        usernameElement.appendChild(usernameText)
        messageElement.appendChild(usernameElement);
    }

    var textElement = document.createElement('p');
    var messageText = document.createTextNode(message.content);
    textElement.appendChild(messageText);

    messageElement.appendChild(textElement);

    messageArea.appendChild(messageElement);
    messageArea.scrollTop=messageArea.scrollHeight;
}


function getAvatarColor(messageSender){
    var hash =0;
    for (var i =0;i<messageSender.length;i++){
        hash=31*hash+messageSender.charCodeAt(i);
    }

    var index = Math.abs(hash % colors.length);
    return colors[index];
}

usernameForm.addEventListener('submit', connect, true);
messageForm.addEventListener('submit', send, true);