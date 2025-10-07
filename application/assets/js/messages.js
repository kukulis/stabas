function testMessage(msg) {
    alert(msg);
}

class Message {
    Message(text) {
        this.text = text;
        // TODO
    }


}

class MessagesList {
    MessagesList() {
        this.messages = [];
        this.messagesListId = 'messagesList';
    }


    addMessage(message) {
        this.messages.add(message)
    }

    /**
     * @return dom element with the messages containing data
     */
    renderMessages() {
        let messagesList =  document.createElement('div');
        messagesList.setAttribute('id', this.messagesListId );

    }

}