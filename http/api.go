package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	"plex-torrent/app"
	"plex-torrent/http/host"
	list_render "plex-torrent/http/render"
	"plex-torrent/playlist"
)

const (
	paramMagnet     = "magnet"
	paramHash       = "hash"
	paramPath       = "path"
	paramWhitelist  = "whitelist"
	paramIgnoretags = "ignoretags"
)

var (
	patternList = fmt.Sprintf("%s:[json,m3u,html]+", list_render.ParamContentType)
)

type handle struct {
	app *app.App
}

func RouteApp(r chi.Router, app *app.App) {
	var h = handle{
		app: app,
	}

	r.Route("/list/{"+patternList+"}/{"+paramWhitelist+"}/{"+paramIgnoretags+"}/", func(r chi.Router) {
		r.Use(
			whitelist,
			ignoretags,
			host.Host,
			list_render.ListContentType,
			list_render.SetListResponder,
		)

		r.With(hash).Get("/hash/{hash}", h.hash)
		r.With(magnet).Get("/magnet/*", h.magnet)
	})

	r.With(hash, path).Get("/content/{"+paramHash+"}/*", h.content)
}

func (h *handle) hash(w http.ResponseWriter, r *http.Request) {
	var hash = r.Context().Value(paramHash).(string)
	var whitelist = r.Context().Value(paramWhitelist).(map[string]struct{})
	var ignoretags = r.Context().Value(paramIgnoretags).(map[string]struct{})

	var t, err = h.app.TrackHashContext(r.Context(), metainfo.NewHashFromHex(hash))
	if err != nil {
		log.Error().Err(err).Msg("track by hash")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	render.Render(w, r, &playlist.PlayList{Torr: t, Whitelist: whitelist, IgnoreTags: ignoretags})
}

func (h *handle) magnet(w http.ResponseWriter, r *http.Request) {
	var magnet = r.Context().Value(paramMagnet).(*metainfo.Magnet)
	var whitelist = r.Context().Value(paramWhitelist).(map[string]struct{})
	var ignoretags = r.Context().Value(paramIgnoretags).(map[string]struct{})

	var t, err = h.app.TrackMagnetContext(r.Context(), magnet)
	if err != nil {
		log.Error().Err(err).Msg("track by hash")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	render.Render(w, r, &playlist.PlayList{Torr: t, Whitelist: whitelist, IgnoreTags: ignoretags})
}

func (h *handle) content(w http.ResponseWriter, r *http.Request) {
	var err error

	var hash = r.Context().Value(paramHash).(string)
	var path = r.Context().Value(paramPath).(string)

	var t, ok = h.app.Client().Torrent(metainfo.NewHashFromHex(hash))

	if !ok {
		t, ok = addNewTorrentHash(r.Context(), h.app, hash)
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	}

	select {
	case <-r.Context().Done():
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		return

	case <-t.GotInfo():
	}

	if t.Info().IsDir() && strings.Count(path, "/") == 0 {
		err = serveTorrentDir(w, r, t, path)
	} else {
		err = serveTorrentFile(w, r, t, path)
	}

	if err != nil {
		log.Warn().Err(err).Msg("serve content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//func fileInfoHeader(fi *torrent.File) (*zip.FileHeader, error) {

//}
