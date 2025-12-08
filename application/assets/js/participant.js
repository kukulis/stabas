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

    /**
     * Handles delete button click, shows confirmation and dispatches delete event
     * @param {Event} event - Click event
     */
    deleteLineCalled(event) {
        if (confirm('Are you sure you want to delete participant "' + this.name + '"?')) {
            this.dispatcher.dispatch('onDeleteParticipant', [event, this])
        }
    }

    /**
     * Returns the DOM element ID for this participant's line
     * @returns {string}
     */
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

    /**
     * Returns the DOM element ID for this participant's hidden fields container
     * @returns {string}
     */
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
        displayFieldsButton.appendChild(document.createTextNode('👁'))
        displayFieldsButton.classList.add('display-participant-fields-button')
        displayFieldsButton.addEventListener('click', () => {
            this.toggleDisplayHiddenFields()
        })
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

    /**
     * Loads and displays the participant's token and password in the hidden fields area
     * @param {Object} participantDto - Participant data with token and password
     */
    loadParticipantHiddenFields(participantDto) {
        let hiddenFieldsDiv = document.getElementById(this.getHiddenFieldsDivId())
        clearTag(hiddenFieldsDiv)

        if (participantDto === null) {
            console.log('failed to reload participant password and token' + this.id)
            return
        }
        let tokenSpan = document.createElement('span')
        tokenSpan.classList.add('participant-token')
        let token = participantDto.token
        tokenSpan.appendChild(document.createTextNode(token))

        hiddenFieldsDiv.appendChild(tokenSpan)

        let passwordSpan = document.createElement('participant-password')
        passwordSpan.classList.add('participant-password')
        let password = participantDto.password
        passwordSpan.appendChild(document.createTextNode(password))
        hiddenFieldsDiv.appendChild(passwordSpan)

        let regenerateButton = document.createElement('button')
        regenerateButton.appendChild(document.createTextNode('regenerate'))
        regenerateButton.classList.add('participant-regenerate-button')
        regenerateButton.addEventListener('click', () => {
            this.handleRegenerate()
        })
        hiddenFieldsDiv.appendChild(regenerateButton)
    }

    /**
     * Toggles the visibility of hidden fields (token and password)
     * Dispatches event to load participant fields when showing
     */
    toggleDisplayHiddenFields() {
        let hiddenFieldsDiv = document.getElementById(this.getHiddenFieldsDivId())
        clearTag(hiddenFieldsDiv)

        let thisParticipant = this
        if (this.passwordHidden) {
            this.dispatcher.dispatch('loadParticipantFields',
                [thisParticipant, (participantDto) => this.loadParticipantHiddenFields(participantDto)]
            )

            this.passwordHidden = false
        } else {
            // hiddenFieldsDiv.classList.add('hidden')
            this.passwordHidden = true
        }
    }

    /**
     * Handles password regeneration request, shows confirmation and dispatches regenerate event
     */
    handleRegenerate() {
        if (!confirm('Do you want to regenerate password for a participant ' + this.id + '?')) {
            return
        }

        this.dispatcher.dispatch('regenerateParticipantPassword',
            [this, (participantDto) => this.loadParticipantHiddenFields(participantDto)]
        );
    }

    /**
     * Initiates edit mode for the participant name, renders line with input field
     * @param {Event} event - Click event
     */
    initiateEditParticipant(event) {
        let lineDiv = document.getElementById(this.getLineElementId())
        this.renderLineInside(lineDiv, true)
    }

    /**
     * Finishes editing participant name, saves changes via API and re-renders line
     * @param {Event} event - Input blur or keyup event
     */
    initiateFinishEditParticipant(event) {
        this.name = event.target.value

        let dataToUpdate = {
            id: Number.parseInt(this.id),
            name: this.name,
        }

        let thisParticipant = this

        const callBackAfterUpdate = () => {
            let lineDiv = document.getElementById(thisParticipant.getLineElementId())
            thisParticipant.renderLineInside(lineDiv, false)
        }

        this.dispatcher.dispatch('onUpdateParticipant', [
            event,
            thisParticipant,
            dataToUpdate,
            callBackAfterUpdate
        ])
    }
}