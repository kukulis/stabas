const TASK_SAVE_BUTTON = 'task_save_button';


class Task {


    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /** @type {int} */
    id = 0;
    /** @type {String} */
    message = "";
    /** @type {String} */
    result = null;
    /** @type {int} */
    status = 0;
    /** @type {int} */
    sender = 0;
    /** @type {[int]} */
    receivers = [];

    /** @type {Date} */
    createdAt = null;
    /** @type {Date} */
    sentAt = null;
    /** @type {Date} */
    receivedAt = null;
    /** @type {Date} */
    executingAt = null;
    /** @type {Date} */
    finishedAt = null;
    /** @type {Date} */
    closedAt = null;

    /** @type {boolean} */
    modified = false;

    constructor(message, id, date) {
        this.message = message;
        this.id = id;
        this.createdAt = date;
    }

    /**
     *
     * @returns {HTMLDivElement}
     */
    renderTaskLine(participantLoader) {

        // TODO more details : status, sender, receivers, current status date

        let taskElement = document.createElement('div')
        taskElement.style.border = "solid thin black";

        let messageDiv = document.createElement('div')
        messageDiv.appendChild(document.createTextNode(this.message))

        taskElement.appendChild(messageDiv);

        messageDiv.addEventListener('click', (e) => {
            let taskDetailsDiv = this.renderTaskDetailsFull(e, participantLoader);
            this.dispatcher.dispatch('taskDetailsRendered', taskDetailsDiv)
        });

        let deleteButton = document.createElement('button');
        deleteButton.appendChild(document.createTextNode('-'));
        const thisTask = this;
        deleteButton.addEventListener('click', (event) => {
                console.log('deleteButton click event, thisTask: ', thisTask)
                thisTask.dispatcher.dispatch('deleteTaskPressed', [event, this.id])
            }
        );

        taskElement.appendChild(deleteButton);

        return taskElement;
    }

    saveAction(event) {
        console.log('save button, event', event)

        this.modified = false;

        let messageInput = document.getElementById('message');
        this.message = messageInput.value;

        let resultInput = document.getElementById('result');
        this.result = resultInput.value;

        let statusSelect = document.getElementById('status');
        this.status = statusSelect.value;

        let senderSelect = document.getElementById('sender');
        this.sender = parseInt(senderSelect.value);

        let receiversSelect = document.getElementById('receivers');
        this.receivers = Array.from(receiversSelect.selectedOptions).map((option) => parseInt(option.value));

        //
        // console.log('receivers after cycle', this.receivers);
        // TODO other fields ( dates ? )

        this.dispatcher.dispatch('taskSaved', this)
    }

    /*******************************************************************************
     * renderTaskDetailsFull
     *
     *******************************************************************************
     */
    renderTaskDetailsFull(event, participantsLoader) {
        let innerDetailsDiv = document.createElement('div')

        // clearTag(parentDiv);
        let tableDiv = document.createElement('table');
        // parentDiv.appendChild(tableDiv);

        this.renderTrId(tableDiv)
        this.renderTrMessage(tableDiv)
        this.renderTrResult(tableDiv)
        this.renderStatus(tableDiv)
        this.renderSender(tableDiv, participantsLoader)
        this.renderReceivers(tableDiv, participantsLoader)
        this.renderDates(tableDiv)

        let saveButton = document.createElement('button')
        saveButton.appendChild(document.createTextNode('save'))
        saveButton.addEventListener('click', (e) => this.saveAction(e))
        saveButton.setAttribute('id', TASK_SAVE_BUTTON)
        saveButton.disabled = true;

        innerDetailsDiv.appendChild(tableDiv)
        innerDetailsDiv.appendChild(saveButton)

        return innerDetailsDiv
    }

    renderTrId(tableDiv) {
        let trId = document.createElement('tr');
        let tdIdLabel = document.createElement('td');
        tdIdLabel.appendChild(document.createTextNode('ID'));
        let tdIdValue = document.createElement('td');
        tdIdValue.appendChild(document.createTextNode(this.id.toString()));
        trId.appendChild(tdIdLabel)
        trId.appendChild(tdIdValue)

        tableDiv.appendChild(trId)
    }

    renderTrMessage(tableDiv) {
        let tr = document.createElement('tr');
        let tdLabel = document.createElement('td');

        let label = document.createElement('label');
        label.appendChild(document.createTextNode('Message'));
        label.setAttribute('for', 'message')
        tdLabel.appendChild(label);

        let tdValue = document.createElement('td');
        let inputField = document.createElement('input')
        tdValue.appendChild(inputField);
        inputField.setAttribute('type', 'text');
        inputField.setAttribute('name', 'message');
        inputField.setAttribute('id', 'message');
        inputField.setAttribute('size', '80');
        inputField.value = this.message;

        inputField.addEventListener('input', (e) => dispatcher.dispatch('inputMessage', [e, this.id]))

        tr.appendChild(tdLabel)
        tr.appendChild(tdValue)

        tableDiv.appendChild(tr)
    }

    renderTrResult(tableDiv) {
        let tr = document.createElement('tr');
        let tdLabel = document.createElement('td');

        let label = document.createElement('label');
        label.appendChild(document.createTextNode('Result'));
        label.setAttribute('for', 'result')
        tdLabel.appendChild(label);

        let tdValue = document.createElement('td');
        let inputField = document.createElement('input')
        tdValue.appendChild(inputField);
        inputField.setAttribute('type', 'text');
        inputField.setAttribute('name', 'result');
        inputField.setAttribute('id', 'result');
        inputField.setAttribute('size', '80');
        inputField.value = this.result;

        inputField.addEventListener('input', (e) => dispatcher.dispatch('inputResult', [e, this.id]))

        tr.appendChild(tdLabel)
        tr.appendChild(tdValue)

        tableDiv.appendChild(tr)
    }

