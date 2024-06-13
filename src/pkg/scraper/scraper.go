package scraper

import (
	"fmt"
	"log"
	"time"

	"app/config"

	"github.com/gocolly/colly/v2"
)

// スクレイピングした記事の情報を保持する構造体
type Article struct {
	Title string
	URL   string
}

// スクレイパーの構造体
type Scraper struct {
	TargetURL    string
	TargetDomain string
	Categories   []string
	Order        string
	MaxArticles  int
	Results      []Article
}

// スクレイパーの作成
func NewScraper(cfg *config.Config) *Scraper {
	return &Scraper{
		TargetURL:    cfg.TargetURL,
		TargetDomain: cfg.TargetDomain,
		Categories:   cfg.Categories,
		Order:        cfg.Order,
		MaxArticles:  cfg.MaxArticles,
	}
}

// スクレイピングの実行
func (s *Scraper) Scrape(category string) {

	// スクレイピングの設定
	c := colly.NewCollector(
		colly.AllowedDomains(s.TargetDomain),
		colly.CacheDir("./cache"),
		colly.MaxDepth(2),
		colly.UserAgent("Sample-Scraper"),
	)

	// 対象サイトにかかる負荷を配慮し
	// リクエスト間で1~2秒の時間を空ける
	c.Limit(&colly.LimitRule{
		DomainGlob:  s.TargetDomain,
		Delay:       time.Second,
		RandomDelay: time.Second,
	})

	//  OnHTML 要素発見時に何を実施するかあらかじめ指定する
	// article 要素が見つかったら実行される関数
	count := 0
	c.OnHTML("article", func(e *colly.HTMLElement) {
		if count < s.MaxArticles {
			title := e.DOM.Find("h2").Text()
			href, ok := e.DOM.Find("a").Attr("href")
			if !ok {
				href = "Not Link!"
			}
			article := Article{
				Title: title,
				URL:   s.TargetURL + href,
			}
			s.Results = append(s.Results, article)
			log.Printf("Title: %s, URL: %s", article.Title, article.URL)
			count++
		}
	})

	// エラー発生時に実行される関数
	c.OnError(func(r *colly.Response, err error) {
		log.Fatalln("Request URL:", r.Request.URL, "が失敗しました。:", r, "\nError:", err)
	})

	// 指定した URL にアクセスしてスクレイピングの開始
	url := fmt.Sprintf("%s/topics/%s?order=%s",
		s.TargetURL,
		category,
		s.Order,
	)
	c.Visit(url)
}
