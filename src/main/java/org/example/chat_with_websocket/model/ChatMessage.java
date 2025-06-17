package org.javatechie.spring.ws.api.model;


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
