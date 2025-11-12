const TASK_SAVE_BUTTON = 'task_save_button';

let statusesI2A = new Map();
statusesI2A.set(1, 'new')
statusesI2A.set(2, 'sent')
statusesI2A.set(3, 'received')
statusesI2A.set(4, 'executing')
statusesI2A.set(5, 'finished')
statusesI2A.set(6, 'closed')

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
     * @returns {Task}
     */
    static createFromDto(taskDTO) {
        let task = new Task(taskDTO.message, taskDTO.id, parseDate(taskDTO.created_at))

        return task.updateFromDTO(taskDTO)
    }

    /**
     * @param taskDTO
     * @returns {Task}
     */
    updateFromDTO(taskDTO) {
        this.setStatus(taskDTO.status)
        this.setSender(taskDTO.sender)
        this.setReceivers(taskDTO.receivers)
        this.setResult(taskDTO.result)
        this.sentAt = parseDate(taskDTO.sent_at)
        this.receivedAt = parseDate(taskDTO.received_at)
        this.executingAt = parseDate(taskDTO.executing_at)
        this.finishedAt = parseDate(taskDTO.finished_at)
        this.closedAt = parseDate(taskDTO.closed_at)

        return this;
    }

    /**
     *
     * @returns {HTMLDivElement}
     */
    renderTaskLine(participantLoader) {

        // TODO more details : sender, receivers, current status date
        let taskElement = document.createElement('div')
        taskElement.setAttribute('class', 'task-line')

        let deleteButton = document.createElement('button');
        deleteButton.appendChild(document.createTextNode('-'));
        const thisTask = this;
        deleteButton.addEventListener('click', (event) => {
                thisTask.dispatcher.dispatch('deleteTaskPressed', [event, this.id])
            }
        );
        deleteButton.setAttribute('class', 'delete-button')

        taskElement.appendChild(deleteButton);

        let messageDiv = document.createElement('div')
        messageDiv.setAttribute('class', 'message')
        messageDiv.appendChild(document.createTextNode('' + this.id + ': ' + this.message))
        messageDiv.addEventListener('click', (e) => {
            let taskDetailsDiv = this.renderTaskDetailsFull(e, participantLoader);
            this.dispatcher.dispatch('taskDetailsRendered', taskDetailsDiv)
        });

        taskElement.appendChild(messageDiv);


        if (this.status < 6) {
            let newStatus = getNextStatus(this.status);
            let newStatusName = statusesI2A.get(newStatus);
            let changeStatusButton = document.createElement('button')
            changeStatusButton.appendChild(document.createTextNode('change to ' + newStatusName))
            changeStatusButton.addEventListener('click', (e) => this.changeTaskStatus(e, this, newStatus));
            changeStatusButton.setAttribute('class', 'change-status-button')

            taskElement.appendChild(changeStatusButton)
        }

        let statusDiv = document.createElement('div')
        statusDiv.setAttribute('class', 'status')
        statusDiv.appendChild(document.createTextNode('[' + this.status + ': ' + statusesI2A.get(this.status) + ']'))

        taskElement.appendChild(statusDiv)

        let clearDiv = document.createElement('div')
        clearDiv.setAttribute('class', 'clear')

        taskElement.appendChild(clearDiv)

        return taskElement;
    }

    changeTaskStatus(event, task, newStatus) {
        fetch('/api/tasks/' + task.id + '/change-status?status=' + newStatus, {
            method: 'POST',
        }).then((response) => {
            console.log('response received after changing status ', response)
            response.text().then((text) => {
                if (response.status === 200) {
                    task.status = newStatus;
                    this.dispatcher.dispatch('afterChangeStatus', [event, task.id]);
                    return;
                }
                alert(text)
            }).catch((error) => console.log('error getting response after changing status', error))
        }).catch((error) => console.log('error changing status', error))
    }

    buildObjectForJson() {
        return {
            id: this.id,
            message: this.message,
            result: this.result,
            status: this.status,
            sender: this.sender,
            receivers: this.receivers,
        }
    }

    saveAction(event) {
        console.log('save button, event', event)

        this.modified = false;

        let messageInput = document.getElementById('message');
        this.message = messageInput.value;

        let resultInput = document.getElementById('result');
        this.result = resultInput.value;

        let statusSelect = document.getElementById('status');
        this.status = parseInt(statusSelect.value);

        let senderSelect = document.getElementById('sender');
        this.sender = parseInt(senderSelect.value);

        let receiversSelect = document.getElementById('receivers');
        this.receivers = Array.from(receiversSelect.selectedOptions).map((option) => parseInt(option.value));

        fetch('/api/tasks/' + this.id, {
            method: 'POST',
            body: JSON.stringify(this.buildObjectForJson())
        })
            .catch((error) => console.log('error updating task to backend', error))
            .then((taskResponse) => {
                taskResponse.json().then((taskDTO) => {
                    this.updateFromDTO(taskDTO)
                    this.dispatcher.dispatch('taskSaved', this)
                })
            })
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

        // console.log('task.js[255]: statusesI2A', statusesI2A)
        statusesI2A.forEach((value, key) => {
            let optionTag = document.createElement('option');
            optionTag.setAttribute('value', key)
            optionTag.appendChild(document.createTextNode(value))
            if (parseInt(key) === this.status) {
                optionTag.setAttribute('selected', 'selected')
            }
            selectField.appendChild(optionTag)
        })

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
        this.renderTextTr(tableDiv, 'Created at', formatDate(this.createdAt))
        this.renderTextTr(tableDiv, 'Sent at', formatDate(this.sentAt))
        this.renderTextTr(tableDiv, 'Received at', formatDate(this.receivedAt))
        this.renderTextTr(tableDiv, 'Executing at', formatDate(this.executingAt))
        this.renderTextTr(tableDiv, 'Finished at', formatDate(this.finishedAt))
        this.renderTextTr(tableDiv, 'Closed at', formatDate(this.closedAt))
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


function getNextStatus(status) {
    if (status === 6) {
        return status;
    }

    return status + 1;
}
