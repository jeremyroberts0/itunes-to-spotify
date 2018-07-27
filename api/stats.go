package api

import (
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/jeremyroberts0/itunes-to-spotify/itunes"
)

type Correlated struct {
	Name              string  `json:"name"`
	SongCount         int     `json:"song_count"`
	PercentageOfSongs float32 `json:"percentage_of_songs"`
}

type PlaylistStats struct {
	TotalSongs      int            `json:"total_songs"`
	UniqueArtists   int            `json:"unique_artists"`
	UniqueSongs     int            `json:"unique_songs"`
	UniqueAlbums    int            `json:"unique_albums"`
	ArtistFrequency map[string]int `json:"-"`
	SongFrequency   map[string]int `json:"-"`
	AlbumFrequency  map[string]int `json:"-"`
	Artists         []Correlated   `json:"artists"`
	Songs           []Correlated   `json:"songs"`
	Albums          []Correlated   `json:"albums"`
}

func GetPlaylistStats(c *gin.Context) {
	songs := itunes.ParsePlaylist(c.Request.Body)

	stats := PlaylistStats{
		ArtistFrequency: make(map[string]int),
		SongFrequency:   make(map[string]int),
		AlbumFrequency:  make(map[string]int),
	}
	for _, song := range songs {
		if _, ok := stats.ArtistFrequency[song.Artist]; !ok {
			stats.ArtistFrequency[song.Artist] = 0
		}
		if _, ok := stats.AlbumFrequency[song.Album]; !ok {
			stats.AlbumFrequency[song.Album] = 0
		}
		if _, ok := stats.SongFrequency[song.Name]; !ok {
			stats.SongFrequency[song.Name] = 0
		}

		stats.ArtistFrequency[song.Artist]++
		stats.AlbumFrequency[song.Album]++
		stats.SongFrequency[song.Name]++

		stats.TotalSongs++
	}

	stats.UniqueAlbums = len(stats.AlbumFrequency)
	stats.UniqueArtists = len(stats.ArtistFrequency)
	stats.UniqueSongs = len(stats.SongFrequency)

	for artist, frequency := range stats.ArtistFrequency {
		stats.Artists = append(stats.Artists, Correlated{
			Name:              artist,
			SongCount:         frequency,
			PercentageOfSongs: float32(frequency) / float32(stats.TotalSongs),
		})
	}
	for album, frequency := range stats.AlbumFrequency {
		stats.Albums = append(stats.Albums, Correlated{
			Name:              album,
			SongCount:         frequency,
			PercentageOfSongs: float32(frequency) / float32(stats.TotalSongs),
		})
	}
	for song, frequency := range stats.SongFrequency {
		stats.Songs = append(stats.Songs, Correlated{
			Name:              song,
			SongCount:         frequency,
			PercentageOfSongs: float32(frequency) / float32(stats.TotalSongs),
		})
	}

	sort.Slice(stats.Artists, func(i, j int) bool {
		return stats.Artists[i].SongCount > stats.Artists[j].SongCount
	})
	sort.Slice(stats.Albums, func(i, j int) bool {
		return stats.Albums[i].SongCount > stats.Albums[j].SongCount
	})
	sort.Slice(stats.Songs, func(i, j int) bool {
		return stats.Songs[i].SongCount > stats.Songs[j].SongCount
	})

	stats.Artists = stats.Artists[:10]
	stats.Albums = stats.Albums[:10]
	stats.Songs = stats.Songs[:10]

	c.IndentedJSON(200, stats)
}

func ApplyStatsRoutes(router *gin.Engine) {
	router.POST("/stats", GetPlaylistStats)
}
