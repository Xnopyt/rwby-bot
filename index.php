<?php
    $json = file_get_contents('./rwby_info.json');
    $json_data = json_decode($json,true);
    $magic_long = $json_data["magic_long"];
    $magic_short = $json_data["magic_short"];
    $epnum = $json_data["epnum"];
    $eptitle = $json_data["title"];
?>
<link href="https://how2trianglemuygud.com/main.css" rel="stylesheet">
<link href="https://unpkg.com/video.js@6.4.0/dist/video-js.css" rel="stylesheet">
<link href="https://unpkg.com/silvermine-videojs-quality-selector@1.1.2/dist/css/quality-selector.css" rel="stylesheet">
<script src="https://unpkg.com/video.js@6.4.0/dist/video.js"></script>
<script src="https://unpkg.com/videojs-flash@2.0.1/dist/videojs-flash.js"></script>
<script src="https://unpkg.com/videojs-contrib-hls@5.12.2/dist/videojs-contrib-hls.js"></script>
<script src="https://unpkg.com/silvermine-videojs-quality-selector/dist/js/silvermine-videojs-quality-selector.min.js"></script>

<html>
    <head>
        <title>H2TMG - RWBY Vol 6</title>
        <meta charset="utf-8">
    </head>
    <body class="w3-pale-blue">
        <div class="w3-cell-row">
            <div class="w3-container w3-light-grey w3-cell">
                <h1>How2TriangleMuyGud</h1>
            </div>
            <div class="w3-bar w3-container w3-grey w3-cell">
                <a href="https://how2trianglemuygud.com/" class="w3-bar-item w3-black w3-hover-teal" style="height:72px"><h3><b>Home</b></h3></a>
                <a href="https://how2trianglemuygud.com/nextcloud/" class="w3-bar-item w3-black w3-hover-teal" style="height:72px"><h4>NextCloud</h4></a>
                <a href="https://how2trianglemuygud.com/rainloop/" class="w3-bar-item w3-black w3-hover-teal" style="height:72px"><h4>Rainloop</h4></a>
            </div>
        </div>
        <div class="w3-cell-row">
                <div class="w3-cell w3-pale-blue" style="width:25%"></div>
                <div class="w3-container w3-teal w3-cell">
                    <h2 align="center">RWBY Volume 6 Episode <?=$epnum?> - <?=$eptitle?></a></h2>
                </div>
                <div class="w3-cell w3-pale-blue" style="width:25%"></div>
            </div>
        </br></br></br>
            <center><video id="my_video_1" class="video-js vjs-default-skin" controls preload="auto" width="800" height="450" autoplay="true" data-setup='{}'>
                <source src="https://rtv3-video.roosterteeth.com/store/<?=$magic_long?>-<?=$magic_short?>/ts/<?=$magic_short?>-hls_1080p-store-<?=$magic_long?>.m3u8" type="application/x-mpegURL" label="1080p" selected="true">
                <source src="https://rtv3-video.roosterteeth.com/store/<?=$magic_long?>-<?=$magic_short?>/ts/<?=$magic_short?>-hls_720p-store-<?=$magic_long?>.m3u8" type="application/x-mpegURL" label="720p">
                <source src="https://rtv3-video.roosterteeth.com/store/<?=$magic_long?>-<?=$magic_short?>/ts/<?=$magic_short?>-hls_480p-store-<?=$magic_long?>.m3u8" type="application/x-mpegURL" label="480p">
                <source src="https://rtv3-video.roosterteeth.com/store/<?=$magic_long?>-<?=$magic_short?>/ts/<?=$magic_short?>-hls_360p-store-<?=$magic_long?>.m3u8" type="application/x-mpegURL" label="360p">
                <source src="https://rtv3-video.roosterteeth.com/store/<?=$magic_long?>-<?=$magic_short?>/ts/<?=$magic_short?>-hls_240p-store-<?=$magic_long?>.m3u8" type="application/x-mpegURL" label="240p">
            </video></center>
        </br></br></br></br></br></br></br></br></br></br></br></br></br></br>
        <div id="footer-pos" class="w3-container w3-bar w3-grey">
            <div class="w3-bar-item w3-grey" style="width:50%">
                <p>Contact: <a href="mailto:admin@how2trianglemuygud.com">admin@how2trianglemuygud.com</a></p>
            </div>
            <div class="w3-bar-item w3-grey" style="width:50%">
                <p align="right">How2TriangleMuyGud.com</p>
            </div>
        </div>
        <script>
            let videojshls= videojs('my_video_1', {  html5: {  
            nativeAudioTracks: false,
            nativeVideoTracks: false,
            hls: {
              debug: true,
              overrideNative: false
            }
           }});
           videojshls.controlBar.addChild('QualitySelector');
        </script>
    </body>
</html>