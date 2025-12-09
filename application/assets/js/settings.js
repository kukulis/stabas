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
        lineDiv.style['border'] = "solid thin black";

        clearTag(lineDiv);

        let idDiv = document.createElement('div');
        idDiv.appendChild(document.createTextNode("Nustatymai"));
        lineDiv.appendChild(idDiv);

        lineDiv.appendChild(document.createElement("br"))

        lineDiv.appendChild(document.createTextNode("newStatusDelay"));
        let newStatusDelayInput = document.createElement('input');
        newStatusDelayInput.value = this.newStatusDelay.toString();
        lineDiv.appendChild(newStatusDelayInput)

        lineDiv.appendChild(document.createElement("br"))

        lineDiv.appendChild(document.createTextNode("newStatusDelaySevere"));
        let newStatusDelaySevereInput = document.createElement('input');
        newStatusDelaySevereInput.value = this.newStatusDelaySevere.toString();
        lineDiv.appendChild(newStatusDelaySevereInput)

        lineDiv.appendChild(document.createElement("br"))

        lineDiv.appendChild(document.createTextNode("sentStatusDelay"));
        let sentStatusDelayInput = document.createElement('input');
        sentStatusDelayInput.value = this.sentStatusDelay.toString();
        lineDiv.appendChild(sentStatusDelayInput)

        lineDiv.appendChild(document.createElement("br"))

        lineDiv.appendChild(document.createTextNode("sentStatusDelaySevere"));
        let sentStatusDelaySevereInput = document.createElement('input');
        sentStatusDelaySevereInput.value = this.sentStatusDelaySevere.toString();
        lineDiv.appendChild(sentStatusDelaySevereInput)

        let saveButton = document.createElement('button')
        saveButton.appendChild(document.createTextNode('save settings'))
        lineDiv.appendChild(saveButton)

        return lineDiv;
    }

    copyFromDto(settingsDto) {
        this.id = settingsDto.id;
        this.newStatusDelay = settingsDto.newStatusDelay;
        this.newStatusDelaySevere = settingsDto.newStatusDelaySevere;
        this.sentStatusDelay = settingsDto.sentStatusDelay;
        this.sentStatusDelaySevere = settingsDto.sentStatusDelaySevere;
        this.receivedStatusDelay = settingsDto.receivedStatusDelay;
        this.receivedStatusDelaySevere = settingsDto.receivedStatusDelaySevere;
        this.executingStatusDelay = settingsDto.executingStatusDelay;
        this.executingStatusDelaySevere = settingsDto.executingStatusDelaySevere;
        this.finishedStatusDelay = settingsDto.finishedStatusDelay;
        this.finishedStatusDelaySevere = settingsDto.finishedStatusDelaySevere;
        this.reloadTasksOnTimer = settingsDto.reloadTasksOnTimer;

        return this;
    }
}