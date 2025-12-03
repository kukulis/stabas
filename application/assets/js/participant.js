class Participant {
    id = 0;
    name = "";

    /**
     * @type {Dispatcher}
     */
    dispatcher;


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
        this.dispatcher.dispatch('onDeleteParticipant', [event, this])
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
        lineDiv.style['border'] = "solid thin black"

        this.renderLineInside(lineDiv)

        return lineDiv;
    }

    /**
     * @param {HTMLDivElement} lineDiv
     * @param {boolean} withEditor
     */
    renderLineInside(lineDiv, withEditor = false) {
        clearTag(lineDiv);

        let deleteDiv = document.createElement('button');
        deleteDiv.appendChild(document.createTextNode('-'))
        deleteDiv.addEventListener('click', (e) => this.deleteLineCalled(e))
        lineDiv.appendChild(deleteDiv)

        let idDiv = document.createElement('div');
        idDiv.appendChild(document.createTextNode(this.id.toString()));
        lineDiv.appendChild(idDiv)

        if (withEditor) {
            let nameInput = document.createElement('input');
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
            nameDiv.appendChild(document.createTextNode(this.name));
            nameDiv.addEventListener('click', (e) => this.initiateEditParticipant(e))
            lineDiv.appendChild(nameDiv)
        }
    }

    initiateEditParticipant(event) {
        let lineDiv = document.getElementById(this.getLineElementId())
        this.renderLineInside(lineDiv, true)
    }

    initiateFinishEditParticipant(event) {
        this.name = event.target.value

        let dataToUpdate = {
            id: Number.parseInt( this.id ),
            name: this.name,
        }

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