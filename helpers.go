package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func combineData(artist Artist, rel Relation) ParsedArtist {
	var parsed ParsedArtist

	parsed.ID = artist.ID
	parsed.Image = artist.Image
	parsed.Name = artist.Name
	parsed.Members = artist.Members
	parsed.CreationDate = artist.CreationDate

	// parse dates
	newDates := map[string]string{}
	for location, dates := range rel.DatesLocations {
		newLoc := beautifyLocation(location)
		for _, v := range dates {
			newDates[beautifyDate(v)] = newLoc
		}
	}
	parsed.DatesLocations = newDates

	// parse first album
	if artist.ID == 1 {
		parsed.FirstAlbum = "July 13, 1973"
	} else {
		fa := strings.Trim(strings.TrimPrefix(artist.FirstAlbum, "*"), " ")
		n, _ := strconv.Atoi(fa[3:5])
		mon := months[(n-1)%12]
		parsed.FirstAlbum = mon + " " + strings.TrimPrefix(fa[:2], "0") + ", " + fa[6:]
	}

	parsed.Slug = encodeURL(artist.Name)

	return parsed
}

func beautifyLocation(location string) string {
	l := strings.SplitN(location, "-", 2)

	for i, word := range l {
		l[i] = strings.Join(strings.Split(word, "_"), " ")
	}

	newLoc := capitalize(l[0])

	if len(l) > 1 {
		newLoc += capitalize(", " + l[1])
	}

	return newLoc
}

func beautifyDate(s string) string {

	s = strings.TrimPrefix(s, " ")
	s = strings.TrimPrefix(s, "*")

	if len(s) != len("dd-mm-yyyy") {
		return "Invalid Date"
	}

	n, _ := strconv.Atoi(s[3:5])

	mon := months[n-1]
	day := strings.TrimPrefix(s[:2], "0")
	year := s[6:]

	newDate := capitalize(mon + " " + day + ", " + year)
	return newDate
}

func capitalize(s string) string {
	words := strings.Split(s, " ")
	for i, word := range words {

		if word == "usa" {
			words[i] = "USA"
			continue
		}

		if word == "uk" {
			words[i] = "UK"
			continue
		}

		if len(word) > 0 && word[0] >= 'a' && word[0] <= 'z' {
			words[i] = string(word[0]+'A'-'a') + word[1:]
		}
	}

	return strings.Join(words, " ")
}

// parse json using get request to url
func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// find artist by slug
func getArtist(s string) (ParsedArtist, bool) {
	for _, artist := range parsed {
		if s == artist.Slug {
			return artist, true
		}
	}

	return ParsedArtist{}, false
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func indexOfInt(element int, data []int) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func seachKeyword(keyword string) SearchData {
	found := SearchData{
		Names:         make(map[string]ParsedArtist),
		Members:       make(map[string]ParsedArtist),
		Locations:     make([]Loc, 0),
		FirstAlbums:   make(map[string]ParsedArtist),
		CreationDates: make(map[int]ParsedArtist),
	}

	keyword = strings.ToLower(keyword)

	for _, artist := range parsed {
		name := strings.ToLower(artist.Name)
		if strings.Index(name, keyword) == 0 {
			found.Names[artist.Name] = artist
		}

		for _, member := range artist.Members {
			if strings.Index(strings.ToLower(member), keyword) == 0 {
				found.Members[member] = artist
			}
		}

		for date, loc := range artist.DatesLocations {
			if strings.Index(strings.ToLower(loc), keyword) == 0 {
				if x := indexOfLocation(found.Locations, loc); x < 0 {
					found.Locations = append(found.Locations, Loc{Location: loc, DateArtists: map[string]ParsedArtist{date: artist}})
				} else {
					found.Locations[x].DateArtists[date] = artist
				}
			}
		}

		f := strings.ToLower(artist.FirstAlbum)
		if strings.Index(f, keyword) == 0 {
			found.FirstAlbums[artist.FirstAlbum] = artist
		}

		c := fmt.Sprint(artist.CreationDate)
		if strings.Index(c, keyword) == 0 {
			found.CreationDates[artist.CreationDate] = artist
		}
	}

	return found
}

func trimWhitespaces(s string) string {
	if len(s) == 0 {
		return s
	}

	if s[0] == ' ' {
		return trimWhitespaces(s[1:])
	}

	if s[len(s)-1] == ' ' {
		return trimWhitespaces(s[:len(s)-1])
	}

	return s
}

func leaveOneSpace(s string) string {
	for i := 1; i < len(s); i++ {
		if s[i-1] == ' ' && s[i] == ' ' {
			s = s[:i-1] + s[i:]
			i--
		}
	}

	return s
}

func indexOfLocation(locations []Loc, locName string) int {
	for i, loc := range locations {
		if loc.Location == locName {
			return i
		}
	}

	return -1
}

func buildSearchData() {
	searchData = SearchData{
		Names:            make(map[string]ParsedArtist),
		Members:          make(map[string]ParsedArtist),
		Locations:        make([]Loc, 0),
		FirstAlbums:      make(map[string]ParsedArtist),
		CreationDates:    make(map[int]ParsedArtist),
		MembersRange:     Range{math.MaxInt32, math.MinInt32},
		FirstAlbumsRange: Range{math.MaxInt32, math.MinInt32},
		CreationRange:    Range{math.MaxInt32, math.MinInt32},
	}

	for _, artist := range parsed {
		// band/artist name
		searchData.Names[artist.Name] = artist

		// members
		for _, member := range artist.Members {
			searchData.Members[member] = artist
		}

		for date, loc := range artist.DatesLocations {
			if x := indexOfLocation(searchData.Locations, loc); x < 0 {
				searchData.Locations = append(searchData.Locations, Loc{Location: loc, DateArtists: map[string]ParsedArtist{date: artist}})
			}
		}

		// first album dates
		searchData.FirstAlbums[artist.FirstAlbum] = artist

		// creation dates
		searchData.CreationDates[artist.CreationDate] = artist
	}

	searchData.Countries = getCountries(searchData.Locations)

	for _, artist := range parsed {
		m := len(artist.Members)
		if m < searchData.MembersRange.Min {
			searchData.MembersRange.Min = m
		} else if m > searchData.MembersRange.Max {
			searchData.MembersRange.Max = m
		}
	}

	for date := range searchData.FirstAlbums {
		year, err := getYear(date)
		if err == nil {
			if year < searchData.FirstAlbumsRange.Min {
				searchData.FirstAlbumsRange.Min = year
			}
			if year > searchData.FirstAlbumsRange.Max {
				searchData.FirstAlbumsRange.Max = year
			}
		}
	}

	for year := range searchData.CreationDates {
		if year < searchData.CreationRange.Min {
			searchData.CreationRange.Min = year
		}
		if year > searchData.CreationRange.Max {
			searchData.CreationRange.Max = year
		}
	}
}

func getYear(date string) (int, error) {
	if len(date) >= 4 {
		return strconv.Atoi(date[len(date)-4:])
	}

	return 0, errors.New("too short input")
}

func encodeURL(s string) string {
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ToLower(s)
	return s
}

func getCountries(locs []Loc) []string {
	res := []string{}
	for _, loc := range locs {
		country := strings.Split(loc.Location, ", ")[1]
		if indexOf(country, res) < 0 {
			res = append(res, country)
		}
	}
	return res
}
