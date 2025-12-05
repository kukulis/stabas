class ApiClient {

    /**
     * loadGroups is used lately.
     *
     * This endpoint may be useful still.
     *
     * @returns {Promise<any>}
     */
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

        return await response.json();
    }


    /**
     * Groups are same tasks but with children
     * @returns {Promise<TaskGroup[]>} It is actually not taskGroups, but has same fields
     */
    async loadGroups() {
        let response = await fetch("/api/groups", {
            method: "GET",
        })
            .catch((error) => {
                console.log('error fetching groups', error)
            });

        if (response === undefined) {
            console.log('loadGroups response is undefined')
            return null;
        }

        return await response.json();
    }

    async loadTask(id) {
        let response = await fetch("/api/tasks/" + id, {
            method: "GET",
        })
            .catch((error) => {
                console.log('error fetching tasks', error)
            });

        if (response === undefined) {
            console.log('loadTasks response is undefined')
            return;
        }

        return await response.json().catch((error) => console.log('Json error', error));
    }

    /**
     *
     * @returns {Promise<Settings>}
     */
    async loadSettings() {
        // TODO load from api
        return new Settings()
    }

    async deleteTask(taskId) {
        await fetch('/api/tasks/' + taskId, {method: 'DELETE'}).catch((error) => {
            console.log('error deleting task ' + taskId, error)
            return false
        })

        return true
    }

    /**
     * Create a new task
     * @param {Object} taskData - Task data to create
     * @returns {Promise<any>}
     */
    async createTask(taskData) {
        let response = await fetch('/api/tasks', {
            method: 'POST',
            body: JSON.stringify(taskData)
        }).catch((error) => {
            console.log('backend error creating task', error)
            return undefined
        })

        if (response === undefined) {
            return null
        }

        return await response.json()
    }

    /**
     * Update an existing task
     * @param {number} taskId - Task ID to update
     * @param {Object} taskData - Updated task data
     * @returns {Promise<[Response, data]>}
     */
    async updateTask(taskId, taskData) {
        let response = await fetch('/api/tasks/' + taskId, {
            method: 'PUT',
            body: JSON.stringify(taskData)
        }).catch((error) => {
            console.log('error updating task to backend', error)
            return undefined
        })

        if (response === undefined) {
            return null
        }

        let data = await response.json()
        return [ response, data ]
    }

    /**
     * Change task status
     * @param {number} taskId - Task ID
     * @param {number} newStatus - New status value
     * @returns {Promise<Response>}
     */
    async changeTaskStatus(taskId, newStatus) {
        return await fetch('/api/tasks/' + taskId + '/change-status?status=' + newStatus, {
            method: 'POST',
        }).catch((error) => {
            console.log('error changing status', error)
            return undefined
        })
    }

    /**
     * Load all participants
     * @returns {Promise<any>}
     */
    async loadParticipants() {
        let response = await fetch('/api/participants', {
            method: 'GET',
        }).catch((error) => {
            console.log('error fetching participants', error)
            return undefined
        })

        if (response === undefined) {
            console.log('loadParticipants response is undefined')
            return null
        }

        return await response.json()
    }

    /**
     * Create a new participant
     * @param {Object} participantData - Participant data to create
     * @returns {Promise<{response: Response, data: any}>}
     */
    async createParticipant(participantData) {
        let response = await fetch('api/participants', {
            method: 'POST',
            body: JSON.stringify(participantData)
        }).catch((error) => {
            console.log('Error adding participant to api', error)
            return undefined
        })

        if (response === undefined) {
            return null
        }

        let data = await response.json()
        return { response, data }
    }

    /**
     * Update an existing participant
     * @param {number} participantId - Participant ID to update
     * @param {Object} participantData - Updated participant data
     * @returns {Promise<Response>}
     */
    async updateParticipant(participantId, participantData) {
        return await fetch('/api/participants/' + participantId, {
            method: 'PUT',
            body: JSON.stringify(participantData),
        }).catch((error) => {
            console.log('updating participant api error', error)
            return undefined
        })
    }

    /**
     * Delete a participant
     * @param {number} participantId - Participant ID to delete
     * @returns {Promise<Response>}
     */
    async deleteParticipant(participantId) {
        return await fetch('/api/participants/' + participantId, {
            method: 'DELETE',
        }).catch((error) => {
            console.log('error deleting participant', error)
            return undefined
        })
    }
}