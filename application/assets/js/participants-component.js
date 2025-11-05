class ParticipantsComponent {
    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @type {[Participant]}
     */
    participants = [];

    /**
     *
     * @param dispatcher {Dispatcher}
     */
    constructor(dispatcher) {
        this.dispatcher = dispatcher;
    }

    loadParticipants() {
        console.log('TODO loadParticipants')
    }

    renderParticipants() {
        let participantsDiv = document.createElement('div');
        participantsDiv.appendChild(document.createTextNode('TODO participants div'))

        return participantsDiv;
    }

    static async initialize(dispatcher) {
        let participantsComponent = new ParticipantsComponent(dispatcher)
        participantsComponent.loadParticipants();

        return participantsComponent
    }
}