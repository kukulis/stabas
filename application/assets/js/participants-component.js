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

    async loadParticipants() {
        let response = await fetch('/api/participants')
        let participantsDto = await response.json()
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

        let newParticipant = new Participant(newId, 'Participant ' + newId, this.dispatcher);
        this.participants.push(newParticipant)
        this.dispatcher.dispatch('afterAddParticipant', [e, newParticipant])
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

    static
    async initialize(dispatcher) {
        let participantsComponent = new ParticipantsComponent(dispatcher)
        await participantsComponent.loadParticipants();

        return participantsComponent
    }

    removeParticipant(id) {
        this.participants = this.participants.filter((participant) => participant.id !== id)
    }
}