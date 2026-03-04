const scheme = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
const wsUrl = new URL("websocket", scheme + window.location.host + window.location.pathname).toString()

/**
 * handleNewMarkup takes a string of markup. It parses the html and expects a top level node with an id that is already present in the DOM.
 * It then replaces the element in the DOM with the new one.
 * @param {string} selector 
 * @param {string} markup 
 */
function handleNewMarkup(selector, markup) {
    promotionDialog.close()

    const newEl = /** @type HTMLElement */ (new DOMParser().parseFromString(markup, 'text/html').body.firstChild)

    if (!newEl) {
        console.log("could not parse DOM")
        return
    }

    console.log(newEl)

    const target = document.querySelector(selector)

    if (!target) {
        console.log("target not found")
        return
    }

    target.replaceWith(newEl)
    console.log("updated", selector)
}

const s = new EventSource("/events")

/**
 * @typedef {{selector: string, markup: string}} MarkupEvent
 */

s.addEventListener("markup", function (ev) {
    const data = /** @type {string}*/ (ev.data)

    const payload = /** @type {MarkupEvent} */ (JSON.parse(data))

    handleNewMarkup(payload.selector, payload.markup)
})

document.body.addEventListener("click", function (ev) {
    const target = /** @type {HTMLElement} */ (ev.target)

    if (target.classList.contains("tile")) {
        handleSquareClick(target)
        return
    }
})

const promotionDialog = /** @type {HTMLDialogElement}*/ (document.getElementById("promotion"))

/**
 * 
 * @param {HTMLElement} tile
 * @returns 
 */
function handleSquareClick(tile) {
    if (tile.classList.contains("promotion")) {
        handlePromotion(tile.dataset.square, tile.dataset.move.split(","))
        return
    }

    if (tile.dataset.move) {
        fetch(`move/${tile.dataset.move}`, {
            method: "PUT",
        })
    } else {
        // (de-) select square
        fetch(`square/${tile.dataset.square}`, {
            method: "PUT",
        })
    }
}

/**
 * 
 * @param {string} square 
 * @param {string[]} moves 
 */
function handlePromotion(square, moves) {
    promotionDialog.showModal()

    promotionDialog.addEventListener("close", ev => {
        console.log(promotionDialog.returnValue)

        if (!promotionDialog.returnValue) return

        const selectedMove = moves.find(m => m.endsWith(promotionDialog.returnValue))

        fetch(`move/${selectedMove}`, {
            method: "PUT",
        })
    }, {
        once: true
    })
}

document.getElementById("closeModal").addEventListener("click", ev => {
    promotionDialog.close()
})
