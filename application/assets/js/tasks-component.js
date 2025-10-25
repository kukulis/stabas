
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

    // TODO participants list

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
            // TODO participants list
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

    loadParticipants() {
        // TODO from backend

        // TODO remove after backend

        // --

    }
}