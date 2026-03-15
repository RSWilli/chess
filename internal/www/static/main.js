const s = new EventSource("/events")

s.addEventListener("change", function (ev) {
    window.location.reload()
})
