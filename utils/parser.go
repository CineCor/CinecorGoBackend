package utils

import (
	"fmt"
	Models "github.com/CineCor/CinecorGoBackend/models"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

var PARSE_TIMEOUT = 60000
var PARSE_USER_AGENT = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"

func getHoursDateFromText(datesString string) []string {
	var utcDates []string
	rDates, _ := regexp.Compile(`(\d{2}:\d{2})`)
	rDate, _ := regexp.Compile(`(?P<Hour>\d{2}):(?P<Minutes>\d{2})`)
	dates := rDates.FindAllString(datesString, -1)
	var lastHour int = -1
	for _, dateString := range dates {
		arrayString := rDate.FindStringSubmatch(dateString)
		if len(arrayString) == 3 {
			_actualHour, err1 := strconv.ParseInt(arrayString[1], 10, 64)
			_actualMinutes, err2 := strconv.ParseInt(arrayString[2], 10, 64)
			actualHour := int(_actualHour)
			actualMinutes := int(_actualMinutes)
			if err1 != nil {
				continue
			}
			if err2 != nil {
				continue
			}
			nowDate := time.Now().Truncate(24 * time.Hour)

			if actualHour < lastHour {
				nowDate = time.Date(nowDate.Year(), nowDate.Month(), nowDate.Day(), actualHour, actualMinutes, 0, 0, nowDate.Location())
				nowDate = nowDate.AddDate(0, 0, 1)
			} else {
				nowDate = time.Date(nowDate.Year(), nowDate.Month(), nowDate.Day(), actualHour, actualMinutes, 0, 0, nowDate.Location())

			}
			utcDates = append(utcDates, nowDate.UTC().Format(time.RFC3339))

			lastHour = actualHour
		}
	}
	return utcDates
}

func GetCinemas() (error, []Models.Cinema) {
	err, cinemas := parseWeb()
	if err != nil {
		return err, nil
	}
	sort.Slice(cinemas, func(i, j int) bool {
		return cinemas[i].Id > cinemas[j].Id
	})
	return err, cinemas
}

func parseWeb() (error, []Models.Cinema) {
	cinemas := make([]Models.Cinema, 0)

	client := &http.Client{}

	baseURL, err := url.Parse(os.Getenv("PARSE_URL"))
	if err != nil {
		return err, nil
	}

	req, err := http.NewRequest("GET", os.Getenv("PARSE_URL"), nil)

	if err != nil {
		return err, nil
	}

	req.Header.Set("User-Agent", PARSE_USER_AGENT)

	resp, err := client.Do(req)

	if err != nil {
		return err, nil
	}

	defer resp.Body.Close()
	charset := "windows-1252"
	utfBody, err := iconv.NewReader(resp.Body, charset, "utf-8")
	if err != nil {
		return err, nil
	}

	doc, _ := goquery.NewDocumentFromReader(utfBody)
	cinemasElements := doc.Find("div#bloqueportadaa")
	if len(cinemasElements.Nodes) > 0 {
		cinemasElements.Each(func(i int, cinemaElement *goquery.Selection) {
			cinema := new(Models.Cinema)
			cinema.Name = cinemaElement.Find("h1 a").Text()

			urlCinema, _ := cinemaElement.Find("a").First().Attr("href")
			u, err := url.Parse(urlCinema)
			if err != nil {
				return
			}
			cinema.Id, err = strconv.Atoi(u.Query()["id"][0])
			if err != nil {
				return
			}
			moviesElements := cinemaElement.Find("div.pildora")
			cinema.Movies = make([]Models.Movie, 0)
			moviesElements.Each(func(j int, moviesElement *goquery.Selection) {
				movieLinks := moviesElement.Find("a")
				if len(movieLinks.Nodes) == 0 {
					return
				}
				movie := new(Models.Movie)

				urlMovie, _ := moviesElement.Find("a").First().Attr("href")
				u, err := url.Parse(urlMovie)
				if err != nil {
					return
				}
				movie.Id, err = strconv.Atoi(u.Query()["id"][0])
				if err != nil {
					return
				}
				movie.Title = moviesElement.Find("a").Text()
				u, err = url.Parse("gestor/ficheros/imagen" + u.Query()["id"][0] + ".jpeg")
				if err != nil {
					return
				}
				movie.Hours = getHoursDateFromText(moviesElement.Find("h5").Text())

				movie.Images = make(map[string]string)

				movie.Images["POSTER"] = baseURL.ResolveReference(u).String()
				moviesElementURL, _ := moviesElement.Find("a").Attr("href")
				if err != nil {
					return
				}
				u, err = url.Parse(moviesElementURL)
				if err != nil {
					return
				}

				movie.Url = baseURL.ResolveReference(u).String()

				cinema.Movies = append(cinema.Movies, *movie)
			})
			if len(cinema.Movies) == 0 {
				return
			}
			cinemas = append(cinemas, *cinema)

		})
	}

	return nil, cinemas
}

func fillDataWithOriginalWeb(movie *Models.Movie) {
	fmt.Printf("Fetching Movie %s from %s\n", movie.Title, movie.Url)
	client := &http.Client{}

	req, err := http.NewRequest("GET", movie.Url, nil)

	if err != nil {
		return
	}

	req.Header.Set("User-Agent", PARSE_USER_AGENT)

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	charset := "windows-1252"
	utfBody, err := iconv.NewReader(resp.Body, charset, "utf-8")
	if err != nil {
		return
	}

	doc, _ := goquery.NewDocumentFromReader(utfBody)
	h5Elements := doc.Find("div#sobrepelicula h5")
	if len(h5Elements.Nodes) > 0 {
		h5Elements.Each(func(i int, h5Element *goquery.Selection) {
			if i == 0 {
				movie.RawDescription = h5Element.Text()
			} else if i == 1 {
				movie.Overview = h5Element.Text()
			}

		})
	}
	return

}
