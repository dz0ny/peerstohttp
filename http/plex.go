package http

import (
	"fmt"
	"net/http"
	"plex-torrent/app"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/go-chi/chi"
)

func RoutePlex(r chi.Router, app *app.App) {

	r.Get("/library/sections", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<MediaContainer friendlyName="Torrent Library">
		   <Directory title="Torrents" key="1"></Directory>
		</MediaContainer>`))
	})

	r.Get("/library/sections/1", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`<MediaContainer friendlyName="Torrents">`))
		for _, t := range app.Client().Torrents() {
			w.Write([]byte(`<Video ratingKey="` + t.InfoHash().HexString() + `" key="/library/metadata/` + t.InfoHash().HexString() + `" title="` + t.Name() + `"></Video> `))
		}

		w.Write([]byte(`</MediaContainer>`))
	})

	r.With(hash).Get("/library/metadata/{hash}", func(w http.ResponseWriter, r *http.Request) {
		var hash = r.Context().Value(paramHash).(string)
		var t, ok = app.Client().Torrent(metainfo.NewHashFromHex(hash))

		if ok {
			w.Write([]byte(`<MediaContainer>
			<Video ratingKey="` + hash + `" key="/list/m3u/mp4/-/hash/` + hash + `" title="` + t.Name() + `">
				<Media audioCodec="aac" videoCodec="h264">
						<Part key="/list/m3u/mp4,mkv/-/hash/` + hash + `" duration="" size="` + fmt.Sprint(t.Info().Length) + `"></Part>  
				</Media>
			</Video> 
		 </MediaContainer>`))
		}

	})

}
