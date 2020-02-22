package main

var months = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

var customImages = map[int]string{1: "https://cdn.vox-cdn.com/thumbor/XxKTJTVlOD3cgza6hL5TgYqe2kU=/0x0:1511x1500/1200x0/filters:focal(0x0:1511x1500):no_upscale()/cdn.vox-cdn.com/uploads/chorus_asset/file/13355413/Queen_II.jpg"}

// Artist struct
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// Relation q
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Relations q
type Relations struct {
	Index []Relation `json:"index"`
}

// ParsedArtist for parsed data
type ParsedArtist struct {
	ID             int               `json:"id"`
	Image          string            `json:"image"`
	Name           string            `json:"name"`
	Members        []string          `json:"members"`
	CreationDate   int               `json:"creationDate"`
	FirstAlbum     string            `json:"firstAlbum"`
	DatesLocations map[string]string `json:"datesLocations"`
	Slug           string            `json:"slug"`
}

// SearchData arrays of search data
type SearchData struct {
	Names            map[string]ParsedArtist
	Members          map[string]ParsedArtist
	Locations        []Loc
	FirstAlbums      map[string]ParsedArtist
	CreationDates    map[int]ParsedArtist
	Keyword          string
	Countries        []string
	MembersRange     Range
	FirstAlbumsRange Range
	CreationRange    Range
}

// Range q
type Range struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// Loc q
type Loc struct {
	Location    string
	DateArtists map[string]ParsedArtist
}

// HomeData send for home page
type HomeData struct {
	Artists    []ParsedArtist
	SearchData SearchData
}

// FilterInput q
type FilterInput struct {
	Countries        []string `json:"countries"`
	CreationRange    Range    `json:"creationRange"`
	MembersRange     Range    `json:"membersRange"`
	FirstAlbumsRange Range    `json:"firstAlbumsRange"`
}

// const apikey = "pk.eyJ1IjoicnRhYnVsb3YiLCJhIjoiY2s2dDRpZTZjMDNyczNlbXlubDN2cHBlaiJ9.YkwacBJwhSZRMt7ycbjzLQ"

// var mapp = "https://api.mapbox.com/styles/v1/mapbox/dark-v10/static/0,0,1/1280x750@2x?access_token=pk.eyJ1IjoicnRhYnVsb3YiLCJhIjoiY2s2dDRpZTZjMDNyczNlbXlubDN2cHBlaiJ9.YkwacBJwhSZRMt7ycbjzLQ"

// var marker2 = "https://api.mapbox.com/styles/v1/mapbox/dark-v10/static/url-https%3A%2F%2Fwww.mapbox.com%2Fimg%2Frocket.png(-46.6334,-23.5507)(-120,37)/0,0,1/1280x750?access_token=pk.eyJ1IjoicnRhYnVsb3YiLCJhIjoiY2s2dDRpZTZjMDNyczNlbXlubDN2cHBlaiJ9.YkwacBJwhSZRMt7ycbjzLQ"

// var marker = "https://www.mapbox.com/img/rocket.png"

// var loc = "https://api.mapbox.com/geocoding/v5/mapbox.places/Sao Paulo, Brazil.json?access_token=pk.eyJ1IjoicnRhYnVsb3YiLCJhIjoiY2s2dDRpZTZjMDNyczNlbXlubDN2cHBlaiJ9.YkwacBJwhSZRMt7ycbjzLQ"
