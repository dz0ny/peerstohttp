# Plextor

Simple torrent proxy to Plex http API.

Built top on [anacrolix/torrent](https://github.com/anacrolix/torrent) lib.

## HTTP API
##### URL Parameters:
* **playlist** - output file format, one of these values: `m3u`,`html`,`json`
* **hash** - torrent info hash. Example: `08ada5a7a6183aae1e09d831df6748d566095a10`
* **extsWhitelist** - list of whitelisted file extensions. Possible values: "-" (any) or list extension names divided by comma. Examples: "`-`", "`mp3,mp4a`"
* **tagsBlacklist** - list of blacklisted tags, extracted from file names. Possible values: "-" (no filter) or list tags divided by comma. See /playlist/tags.go for full list of possible tags. Examples: "`-`", "`remix,interview`"

Get list of files by magnet uri:
```
GET http://localhost/list/{playlist}/{extsWhitelist}/{tagsBlacklist}/magnet/{magnetURI}
```

Get list of files by info hash:
```
GET http://localhost/list/{playlist}/{extsWhitelist}/{tagsBlacklist}/hash/{hash}
```

Download file:
```
GET http://localhost/content/{hash}/{filePath}
```

## Examples
Get HTML links list for Sintel by torrent hash:
```
http://localhost/list/html/mp4/-/hash/08ada5a7a6183aae1e09d831df6748d566095a10
```

or by magnet URI:
```
http://localhost/list/html/mp4/-/magnet/magnet:?xt=urn:btih:08ada5a7a6183aae1e09d831df6748d566095a10&dn=Sintel&tr=udp%3A%2F%2Fexplodie.org%3A6969&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.empire-js.us%3A1337&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=wss%3A%2F%2Ftracker.btorrent.xyz&tr=wss%3A%2F%2Ftracker.fastcast.nz&tr=wss%3A%2F%2Ftracker.openwebtorrent.com&ws=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2F&xs=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2Fsintel.torrent
```

Then watch Sintel.mp4 using VLC video player:
```
$ vlc http://localhost/content/08ada5a7a6183aae1e09d831df6748d566095a10/Sintel/Sintel.mp4
```

Or open m3u playlist in VLC video player:
```
$ vlc http://localhost/list/m3u/mp4/-/hash/08ada5a7a6183aae1e09d831df6748d566095a10
```


## Build
```
$ make build
```
