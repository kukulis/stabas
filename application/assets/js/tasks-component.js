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
    }

    addTask(task) {
        this.tasks.push(task)
    }

    /**
     * @returns {Promise<Task|null>}
     */
    async newTask() {
        let maxId = this.tasks.reduce((prev, curr) => {
            return Math.max(prev, curr.id);
        }, 0)

        console.log('newTask, maxId', maxId)
        let response = await fetch('/api/tasks', {
            method: 'PUT',
            body: JSON.stringify({
                message: 'Task ' + (maxId + 1),
                sender: 1,
                receivers: [],
            })
        }).catch((error) => console.log('backend error creating task ', error));

        if (response === undefined) {
            return null;
        }

        // TODO better to return the new created task
        let taskId = await response.json()

        let taskResponse = await fetch('/api/tasks/' + taskId,
            {method: 'GET'}
        ).catch((error) => console.log('error getting task by id ' + taskId, error));

        if (taskResponse === undefined) {
            return null;
        }

        let taskDto = await taskResponse.json();

        let task = Task.createFromDto(taskDto).setDispatcher(this.dispatcher)

        this.tasks.push(task);

        return task;
    }

    deleteTask(event, taskId) {
        event.stopPropagation();

        let foundTask = this.tasks.find((task) => task.id === taskId)

        let confirmDelete = true;
        if (foundTask.status !== STATUS_CLOSED) {
            confirmDelete = confirm("Delete task " + taskId + " ?")
        }

        if (!confirmDelete) {
            return;
        }

        fetch('/api/tasks/' + taskId, {method: 'DELETE'}).then(
            (response) => {
                this.tasks = this.tasks.filter((task) => (task.id !== taskId))
                this.dispatcher.dispatch('afterDeleteTask', taskId);
            })
            .catch((error) => console.log('error deleting task ' + taskId, error))
    }

    /**
     * @return dom element with the messages containing data
     */
    renderTasks() {
        let tasksListElement = document.createElement('div');

        let addButton = document.createElement('button');
        addButton.appendChild(document.createTextNode('+'));
        addButton.setAttribute('class', 'add-task')
        addButton.addEventListener('click', () => {
            this.newTask()
                .then((task) => this.dispatcher.dispatch('onAddTask', task));
        });

        tasksListElement.appendChild(addButton);

        for (let task of this.tasks) {
            tasksListElement.appendChild(task.renderTaskLine(() => this.participants));
        }

        return tasksListElement;
    }

    async loadParticipants() {
        let response = await fetch("/api/participants", {
            method: "GET",
        });

        let participantsDTO = await response.json();

        for (let p of participantsDTO) {
            this.participants.push(new Participant(p.id, p.name))
        }
    }

    async loadTasks() {
        let response = await fetch("/api/tasks", {
            method: "GET",
        })
            .catch((error) => {
                console.log('error fetching tasks', error)
            });

        if (response === undefined) {
            console.log('loadTasks response is undefined')
            return;
        }

        let tasksDto = await response.json();
        // console.log('tasksDto', tasksDto)


        for (let taskDto of tasksDto) {
            // console.log('loadTasks: adding task from backend, status ', taskDto.status)
            this.addTask(
                Task.createFromDto(taskDto).setDispatcher(this.dispatcher)
            );
        }

        // console.log ('tasks after loadTasks ',  this.tasks );
    }

    static async initialize(dispatcher) {
        let tasksComponent = new TasksComponent(dispatcher);

        await tasksComponent.loadParticipants();
        await tasksComponent.loadTasks();

        tasksComponent.startTimer();

        return tasksComponent;
    }

    enableSaveButton() {
        let saveButton = document.getElementById(TASK_SAVE_BUTTON);
        saveButton.disabled = false;
    }

    disableSaveButton() {
        let saveButton = document.getElementById(TASK_SAVE_BUTTON);
        saveButton.disabled = true;
    }

    startTimer() {
        var intervalId = setInterval(()=> {
            this.dispatcher.dispatch('timerTick')
            // clearInterval(intervalId)
        }, 1000)
    }

    setTimers(now) {
        for(let task of this.tasks) {
            task.setTimer(now)
        }
    }
}