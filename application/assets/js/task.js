const TASK_SAVE_BUTTON = 'task_save_button';

const STATUS_NEW = 1
const STATUS_SENT = 2
const STATUS_RECEIVED = 3
const STATUS_EXECUTING = 4
const STATUS_FINISHED = 5
const STATUS_CLOSED = 6

let statusesI2A = new Map();
statusesI2A.set(STATUS_NEW, 'new')
statusesI2A.set(STATUS_SENT, 'sent')
statusesI2A.set(STATUS_RECEIVED, 'received')
statusesI2A.set(STATUS_EXECUTING, 'executing')
statusesI2A.set(STATUS_FINISHED, 'finished')
statusesI2A.set(STATUS_CLOSED, 'closed')

function getStatusClass(statusId) {
    if (!statusesI2A.has(statusId)) {
        return '';
    }
    return 'status-' + statusesI2A.get(statusId);
}

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

    // TODO single receiver
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

    version = 0;

    taskGroup = 0;

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
        this.message = taskDTO.message
        this.sentAt = parseDate(taskDTO.sent_at)
        this.receivedAt = parseDate(taskDTO.received_at)
        this.executingAt = parseDate(taskDTO.executing_at)
        this.finishedAt = parseDate(taskDTO.finished_at)
        this.closedAt = parseDate(taskDTO.closed_at)
        this.taskGroup = parseDate(taskDTO.task_group)
        this.version = taskDTO.version

        return this;
    }
    updateFromDTOMerged(taskDTO, myVersionTaskDTO) {

        this.status = selectValue(this.status, taskDTO.status, myVersionTaskDTO.status)
        this.sender = selectValue(this.sender, taskDTO.sender, myVersionTaskDTO.sender)
        this.receivers = selectValue(this.receivers, taskDTO.receivers, myVersionTaskDTO.receivers)
        this.message =  selectValue(this.message, taskDTO.message, myVersionTaskDTO.message)
        this.result =  selectValue(this.result, taskDTO.result, myVersionTaskDTO.result)
        this.sentAt = parseDate(taskDTO.sent_at)
        this.receivedAt = parseDate(taskDTO.received_at)
        this.executingAt = parseDate(taskDTO.executing_at)
        this.finishedAt = parseDate(taskDTO.finished_at)
        this.closedAt = parseDate(taskDTO.closed_at)
        this.version = taskDTO.version

        return this;
    }

    getTimerDivId() {
        return 'task-timer-' + this.id.toString();
    }

    /**
     * @param participantLoader {function}
     * @param settings {Settings}
     * @returns {HTMLDivElement}
     */
    renderTaskLine(participantLoader, settings) {

        let taskElement = document.createElement('div')
        // taskElement.setAttribute('class', 'task-line')
        taskElement.classList.add('task-line')
        taskElement.classList.add(getStatusClass(this.status))

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
            let newStatusClass = getStatusClass(newStatus);
            let newStatusName = statusesI2A.get(newStatus);
            let changeStatusButton = document.createElement('button')
            changeStatusButton.appendChild(document.createTextNode('change to ' + newStatusName))
            changeStatusButton.addEventListener('click', (e) => this.changeTaskStatus(e, this, newStatus));
            changeStatusButton.classList.add('change-status-button')
            changeStatusButton.classList.add(newStatusClass)
            // changeStatusButton.setAttribute('class', 'change-status-button')

            taskElement.appendChild(changeStatusButton)
        }

        let statusDiv = document.createElement('div')
        statusDiv.setAttribute('class', 'status')
        statusDiv.appendChild(document.createTextNode('[' + this.status + ': ' + statusesI2A.get(this.status) + ']'))

        taskElement.appendChild(statusDiv)


        let dateDiv = document.createElement('div')
        dateDiv.classList.add('task-date')
        let currentStatusDate = this.getCurrentStatusDate();
        dateDiv.appendChild(document.createTextNode(formatTimer(currentStatusDate)))
        taskElement.appendChild(dateDiv)

        let now = new Date();

        let timerDiv = document.createElement('div')
        timerDiv.classList.add('task-timer')

        let late = 0

        if ( now !== null && currentStatusDate !== null ) {
            late = now.getTime() - currentStatusDate.getTime();
        }

        let criticality = settings.calculateCriticality(late, this.status)

        let duration = this.calculateIntervalFromTheCurrentStatusDate(now);

        timerDiv.appendChild(document.createTextNode(duration))
        timerDiv.setAttribute('id', this.getTimerDivId())
        timerDiv.classList.add('late-' + criticality)
        taskElement.appendChild(timerDiv)

        let clearDiv = document.createElement('div')
        clearDiv.setAttribute('class', 'clear')

        taskElement.appendChild(clearDiv)

        // TODO sender and receivers

        return taskElement;
    }

    changeTaskStatus(event, task, newStatus) {
        fetch('/api/tasks/' + task.id + '/change-status?status=' + newStatus, {
            method: 'POST',
        }).then((response) => {
            // console.log('response received after changing status ', response)
            response.text().then((text) => {
                if (response.status === 200) {
                    task.status = newStatus;
                    console.log(text)
                    this.dispatcher.dispatch('afterChangeStatus', [event, task]);
                }
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
            version: this.version+1,
        }
    }

    saveAction(event) {
        // console.log('save button, event', event)

        this.modified = false;
        let messageInput = document.getElementById('message');
        let resultInput = document.getElementById('result');
        let statusSelect = document.getElementById('status');
        let senderSelect = document.getElementById('sender');
        let receiversSelect = document.getElementById('receivers');

        let myVersionDto = {
            message: messageInput.value,
            result: resultInput.value,
            status: parseInt(statusSelect.value),
            sender: parseInt(senderSelect.value),
            receivers: Array.from(receiversSelect.selectedOptions).map((option) => parseInt(option.value)),
            version: this.version +1,
        }

        //= this.buildObjectForJson()

        fetch('/api/tasks/' + this.id, {
            method: 'POST',
            body: JSON.stringify(myVersionDto)
        })
            .catch((error) => console.log('error updating task to backend', error))
            .then((taskResponse) => {
                taskResponse.json().then((taskDTO) => {
                    if ( taskResponse.status === 409 ) {
                        this.updateFromDTOMerged(taskDTO, myVersionDto)
                        this.dispatcher.dispatch('taskSavedPartially', this)
                    }
                    else {
                        this.updateFromDTO(taskDTO)
                        this.dispatcher.dispatch('taskSaved', this)
                    }


                    // if ( taskResponse.status === 409 ) {
                    //     this.dispatcher.dispatch('inputMessage', this)
                    // }
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

        let hideButton = document.createElement('button')
        hideButton.appendChild(document.createTextNode('hide'))
        hideButton.addEventListener('click', (event) => dispatcher.dispatch('hideDetailsPressed', event))

        innerDetailsDiv.appendChild(tableDiv)
        innerDetailsDiv.appendChild(saveButton)
        innerDetailsDiv.appendChild(hideButton)

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

    getCurrentStatusDate() {
        if (this.status === STATUS_NEW) {
            return this.createdAt
        }
        if (this.status === STATUS_SENT) {
            return this.sentAt
        }
        if (this.status === STATUS_RECEIVED) {
            return this.receivedAt
        }
        if (this.status === STATUS_EXECUTING) {
            return this.executingAt
        }
        if (this.status === STATUS_FINISHED) {
            return this.finishedAt
        }
        if (this.status === STATUS_CLOSED) {
            return this.closedAt
        }
    }

    /**
     *
     * @param date {Date}
     * @return {string}
     */
    // TODO use distance instead of calculating distance inside the function
    calculateIntervalFromTheCurrentStatusDate(date) {

        let statusDate = this.getCurrentStatusDate();
        if (statusDate === null) {
            return '-';
        }

        let distance = date.getTime() - statusDate.getTime()
        let secondsDistance = Math.floor(distance / 1000);
        let minutesDistance = Math.floor(secondsDistance / 60);
        let remainingSecondsDistance = secondsDistance % 60;
        let hoursDistance = Math.floor(minutesDistance / 60);
        let remainingMinutesDistance = minutesDistance % 60;

        return hoursDistance.toString() + ':' + remainingMinutesDistance.toString() + ':' + remainingSecondsDistance.toString();
    }
}


function getNextStatus(status) {
    if (status === 6) {
        return status;
    }

    return status + 1;
}

/**
 *
 * @param date {Date}
 */
function formatTimer(date) {
    if (date === null) {
        return '-';
    }
    return date.getHours().toString() + ':' + date.getMinutes().toString() + ':' + date.getSeconds().toString()
}

function selectValue(originalValue, newValue, anotherValue ) {
    if ( originalValue !== newValue ) {
        return newValue;
    }
    return anotherValue;
}


class TaskGroup extends Task {

    /**
     * @type {[Task]}
     */
    children = [];

}


