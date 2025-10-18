function testMessage(msg) {
    alert(msg);
}

// TODO Task
class Message {
    constructor(text, id, date) {
        this.text = text;
        this.id = id;
        this.date = date;
    }
}

// TODO TasksList
class MessagesList {
    constructor(messageDisplayContainerId) {
        this.messages = [];
        this.messagesListId = 'messagesList';
        this.messageDisplayContainerId = messageDisplayContainerId;
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

        let addButton = document.createElement('button');
        addButton.appendChild(document.createTextNode('+'));
        messagesListElement.appendChild(addButton);

        addButton.addEventListener('click', (e)=>this.addMessagePressed(e));

        return messagesListElement;
    }

    addMessagePressed(event) {
        alert('TODO add message');
    }

    renderMessageDetails(message) {

    }

}