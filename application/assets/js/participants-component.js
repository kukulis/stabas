class ParticipantsComponent {
    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     * @type {[Participant]}
     */
    participants = [];

    /**
     *
     * @param dispatcher {Dispatcher}
     * @param apiClient {ApiClient}
     */
    constructor(dispatcher, apiClient) {
        this.dispatcher = dispatcher;
        this.apiClient = apiClient;
    }

    async loadParticipants() {
        let participantsDto = await this.apiClient.loadParticipants()

        if (participantsDto === null) {
            throw new Error('Failed to load participants')
        }

        // console.log('Loaded participants ' + participantsDto.length);
        this.participants = [];

        for (let participantDto of participantsDto) {
            this.participants.push(new Participant(participantDto.id, participantDto.name, dispatcher))
            dispatcher.dispatch('afterLoadParticipants')
        }
    }

    addParticipantCalled(e) {

        let maxId = this.participants.reduce(
            (previousId, participant) => Math.max(previousId, participant.id), 0
        );

        let newId = maxId + 1;

        let participantData = {
            id: 0,
            name: "Participant " + newId
        }

        // till here
        this.apiClient.createParticipant(participantData)
            .then((result) => {
                if (result === null) {
                    throw new Error('Failed to create participant')
                }

                let { response, data } = result
                let id = data
                let participant = new Participant(Number.parseInt(id), participantData.name, this.dispatcher)
                this.participants.push(participant)

                this.dispatcher.dispatch('afterAddParticipant', [e, participant])
            })
            .catch((error) => {
                console.log('Error adding participant', error)
            })

    }

    renderParticipants() {
        let participantsDiv = document.createElement('div');

        for (let participant of this.participants) {
            participantsDiv.appendChild(participant.renderLine());
        }

        let addParticipantButton = document.createElement('button')
        addParticipantButton.appendChild(document.createTextNode('+'))
        addParticipantButton.addEventListener('click', (e) => this.addParticipantCalled(e))
        participantsDiv.appendChild(addParticipantButton)

        return participantsDiv;
    }

    static async initialize(dispatcher, apiClient) {
        let participantsComponent = new ParticipantsComponent(dispatcher, apiClient)
        await participantsComponent.loadParticipants();

        return participantsComponent
    }

    removeParticipant(id) {
        // this.participants = this.participants.filter((participant) => participant.id !== id)

        this.apiClient.deleteParticipant(id)
            .then((response) => {
                if (response === undefined) {
                    throw new Error('Failed to delete participant')
                }
                return this.loadParticipants()
            })
            .then(() => {
                this.dispatcher.dispatch('afterDeletingParticipant')
            })
    }
}