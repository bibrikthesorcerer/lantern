import { updatePlayer } from "./utils.js";

const tracklist = document.getElementById("tracklist")

function createTrack(track) {
    const root = document.createElement("div")
    root.classList.add("trackCard")
    root.setAttribute("data-title", `${track.title}`)
    root.setAttribute("data-artist", `${track.artist}`)
    root.setAttribute("data-id", `${track.id}`)

    const img = document.createElement("img")
    img.setAttribute("src", `/api/cover/${track.id}`)
    img.style.width = "10vw"
    img.setAttribute("loading", "lazy")
    root.appendChild(img)

    const title = document.createElement("span")
    title.classList.add("trackTitle")
    title.textContent = track.title
    root.appendChild(title)

    const info = document.createElement("span")
    info.textContent = `by ${track.artist}, from ${track.album}`
    root.appendChild(info)
    return root
}

async function loadTracks() {
    const res = await fetch("/api/tracks")
    const tracks = await res.json()

    for (const track of tracks) {
        const elem = createTrack(track)
        elem.addEventListener("click", function(e) {updatePlayer(e.currentTarget)})
        tracklist.appendChild(elem)
    }
}

loadTracks()