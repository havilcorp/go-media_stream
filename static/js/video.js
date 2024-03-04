function turnFullscreen(fullscreen) {
	if (fullscreen) {
		if (video.requestFullScreen) video.requestFullScreen()
		else if (video.webkitRequestFullScreen) video.webkitRequestFullScreen()
		else if (video.mozRequestFullScreen) video.mozRequestFullScreen()
	} else {
		if (document.cancelFullScreen) document.cancelFullScreen()
		else if (document.mozCancelFullScreen) document.mozCancelFullScreen()
		else if (document.webkitCancelFullScreen) document.webkitCancelFullScreen()
		else if (document.msCancelFullScreen) document.msCancelFullScreen()
	}
}

function addTimeOffset(time) {
	video.pause()
	video.currentTime += time
	// audio.currentTime = video.currentTime
	timeRange.value = video.currentTime
	video.play()
	// setTimeout(() => {
	// 	audio.currentTime = video.currentTime
	// }, 1000)
}

function setTime(second) {
	video.currentTime = second
	audio.currentTime = second
	timeRange.value = second
	setTimeout(() => {
		audio.currentTime = video.currentTime
	}, 1000)
}

function playPause() {
	isPlaying = !isPlaying
	if (isPlaying) {
		btnPlay.textContent = '⏸'
		video.play()
		audio.muted = false
		audio.play()
	} else {
		btnPlay.textContent = '▶'
		video.pause()
		audio.pause()
	}
}

function addVideoSrc(src, onLoad) {
	video.addEventListener('loadeddata', onLoad, false)
	video.src = src
}

function addAudioSrc(src, onLoad) {
	audio.addEventListener('loadeddata', onLoad, false)
	audio.src = src
}
