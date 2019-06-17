package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

/**
用户身份
*/
var _ga string
var _gid string
var GCID string
var GCESS string

/**
课程id
*/
var cid string

var hostname = "https://time.geekbang.org"
var cookie = "_ga=%s; _gid=%s; GCID=%s; GCESS=%s"

func main() {

	loadConfig()

	articles := getArticles()

	articlesLen := len(articles)
	for i := 0; i < articlesLen; i++ {
		download(articles[i])
	}

}

func getArticles() []Article {
	client := http.Client{}
	reqBody := fmt.Sprintf("{\"cid\":\"%s\",\"size\":200,\"prev\":0,\"order\":\"earliest\",\"sample\":true}", cid)

	req, _ := http.NewRequest(http.MethodPost, hostname+"/serv/v1/column/articles", strings.NewReader(reqBody))

	header := req.Header

	header.Set("Origin", hostname)
	header.Set("Referer", hostname)
	header.Set("User-Agent", randomUserAgent())
	header.Set("X-Real-IP", randomIpAddress())
	header.Set("Cookie", fmt.Sprintf(cookie, _ga, _gid, GCID, GCESS))
	header.Set("Connection", "keep-alive")
	header.Set("Content-Type", "application/json")

	response, _ := client.Do(req)

	dataByte, _ := ioutil.ReadAll(response.Body)

	dataStr := string(dataByte)
	if len(dataStr) == 0 {
		panic("cookie is not valid! please check your config")
	}

	articlesMap := make(map[string]interface{})
	jsonErr := json.Unmarshal(dataByte, &articlesMap)
	if jsonErr != nil {
		panic(jsonErr)
	}

	var videos = articlesMap["data"].(map[string]interface{})["list"].([]interface{})

	var articles []Article
	for _, video := range videos {
		var url = video.(map[string]interface{})["video_media_map"].(map[string]interface{})["hd"].(map[string]interface{})["url"].(string)
		var title = strings.ReplaceAll(strings.ReplaceAll(video.(map[string]interface{})["article_title"].(string), " ", ""), "|", "")
		article := Article{url, title}
		articles = append(articles, article)
	}
	return articles
}

func download(article Article) {
	path, _ := os.Getwd()
	articleSavePath := path + "/" + article.title + ".mp4"
	file, _ := os.Stat(articleSavePath) //os.Stat获取文件信息
	if file != nil {
		fmt.Println("file exists : " + article.title)
		return
	}

	fmt.Println("start download : ", article.title)
	cmd := exec.Command("ffmpeg", "-i", article.url, "-c", "copy", "-bsf:a", "aac_adtstoasc", articleSavePath)

	fmt.Println(cmd.Args)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		os.Remove(articleSavePath)
		panic("download failed: " + err.Error() + ", " + article.url + ", " + article.title)
	}

	fmt.Println("end download : ", article.title)
}

var userAgentList = [13]string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
	"Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 6 Build/LYZ28E) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Mobile/14F89;GameHelper",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.4",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:46.0) Gecko/20100101 Firefox/46.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/13.10586",
	"Mozilla/5.0 (iPad; CPU OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1",
}

func randomUserAgent() string {
	r := rand.Intn(13)
	return userAgentList[r]
}

func randomIpAddress() string {
	r := rand.Intn(254)
	return fmt.Sprintf("211.161.244.%d", r)
}

func loadConfig() {
	f, err := os.Open("./config.json")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	config := make(map[string]string)

	dataByte, _ := ioutil.ReadAll(f)

	jsonErr := json.Unmarshal(dataByte, &config)
	if jsonErr != nil {
		panic(jsonErr)
	}

	_ga = config["_ga"]
	if _ga == "" {
		panic("config['_ga'] can not be null")
	}
	_gid = config["_gid"]
	if _gid == "" {
		panic("config['_gid'] can not be null")
	}
	GCID = config["GCID"]
	if GCID == "" {
		panic("config['GCID'] can not be null")
	}
	GCESS = config["GCESS"]
	if GCESS == "" {
		panic("config['GCESS'] can not be null")
	}
	cid = config["cid"]
	if cid == "" {
		panic("config['cid'] can not be null")
	}
}

type Article struct {
	url   string
	title string
}
