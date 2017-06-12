package utils

import (
	"fmt"
	Models "github.com/CineCor/CinecorGoBackend/models"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ryanbradynd05/go-tmdb"
)

type TMDBManager struct {
	tmdb *tmdb.TMDb
}

var TMDB_LANGUAGE = "es-ES"

var tmdbmutex sync.Mutex
var tmdbmanager *TMDBManager

func NewTMDBManager() *TMDBManager {
	manager := new(TMDBManager)
	manager.tmdb = nil
	manager.tmdb = tmdb.Init(os.Getenv("TMDB_API_KEY"))

	return manager
}

func GetTMDBManagerInstance() *TMDBManager {
	if tmdbmanager == nil {
		tmdbmutex.Lock()
		defer tmdbmutex.Unlock()
		if tmdbmanager == nil {
			tmdbmanager = NewTMDBManager()
		}
	}
	return tmdbmanager
}
func getMovieColorsFromUrl(url string) map[string]string {
	var colors map[string]string

	err, palette := getPalettefromURL(url)
	if err == nil {
		swatch := palette.VibrantSwatch()
		if swatch == nil {
			swatch = palette.MutedSwatch()
		}
		if swatch == nil {
			return colors
		}
		colors = map[string]string{}
		_, color := getTextColorForBackground(uint32(swatch.RGBAInt()), MIN_CONTRAST_TITLE_TEXT)
		colors["TITLE"] = fmt.Sprintf("#%06s", strings.ToLower(strconv.FormatUint(uint64(color), 16))[2:])
		colors["MAIN"] = fmt.Sprintf("#%06s", strings.ToLower(strconv.FormatUint(uint64(swatch.RGBAInt()), 16))[2:])
	}
	return colors
}

func fillColors(movie *Models.Movie) {
	images := movie.Images
	if len(images) == 0 {
		return
	}
	var url string
	if value, ok := images["BACKDROP"]; ok {
		url = value
	} else {
		url = images["POSTER"]
	}

	movie.Colors = getMovieColorsFromUrl(url)
}

func FillMovieData(cinemas *[]Models.Cinema) {
	for idx, cinema := range *cinemas {
		for idy, movie := range cinema.Movies {
			if !fillDataWithExistingMovie(cinemas, &movie) {
				if !fillDataWithExternalApi(&movie) {
					fillDataWithOriginalWeb(&movie)
				}
				fillColors(&movie)
			}
			cinema.Movies[idy] = movie
		}
		(*cinemas)[idx] = cinema
	}
}
func fillDataWithExistingMovie(cinemas *[]Models.Cinema, movieOriginal *Models.Movie) bool {
	for _, cinema := range *cinemas {
		for _, movie := range cinema.Movies {
			if movie.Id == movieOriginal.Id && len(movie.Overview) != 0 {
				(*movieOriginal).Copy(movie)
				return true
			}
		}
	}
	return false
}
func fillDataWithExternalApi(movie *Models.Movie) bool {
	fmt.Printf("Fetching Movie %s from IMDB API\n", movie.Title)

	var options = make(map[string]string)
	movieResults, err := searchMovie(movie.Title, time.Now().Year())
	if err != nil || movieResults.TotalResults == 0 {
		movieResults, err = searchMovie(movie.Title, time.Now().Year()-1)

	}
	if err != nil || movieResults.TotalResults == 0 {
		movieResults, err = searchMovie(movie.Title, 0)
	}
	if err != nil || movieResults.TotalResults > 0 {
		options["language"] = TMDB_LANGUAGE
		options["append_to_response"] = "videos"
		result, _ := GetTMDBManagerInstance().tmdb.GetMovieInfo(movieResults.Results[0].ID, options)
		if result != nil {
			(*movie).Copy(*result)
			return true
		}
	}

	return false
}

func searchMovie(title string, year int) (*tmdb.MovieSearchResults, error) {

	var options = make(map[string]string)

	options["language"] = TMDB_LANGUAGE
	if year != 0 {
		options["year"] = fmt.Sprintf("%d", year)
	}
	options["include_adult"] = "true"
	options["page"] = "1"

	return GetTMDBManagerInstance().tmdb.SearchMovie(title, options)
}
