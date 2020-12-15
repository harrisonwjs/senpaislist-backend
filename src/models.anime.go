package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
)

type anime struct {
	ID              int64    `json:"id"`
	Title           string   `json:"title"`
	ContentType     string   `json:"contentType"`
	BeginDate       string   `json:"beginDate"`
	PremieredSeason string   `json:"premieredSeason"`
	Airing          bool     `json:"airing"`
	CurrentStatus   string   `json:"currentStatus"`
	NumEpisodes     int64    `json:"numEpisodes"`
	EpisodeDuration string   `json:"episodeDuration"`
	BroadcastTime   string   `json:"string"`
	Genres          []string `json:"genres"`
}

var animeList = []anime{}

func getAnimeIds() []int64 {
	client := &http.Client{}

	// the request
	req, _ := http.NewRequest("GET", "https://api.jikan.moe/v3/season/2020/fall", nil)

	// the response / any error
	resp, err := client.Do(req)

	// check for any errors
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return []int64{}
	}

	defer resp.Body.Close()
	// uint8 array
	resp_body, _ := ioutil.ReadAll(resp.Body)

	var animeIds []int64
	jsonparser.ArrayEach(resp_body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		animeId, err := jsonparser.GetInt(value, "mal_id")
		animeIds = append(animeIds, animeId)
	}, "anime")

	return animeIds
}

func getAllAnimes() []anime {
	client := &http.Client{}
	animeIds := getAnimeIdsT()
	var animeListTemp []anime
	// loop through the list of anime ids and retrieve info
	for i := 0; i < len(animeIds); i++ {
		animeId := animeIds[i]

		// the request
		req, _ := http.NewRequest("GET", "https://api.jikan.moe/v3/anime/"+strconv.Itoa(int(animeId)), nil)
		// the response / any error
		resp, err := client.Do(req)
		// check for any errors
		if err != nil {
			fmt.Println("Errored when sending request to the server")
			return []anime{}
		}

		// wait 4 seconds
		time.Sleep(4 * time.Second)

		// uint8 array
		resp_body, _ := ioutil.ReadAll(resp.Body)

		// get data
		currentAnime := anime{}
		paths := [][]string{
			{"mal_id"},
			{"title"},
			{"type"},
			{"aired", "from"},
			{"premiered"},
			{"airing"},
			{"status"},
			{"episodes"},
			{"duration"},
			{"broadcast"},
			// {"genres"},
		}
		jsonparser.EachKey(resp_body, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
			switch idx {
			case 0: // []string{"mal_id"}
				id, _ := jsonparser.ParseInt(value)
				currentAnime.ID = id
			case 1: // []string{"title"}
				title, _ := jsonparser.ParseString(value)
				currentAnime.Title = title
			case 2: // []string{"type"}
				contentType, _ := jsonparser.ParseString(value)
				currentAnime.ContentType = contentType
			case 3: // []string{"aired", "from"}
				beginDate, _ := jsonparser.ParseString(value)
				currentAnime.BeginDate = beginDate
			case 4: // []string{"premiered"}
				premieredSeason, _ := jsonparser.ParseString(value)
				currentAnime.PremieredSeason = premieredSeason
			case 5: // []string{"airing"}
				airing, _ := jsonparser.ParseBoolean(value)
				currentAnime.Airing = airing
			case 6: // []string{"status"}
				currentStatus, _ := jsonparser.ParseString(value)
				currentAnime.CurrentStatus = currentStatus
			case 7: // []string{"episodes"}
				numEpisodes, _ := jsonparser.ParseInt(value)
				currentAnime.NumEpisodes = numEpisodes
			case 8: // []string{"duration"}
				episodeDuration, _ := jsonparser.ParseString(value)
				currentAnime.EpisodeDuration = episodeDuration
			case 9: // []string{"broadcast"}
				broadcastTime, _ := jsonparser.ParseString(value)
				currentAnime.BroadcastTime = broadcastTime
				// case 10: // []string{"genres"}
				// 	var genreArray []string
				// 	jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				// 		genreArray = append(genreArray, )
				// 	}, "genres")
			}
		}, paths...)
		// get genres
		var genreArray []string
		jsonparser.ArrayEach(resp_body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			genreName, _ := jsonparser.GetString(value, "name")
			genreArray = append(genreArray, genreName)
		}, "genres")
		currentAnime.Genres = genreArray

		// close the response body
		resp.Body.Close()

		animeListTemp = append(animeListTemp, currentAnime)
	} // for

	animeList = animeListTemp
	return animeListTemp
}

func getAnimeByID(id int) (*anime, error) {
	for _, a := range animeList {
		// fmt.Println(a.ID, id)
		// fmt.Println(reflect.TypeOf(a.ID), reflect.TypeOf(id))
		if int(a.ID) == id {
			return &a, nil
		}
	}
	return nil, errors.New("Anime not found")
}