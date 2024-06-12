package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	targetDomain = "zenn.dev"
	targetURL    = "https://" + targetDomain
)

func main() {
	c := colly.NewCollector(
		// Zenn 以外のアクセスを許可しない
		colly.AllowedDomains(targetDomain),
		// ./cache でレスポンスをキャッシュする
		colly.CacheDir("./cache"),
		// アクセスするページの再帰の深さを設定
		colly.MaxDepth(2),
		// ユーザーエージェントを設定
		colly.UserAgent("Sample-Scraper"),
	)

	// リクエスト間で1~2秒の時間を空ける
	c.Limit(&colly.LimitRule{
		DomainGlob:  targetDomain,
		Delay:       time.Second,
		RandomDelay: time.Second,
	})

	// 全ての article 要素に対して実行される関数
	c.OnHTML("article", func(e *colly.HTMLElement) {
		title := e.DOM.Find("h2").Text()
		href, ok := e.DOM.Find("a").Attr("href")
		if !ok {
			href = "Not Link!"
		}
		log.Printf("Title: %s, URL: %s%s", title, targetURL, href)
	})

	c.OnHTML("a.Button_secondary__cM38g.Button_baseStyle__Vhn6Y.Button_medium__STW9Z.Button_shadow__3xqQY", func(e *colly.HTMLElement) {
		// 次のページに遷移する
		e.Request.Visit(e.Attr("href"))
	})

	// エラー発生時に実行される関数
	c.OnError(func(r *colly.Response, err error) {
		log.Fatalln("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	url := fmt.Sprintf("%s/topics/go?order=latest", targetURL)
	c.Visit(url)
}
