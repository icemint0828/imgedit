// see:http://bashalog.c-brains.jp/20/03/30-170110.php
window.onload = function() {
    let dropZone = document.getElementById('drop-zone')
    let preview = document.getElementById('preview')
    let fileInput = document.getElementById('file-input')

    dropZone.addEventListener('dragover', function(e) {
        e.stopPropagation()
        e.preventDefault()
        this.style.background = '#e1e7f0'
    }, false)

    dropZone.addEventListener('dragleave', function(e) {
        e.stopPropagation()
        e.preventDefault()
        this.style.background = '#ffffff'
    }, false)

    fileInput.addEventListener('change', function () {
        previewFile(this.files[0])
    })

    dropZone.addEventListener('drop', function(e) {
        e.stopPropagation()
        e.preventDefault()
        this.style.background = '#ffffff' //背景色を白に戻す
        let files = e.dataTransfer.files //ドロップしたファイルを取得
        if (files.length > 1) return alert('you can only one file to upload ')
        fileInput.files = files //inputのvalueをドラッグしたファイルに置き換える。
        previewFile(files[0])
    }, false)

    function previewFile(file) {
        /* FileReaderで読み込み、プレビュー画像を表示。 */
        let fr = new FileReader()
        fr.readAsDataURL(file)
        fr.onload = function() {
            let img = document.createElement('img')
            img.setAttribute('src', fr.result)
            img.setAttribute('class', 'image')
            preview.innerHTML = ''
            preview.appendChild(img)
        }
    }
}