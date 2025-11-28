class SettingsComponent {
    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @type {[Settings]}
     */
    settings;

    /**
     *
     * @param dispatcher {Dispatcher}
     */
    constructor(dispatcher) {
        this.dispatcher = dispatcher;
    }

    async loadSettings() {
        let response = await fetch('/api/settings')
        let settingsDto = await response.json()
        console.log('Loaded settings ' + settingsDto.id, settingsDto.newStatusDelay);
        this.settings;

        // for (let settingDto of settingsDto) {
        this.settings = (new Settings(
            settingsDto.id,
            settingsDto.newStatusDelay,
            settingsDto.newStatusDelaySevere,
            settingsDto.sentStatusDelay,
            settingsDto.sentStatusDelaySevere,
            settingsDto.receivedStatusDelay,
            settingsDto.receivedStatusDelaySevere,
            settingsDto.executingStatusDelay,
            settingsDto.executingStatusDelaySevere,
            settingsDto.finishedStatusDelaySevere,
            settingsDto.finishedStatusDelaySevere,
        ))
        dispatcher.dispatch('afterLoadSettings')
    }

    // dispatcher?
    static async initialize(dispatcher) {
        let settingsComponent = new SettingsComponent(dispatcher)
        await settingsComponent.loadSettings();

        return settingsComponent
    }

    renderSettings() {
        let settingsDiv = document.createElement('div');

        console.log('render settingsComponent', this.settings) 
        settingsDiv.appendChild(this.settings.renderLine());

        return settingsDiv;
    }

}