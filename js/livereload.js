function newConnection(ws) {
	let loc = window.location
	let url = (loc.protocol === 'https:') ? 'wss:': 'ws:'
	url += `//${loc.host}/__livereload__`

	ws = new WebSocket(url)
	let intervalHandle = 0

	ws.onopen = function() {
		console.log('livereload: connected')
	}

	ws.onmessage = function(evt) {
		console.log(evt.data)
		if (evt.data === 'reload') {
			window.location.reload()
		}
	}

	ws.onclose = function(event) {
		console.log('livereload: connection died')

		clearInterval(intervalHandle)
		ws.close()
		setTimeout(function() {
			newConnection(ws)
		}, 1000)
	}

	ws.onerror = function(error) {
		console.log(`livereload: connection error`)
	}

	intervalHandle = setInterval(function() {
		ws.send('livereload: ping')
	}, 20000)

	return ws
}

let ws = null
newConnection(ws)
