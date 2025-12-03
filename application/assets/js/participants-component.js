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
        fetch('api/participants', {
            method: 'POST',
            body: JSON.stringify(participantData)
        })
            .then((response) => {
                response.json().then((id) => {
                        let participant = new Participant(Number.parseInt(id), participantData.name, this.dispatcher)
                        this.participants.push(participant)

                        this.dispatcher.dispatch('afterAddParticipant', [e, participant])
                    }
                ).catch((error) => {
                    console.log('Error parsing response after adding new participant', error)
                })
            })
            .catch((error) => {
                console.log('Error adding participant to api', error)
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

    static async initialize(dispatcher) {
        let participantsComponent = new ParticipantsComponent(dispatcher)
        await participantsComponent.loadParticipants();

        return participantsComponent
    }

    removeParticipant(id) {
        // this.participants = this.participants.filter((participant) => participant.id !== id)

        fetch('/api/participants/' + id, {
            method: 'DELETE',
        })
            .catch((error) => {
                console.log('error deleting participant ', error)
            })
            .then((response) => {
                this.loadParticipants().then(() => {
                        this.dispatcher.dispatch('afterDeletingParticipant')
                    }
                )
            })
    }
}