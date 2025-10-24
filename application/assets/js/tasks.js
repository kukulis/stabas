class Task {
    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    // TODO all fields as variables


    constructor(message, id, date) {
        this.message = message;
        this.id = id;
        // deprecated
        this.date = date;
        this.createdAt = date;

        this.result = null;
        this.status = null;
        this.sender = null;
        this.receivers = [];

        this.sentAt = null;
        this.receivedAt = null;
        this.executingAt = null;
        this.finishedAt = null;
        this.closedAt = null;

    }

    renderTaskLine() {
        let taskElement = document.createElement('div')
        taskElement.appendChild(document.createTextNode(this.message))

        taskElement.addEventListener('click', (e) => {
            let taskDetailsDiv = this.renderTaskDetailsFull(e);
            this.dispathcer.dispatch('taskDetailsRendered', taskDetailsDiv)
        });
        taskElement.style.border = "solid thin black";

        let deleteButton = document.createElement('button');
        deleteButton.appendChild(document.createTextNode('-'));
        deleteButton.addEventListener('click', (event) =>
            dispatcher.dispatch('deleteTaskPressed', [event, this.id])
        );

        taskElement.appendChild(deleteButton);

        return taskElement;
    }

    saveAction(event) {
        // TODO unset 'modified' flag, which is set when any of the detail elements are modified
        // alert('save button pressed')
        console.log('save button, event', event)

        let messageInput = document.getElementById('message');
        this.message = messageInput.value;

        let resultInput = document.getElementById('result');
        this.result = resultInput.value;

        let statusSelect = document.getElementById('status');
        this.status = statusSelect.value;

        let senderSelect = document.getElementById('sender');
        this.sender = senderSelect.value;

        let receiversSelect = document.getElementById('receivers');
        // receiversSelect.children

        this.receivers = Array.from(receiversSelect.selectedOptions).map((option) => option.value);
        //
        // console.log('receivers after cycle', this.receivers);
        // TODO other fields ( dates ? )

        this.dispathcer.dispatch('taskSaved', this)
    }

    /*******************************************************************************
     * renderTaskDetailsFull
     *
     *******************************************************************************
     */
    renderTaskDetailsFull(event) {
        let innerDetailsDiv = document.createElement('div')

        // clearTag(parentDiv);
        let tableDiv = document.createElement('table');
        // parentDiv.appendChild(tableDiv);

        this.renderTrId(tableDiv)
        this.renderTrMessage(tableDiv)
        this.renderTrResult(tableDiv)
        this.renderStatus(tableDiv)
        this.renderSender(tableDiv)
        this.renderReceivers(tableDiv)
        this.renderDates(tableDiv)

        let saveButton = document.createElement('button')
        saveButton.appendChild(document.createTextNode('save'))
        saveButton.addEventListener('click', (e) => this.saveAction(e))

        // parentDiv.appendChild(saveButton);

        innerDetailsDiv.appendChild(tableDiv)
        innerDetailsDiv.appendChild(saveButton)

        return innerDetailsDiv
    }

    renderTrId(tableDiv) {
        let trId = document.createElement('tr');
        let tdIdLabel = document.createElement('td');
        tdIdLabel.appendChild(document.createTextNode('ID'));
        let tdIdValue = document.createElement('td');
        tdIdValue.appendChild(document.createTextNode(this.id));
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

            if (key === this.status) {
                optionTag.setAttribute('selected', 'selected')
            }

            selectField.appendChild(optionTag)
        }

        tr.appendChild(tdLabel)
        tr.appendChild(tdValue)

        tableDiv.appendChild(tr)
    }

    // TODO all possible participants should be passed through parameters
    renderSender(tableDiv) {
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

        // TODO take values from parameters
        let options = {
            "1": "Participant1",
            "2": "Participant2",
            "3": "Participant3",
        }

        for (const [key, value] of Object.entries(options)) {
            let optionTag = document.createElement('option');
            optionTag.setAttribute('value', key)
            optionTag.appendChild(document.createTextNode(value))

            if (key === ('' + this.sender)) {
                optionTag.setAttribute('selected', 'selected')
            }

            selectField.appendChild(optionTag)
        }

        tr.appendChild(tdLabel)
        tr.appendChild(tdValue)

        tableDiv.appendChild(tr)
    }

    // TODO pass available participants through parameters
    renderReceivers(tableDiv) {
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
        selectField.value = this.receivers;

        // TODO take values from parameters
        let options = {
            "1": "participant1",
            "2": "participant2",
            "3": "participant3",
        }

        let receiversMap = new Map();
        for (let receiver of this.receivers) {
            receiversMap.set(receiver.toString(), receiver)
        }


        for (const [key, value] of Object.entries(options)) {
            let optionTag = document.createElement('option');
            optionTag.setAttribute('value', key)
            optionTag.appendChild(document.createTextNode(value))

            if (receiversMap.has(key)) {
                optionTag.setAttribute('selected', 'selected')
            }

            selectField.appendChild(optionTag)
        }

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

    setId(id) {
        this.id = id;

        return this;
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
        this.dispathcer = dispatcher;

        return this;
    }

}

function clearTag(tag) {
    while (tag.firstChild) {
        tag.removeChild(tag.firstChild);
    }
}

/***********************************************************************
 *
 * TaskList
 *
 ***********************************************************************/
class TasksList {

    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @type {[Task]}
     */
    tasks = [];

    constructor(dispatcher) {
        this.dispatcher = dispatcher;

        // for debuging only
        this.maxTaskId = 3;
    }

    addTask(task) {
        this.tasks.push(task)
    }

    newTask() {
        // TODO with backend later
        this.maxTaskId++;
        let task = new Task("Task " + this.maxTaskId, this.maxTaskId, new Date())
        task.setDispatcher(this.dispatcher)
        this.tasks.push(task);

        return task;
    }

    deleteTask(event, taskId) {
        this.tasks = this.tasks.filter((task) => (task.id !== taskId))
        event.stopPropagation();
        dispatcher.dispatch('afterDeleteTask', taskId);
    }

    /**
     * @return dom element with the messages containing data
     */
    renderTasks() {
        let tasksListElement = document.createElement('div');

        for (let task of this.tasks) {
            tasksListElement.appendChild(task.renderTaskLine());
        }

        let addButton = document.createElement('button');
        addButton.appendChild(document.createTextNode('+'));
        tasksListElement.appendChild(addButton);

        addButton.addEventListener('click', () => {
            let task = this.newTask();
            this.dispatcher.dispatch('onAddTask', task)
        });

        return tasksListElement;
    }
}