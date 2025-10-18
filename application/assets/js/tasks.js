function testMessage(msg) {
    alert(msg);
}

class Task {
    constructor(message, id, date) {
        this.message = message;
        this.id = id;
        this.date = date;

        // TODO other fields
    }

    renderTask() {
        let taskElement = document.createElement('div')
        taskElement.appendChild( document.createTextNode(this.message ))

        taskElement.addEventListener('click', (e) => this.displayMessageDetails(e))
        taskElement.style.border  = "solid thin black";

        return taskElement;
    }

    displayMessageDetails(event) {
        // console.log('displayMessageDetails this', this)
        alert('TODO displayMessageDetails, id='+this.id);
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
            tasksListElement.appendChild(task.renderTask())
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