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
        taskElement.appendChild( document.createTextNode(this.message ))

        taskElement.addEventListener('click', (e) => {
            this.renderTaskDetails(e);
            // TODO pass through parameters this.details2Div
            this.renderTaskDetailsFull(e, this.details2Div);
        });
        taskElement.style.border  = "solid thin black";

        return taskElement;
    }

    renderTaskDetails(event) {
        // console.log('displayMessageDetails this', this)
        // alert('TODO displayMessageDetails, id='+this.id);
        // will set to available tags first

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
        this.setReceiversTag( receiversTag, this.receivers);

        setTextNodeValueToTag(createdAtTag, this.createdAt);
        setTextNodeValueToTag(sentAtTag, this.sentAt);
        setTextNodeValueToTag(receivedAtTag, this.receivedAt);
        setTextNodeValueToTag(executingAtTag, this.executingAt);
        setTextNodeValueToTag(finishedAtTag, this.finishedAt);
        setTextNodeValueToTag(closedAtTag, this.closedAt);

    }

    setReceiversTag(receiversTag, receivers) {
        let receiversSet = new Map();
        for ( let receiver of receivers ) {
            receiversSet.set( receiver.toString(), receiver);
        }
        Array.from(receiversTag.options).forEach(function (option) {
            option.selected = receiversSet.has(option.value);
        });
    }

    renderTaskDetailsFull( event, parentDiv ) {
        clearTag(parentDiv);
        let tableDiv = document.createElement('table');
        parentDiv.appendChild(tableDiv);


        let trId = document.createElement('tr');
        tableDiv.appendChild(trId)

        let tdIdLabel = document.createElement('td');
        tdIdLabel.appendChild(document.createTextNode('ID'));
        let tdIdValue = document.createElement('td');
        tdIdValue.appendChild(document.createTextNode(this.id));
        trId.appendChild(tdIdLabel)
        trId.appendChild(tdIdValue)
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

}

function setTextNodeValueToTag(tag, value) {
    clearTag(tag);

    tag.appendChild(document.createTextNode(value));
}

function clearTag(tag) {
    while(tag.firstChild){
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
        let tasksListElement =  document.createElement('div');
        tasksListElement.setAttribute('id', this.tasksListId );

        for ( let task of this.tasks ) {
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