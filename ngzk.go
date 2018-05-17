package main

import (
	"github.com/valyala/fasthttp"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"regexp"
	"strconv"
	"time"
	"os"
	"fmt"
)

func main() {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI("http://blog.nogizaka46.com/")
	fasthttp.Do(req, resp)
	content, _ := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
	for id := 1; id <33; id++ {
		url, _ := content.Find("div.unit:nth-child("+strconv.Itoa(id)+") > a:nth-child(1)").Attr("href")
		idolName:=strings.Replace(url, "./", "", -1)
/*		req.SetRequestURI(strings.Replace(url, "./", "http://www.nogizaka46.com/member/", -1))
		fmt.Println(strings.Replace(url, "./", "http://www.nogizaka46.com/member/", -1))
		fasthttp.Do(req, resp)
		Blogcontent, _ := goquery.NewDocumentFromReader(strings.NewReader(string(resp.Body())))
		blogUrl, _ := Blogcontent.Find("li.clearfix:nth-child(1) > span:nth-child(2) > span:nth-child(3) > strong:nth-child(1) > a:nth-child(1)").Attr("href")
		if blogUrl != "" {*/
			//idolName := strings.Split(blogUrl, "/")[3]
			x := time.Now()
			h, _ := time.ParseDuration("-1h")
			for i := 0; i < 24; i++ {
				x = x.Add(24 * 30 * h)
				y := strings.Split(x.Format("2006-01-02"), "-")[0] + strings.Split(x.Format("2006-01-02"), "-")[1]
				blog := "http://blog.nogizaka46.com/" + idolName + "/?d=" + y
				fmt.Println(blog)
				req.SetRequestURI(blog)
				fasthttp.Do(req, resp)
				reg := regexp.MustCompile(`http...img.nogizaka46.com.blog.*?img.*?.jpeg`)
				reg2 := regexp.MustCompile(`http...img.nogizaka46.com.blog.*?img.*?.jpg`)
				imageurls := reg.FindAllString(string(resp.Body()), -1)
				imageurls2 := reg2.FindAllString(string(resp.Body()), -1)
				imageurls = append(imageurls, imageurls2...)
				for i, x := range imageurls {
					if len(x) > 75 {
						reg3 := regexp.MustCompile(`.http...img.nogizaka46.com.blog.*?img.*?.jpeg`)
						y := strings.Replace(reg3.FindString(x), "src=", "", -1)
						y = strings.Replace(y, `"`, "", -1)
						if len(y) < 100 {
							work(y, url, i)
						}
					} else  {
						work(x, url, i)
					}
				}
			}
		}
	//}
}

func work(x, url string, i int) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	list := strings.Split(x, "/")
	if len(list) > 9 {
		fmt.Println(x)
		tm := strings.Join([]string{strings.Split(x, "/")[6] + strings.Split(x, "/")[7] + strings.Split(x, "/")[8]}, "-")
		_, err := os.Stat(strings.Replace(strings.Replace(url, "./detail/", "", -1), ".php", "", -1) + "/" + tm)
		if err != nil {
			os.MkdirAll(strings.Replace(strings.Replace(url, "./detail/", "", -1), ".php", "", -1)+"/"+tm, os.ModePerm)
		}
		req.SetRequestURI(x)
		fasthttp.Do(req, resp)
		if len(resp.Body()) < 1000 {
			x = strings.Replace(x, "img", "blog", 1)
			x = strings.Replace(x, "blog/", "", 1)
			req.SetRequestURI(x)
			fasthttp.Do(req, resp)
		}
		fd, _ := os.OpenFile(strings.Replace(strings.Replace(url, "./detail/", "", -1), ".php", "", -1)+"/"+tm+"/"+strconv.Itoa(i)+".jpg", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		buf := []byte(resp.Body())
		fd.Write(buf)
		fd.Close()
	}
}