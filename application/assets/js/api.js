class ApiClient {

    /**
     * Get headers with auth token from localStorage
     * @returns {Object}
     */
    getHeaders() {
        const authToken = localStorage.getItem('authToken');
        return {
            'auth_token': authToken || ''
        };
    }

    /**
     * Handle response and check for 401 unauthorized
     * @param {Response|undefined} response - Fetch response
     * @returns {Response|null} - Returns response if valid, null if 401 or undefined
     */
    handleResponse(response) {
        if (response === undefined) {
            return null;
        }

        if (response.status === 401) {
            window.location.href = '/login';
            return null;
        }

        return response;
    }

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
            headers: this.getHeaders(),
        })
            .catch((error) => {
                console.log('error fetching tasks', error)
            });

        response = this.handleResponse(response);
        if (!response) {
            console.log('loadTasks response is invalid')
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
            headers: this.getHeaders(),
        })
            .catch((error) => {
                console.log('error fetching groups', error)
            });

        response = this.handleResponse(response);
        if (!response) {
            console.log('loadGroups response is invalid')
            return null;
        }

        return await response.json();
    }

    async loadTask(id) {
        let response = await fetch("/api/tasks/" + id, {
            method: "GET",
            headers: this.getHeaders(),
        })
            .catch((error) => {
                console.log('error fetching tasks', error)
            });

        response = this.handleResponse(response);
        if (!response) {
            console.log('loadTask response is invalid')
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
        let response = await fetch('/api/tasks/' + taskId, {
            method: 'DELETE',
            headers: this.getHeaders(),
        }).catch((error) => {
            console.log('error deleting task ' + taskId, error)
        })

        return this.handleResponse(response);
    }

    /**
     * Create a new task
     * @param {Object} taskData - Task data to create
     * @returns {Promise<any>}
     */
    async createTask(taskData) {
        let response = await fetch('/api/tasks', {
            method: 'POST',
            headers: {
                ...this.getHeaders(),
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(taskData)
        }).catch((error) => {
            console.log('backend error creating task', error)
        })

        response = this.handleResponse(response);
        if (!response) {
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
            headers: {
                ...this.getHeaders(),
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(taskData)
        }).catch((error) => {
            console.log('error updating task to backend', error)
        })

        response = this.handleResponse(response);
        if (!response) {
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
        let response = await fetch('/api/tasks/' + taskId + '/change-status?status=' + newStatus, {
            method: 'POST',
            headers: this.getHeaders(),
        }).catch((error) => {
            console.log('error changing status', error)
        })

        return this.handleResponse(response);
    }

    /**
     * Load all participants
     * @returns {Promise<any>}
     */
    async loadParticipants() {
        let response = await fetch('/api/participants', {
            method: 'GET',
            headers: this.getHeaders(),
        }).catch((error) => {
            console.log('error fetching participants', error)
        })

        response = this.handleResponse(response);
        if (!response) {
            console.log('loadParticipants response is invalid')
            return null
        }

        return await response.json()
    }

    /**
     * Load a single participant by ID
     * @param {number} id - Participant ID
     * @returns {Promise<any>}
     */
    async loadParticipant(id) {
        let response = await fetch('/api/participants/' + id, {
            method: 'GET',
            headers: this.getHeaders(),
        }).catch((error) => {
            console.log('error fetching participant', error)
        })

        response = this.handleResponse(response);
        if (!response) {
            console.log('loadParticipant response is invalid')
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
            headers: {
                ...this.getHeaders(),
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(participantData)
        }).catch((error) => {
            console.log('Error adding participant to api', error)
        })

        response = this.handleResponse(response);
        if (!response) {
            return null
        }

        let data = await response.json()
        return { response, data }
    }

    /**
     * Update an existing participant
     * @param {number} participantId - Participant ID to update
     * @param {Object} participantData - Updated participant data
     * @returns {Promise<any>}
     */
    async updateParticipant(participantId, participantData) {
        let response = await fetch('/api/participants/' + participantId, {
            method: 'PUT',
            headers: {
                ...this.getHeaders(),
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(participantData),
        }).catch((error) => {
            console.log('updating participant api error', error)
        })

        response = this.handleResponse(response);
        if (!response) {
            return null
        }

        return await response.json()
    }

    /**
     * Delete a participant
     * @param {number} participantId - Participant ID to delete
     * @returns {Promise<Response>}
     */
    async deleteParticipant(participantId) {
        let response = await fetch('/api/participants/' + participantId, {
            method: 'DELETE',
            headers: this.getHeaders(),
        }).catch((error) => {
            console.log('error deleting participant', error)
        })

        return this.handleResponse(response);
    }

    /**
     * Regenerate password for a participant
     * @param {number} participantId - Participant ID
     * @returns {Promise<any>}
     */
    async regenerateParticipantPassword(participantId) {
        let response = await fetch('/api/participants/' + participantId + '/regenerate-password', {
            method: 'POST',
            headers: this.getHeaders(),
        }).catch((error) => {
            console.log('error regenerating password', error)
        })

        response = this.handleResponse(response);
        if (!response) {
            return null
        }

        return await response.json()
    }
}