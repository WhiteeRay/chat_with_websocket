package org.example.chat_with_websocket.model;


import lombok.Getter;
import lombok.Setter;



@Getter
@Setter
public class ChatMessage {
    private String content;
    private String sender;
    private MessageType type;

    public enum MessageType{
        CHAT,LEAVE,JOIN
    }
}


