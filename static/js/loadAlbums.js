import { className, updatePlayer } from "./utils.js";

const albumList = document.getElementById("albumlist")

function createTrack(track) {
    const root = document.createElement("div")
    className(root, "albumTrackEntryBG flex flex-row justify-between px-4 py-2 items-center mb-2 rounded transition-colors duration-300")
    root.setAttribute("data-title", `${track.title}`)
    root.setAttribute("data-artist", `${track.artist}`)
    root.setAttribute("data-id", `${track.id}`)
    
    const entry = document.createElement("div")
    className(entry, "albumTrackEntry")
    root.appendChild(entry)

    const trackNum = document.createElement("span")
    trackNum.textContent = track.track
    entry.appendChild(trackNum)
    
    const trackInfo = document.createElement("div")
    entry.appendChild(trackInfo)

    const title = document.createElement("span")
    title.classList.add("trackTitle")
    title.textContent = track.title
    trackInfo.appendChild(title)

    const artist = document.createElement("span")
    artist.textContent = `${track.artist}`
    trackInfo.appendChild(artist)
    
    const downloadButton = document.createElement("a")
    className(downloadButton, "font-bold text-xl px-4 py-2 bg-sky-700 rounded-xl")
    downloadButton.href = `/api/tracks/${track.id}/download`
    const downloadIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" fill="currentColor" class="bi bi-download" viewBox="0 0 18 18">
  <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5"/>
  <path d="M7.646 11.854a.5.5 0 0 0 .708 0l3-3a.5.5 0 0 0-.708-.708L8.5 10.293V1.5a.5.5 0 0 0-1 0v8.793L5.354 8.146a.5.5 0 1 0-.708.708z"/>
</svg>`
    downloadButton.innerHTML = downloadIcon
    root.appendChild(downloadButton)

    return root
}

function createAlbum(album) {
    const root = document.createElement("div")
    root.classList.add("albumContainer")

    const albumHeader = document.createElement("div")
    albumHeader.classList.add("albumHeader")
    root.appendChild(albumHeader)

    const img = document.createElement("img")
    img.setAttribute("src", `/api/cover/${album.tracks[0].id}`)
    img.style.width = "10vw"
    img.setAttribute("loading", "lazy")
    albumHeader.appendChild(img)
    
    const albumInfo = document.createElement("div")
    className(albumInfo, "flex flex-col")
    albumHeader.appendChild(albumInfo)

    const title = document.createElement("span")
    className(title, "text-3xl font-bold")
    title.textContent = album.title
    albumInfo.appendChild(title)
    
    const albumDesc = document.createElement("div")
    className(albumDesc, "flex flex-row gap-2")
    albumInfo.appendChild(albumDesc)
    
    const desc = document.createElement("span") 
    desc.textContent = `${album.artist} • ${album.year}`
    albumDesc.appendChild(desc)

    const tracksContainer = document.createElement("div")
    tracksContainer.classList.add("tracksContainer")
    root.appendChild(tracksContainer)

    for (const track of album.tracks) {
        const elem = createTrack(track)
        elem.addEventListener("click", function(e) {updatePlayer(e.currentTarget)})
        tracksContainer.appendChild(elem)
    }
    
    return root
}

async function loadAlbums() {
    const res = await fetch("/api/albums")
    const albums = await res.json()

    for (const al of albums) {
        const albumElem = createAlbum(al)
        albumList.appendChild(albumElem)
    }
}

loadAlbums()