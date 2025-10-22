function testMessage(msg) {
    alert(msg);
}

class Task {
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

    setDetails2Div(details2Div) {
        this.details2Div = details2Div;

        return this;
    }

    renderTaskLine() {
        let taskElement = document.createElement('div')
        taskElement.appendChild(document.createTextNode(this.message))

        taskElement.addEventListener('click', (e) => {
            this.renderTaskDetails(e);
            // TODO pass through parameters this.details2Div
            this.renderTaskDetailsFull(e, this.details2Div);
        });
        taskElement.style.border = "solid thin black";

        return taskElement;
    }

    renderTaskDetails(event) {
        let identifierTag = document.getElementById('identifier');
        let messageTag = document.getElementById('message');
        let createdAtTag = document.getElementById('createdAt');
        let resultTag = document.getElementById('result');
        let statusTag = document.getElementById('status');
        let senderTag = document.getElementById('sender');
        let receiversTag = document.getElementById('receivers');

        let sentAtTag = document.getElementById('sentAt');
        let receivedAtTag = document.getElementById('receivedAt');
        let executingAtTag = document.getElementById('executingAt');
        let finishedAtTag = document.getElementById('finishedAt');
        let closedAtTag = document.getElementById('closedAt');

        // put values
        setTextNodeValueToTag(identifierTag, this.id);
        messageTag.value = this.message;
        resultTag.value = this.result;
        // TODO select box assignment
        // statusTag.setAttribute('value', this.status);
        statusTag.value = this.status;


        senderTag.value = this.sender;
        this.setReceiversTag(receiversTag, this.receivers);

        setTextNodeValueToTag(createdAtTag, this.createdAt);
        setTextNodeValueToTag(sentAtTag, this.sentAt);
        setTextNodeValueToTag(receivedAtTag, this.receivedAt);
        setTextNodeValueToTag(executingAtTag, this.executingAt);
        setTextNodeValueToTag(finishedAtTag, this.finishedAt);
        setTextNodeValueToTag(closedAtTag, this.closedAt);

    }

    setReceiversTag(receiversTag, receivers) {
        let receiversSet = new Map();
        for (let receiver of receivers) {
            receiversSet.set(receiver.toString(), receiver);
        }
        Array.from(receiversTag.options).forEach(function (option) {
            option.selected = receiversSet.has(option.value);
        });
    }

    renderTaskDetailsFull(event, parentDiv) {
        clearTag(parentDiv);
        let tableDiv = document.createElement('table');
        parentDiv.appendChild(tableDiv);

        this.renderTrId(tableDiv)
        this.renderTrMessage(tableDiv)
        this.renderTrResult(tableDiv)
        this.renderStatus(tableDiv)
        this.renderSender(tableDiv)
        this.renderReceivers(tableDiv)
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

    // TODO render createdAt
    // TODO render sentAt
    // TODO render receivedAt
    // TODO render executingAt
    // TODO render finishedAt
    // TODO render closedAt
    // TODO save button

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

}

function setTextNodeValueToTag(tag, value) {
    clearTag(tag);

    tag.appendChild(document.createTextNode(value));
}

function clearTag(tag) {
    while (tag.firstChild) {
        tag.removeChild(tag.firstChild);
    }
}

class TasksList {
    constructor() {
        this.tasks = [];
        // inner id for a generated dom element
        this.tasksListId = 'tasksList';
    }

    addTask(task) {
        this.tasks.push(task)
    }

    /**
     * @return dom element with the messages containing data
     */
    renderTasks() {
        let tasksListElement = document.createElement('div');
        tasksListElement.setAttribute('id', this.tasksListId);

        for (let task of this.tasks) {
            // console.log( 'message:', message )
            tasksListElement.appendChild(task.renderTaskLine())
        }

        let addButton = document.createElement('button');
        addButton.appendChild(document.createTextNode('+'));
        tasksListElement.appendChild(addButton);

        addButton.addEventListener('click', this.addTaskPressed);

        return tasksListElement;
    }

    addTaskPressed(event) {
        alert('TODO add task');
    }
}