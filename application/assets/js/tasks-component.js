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
     * @type {[Task]|[TaskGroup]}
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
        let task = TaskGroup.createFromDto(taskDto).setDispatcher(this.dispatcher)

        this.tasks.push(task);

        return task;
    }

    findTask(taskId) {
        for ( let task of this.tasks) {
            let foundTask = task.findTask(taskId)

            if ( foundTask !== null ) {
                return foundTask
            }
        }

        return null;
    }

    /**
     *
     * @param event
     * @param taskId
     */
    deleteTask(event, taskId) {
        event.stopPropagation();

        let foundTask = this.findTask(taskId)

        if ( foundTask === null) {
            console.log('error: no task found by taskId '+taskId)
            return;
        }

        let confirmDelete = true;
        if (foundTask.status !== STATUS_CLOSED && foundTask.status !== STATUS_NEW ) {
            confirmDelete = confirm("Delete task " + taskId + " ?")
        }

        if (!confirmDelete) {
            return;
        }

        this.apiClient.deleteTask(taskId).then(
            (success) => {
                if ( !success) {
                    return
                }
                this.tasks = this.tasks.filter((task) => (task.id !== taskId))
                // delete task from each child
                this.tasks.forEach((task)=>task.deleteTask(taskId))
                this.dispatcher.dispatch('afterDeleteTask', taskId);
            })
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

        this.tasks = [];
        let groupsDto = await this.apiClient.loadGroups();

        for (let groupDto of groupsDto) {
            // console.log('loadTasks: adding task from backend, status ', taskDto.status)
            let taskGroup = TaskGroup.createFromDto(groupDto).setDispatcher(this.dispatcher)
            this.addTask( taskGroup )
        }
        console.log ('setting dispatcher in the loadTasks ',  this.tasks );
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

        let tasksDto = await this.apiClient.loadGroups()

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
                task.updateFromDTO(taskDto).setDispatcher(dispatcher)
            }

            if ( task === null ) {
                task = TaskGroup.createFromDto(taskDto).setDispatcher(this.dispatcher)
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