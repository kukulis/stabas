class SettingsComponent {
    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     * @type {Settings}
     */
    settings;

    /**
     *
     * @param dispatcher {Dispatcher}
     * @param apiClient {ApiClient}
     */
    constructor(dispatcher, apiClient) {
        this.dispatcher = dispatcher;
        this.apiClient = apiClient;
    }

    async loadSettings() {
        let settingsDto = await this.apiClient.loadSettings()
        console.log('Loaded settings ' + settingsDto.id);
        this.settings = (new Settings(this.dispatcher)).copyFromDto(settingsDto)
        dispatcher.dispatch('afterLoadSettings')
    }

    // dispatcher?
    static async initialize(dispatcher, apiClient) {
        let settingsComponent = new SettingsComponent(dispatcher, apiClient)
        await settingsComponent.loadSettings();

        return settingsComponent
    }

    renderSettings() {
        let settingsDiv = document.createElement('div');

        console.log('render settingsComponent', this.settings)
        settingsDiv.appendChild(this.settings.renderLine());

        return settingsDiv;
    }

    updateSettings(settingsData) {
        this.apiClient.updateSettings(settingsData).then(
            (data) => {
                if (data) {
                    console.log('Settings saved successfully:', data);
                }
            });
    }

}