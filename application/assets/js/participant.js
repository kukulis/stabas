class Participant {
    id = 0;
    name = "";

    /**
     * @type {Dispatcher}
     */
    dispatcher = null;


    /**
     * @param id {number}
     * @param name {string}
     * @param dispatcher {Dispatcher}
     */
    constructor(id, name, dispatcher) {
        this.id = id;
        this.name = name;
    }

    /**
     *
     * @returns {HTMLDivElement}
     */
    renderLine() {
        let lineDiv = document.createElement('div');
        lineDiv.setAttribute('id', 'participant-id-'+this.id );

        this.renderLineInside(lineDiv)

        return lineDiv;
    }

    /**
     * @param {HTMLDivElement} lineDiv
     */
    renderLineInside(lineDiv ) {
        clearTag(lineDiv);

        let nameDiv = document.createElement('div');
        nameDiv.appendChild(document.createTextNode(this.name ));

        lineDiv.appendChild(nameDiv)
    }
}