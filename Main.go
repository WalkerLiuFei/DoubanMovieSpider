package main

import (
	"./dao"
	"./myconstant"
	"./utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hu17889/go_spider/core/common/page"
	"github.com/hu17889/go_spider/core/common/request"
	"github.com/hu17889/go_spider/core/pipeline"
	"github.com/hu17889/go_spider/core/spider"
	"net/http"
	"strings"
	"log"
	"strconv"
)

type CurrentRequest struct{
	url string
	currentPage int
}
var currentRequest CurrentRequest
var DoubanSpider *spider.Spider

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

// Parse html dom here and record the parse result that we want to crawl.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.Errormsg())
		return
	}

	q := p.GetHtmlParser()

	subjects := q.Find("div.doulist-subject")
	fmt.Println("\n how much items ? :" + strconv.Itoa(len(subjects.Nodes)))

	if len(subjects.Nodes) > 0 { //如果有匹配的模式,就进处理
		subjects.Each(func(i int, element *goquery.Selection) {
			link, _ := element.Find("div.title a").Attr("href")
			postImg, _ := element.Find("div.post a img").Attr("src")
			ratings := element.Find("div.rating span")
			rating := ratings.Eq(1).Text()
			judgeNum := ratings.Eq(2).Text()
			sourceAbstracts := strings.Split(element.Find("div.abstract").Text(), "\n")
			afterDeakAbstract := make([]string, 0)

			for _, str := range sourceAbstracts {
				str = strings.TrimSpace(str)
				if len(str) > 1 {
					str = strings.Split(str, ":")[1]
					afterDeakAbstract = append(afterDeakAbstract, str)
				}
			}
			m := dao.Movie{
				Name:        strings.TrimSpace(element.Find("div.title a").Text()),
				Link:        strings.TrimSpace(link),
				PostImgLink: strings.TrimSpace(postImg),
				Score:       strings.TrimSpace(rating),
				JudgeNumber: strings.Trim(strings.Trim(judgeNum, "("), ")"),
				Director:    afterDeakAbstract[0],
				LeadingStar: afterDeakAbstract[1],
				MovieType:   afterDeakAbstract[2],
				Country:     afterDeakAbstract[3],
				PostYear:    afterDeakAbstract[4],
			}
			m.ID = Utils.MD5(m.PostImgLink)

			dao.WriteMovie(m)
		})
		currentRequest.currentPage += 25
		dealCurrentRequest()
	}else{
		//TODO: NOTHING....
	}

}
func (this *MyPageProcesser) Finish() {
	fmt.Printf("TODO:before end spider \r\n")
}

const (
	up     = "https://www.douban.com/doulist/240962/"  //上
	middle = "https://www.douban.com/doulist/243559/" //中
	foot   = "https://www.douban.com/doulist/248893/"  //下
)



func main() {

	DoubanSpider = spider.NewSpider(NewMyPageProcesser(), "test")
	DoubanSpider.AddPipeline(pipeline.NewPipelineConsole())
	urls := []string{
		up, middle, foot,
	}
	dao.CheckTableAndColumn("root", "3252860")
	for index, item := range urls {
		currentRequest = CurrentRequest{
			url:item,
			currentPage: 1,
		}
		dealCurrentRequest()
		log.Printf("进行第%d个",index+1)
	}

}
func dealCurrentRequest() {
	//requests := make([]*request.Request,0) //
	rawUrl := currentRequest.url + "?start=" + strconv.Itoa(currentRequest.currentPage)
	strings.TrimSpace(rawUrl)
	headers := make(http.Header)
	headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/57.0.2987.98 Chrome/57.0.2987.98 Safari/537.36")
	headers.Set("Cookie",
		"bid=CkhgPN5OMCg; __yadk_uid=WBXysLWOtt7HNArIetA72lvQn8qMltk8; ct=y; ll=\"118172\"; gr_user_id=13ed0786-79f7-4f40-be3f-fbc80d29d9fb; ap=1; viewed=\"25779298_25900156_25742200_26327897_1230036\"; _vwo_uuid_v2=3FCC592C6144551DEFD7378DBDBF32B4|184602a4eaf4b8eb607a1e245ad4f33f; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1494484687%2C%22https%3A%2F%2Fwww.google.nl%2F%22%5D; _pk_id.100001.8cb4=9b1613a79a527a52.1493702857.13.1494487624.1494482273.; _pk_ses.100001.8cb4=*; __utmt=1; __utma=30149280.1043816372.1493702870.1494482274.1494484687.16; __utmb=30149280.5.10.1494484687; __utmc=30149280; __utmz=30149280.1494484687.16.13.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided)")
	re := request.NewRequest(rawUrl, myconstant.HTML,
		"", myconstant.GET,
		"", headers, nil, nil, nil)
	pageItem := DoubanSpider.SetThreadnum(1).
		SetSleepTime("fixed", 1000, 0).
		GetByRequest(re)
	if (pageItem != nil){
		request := pageItem.GetRequest()
		if (request != nil ){
			fmt.Println("rawUrl\t:\t" + request.Url)
		}
	}

}
