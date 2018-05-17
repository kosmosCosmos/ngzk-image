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
	for id := 1; id < 2; id++ {
		url, _ := content.Find(".unit2 > a:nth-child(3)").Attr("href")
		idolName := strings.Replace(url, "./", "", -1)
		x := time.Now()
		h, _ := time.ParseDuration("-1h")
		for i := 0; i < 24; i++ {
			x = x.Add(24 * 30 * h)
			y := strings.Split(x.Format("2006-01-02"), "-")[0] + strings.Split(x.Format("2006-01-02"), "-")[1]
			blog := "http://blog.nogizaka46.com/" + idolName + "/?d=" + y
			req.SetRequestURI(blog)
			fasthttp.Do(req, resp)
			reg := regexp.MustCompile(`http...img.nogizaka46.com.blog.*?img.\d{4}.\d{2}.\d{2}.\d{7}.\d{4}.\w{3,4}`)
			imageurls := reg.FindAllString(string(resp.Body()), -1)
			for _, x := range imageurls {
				work(x, url, i)
			}
		}
	}
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