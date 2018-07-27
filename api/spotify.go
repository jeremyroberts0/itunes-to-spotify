package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jeremyroberts0/itunes-to-spotify/itunes"
	"github.com/jeremyroberts0/itunes-to-spotify/types"
	"github.com/jeremyroberts0/pool"
	"github.com/zmb3/spotify"

	"github.com/gin-gonic/gin"
)

type matchErr struct {
	Song types.ItunesSong `json:"song,omitempty"`
	Err  error            `json:"err,omitempty"`
}

type matchCompletedResponse struct {
	PlaylistName     string     `json:"playlist_name,omitempty"`
	UnmatchedLines   []matchErr `json:"unmatched_lines,omitempty"`
	TotalUnmatched   int        `json:"total_unmatched,omitempty"`
	TotalMatched     int        `json:"total_matched,omitempty"`
	Total            int        `json:"total,omitempty"`
	MatchSuccessRate float32    `json:"match_success_rate,omitempty"`
}

func makeSearchJob(client spotify.Client, song types.ItunesSong, index, total int, idChan chan spotify.ID, matchErrChan chan matchErr) func() error {
	return func() error {
		fmt.Printf("Searching for song %v of %v\n", index, total)

		results, err := client.Search(
			fmt.Sprintf("track:%v artist:%v", song.Name, song.Artist),
			spotify.SearchTypeTrack,
		)

		if err != nil {
			matchErrChan <- matchErr{
				Song: song,
				Err:  err,
			}
		} else if results.Tracks == nil || len(results.Tracks.Tracks) == 0 {
			matchErrChan <- matchErr{
				Song: song,
				Err:  errors.New("Song not found on spotify"),
			}
		} else {
			track := results.Tracks.Tracks[0]
			idChan <- track.ID
		}

		return nil
	}
}

func MatchItunesPlaylistToSpotify(router *gin.Engine) {
	router.POST("/itunes-to-spotify", func(c *gin.Context) {
		playlistName := c.Query("name")
		token, err := getCookieToken(c)

		if err != nil || token.AccessToken == "" {
			c.IndentedJSON(http.StatusUnauthorized, createError("Auth cookie missing", err))
			return
		}

		if playlistName == "" {
			c.IndentedJSON(http.StatusBadRequest, createError("Missing 'name' query parameter", errors.New("Missing 'name' query parameter")))
			return
		}

		client := auth.NewClient(token)

		// AutoRetry automatically tried requests again when they fail due to rate limiting
		// It will wait the ms specified in the Retry-After header, as per the Spotify REST API Guidelines
		client.AutoRetry = true

		songs := itunes.ParsePlaylist(c.Request.Body)

		trackIds := []spotify.ID{}
		matchErrs := []matchErr{}

		idChan := make(chan spotify.ID)
		matchErrChan := make(chan matchErr)

		searchPool := pool.New(1)
		for index, song := range songs {
			searchPool.Add(makeSearchJob(
				client,
				song,
				index,
				len(songs),
				idChan,
				matchErrChan,
			))
		}

		start := time.Now()
		fmt.Printf("Starting matching on %v songs", len(songs))
		doneChan, errsChan := searchPool.MakeAsyncChans()
		go searchPool.RunAsync(doneChan, errsChan)

	ChanBlocker:
		for {
			select {
			case <-errsChan:
				// Don't care about errors in pool, we've captured them another way
			case id := <-idChan:
				trackIds = append(trackIds, id)
			case matchErr := <-matchErrChan:
				matchErrs = append(matchErrs, matchErr)
				fmt.Println(matchErr.Err.Error())
			case <-doneChan:
				break ChanBlocker
			}
		}

		duration := time.Now().Sub(start)
		fmt.Printf("Matching completed in %v seconds\n", duration.Seconds())
		fmt.Printf("%v of %v songs matched\n", len(trackIds), len(songs))

		// Wait a second so Spotify's rate limiting chills out
		time.Sleep(time.Second)

		user, err := client.CurrentUser()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, createError("Error creating playlist", err))
			return
		}
		userID := user.ID
		playlist, err := client.CreatePlaylistForUser(userID, playlistName, false)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, createError("Error creating playlist", err))
			return
		}

		trackIDGroups := [][]spotify.ID{[]spotify.ID{}}
		currentIndex := 0

		for _, trackID := range trackIds {
			trackIDGroups[currentIndex] = append(trackIDGroups[currentIndex], trackID)
			if len(trackIDGroups[currentIndex]) == 100 {
				currentIndex++
				trackIDGroups = append(trackIDGroups, []spotify.ID{})
			}
		}
		for _, trackIDGroup := range trackIDGroups {
			_, err = client.AddTracksToPlaylist(userID, playlist.ID, trackIDGroup...)
			if err != nil {
				c.IndentedJSON(
					http.StatusInternalServerError,
					createError(
						fmt.Sprintf("Error adding tracks to playlist, you probably have an incomplete playlist in your account called %v", playlistName),
						err,
					),
				)
				return
			}
		}

		fmt.Println(trackIds, matchErrs)
		c.IndentedJSON(http.StatusOK, matchCompletedResponse{
			PlaylistName:     playlistName,
			UnmatchedLines:   matchErrs,
			Total:            len(songs),
			TotalMatched:     len(trackIds),
			TotalUnmatched:   len(matchErrs),
			MatchSuccessRate: float32(len(trackIds)) / float32(len(songs)),
		})
	})
}
