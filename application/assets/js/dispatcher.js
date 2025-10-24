class Dispatcher {
    constructor() {
        this.listeners = new Map();
    }

    addListener(eventName, listener) {
        if (!this.listeners.has(eventName)) {
            this.listeners.set(eventName, []);
        }

        let eventListeners = this.listeners.get(eventName);
        eventListeners.push(listener);

        this.listeners.set(eventName, eventListeners);
    }

    dispatch(eventName, parameters) {
        if (!this.listeners.has(eventName)) {
            return 0;
        }

        let count = 0;
        for (let listener of this.listeners.get(eventName)) {
            if (listener(parameters)) {
                count++;
            }
        }

        return count;
    }
}