export {updatePlayer, className};

const player = document.getElementById("player")
const playerTrackTitle = document.getElementById('track-title');
const playerTrackArtist = document.getElementById('track-artist');
const playerTrackArt = document.getElementById('track-art');

async function updatePlayer(track) {
    player.src = `/api/stream/${track.getAttribute("data-id")}`
    playerTrackTitle.textContent = track.getAttribute("data-title")
    playerTrackArtist.textContent = track.getAttribute("data-artist")
    playerTrackArt.replaceChildren(getTrackCover(track.getAttribute("data-id")))
    player.play()
}

function getTrackCover(id) {
    const img = document.createElement("img")
    img.src = `/api/cover/${id}`
    return img
}

function className(elem, cls) {
    const classes = cls.split(" ")
    for (const cl of classes) {
        elem.classList.add(cl)
    }
}