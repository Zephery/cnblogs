package main

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"encoding/json"
	"strconv"
	"time"
)

type CnBlog struct {
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
	//fmt.Print(string(blogList))

	//https://blog.csdn.net/u011304970/article/details/70769949

	var dbgInfos []CnBlog
	json.Unmarshal(blogList, &dbgInfos)
	fmt.Print(dbgInfos)

	db, err := sql.Open("mysql", "root:helloroot@tcp(119.23.46.71:3306)/myblog?charset=utf8")
	if err != nil {
		fmt.Print(err)
	}
	for _, value := range dbgInfos {
		fmt.Print(value)
		var id = value.ID
		var Title = value.Title
		var PostDate = value.PostDate
		fmt.Print(id, Title, PostDate)
		getUrl := "https://api.cnblogs.com/api/blogposts/" + strconv.Itoa(id) + "/body"
		reqest, _ := http.NewRequest("GET", getUrl, nil)
		reqest.Header.Add("Authorization", "Bearer "+accessToken)
		response, _ := client.Do(reqest)
		temp, _ := ioutil.ReadAll(response.Body)
		blog := string(temp)
		fmt.Print(blog)

		stmt, _ := db.Prepare("insert sync_blog set cn_blogs_id=?,blog_id=?,create_time=?")
		res, _ := stmt.Exec(id, 23, time.Now())
		fmt.Println(res)
	}

}

func main() {
	ExampleScrape()
}
