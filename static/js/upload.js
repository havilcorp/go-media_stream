function createUploadCard(i, filmName, onClickDelete) {
	let div = document.createElement('div')
	div.classList.add('upload_file')
	div.classList.add('file_item')

	let btnDel = document.createElement('span')
	btnDel.textContent = '✖️'
	btnDel.classList.add('btn_del')
	btnDel.addEventListener('click', function (e) {
		onClickDelete(e)
		// const index = films.findIndex(v => v.i == i)
		// films.splice(index, 1)
		// div.remove()
		// console.log(films)
		// fileCount.textContent = films.length
		// if (films.length == 0) {
		// 	loadImageBlock.style.display = 'flex'
		// 	imageLoadedBlock.style.display = 'none'
		// 	uploadFiles.value = ''
		// }
	})
	div.appendChild(btnDel)

	const image = new Image()
	image.classList.add('image')
	image.src = ''

	let inputImageFile = document.createElement('input')
	inputImageFile.classList.add('upload_file_img')
	inputImageFile.type = 'file'
	inputImageFile.style.display = 'none'
	image.appendChild(inputImageFile)
	image.onclick = function () {
		inputImageFile.click()
	}

	div.appendChild(image)

	inputImageFile.addEventListener('change', function (e) {
		image.src = window.URL.createObjectURL(e.target.files[0])
	})

	let name = document.createElement('input')
	name.id = `file_item_name_${i}`
	name.type = 'text'
	name.value = filmName
	div.appendChild(name)

	let progress = document.createElement('span')
	progress.id = `file_item_progress_${i}`
	div.appendChild(progress)

	return {
		div: div,
		imgFile: inputImageFile,
	}
}
