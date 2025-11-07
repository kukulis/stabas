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

    /**
     *
     * @returns {HTMLDivElement}
     */
    renderLine() {
        let lineDiv = document.createElement('div');
        lineDiv.setAttribute('id', 'participant-id-' + this.id);
        lineDiv.style['border'] = "solid thin black"

        this.renderLineInside(lineDiv)

        return lineDiv;
    }

    /**
     * @param {HTMLDivElement} lineDiv
     */
    renderLineInside(lineDiv) {
        clearTag(lineDiv);

        let deleteDiv = document.createElement('button');
        deleteDiv.appendChild(document.createTextNode('-'))
        deleteDiv.addEventListener('click', (e)=>this.deleteLineCalled(e))
        lineDiv.appendChild(deleteDiv)


        let idDiv = document.createElement('div');
        idDiv.appendChild(document.createTextNode(this.id.toString()));
        lineDiv.appendChild(idDiv)

        let nameDiv = document.createElement('div');
        nameDiv.appendChild(document.createTextNode(this.name));
        lineDiv.appendChild(nameDiv)
    }
}