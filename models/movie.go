package models

import (
	"fmt"

	"github.com/ryanbradynd05/go-tmdb"
)

const POSTER, POSTER_THUMBNAIL, BACKDROP, BACKDROP_THUMBNAIL = 1, 2, 3, 4

const MAIN, TITLE = 1, 2

type Movie struct {
	Id             int
	Colors         map[string]string
	Images         map[string]string
	Hours          []string
	Genres         []string
	RawDescription string
	Imdb           string
	Trailer        string
	Duration       int
	ReleaseDate    string
	Overview       string
	Director       string
	Url            string
	Rating         float64
	Title          string
}

func (self *Movie) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["id"] = self.Id
	m["colors"] = self.Colors
	m["hours"] = self.Hours
	m["images"] = self.Images
	m["genres"] = self.Genres
	if len(self.RawDescription) > 0 {
		m["rawDescription"] = self.RawDescription
	}
	if len(self.Imdb) > 0 {
		m["imdb"] = self.Imdb
	}
	if len(self.Trailer) > 0 {
		m["trailer"] = self.Trailer
	}
	if len(self.ReleaseDate) > 0 {
		m["releaseDate"] = self.ReleaseDate
	}
	m["duration"] = self.Duration
	if len(self.Overview) > 0 {
		m["overview"] = self.Overview
	}
	if len(self.Director) > 0 {
		m["director"] = self.Director
	}
	if len(self.Url) > 0 {
		m["url"] = self.Url
	}
	m["rating"] = self.Rating
	if len(self.Title) > 0 {
		m["title"] = self.Title
	}
	return m

}

func (self *Movie) Copy(v interface{}) {
	switch movie := v.(type) { // HL
	case Movie:
		self.Colors = movie.Colors
		self.Images = movie.Images
		self.Genres = movie.Genres
		self.RawDescription = movie.RawDescription
		self.Imdb = movie.Imdb
		self.Trailer = movie.Trailer
		self.Duration = movie.Duration
		self.ReleaseDate = movie.ReleaseDate
		self.Overview = movie.Overview
		self.Director = movie.Director
		self.Url = movie.Url
		self.Rating = movie.Rating
		self.Title = movie.Title
	case tmdb.Movie:
		self.Title = movie.Title
		self.Imdb = "http://www.imdb.com/title/" + movie.ImdbID
		self.ReleaseDate = movie.ReleaseDate
		self.Overview = movie.Overview
		self.Rating = float64(movie.VoteAverage)
		self.Duration = int(movie.Runtime)
		self.Genres = make([]string, 0)
		if len(movie.PosterPath) > 0 {
			self.Images["POSTER"] = "http://image.tmdb.org/t/p/w780" + movie.PosterPath
			self.Images["POSTER_THUMBNAIL"] = "http://image.tmdb.org/t/p/w92" + movie.PosterPath
		}
		if len(movie.BackdropPath) > 0 {
			self.Images["BACKDROP"] = "http://image.tmdb.org/t/p/w780" + movie.BackdropPath
			self.Images["BACKDROP_THUMBNAIL"] = "http://image.tmdb.org/t/p/w300" + movie.BackdropPath
		}
		for _, genre := range movie.Genres {
			self.Genres = append(self.Genres, genre.Name)
		}
		for _, video := range movie.Videos.Results {
			if video.Type == "Trailer" && video.Site == "YouTube" {
				self.Trailer = "https://www.youtube.com/watch?v=" + fmt.Sprintf("%s", video.Key)
			}
		}
	}
}
