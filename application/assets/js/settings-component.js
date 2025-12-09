class SettingsComponent {
    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @type {Settings}
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
        // TODO move fetch to the ApiClient
        let response = await fetch('/api/settings')
        let settingsDto = await response.json()
        console.log('Loaded settings ' + settingsDto.id);
        this.settings = (new Settings()).copyFromDto(settingsDto)
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