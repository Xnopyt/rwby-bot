<?php
    $json = file_get_contents('./rwby_info.json');
    $json_data = json_decode($json,true);
    $magic_long = $json_data["magic_long"];
    $magic_short = $json_data["magic_short"];
    $epnum = $json_data["epnum"];
    $eptitle = $json_data["title"];
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge"> 
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" type="text/css" href="../styles/home.css">
    <link href="https://fonts.googleapis.com/css?family=Red+Hat+Text&display=swap" rel="stylesheet">
    <link href="https://unpkg.com/video.js@6.4.0/dist/video-js.css" rel="stylesheet">
    <link href="https://unpkg.com/silvermine-videojs-quality-selector@1.1.2/dist/css/quality-selector.css" rel="stylesheet">
    <script src="https://unpkg.com/video.js@6.4.0/dist/video.js"></script>
    <script src="https://unpkg.com/videojs-flash@2.0.1/dist/videojs-flash.js"></script>
    <script src="https://unpkg.com/videojs-contrib-hls@5.12.2/dist/videojs-contrib-hls.js"></script>
    <script src="https://unpkg.com/silvermine-videojs-quality-selector/dist/js/silvermine-videojs-quality-selector.min.js"></script>
    <link href="../assets/favicon.svg" rel="icon">
    <title>Xnopyt - Home</title>
</head>
<body>
    <div class="container">
        <p class="title"><img src="../assets/logo-greyscale.svg">nopyt - RWBY Vol 7</p>
        <br />
        <center><video id="my_video_1" class="video-js vjs-default-skin" controls preload="auto" width="800" height="450" autoplay="true" data-setup='{}'>
            <source src="https://rtv3-roosterteeth.akamaized.net/store/<?=$magic_long?>-<?=$magic_short?>/ts/<?=$magic_short?>-hls_1080p-store-<?=$magic_long?>.m3u8" type="application/x-mpegURL" label="1080p" selected="true">
            <source src="https://rtv3-roosterteeth.akamaized.net/store/<?=$magic_long?>-<?=$magic_short?>/ts/<?=$magic_short?>-hls_720p-store-<?=$magic_long?>.m3u8" type="application/x-mpegURL" label="720p">
            <source src="https://rtv3-roosterteeth.akamaized.net/store/<?=$magic_long?>-<?=$magic_short?>/ts/<?=$magic_short?>-hls_480p-store-<?=$magic_long?>.m3u8" type="application/x-mpegURL" label="480p">
            <source src="https://rtv3-roosterteeth.akamaized.net/store/<?=$magic_long?>-<?=$magic_short?>/ts/<?=$magic_short?>-hls_360p-store-<?=$magic_long?>.m3u8" type="application/x-mpegURL" label="360p">
            <source src="https://rtv3-roosterteeth.akamaized.net/store/<?=$magic_long?>-<?=$magic_short?>/ts/<?=$magic_short?>-hls_240p-store-<?=$magic_long?>.m3u8" type="application/x-mpegURL" label="240p">
        </video></center>
        <br /> <br />
    </div>
</body>
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
</html>