class ApiClient {
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

    // TODO load groups

    async loadTask(id){
        let response = await fetch("/api/tasks/"+id, {
            method: "GET",
        })
            .catch((error) => {
                console.log('error fetching tasks', error)
            });

        if (response === undefined) {
            console.log('loadTasks response is undefined')
            return;
        }

        return await response.json().catch((error)=> console.log('Json error', error));
    }

    /**
     *
     * @returns {Promise<Settings>}
     */
    async loadSettings() {
        // TODO load from api
        return new Settings()
    }
}