    renderStatus(tableDiv) {
        let tr = document.createElement('tr');
        let tdLabel = document.createElement('td');

        let label = document.createElement('label');
        label.appendChild(document.createTextNode('Status'));
        label.setAttribute('for', 'status')
        tdLabel.appendChild(label);

        let tdValue = document.createElement('td');
        let selectField = document.createElement('select')
        tdValue.appendChild(selectField);
        // inputField.setAttribute('type', 'text' );
        selectField.setAttribute('name', 'status');
        selectField.setAttribute('id', 'status');
        selectField.value = this.status;

        let options = {
            new: "New",
            sent: "Sent",
            received: "Received",
            executing: "Executing",
            finished: "Finished",
            closed: "Closed",
        }

        for (const [key, value] of Object.entries(options)) {
            let optionTag = document.createElement('option');
            optionTag.setAttribute('value', key)
            optionTag.appendChild(document.createTextNode(value))

            if (parseInt(key) === this.status) {
                optionTag.setAttribute('selected', 'selected')
            }

            selectField.appendChild(optionTag)
        }
        selectField.addEventListener('input', (e) => dispatcher.dispatch('inputStatus', [e, this.id]))

        tr.appendChild(tdLabel)
        tr.appendChild(tdValue)

        tableDiv.appendChild(tr)
    }

    renderSender(tableDiv, participantsLoader) {
        let tr = document.createElement('tr');
        let tdLabel = document.createElement('td');

        let label = document.createElement('label');
        label.appendChild(document.createTextNode('Sender'));
        label.setAttribute('for', 'sender')
        tdLabel.appendChild(label);

        let tdValue = document.createElement('td');
        let selectField = document.createElement('select')
        tdValue.appendChild(selectField);
        selectField.setAttribute('name', 'sender');
        selectField.setAttribute('id', 'sender');
        selectField.value = this.sender;

        /**
         * @type {[Participant]}
         */
        let participants = participantsLoader();

        for (let participant of participants) {
            let optionTag = document.createElement('option');
            optionTag.setAttribute('value', participant.id.toString())
            optionTag.appendChild(document.createTextNode(participant.name))

            if (participant.id === this.sender) {
                optionTag.setAttribute('selected', 'selected')
            }

            selectField.appendChild(optionTag)
        }

        selectField.addEventListener('input', (e) => dispatcher.dispatch('inputSender', [e, this.id]))

        tr.appendChild(tdLabel)
        tr.appendChild(tdValue)

        tableDiv.appendChild(tr)
    }

    renderReceivers(tableDiv, participantsLoader) {
        let tr = document.createElement('tr');
        let tdLabel = document.createElement('td');

        let label = document.createElement('label');
        label.appendChild(document.createTextNode('Receivers'));
        label.setAttribute('for', 'receivers')
        tdLabel.appendChild(label);

        let tdValue = document.createElement('td');
        let selectField = document.createElement('select')
        tdValue.appendChild(selectField);
        selectField.setAttribute('multiple', 'multiple');
        selectField.setAttribute('name', 'receivers');
        selectField.setAttribute('id', 'receivers');
        // selectField.value = this.receivers;

        /**
         * @type {[Participant]}
         */
        let participants = participantsLoader();

        let receiversMap = new Map();
        for (let receiver of this.receivers) {
            receiversMap.set(receiver, receiver)
        }

        for (let participant of participants) {
            let optionTag = document.createElement('option');
            optionTag.setAttribute('value', '' + participant.id)
            optionTag.appendChild(document.createTextNode(participant.name))

            if (receiversMap.has(participant.id)) {
                optionTag.setAttribute('selected', 'selected')
            }

            selectField.appendChild(optionTag)
        }

        selectField.addEventListener('input', (e) => dispatcher.dispatch('inputReceivers', [e, this.id]))

        tr.appendChild(tdLabel)
        tr.appendChild(tdValue)

        tableDiv.appendChild(tr)
    }

    renderTextTr(tableDiv, labelText, value) {
        let trId = document.createElement('tr');
        let tdIdLabel = document.createElement('td');
        tdIdLabel.appendChild(document.createTextNode(labelText));
        let tdIdValue = document.createElement('td');
        tdIdValue.appendChild(document.createTextNode(value));
        trId.appendChild(tdIdLabel)
        trId.appendChild(tdIdValue)

        tableDiv.appendChild(trId)
    }

    renderDates(tableDiv) {
        this.renderTextTr(tableDiv, 'Created at', this.createdAt)
        this.renderTextTr(tableDiv, 'Sent at', this.sentAt)
        this.renderTextTr(tableDiv, 'Received at', this.receivedAt)
        this.renderTextTr(tableDiv, 'Executing at', this.executingAt)
        this.renderTextTr(tableDiv, 'Finished at', this.finishedAt)
        this.renderTextTr(tableDiv, 'Closed at', this.closedAt)
    }

    setResult(result) {
        this.result = result;

        return this;
    }

    setStatus(status) {
        this.status = status;

        return this;
    }

    setSender(sender) {
        this.sender = sender;

        return this;
    }

    setReceivers(receivers) {
        this.receivers = receivers;

        return this;
    }

    setDispatcher(dispatcher) {
        this.dispatcher = dispatcher;

        return this;
    }


}

function clearTag(tag) {
    while (tag.firstChild) {
        tag.removeChild(tag.firstChild);
    }
}
