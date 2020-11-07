queryString = window.location.search;
urlParams = new URLSearchParams(queryString);
ep = urlParams.get("ep");
title = urlParams.get("title");
tokenlong = urlParams.get("tokenlong");
tokenshort = urlParams.get("tokenshort");
if (tokenlong == null || tokenshort == null) {
    showerror();
    throw '';
}
if (ep != null && title != null) {
    document.getElementById("title").innerHTML = "Episode " + ep + " - " + title
} else {
   document.getElementById("title").parentElement.removeChild(document.getElementById("title"))
}
urlstart = "https://rtv3-roosterteeth.akamaized.net/store/"+tokenlong+"-"+tokenshort+"/ts/"+tokenshort+"-hls_"
urlend = "p-store-"+tokenlong+".m3u8"
document.getElementById("1080p").src = urlstart + "1080" + urlend
document.getElementById("720p").src = urlstart + "720" + urlend
document.getElementById("480p").src = urlstart + "480" + urlend
document.getElementById("360p").src = urlstart + "360" + urlend
document.getElementById("240p").src = urlstart + "240" + urlend
let videojshls= videojs('rwby', {  html5: {  
    nativeAudioTracks: false,
    nativeVideoTracks: false,
    hls: {
        debug: true,
        overrideNative: false
   }
}});
videojshls.on("error", showerror)
videojshls.controlBar.addChild('QualitySelector');
videojshls.play()

function showerror() {
    document.getElementById("player").innerHTML = '<p class="text subhead">Failed to load video</p><p class="text">Did you get this link from RWBY bot?</p>'
    document.getElementById("break").parentElement.removeChild(document.getElementById("break"))
    if (document.getElementById("title") != null) {
        document.getElementById("title").parentElement.removeChild(document.getElementById("title"))
    }
}