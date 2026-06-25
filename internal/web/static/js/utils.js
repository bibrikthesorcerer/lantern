export {updatePlayer, className, prevTrack, nextTrack};

let currentTrackElem = null;
let player;
let playerTrackTitle;
let playerTrackArtist;
let playerTrackArt;

document.addEventListener('DOMContentLoaded', () => {
    player = document.getElementById("player")
    player.addEventListener('ended', nextTrack)
    playerTrackTitle = document.getElementById('track-title');
    playerTrackArtist = document.getElementById('track-artist');
    playerTrackArt = document.getElementById('track-art');
});

async function updatePlayer(track) {
    currentTrackElem = track

    player.src = `/api/stream/${track.getAttribute("data-id")}`
    playerTrackTitle.textContent = track.getAttribute("data-title")
    playerTrackArtist.textContent = track.getAttribute("data-artist")
    playerTrackArt.replaceChildren(getTrackCover(track.getAttribute("data-id")))

    document.querySelectorAll('.albumTrackEntryWrapper').forEach(el => el.classList.remove("active"))
    track.classList.add("active")

    player.play()
}

async function prevTrack() {
    if (!currentTrackElem) return;

    const prevTrack = currentTrackElem.previousElementSibling;
    if (prevTrack) {
        updatePlayer(prevTrack)
    } else {
        currentTrackElem.classList.remove("active")
        currentTrackElem = null;
    }
}

async function nextTrack() {
    if (!currentTrackElem) return;

    const nextTrack = currentTrackElem.nextElementSibling;
    if (nextTrack) {
        updatePlayer(nextTrack)
    } else {
        currentTrackElem.classList.remove("active")
        currentTrackElem = null;
    }
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