/***********************************************************************
 *
 * TasksComponent
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

    /**
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     *
     * @type {Settings}
     */
    settings = null;

    /**
     *
     * @param dispatcher {Dispatcher}
     * @param apiClient {ApiClient}
     */
    constructor(dispatcher, apiClient) {
        this.dispatcher = dispatcher;
        this.apiClient = apiClient;
    }

    addTask(task) {
        this.tasks.push(task)
    }

    /**
     * Called when 'add task' button is pressed (currently with '+' label )
     * @returns {Promise<Task|null>}
     */
    async newTask() {
        let maxId = this.tasks.reduce((prev, curr) => {
            return Math.max(prev, curr.id);
        }, 0)

        console.log('newTask, maxId', maxId)
        // TODO use ApiClient
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

        let taskDto = await response.json()

        // let taskResponse = await fetch('/api/tasks/' + taskId,
        //     {method: 'GET'}
        // ).catch((error) => console.log('error getting task by id ' + taskId, error));
        //
        // if (taskResponse === undefined) {
        //     return null;
        // }
        // let taskDto = await taskResponse.json();

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

        // TODO use ApiClient
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
            tasksListElement.appendChild(task.renderTaskLine(() => this.participants, this.settings));
        }

        return tasksListElement;
    }

    async loadParticipants() {
        // TODO use ApiClient
        let response = await fetch("/api/participants", {
            method: "GET",
        });
        // TODO handle errors

        let participantsDTO = await response.json();

        for (let p of participantsDTO) {
            this.participants.push(new Participant(p.id, p.name))
        }
    }

    async loadTasks() {
        // TODO use Api client
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

    async loadSettings() {
        this.settings = await this.apiClient.loadSettings()
    }

    /**
     * @param dispatcher {Dispatcher}
     * @param apiClient {ApiClient}
     * @returns {Promise<TasksComponent>}
     */
    static async initialize(dispatcher, apiClient) {
        let tasksComponent = new TasksComponent(dispatcher, apiClient);

        await tasksComponent.loadParticipants();
        await tasksComponent.loadTasks();
        await tasksComponent.loadSettings();

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

    // TODO make possibility to start and stop timer
    startTimer() {
        var intervalId = setInterval(()=> {
            this.dispatcher.dispatch('timerTick')
            // clearInterval(intervalId)
        }, 5000)
    }

    async reloadSettings() {
        // TODO
    }

    async reloadParticipants() {
        // TODO
    }

    async reloadTasks() {

        let tasksDto = await this.apiClient.loadTasks()

        let tasksMap = new Map();

        for ( let task of  this.tasks ) {
             tasksMap.set(task.id, task)
        }

        /**
         * @type {Task[]}
         */
        let newTasksList = [];

        for ( let taskDto of tasksDto ) {
            let taskId = taskDto.id;
            let task = null;
            if ( tasksMap.has(taskId)) {
                task = tasksMap.get(taskId)
                task.updateFromDTO(taskDto)
            }

            if ( task === null ) {
                task = Task.createFromDto(taskDto).setDispatcher(this.dispatcher)
            }

            newTasksList.push(task)
        }

        this.tasks = newTasksList;
    }

    /**
     *
     * @param task {Task}
     * @returns {Promise<void>}
     */
    async reloadSingleTask(task) {
         let taskDto = await this.apiClient.loadTask(task.id)

        task.updateFromDTO(taskDto)
    }
}