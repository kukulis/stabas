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
        lineDiv.classList.add('settings-div')

        clearTag(lineDiv);

        let settingsTitleDiv = document.createElement('div');
        settingsTitleDiv.appendChild(document.createTextNode("Settings"));
        lineDiv.appendChild(settingsTitleDiv);

        

        let dl = document.createElement('dl');

        let dtNewStatusDelay = document.createElement('dt');
        dtNewStatusDelay.appendChild(document.createTextNode("New status delay"));
        dl.appendChild(dtNewStatusDelay);

        let ddNewStatusDelay = document.createElement('dd');
        let newStatusDelayInput = document.createElement('input');
        newStatusDelayInput.id = 'newStatusDelayInput';
        newStatusDelayInput.value = this.newStatusDelay.toString();
        ddNewStatusDelay.appendChild(newStatusDelayInput);
        dl.appendChild(ddNewStatusDelay);

        let dtNewStatusDelaySevere = document.createElement('dt');
        dtNewStatusDelaySevere.appendChild(document.createTextNode("New status delay severe"));
        dl.appendChild(dtNewStatusDelaySevere);

        let ddNewStatusDelaySevere = document.createElement('dd');
        let newStatusDelaySevereInput = document.createElement('input');
        newStatusDelaySevereInput.id = 'newStatusDelaySevereInput';
        newStatusDelaySevereInput.value = this.newStatusDelaySevere.toString();
        ddNewStatusDelaySevere.appendChild(newStatusDelaySevereInput);
        dl.appendChild(ddNewStatusDelaySevere);

        let dtSentStatusDelay = document.createElement('dt');
        dtSentStatusDelay.appendChild(document.createTextNode("Sent status delay"));
        dl.appendChild(dtSentStatusDelay);

        let ddSentStatusDelay = document.createElement('dd');
        let sentStatusDelayInput = document.createElement('input');
        sentStatusDelayInput.id = 'sentStatusDelayInput';
        sentStatusDelayInput.value = this.sentStatusDelay.toString();
        ddSentStatusDelay.appendChild(sentStatusDelayInput);
        dl.appendChild(ddSentStatusDelay);

        let dtSentStatusDelaySevere = document.createElement('dt');
        dtSentStatusDelaySevere.appendChild(document.createTextNode("Sent status delay severe"));
        dl.appendChild(dtSentStatusDelaySevere);

        let ddSentStatusDelaySevere = document.createElement('dd');
        let sentStatusDelaySevereInput = document.createElement('input');
        sentStatusDelaySevereInput.id = 'sentStatusDelaySevereInput';
        sentStatusDelaySevereInput.value = this.sentStatusDelaySevere.toString();
        ddSentStatusDelaySevere.appendChild(sentStatusDelaySevereInput);
        dl.appendChild(ddSentStatusDelaySevere);

        let dtReceivedStatusDelay = document.createElement('dt');
        dtReceivedStatusDelay.appendChild(document.createTextNode("Received status delay"));
        dl.appendChild(dtReceivedStatusDelay);

        let ddReceivedStatusDelay = document.createElement('dd');
        let receivedStatusDelayInput = document.createElement('input');
        receivedStatusDelayInput.id = 'receivedStatusDelayInput';
        receivedStatusDelayInput.value = this.receivedStatusDelay.toString();
        ddReceivedStatusDelay.appendChild(receivedStatusDelayInput);
        dl.appendChild(ddReceivedStatusDelay);

        let dtReceivedStatusDelaySevere = document.createElement('dt');
        dtReceivedStatusDelaySevere.appendChild(document.createTextNode("Received status delay severe"));
        dl.appendChild(dtReceivedStatusDelaySevere);

        let ddReceivedStatusDelaySevere = document.createElement('dd');
        let receivedStatusDelaySevereInput = document.createElement('input');
        receivedStatusDelaySevereInput.id = 'receivedStatusDelaySevereInput';
        receivedStatusDelaySevereInput.value = this.receivedStatusDelaySevere.toString();
        ddReceivedStatusDelaySevere.appendChild(receivedStatusDelaySevereInput);
        dl.appendChild(ddReceivedStatusDelaySevere);

        let dtExecutingStatusDelay = document.createElement('dt');
        dtExecutingStatusDelay.appendChild(document.createTextNode("Executing status delay"));
        dl.appendChild(dtExecutingStatusDelay);

        let ddExecutingStatusDelay = document.createElement('dd');
        let executingStatusDelayInput = document.createElement('input');
        executingStatusDelayInput.id = 'executingStatusDelayInput';
        executingStatusDelayInput.value = this.executingStatusDelay.toString();
        ddExecutingStatusDelay.appendChild(executingStatusDelayInput);
        dl.appendChild(ddExecutingStatusDelay);

        let dtExecutingStatusDelaySevere = document.createElement('dt');
        dtExecutingStatusDelaySevere.appendChild(document.createTextNode("Executing status delay severe"));
        dl.appendChild(dtExecutingStatusDelaySevere);

        let ddExecutingStatusDelaySevere = document.createElement('dd');
        let executingStatusDelaySevereInput = document.createElement('input');
        executingStatusDelaySevereInput.id = 'executingStatusDelaySevereInput';
        executingStatusDelaySevereInput.value = this.executingStatusDelaySevere.toString();
        ddExecutingStatusDelaySevere.appendChild(executingStatusDelaySevereInput);
        dl.appendChild(ddExecutingStatusDelaySevere);

        let dtFinishedStatusDelay = document.createElement('dt');
        dtFinishedStatusDelay.appendChild(document.createTextNode("Finished status delay"));
        dl.appendChild(dtFinishedStatusDelay);

        let ddFinishedStatusDelay = document.createElement('dd');
        let finishedStatusDelayInput = document.createElement('input');
        finishedStatusDelayInput.id = 'finishedStatusDelayInput';
        finishedStatusDelayInput.value = this.finishedStatusDelay.toString();
        ddFinishedStatusDelay.appendChild(finishedStatusDelayInput);
        dl.appendChild(ddFinishedStatusDelay);

        let dtFinishedStatusDelaySevere = document.createElement('dt');
        dtFinishedStatusDelaySevere.appendChild(document.createTextNode("Finished status delay severe"));
        dl.appendChild(dtFinishedStatusDelaySevere);

        let ddFinishedStatusDelaySevere = document.createElement('dd');
        let finishedStatusDelaySevereInput = document.createElement('input');
        finishedStatusDelaySevereInput.id = 'finishedStatusDelaySevereInput';
        finishedStatusDelaySevereInput.value = this.finishedStatusDelaySevere.toString();
        ddFinishedStatusDelaySevere.appendChild(finishedStatusDelaySevereInput);
        dl.appendChild(ddFinishedStatusDelaySevere);

        let dtReloadTasksOnTimer = document.createElement('dt');
        dtReloadTasksOnTimer.appendChild(document.createTextNode("Reload tasks on timer"));
        dl.appendChild(dtReloadTasksOnTimer);

        let ddReloadTasksOnTimer = document.createElement('dd');
        let reloadTasksOnTimerInput = document.createElement('input');
        reloadTasksOnTimerInput.id = 'reloadTasksOnTimerInput';
        reloadTasksOnTimerInput.type = 'checkbox';
        reloadTasksOnTimerInput.checked = this.reloadTasksOnTimer;
        ddReloadTasksOnTimer.appendChild(reloadTasksOnTimerInput);
        dl.appendChild(ddReloadTasksOnTimer);

        lineDiv.appendChild(dl);

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