package main

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"fmt"
	"encoding/json"
)

type CnBlog []struct {
	ID           int    `json:"Id"`
	Title        string `json:"Title"`
	URL          string `json:"Url"`
	Description  string `json:"Description"`
	Author       string `json:"Author"`
	BlogApp      string `json:"BlogApp"`
	Avatar       string `json:"Avatar"`
	PostDate     string `json:"PostDate"`
	ViewCount    int    `json:"ViewCount"`
	CommentCount int    `json:"CommentCount"`
	DiggCount    int    `json:"DiggCount"`
}

func ExampleScrape() {
	client := &http.Client{
	}
	// Request the HTML page.
	postValues := url.Values{}
	postValues.Add("client_id", "88cfd065-5321-47f5-bbaa-46ece97e7e0f")
	postValues.Add("client_secret", "Wua3iSHwihoRNvzh2oSPVIfpUppqXPRwm6JtuqPySo8x604NXHbpEIajQOwGcrOni3XdxxK8DJbSh5lt")
	postValues.Add("grant_type", "client_credentials")
	res, _ := http.PostForm("https://oauth.cnblogs.com/connect/token", postValues)
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	fmt.Println(string(data))

	byt := []byte(string(data))
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	accessToken := dat["access_token"].(string)
	fmt.Println(accessToken)

	//start to get BlogList
	req, _ := http.NewRequest("GET", "https://api.cnblogs.com/api/blogs/w1570631036/posts?pageIndex=1", nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, _ := client.Do(req)
	blogList, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Print(string(blogList))

	//

}


func main() {
	ExampleScrape()
}
