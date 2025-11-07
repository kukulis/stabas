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

        for (let participant of this.participants) {
            participantsDiv.appendChild(participant.renderLine());
        }

        let addParticipantButton = document.createElement('button')
        addParticipantButton.appendChild(document.createTextNode('+'))

        addParticipantButton.addEventListener('click', (e) => {
            this.participants.push((new Participant(100, 'Participant 100', this.dispatcher)))
            this.dispatcher.dispatch('afterAddParticipant', [e, this])
        })

        participantsDiv.appendChild(addParticipantButton)

        return participantsDiv;
    }

    static async initialize(dispatcher) {
        let participantsComponent = new ParticipantsComponent(dispatcher)
        participantsComponent.loadParticipants();

        return participantsComponent
    }
}