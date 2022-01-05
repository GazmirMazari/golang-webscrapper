package main

import (
	"math/rand"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var googleDomain = map[string]string{

}

type SearchResult struct {
	ResultRank  int
	Resulturl   string
	ResultTitle string
	ResultDesc  string
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

func randomUserAgent() string {
	randNum := rand.Int() % len(userAgents)

	return userAgents[randNum]
}

func buildGoogleUrls(searchTerm, countryCode, languageCode string, pages, count int)([]string, error){
	toScrape := []string{}
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	if googleBase, found := googleDomains[countryCode]; found{
		for i :=0; i<pages ; i++{
			start := i*count
			scrapeURL := fmt.Sprintf("%s%s&num=%d&hl=%s&start=%d&filter=0",googleBase, searchTerm, count, languageCode, start)
			toScrape = append(toScrape, scrapeURL)
		}
	}else{
		err := fmt.Errorf("country (%s) is currently not supported", countryCode)
		return nil, err
	}
	return toScrape, nil

}

func scrapeClientRequest(searchUrl string, proxyString interface{}) (*http.Response, error) {
	baseClient := getScrapeBaseClient(proxyString)
	req, _ = http.NewRequest("GET", searchUrl, nil)
	req.Header.Set("User-Agent", randomUserAgent)
	res, err := baseClient.Do(req)
	if res.statusCode != 200 {
		err := fmt.Errorf("Scrapper received a non -200 status code")
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

func googleResultParsing(response *http.Response, rank int)([]SearchResult, error){
doc, err := goquery.NewDocumentFromResponse(response)

if err !=nil {
	return nil, err
}

results := []SearchResult{}
sel := doc.Find("div.g")
rank ++
for i := range sel.Nodes{
	item := sel.Eq(i)
	linkTag := item.Find("a")
	link,_ := linkTag.Attr("href")
	titleTag := item.Find("h3.r")
	descTag := item.Find("span.st")
	desc := descTag.Text()
	title := titleTag.Text()
	link = strings.Trim(link, " ")

	if link != "" && link !="#" && !strings.HasPrefix(link,"/"){
		result := SearchResult{
			rank,
			link,
			title,
			desc,
		}
		results = append(results, result)
		rank ++
	}
}
return results, err

}

func getScrapeClient(proxyString interface{}) (*http.Client) {

switch v :=proxyString.(type){

case string: 
	proxyURL, _ :=url.Parse(v)
	return &http.Client{Transport : &http.Transport{Proxy: http.Proxy(proxyUrl)}}
	default: 
		return &http.Client{}
	}
}

func GoogleScrape(searchTerm, countryCode, languageCode string,  pages, count ) ([]SearchResult, err) {
	results := []SearchResult{}
	resultConter := 0
	googlePages, err := buildGoogleUrls(searchTerm, countryCode, languageCode, pages, count)
	if err !=nil {
		return nil, err
	}

	for _, page: range googlePages {
		res, err := scrapeClientRequest(page, proxyString)
		if err != nil{
			return nil, err
		}
		data, err := googleResultParsing(res, resultConter)
		if err != nil {
			return nil, err
		}
		resultCounter += len(data)
		for _, result := range data {
			results = append(results, result)
		}
		time.sleep(time.Duration(backoff) * time.Second)
	}
}

func main() {
	res, err :=GoogleScrape("Gazmir Mazari",  "com", "en", 1, 30)
	if err == nil {
		for_, res := range res{
			fmt.Println(res)
		}
	}
}



