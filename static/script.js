document.addEventListener('DOMContentLoaded', function() {
    const chatContainer = document.getElementById('chat-container');
    const chatBox = document.getElementById('chat-box');
    const chatLink = document.getElementById('chat-link');
    const clientNameDisplay = document.getElementById('client-name-display');
    const messageInput = document.getElementById('message');
    const sendMessageButton = document.querySelector('#message-form button');

    const messagesContainer = document.getElementById('chat-box');

    // Function to get query parameters from URL
    function getQueryParam(param) {
        const urlParams = new URLSearchParams(window.location.search.substring(1));
        return urlParams.get(param);
    }

    // if the roomId is present in the URL, join the room, or else create a room
    const roomId = getQueryParam("roomId");
    if (roomId) {
        joinRoom(roomId);
    } else {
        createRoom();
    }

    function createRoom() {
        fetch('/api/rooms', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => response.json())
        .then(data => {
            const roomId = data.room_id;
            // go to url with roomId
            window.location.href = "?roomId=" + roomId;
        })
        .catch(error => console.error('Error creating room:', error));
    }

    function joinRoom(roomId) {
        // first make sure that the room exists by calling the API
        fetch(`/api/rooms/${roomId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Room not found');
            }
            startListeningForMessages(roomId);
            //roomInfoDisplay.textContent = `Room ID: ${roomId}`;
            chatLink.href = `?roomId=${roomId}`;
            chatContainer.classList.remove("d-none");
        })
        .catch(error => {
            console.error('Error joining room:', error);
            // window.location.href = '/';
            //showToast("Room not found. Please check the room ID and try again. " + error, "errorToast");
            document.getElementById("chat-room-error-message").classList.remove("d-none");
        });
    }

    function startListeningForMessages(roomId) {
        const socket = new WebSocket(`ws://${window.location.host}/api/rooms/${roomId}/join`);
        socket.onmessage = function(event) {
            const data = JSON.parse(event.data);
            if (data.type === 'chat') {
                displayMessage(data.sender, data.content);
            } else if (data.type === 'clientName') {
                updateClientName(data.content);
            } else if (data.type === 'info') {
                displayMessage(data.sender, data.content);
            }
        };
        socket.onerror = function(error) {
            console.error('WebSocket error:', error);
        };

        sendMessageButton.addEventListener('click', function(event) {
            event.preventDefault();
            const message = messageInput.value.trim();
            if (message && socket && socket.readyState === WebSocket.OPEN) {
                // const messageData = {
                //     sender: 'You',
                //     message: message
                // };
                socket.send(JSON.stringify(message));
                messageInput.value = '';
            }
        });
    }

    function displayMessage(sender, message) {
        const messageElement = document.createElement('div');
        messageElement.textContent = `${sender}: ${message}`;
        messagesContainer.appendChild(messageElement);
        chatBox.scrollTop = chatBox.scrollHeight; // Scroll to the bottom
    }

    function updateClientName(clientName) {
        clientNameDisplay.textContent = `You are: ${clientName}`;
    }

    document.getElementById("chat-link-copy-button").addEventListener("click", function() {
        const link = document.getElementById("chat-link").href;

        navigator.clipboard.writeText(link).then(() => {
            showToast("Copied to clipboard!", "successToast");
        }).catch(err => {
            showToast("Failed to copy: " + err, "errorToast");
        });
    });

    function showToast(message, toastId) {
        const toastElement = document.getElementById(toastId);
        toastElement.querySelector(".toast-body").textContent = message;
        new bootstrap.Toast(toastElement).show();
    }

    window.onbeforeunload = function() {
        return "Are you sure you want to leave? You will lose your chat history.";
    }
});