import React, { useEffect, useState } from "react";
import "./App.css";

const socket = new WebSocket("ws://localhost:8000/ws");

function App() {
  const [messages, setMessages] = useState([]);
  const [messageInput, setMessageInput] = useState("");
  const [usernameInput, setUsernameInput] = useState("");

  useEffect(() => {
    // Handle incoming messages from the WebSocket connection
    socket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      setMessages((messages) => [...messages, message]);
    };

    // Clean up the WebSocket connection when component unmounts
    return () => {
      socket.close();
    };
  }, []);

  const handleSend = (event) => {
    event.preventDefault();

    const message = {
      username: usernameInput,
      content: messageInput,
    };

    // Send the message through the WebSocket connection
    socket.send(JSON.stringify(message));

    setMessageInput("");
  };

  return (
    <div className="App">
      <h1>Chat Application</h1>
      <div className="MessageContainer">
        {messages.map((message, index) => (
          <div key={index} className="Message">
            <span className="Username">{message.username}: </span>
            <span className="Content">{message.content}</span>
          </div>
        ))}
      </div>
      <form onSubmit={handleSend} className="InputForm">
        <input
          type="text"
          value={usernameInput}
          onChange={(event) => setUsernameInput(event.target.value)}
          placeholder="Enter your username"
          required
        />
        <input
          type="text"
          value={messageInput}
          onChange={(event) => setMessageInput(event.target.value)}
          placeholder="Enter your message"
          required
        />
        <button type="submit">Send</button>
      </form>
    </div>
  );
}

export default App;
