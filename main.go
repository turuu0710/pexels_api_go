package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	PhotoApi = "https://api.pexels.com/v1"
	VideoApi = "https://api.pexels.com/video"
)

type Client struct {
	Token          string
	hc             http.Client
	RemainingTimes int32
}

func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{Token: token, hc: c}
}

type SearchResults struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}
type Photo struct {
	Id              int32       `json: "id"`
	Width           int32       `json: "width"`
	Height          int32       `json: "height"`
	Url             string      `json: "url"`
	Photographer    string      `json: "photographer"`
	PhotographerUrl string      `json: "photographer_url"`
	Src             PhotoSource `json: "src"`
}

type PhotoSource struct {
	Original  string `json: "original"`
	Large     string `json: "large"`
	Large2x   string `json: "large2x"`
	Medium    string `json: "medium"`
	Small     string `json: "small"`
	Portrait  string `json: "portrait"`
	Square    string `json: "square"`
	Landscape string `json: "landscape"`
	Tiny      string `json: "tiny"`
}

func (c *Client) SearchPhotos(query string, perPage, page int) (*SearchResults, error) {
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	//to make the request have to create function c->client request
	//to capture this with resp
	resp, err := c.requestDoWithAuth("GET", url)
	defer resp.Body.Close()
	//defer function ensure that that this function will be closed when this function closes

	//so access to real all function from ioutil package
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	var result SearchResults

	return &result, err
}

func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Autherization", c.Token)
	resp, err := c.hc.Do(req)
	if err != nil {
		return resp, err
	}
	times, err := strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return resp, nil
	} else {
		c.RemainingTimes = int32(times)
	}
	return resp, nil

}

func main() {
	os.Setenv("PexelsToken", "6pxw7mNcAUgjOTgD46hKCWAAGWMudLBcoEMQwvVsxD0TIIdiHVXUOC3Y")
	var TOKEN = os.Getenv("PexelsToken")

	var c = NewClient(TOKEN)

	result, err := c.SearchPhotos("waves")

	if err != nil {
		fmt.Errorf("Search error: %v", err)
	}

	if result.Page == 0 {
		fmt.Errorf("search result wrong")
	}

	fmt.Println(result)
}
