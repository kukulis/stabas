
/***********************************************************************
 *
 * TaskList
 *
 ***********************************************************************/
class TasksComponent {

    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @type {[Task]}
     */
    tasks = [];

    /**
     * @type {[Participant]}
     */
    participants = [];

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
        event.stopPropagation();

        let confirmDelete = confirm("Delete task "+taskId+" ?")

        if ( !confirmDelete ) {
            return;
        }

        this.tasks = this.tasks.filter((task) => (task.id !== taskId))
        this.dispatcher.dispatch('afterDeleteTask', taskId);
    }

    /**
     * @return dom element with the messages containing data
     */
    renderTasks() {
        let tasksListElement = document.createElement('div');

        for (let task of this.tasks) {
            tasksListElement.appendChild(task.renderTaskLine(()=>this.participants));
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

    loadParticipants() {
        // TODO from backend

        // TODO remove after backend
        this.participants.push( new Participant(1, "loaded Participant 1"))
        this.participants.push( new Participant(2, "loaded Participant 2"))
        this.participants.push( new Participant(3, "L Participant 3"))
        // --

    }

    loadTasks() {
        // TODO from backend

        // TODO remove after backend
        this.addTask(
            (new Task("Prepare to commendant hour squad I", 1, new Date()))
                .setStatus("new")
                .setSender(1)
                .setReceivers([1, 2])
                .setResult('aaa')
                .setDispatcher(this.dispatcher)
        );
        this.addTask(
            (new Task("Prepare to commendant hour squad II", 2, new Date()))
                .setStatus("sent")
                .setSender(2)
                .setReceivers([2, 3])
                .setResult('bbb')
                .setDispatcher(this.dispatcher)
        );
        this.addTask(
            (new Task("Prepare to commendant hour squad III", 3, new Date()))
                .setStatus("received")
                .setSender(3)
                .setReceivers([1, 3])
                .setResult('ccc')
                .setDispatcher(this.dispatcher)
        );
        // --
    }

    static initialize(dispatcher) {
        let tasksComponent = new TasksComponent(dispatcher);

        tasksComponent.loadParticipants();
        tasksComponent.loadTasks();

        return tasksComponent;
    }
    enableSaveButton() {
        let saveButton = document.getElementById('task_save_button' );
        saveButton.disabled = false;
    }

    disableSaveButton() {
        let saveButton = document.getElementById('task_save_button' );
        saveButton.disabled = true;
    }
}