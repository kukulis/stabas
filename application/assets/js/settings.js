const LATE_NONE = 'none'
const LATE_SOFT = 'soft'
const LATE_SEVERE = 'severe'

class Settings {
    id = 0;
    newStatusDelay = 5;
    newStatusDelaySevere = 15;
    sentStatusDelay = 2;
    sentStatusDelaySevere = 6;
    receivedStatusDelay = 5;
    receivedStatusDelaySevere = 15;
    executingStatusDelay = 10;
    executingStatusDelaySevere = 20;
    finishedStatusDelay = 60;
    finishedStatusDelaySevere = 120;

    reloadTasksOnTimer = false;

    // other settings values

    /**
     * @param id {number}
     * @param newStatusDelay {number}
     */
    constructor(id, newStatusDelay, newStatusDelaySevere) {
        this.id = id;
        this.newStatusDelay = newStatusDelay;
        this.newStatusDelaySevere = newStatusDelaySevere;
    }

    calculateCriticality(delay, status) {
        let delayMinutes = delay / (1000 * 60)
        switch (status) {
            case STATUS_NEW:
                return Settings.solveCriticality(delayMinutes, this.newStatusDelay, this.newStatusDelaySevere)
            case STATUS_SENT:
                return Settings.solveCriticality(delayMinutes, this.sentStatusDelay, this.sentStatusDelaySevere)
            case STATUS_RECEIVED:
                return Settings.solveCriticality(delayMinutes, this.receivedStatusDelay, this.receivedStatusDelaySevere)
            case STATUS_EXECUTING:
                return Settings.solveCriticality(delayMinutes, this.executingStatusDelay, this.executingStatusDelaySevere)
            case STATUS_FINISHED:
                return Settings.solveCriticality(delayMinutes, this.finishedStatusDelay, this.finishedStatusDelaySevere)
        }

        return LATE_NONE
    }

    static solveCriticality(delay, softLate, severeLate) {
        if (delay < softLate) {
            return LATE_NONE;
        }
        if (delay < severeLate) {
            return LATE_SOFT
        }

        return LATE_SEVERE
    }
    /**
     *
     * @returns {HTMLDivElement}
     */
    renderLine() {
        console.log("render settings line :)" + this)
        let lineDiv = document.createElement('div');
        lineDiv.setAttribute('id', "settings-div");
        lineDiv.style['border'] = "solid thin black"

        clearTag(lineDiv);

        let idDiv = document.createElement('div');
        idDiv.appendChild(document.createTextNode(this.id.toString()));
        lineDiv.appendChild(idDiv)

        let newStatusDelayInput = document.createElement('input');
        newStatusDelayInput.value = this.newStatusDelay

        lineDiv.appendChild(newStatusDelayInput)


        return lineDiv;
    }



}