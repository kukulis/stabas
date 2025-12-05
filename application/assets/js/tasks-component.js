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
        let taskDto = await this.apiClient.createTask({
            message: 'Task ' + (maxId + 1),
            sender: 1,
            receivers: [],
        })

        if (taskDto === null) {
            throw new Error('Could not create a new task')
        }
        let task = TaskGroup.createFromDto(taskDto).setDispatcher(this.dispatcher)

        this.tasks.push(task);

        return task;
    }

    findTask(taskId) {
        for (let task of this.tasks) {
            let foundTask = task.findTask(taskId)

            if (foundTask !== null) {
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

        if (foundTask === null) {
            console.log('error: no task found by taskId ' + taskId)
            return;
        }

        let confirmDelete = true;
        if (foundTask.status !== STATUS_CLOSED && foundTask.status !== STATUS_NEW) {
            confirmDelete = confirm("Delete task " + taskId + " ?")
        }

        if (!confirmDelete) {
            return;
        }

        this.apiClient.deleteTask(taskId).then(
            (success) => {
                if (!success) {
                    return
                }
                this.tasks = this.tasks.filter((task) => (task.id !== taskId))
                // delete task from each child
                this.tasks.forEach((task) => task.deleteTask(taskId))
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
                .then((task) => this.dispatcher.dispatch('afterAddTask', task));
        });

        tasksListElement.appendChild(addButton);

        for (let task of this.tasks) {
            tasksListElement.appendChild(task.renderTaskLine(() => this.participants, this.settings));
        }

        return tasksListElement;
    }

    async loadParticipants() {
        let participantsDTO = await this.apiClient.loadParticipants();

        if (participantsDTO === null) {
            throw new Error('Failed to load participants')
        }

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
            this.addTask(taskGroup)
        }
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
        var intervalId = setInterval(() => {
            this.dispatcher.dispatch('timerTick')
            // clearInterval(intervalId)
        }, 5000)
    }

    async reloadSettings() {
        // TODO
    }

    async reloadParticipants() {
        this.participants = [];
        await this.loadParticipants();
    }

    /**
     * Tries to assign values, received from API, to the already existing Task objects instead of creating new.
     * @returns {Promise<void>}
     */
    async reloadTasks() {

        let tasksDto = await this.apiClient.loadGroups()

        let tasksMap = new Map();

        for (let task of this.tasks) {
            tasksMap.set(task.id, task)
        }

        /**
         * @type {Task[]}
         */
        let newTasksList = [];

        for (let taskDto of tasksDto) {
            let taskId = taskDto.id;
            let task = null;
            if (tasksMap.has(taskId)) {
                task = tasksMap.get(taskId)
                task.updateFromDTO(taskDto).setDispatcher(dispatcher)
            }

            if (task === null) {
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

    handleUpdateTask(e, task, myVersionDto) {
        apiClient.updateTask(task.id, myVersionDto)
            .catch((error) => console.log('error updating task to backend', error))
            .then((result) => {
                if (result === null) {
                    alert('did not receive any result for updating task ' + task.id)
                }

                let [taskResponse, data] = result

                if (taskResponse.status === 200) {
                    task.updateFromDTO(data)
                    // for possible children
                    task.setDispatcher(this.dispatcher)
                    this.dispatcher.dispatch('taskSaved', task)
                    return
                }
                if (taskResponse.status === 409) {
                    task.updateFromDTOMerged(data, myVersionDto)
                    // for possible children
                    this.setDispatcher(this.dispatcher)
                    this.dispatcher.dispatch('taskSavedPartially', task)
                    return
                }
                if (taskResponse.status === 400) {
                    alert('Error updating task ' + task.id + ' : ' + data.error)
                    return
                }
                alert('Unhandled error ' + taskResponse.status + ' updating task ' + task.id)
                console.log('response data:', data)
            })
    }
}