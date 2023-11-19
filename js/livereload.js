function newConnection() {
	let url = `/__livereload__`

	let ws = new EventSource(url)
	let intervalHandle = 0

	ws.onopen = function() {
		console.log('[livereload] connected')
	}

	ws.onclose = function(event) {
		console.log('[livereload] connection died')
	}

	ws.onerror = function(ev, err) {
		console.log(`[livereload] error`, err)
	}

	return ws
}

let ws = newConnection()
ws.addEventListener("message", (evt) => {
	let msg = evt.data
	console.log("message:", msg)
	if (msg === 'reload') {
		ws.close()
		window.location.reload()
		console.log("window.reload")
	}
})