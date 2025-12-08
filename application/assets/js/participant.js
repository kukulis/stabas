class Participant {
    id = 0;
    name = "";

    /**
     * @type {Dispatcher}
     */
    dispatcher;

    passwordHidden = true;

    /**
     * @param id {number}
     * @param name {string}
     * @param dispatcher {Dispatcher}
     */
    constructor(id, name, dispatcher) {
        this.id = id;
        this.name = name;
        this.dispatcher = dispatcher;
    }

    deleteLineCalled(event) {
        if (confirm('Are you sure you want to delete participant "' + this.name + '"?')) {
            this.dispatcher.dispatch('onDeleteParticipant', [event, this])
        }
    }

    getLineElementId() {
        return 'participant-id-' + this.id
    }

    /**
     *
     * @returns {HTMLDivElement}
     */
    renderLine() {
        let lineDiv = document.createElement('div');
        lineDiv.setAttribute('id', this.getLineElementId(this.getLineElementId()));
        lineDiv.classList.add('participant-line')

        this.renderLineInside(lineDiv)

        return lineDiv;
    }

    getHiddenFieldsDivId() {
        return 'participant-hidden-fields-' + this.id;
    }

    /**
     * @param {HTMLDivElement} lineDiv
     * @param {boolean} withEditor
     */
    renderLineInside(lineDiv, withEditor = false) {
        clearTag(lineDiv);

        let idDiv = document.createElement('div');
        idDiv.classList.add('participant-id')
        idDiv.appendChild(document.createTextNode(this.id.toString()));
        lineDiv.appendChild(idDiv)

        if (withEditor) {
            let nameInput = document.createElement('input');
            nameInput.classList.add('participant-name-input')
            nameInput.value = this.name

            nameInput.addEventListener('keyup', (e) => {
                // console.log('key code : ', e.code )
                if (e.code === "Enter") {
                    this.initiateFinishEditParticipant(e)
                }
            })

            nameInput.addEventListener('focusout', (e) => {
                this.initiateFinishEditParticipant(e)
            })

            lineDiv.appendChild(nameInput)
        } else {
            let nameDiv = document.createElement('div');
            nameDiv.classList.add('participant-name')
            nameDiv.appendChild(document.createTextNode(this.name));
            nameDiv.addEventListener('click', (e) => this.initiateEditParticipant(e))
            lineDiv.appendChild(nameDiv)
        }


        let displayFieldsButton = document.createElement('button')
        displayFieldsButton.appendChild(document.createTextNode('v'))
        displayFieldsButton.classList.add('display-participant-fields-button')
        displayFieldsButton.addEventListener('click', () => { this.toggleDisplayHiddenFields() })
        lineDiv.appendChild(displayFieldsButton)

        let hiddenFieldsDiv = document.createElement('div')
        hiddenFieldsDiv.id = this.getHiddenFieldsDivId()
        hiddenFieldsDiv.classList.add('participant-hidden-fields')
        // hiddenFieldsDiv.classList.add('hidden')

        lineDiv.appendChild(hiddenFieldsDiv)

        let deleteDiv = document.createElement('button');
        deleteDiv.classList.add('delete-button')
        deleteDiv.appendChild(document.createTextNode('✕'))
        deleteDiv.addEventListener('click', (e) => this.deleteLineCalled(e))
        lineDiv.appendChild(deleteDiv)

        let clearDiv = document.createElement('div');
        clearDiv.classList.add('clear');
        lineDiv.appendChild(clearDiv);
    }

    toggleDisplayHiddenFields() {
        let hiddenFieldsDiv = document.getElementById(this.getHiddenFieldsDivId())
        clearTag(hiddenFieldsDiv)

        let thisParticipant = this
        if ( this.passwordHidden ) {
            const loadParticipantHiddenFields = (participantDto) => {
                if ( participantDto === null ) {
                    console.log ('failed to reload participant '+thisParticipant.id )
                    return
                }
                let tokenSpan = document.createElement('span')
                tokenSpan.classList.add('participant-token')
                let token = participantDto.token
                // if ( token == '' ) token = '-'
                tokenSpan.appendChild(document.createTextNode(token))

                hiddenFieldsDiv.appendChild(tokenSpan)

                let passwordSpan = document.createElement('participant-password')
                passwordSpan.classList.add('participant-password')
                let password = participantDto.password
                // if ( password == '' ) password = '-'
                passwordSpan.appendChild(document.createTextNode(password))
                hiddenFieldsDiv.appendChild(passwordSpan)

                let regenerateButton = document.createElement('button')
                regenerateButton.appendChild(document.createTextNode('regenerate'))
                regenerateButton.classList.add('participant-regenerate-button')
                hiddenFieldsDiv.appendChild(regenerateButton)
            }
            this.dispatcher.dispatch('loadParticipantFields', [thisParticipant, loadParticipantHiddenFields])

            this.passwordHidden = false
        }
        else {
            // hiddenFieldsDiv.classList.add('hidden')
            this.passwordHidden = true
        }
    }

    initiateEditParticipant(event) {
        let lineDiv = document.getElementById(this.getLineElementId())
        this.renderLineInside(lineDiv, true)
    }

    initiateFinishEditParticipant(event) {
        this.name = event.target.value

        let dataToUpdate = {
            id: Number.parseInt(this.id),
            name: this.name,
        }

        // TODO use ApiClient.updateParticipant() through Dispatcher
        fetch('/api/participants/' + this.id, {
            method: 'PUT',
            body: JSON.stringify(dataToUpdate),
        })
            .then((response) => {
                if (response.status !== 200) {
                    console.log('updating participant api returned status ' + response.status)
                    response.text().then((text) => console.log('failed update participant response text ' + text))
                    return
                }
                response.json()
                    .then((data) =>
                        this.dispatcher.dispatch('afterUpdateParticipant', data)
                    )
                    .catch((jsonError) => console.log('updating participant json error', jsonError))
            })
            .catch((updateError) => console.log('updating participant api error', updateError))

        let lineDiv = document.getElementById(this.getLineElementId())
        this.renderLineInside(lineDiv, false)
    }
}