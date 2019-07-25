package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"time"
)

func writeHTML(html string, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, html)
}

func webpage(w http.ResponseWriter, r *http.Request) {
	var url, element, html string
	var quality, randerWaitTime = int64(90), int64(0)
	r.ParseForm()

	query := r.URL.Query()
	if len(query["html"]) > 0 {
		html = query["html"][0]
	} else if len(r.PostFormValue("html")) > 0 {
		html = r.PostFormValue("html")
	}
	if len(query["url"]) > 0 {
		url = query["url"][0]
	} else if len(html) > 0 {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, html)
		}))
		defer ts.Close()
		url = ts.URL
	} else {
		fmt.Fprint(w, "缺少参数url或html")
		return
	}
	if len(query["quality"]) > 0 {
		quality, _ = strconv.ParseInt(query["quality"][0], 10, 64)
	}
	if len(query["element"]) > 0 {
		element = query["element"][0]
	}
	if len(query["rander_wait_time"]) > 0 {
		randerWaitTime, _ = strconv.ParseInt(query["rander_wait_time"][0], 10, 64)
	}
	defer r.Body.Close()
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf []byte

	if element == "" {
		if err := chromedp.Run(ctx, fullScreenshot(url, quality, randerWaitTime, &buf)); err != nil {
			fmt.Fprint(w, err)
			return
		}
	} else {
		if err := chromedp.Run(ctx, elementScreenshot(url, element, randerWaitTime, &buf)); err != nil {
			fmt.Fprint(w, err)
			return
		}
	}

	var s = string(buf)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(s)))
	fmt.Fprint(w, s)
}

func echarts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var pageurl string
	var randerWaitTime = int64(0)
	element := "#main"

	params := url.Values{}

	if len(query["config"]) > 0 {
		params.Add("config", query["config"][0])
	} else {
		fmt.Fprint(w, "参数错误")
		return
	}
	pageurl = "http://localhost:80/echarts/?" + params.Encode()
	if len(query["rander_wait_time"]) > 0 {
		randerWaitTime, _ = strconv.ParseInt(query["rander_wait_time"][0], 10, 64)
	}
	defer r.Body.Close()
	draw(pageurl, element, 0, randerWaitTime, w)
}

func draw(url string, element string, quality int64, randerWaitTime int64, w http.ResponseWriter) {
	fmt.Println(url)
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf []byte

	if element == "" {
		if quality == 0 {
			quality = 90
		}
		if err := chromedp.Run(ctx, fullScreenshot(url, quality, randerWaitTime, &buf)); err != nil {
			fmt.Fprint(w, err)
			return
		}
	} else {
		if err := chromedp.Run(ctx, elementScreenshot(url, "#main", randerWaitTime, &buf)); err != nil {
			fmt.Fprint(w, err)
			return
		}
	}

	var s = string(buf)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(s)))
	fmt.Fprint(w, s)
}

func elementScreenshot(urlstr, sel string, randerWaitTime int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(sel, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if randerWaitTime != 0 {
				time.Sleep(time.Duration(randerWaitTime) * time.Second)
			}
			return nil
		}),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByID),
	}
}

func fullScreenshot(urlstr string, quality int64, randerWaitTime int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			if randerWaitTime != 0 {
				time.Sleep(time.Duration(randerWaitTime) * time.Second)
			}
			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

func main() {
	pathPrefix := "/echarts/"
	staticDir := "./echarts"
	http.Handle(pathPrefix, http.StripPrefix(pathPrefix, http.FileServer(http.Dir(staticDir))))

	http.HandleFunc("/echarts/", echarts)
	http.HandleFunc("/webpage/", webpage)
	fmt.Println("Server is at localhost:8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}
