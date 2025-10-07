function testMessage(msg) {
    alert(msg);
}

class Message {
    constructor(text) {
        this.text = text;
    }
}

class MessagesList {
    constructor() {
        this.messages = [];
        this.messagesListId = 'messagesList';
    }


    addMessage(message) {
        this.messages.push(message)
    }

    /**
     * @return dom element with the messages containing data
     */
    renderMessages() {
        let messagesListElement =  document.createElement('div');
        messagesListElement.setAttribute('id', this.messagesListId );

        for ( let message of this.messages ) {
            // console.log( 'message:', message )

            let messageElement = document.createElement('div')
            messageElement.appendChild( document.createTextNode(message.text ))

            messagesListElement.appendChild(messageElement)
        }

        return messagesListElement;
    }

}