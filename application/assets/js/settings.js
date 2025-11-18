const LATE_NONE = 'none'
const LATE_SOFT = 'soft'
const LATE_SEVERE = 'severe'

class Settings {
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

    reloadTasksOnTimer = true;

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

    static solveCriticality ( delay, softLate, severeLate ) {
        if ( delay < softLate) {
            return LATE_NONE;
        }
        if ( delay < severeLate ) {
            return LATE_SOFT
        }

        return LATE_SEVERE
    }
}