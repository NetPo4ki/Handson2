// Generate a random username
const username = `User${Math.floor(Math.random() * 1000)}`;

// WebSocket setup
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = function (event) {
    const message = JSON.parse(event.data);
    appendMessage(message);
};

ws.onopen = function () {
    console.log('Connected to the WebSocket server');
};

async function sendMessage() {
    const messageInput = document.getElementById('messageInput');
    const message = messageInput.value;

    if (message.trim() === '') {
        alert('Message cannot be empty!');
        return;
    }

    try {
        const response = await fetch('http://localhost:8080/messages', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, message }),
        });

        if (response.ok) {
            messageInput.value = '';
        } else {
            alert('Failed to send message.');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error sending message.');
    }
}

async function getMessages() {
    try {
        const response = await fetch('http://localhost:8080/messages');
        const messages = await response.json();
        const chatWindow = document.getElementById('chatWindow');
        chatWindow.innerHTML = ''; // Clear current messages

        messages.forEach(message => appendMessage(message));
    } catch (error) {
        console.error('Error:', error);
        alert('Error retrieving messages.');
    }
}

function appendMessage(message) {
    const chatWindow = document.getElementById('chatWindow');
    const messageDiv = document.createElement('div');
    messageDiv.className = 'message';
    messageDiv.innerText = `[${message.username}] ${message.content}`;
    chatWindow.appendChild(messageDiv);
    chatWindow.scrollTop = chatWindow.scrollHeight; // Scroll to the bottom
}

// Initial messages retrieval
getMessages();