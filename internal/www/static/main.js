const scheme = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
const wsUrl = new URL("websocket", scheme + window.location.host + window.location.pathname).toString()

function connectWebSocket() {
    console.log("connecting websocket to ", wsUrl)
    const w = new WebSocket(wsUrl)

    // continuously ping the server to keep the connection alive and detect dead connections
    let interval

    w.addEventListener('open', (event) => {
        console.log("ws connection established")

        interval = setInterval(() => {
            w.send('ping')
        }, 1_000)
    })

    w.addEventListener('message', (event) => {
        const newEl = new DOMParser().parseFromString(event.data, 'text/html').body.firstChild
        newEl && document.getElementById(newEl.id)?.replaceWith(newEl)
        console.log("updated", newEl.id)
    })

    w.addEventListener('close', (event) => {
        clearInterval(interval)

        console.log("websocket closed")

        setTimeout(connectWebSocket, 2000)
    })

    w.addEventListener('error', (error) => {
        console.error('WebSocket error:', error)
        w.close()
    })
}

connectWebSocket()

document.body.addEventListener("click", ev => {
    /**
     * @type {HTMLElement | null}
     */
    const el = ev.target

    if (!el) return

    fetch(`square/${el.dataset.square}`, {
        method: "PUT",
    })
})