package scraper

import (
	"fmt"
	"log"
	"time"

	"app/config"

	"github.com/gocolly/colly/v2"
)

// Article represents the scraped data structure
type Article struct {
	Title string
	URL   string
}

type Scraper struct {
	TargetURL      string
	TargetDomain   string
	TargetCategory string
	Order          string
	MaxArticles    int
	Results        []Article
}

func NewScraper(cfg *config.Config) *Scraper {
	return &Scraper{
		TargetURL:      cfg.TargetURL,
		TargetDomain:   cfg.TargetDomain,
		TargetCategory: cfg.TargetCategory,
		Order:          cfg.Order,
		MaxArticles:    cfg.MaxArticles,
	}
}

func (s *Scraper) Scrape() {
	c := colly.NewCollector(
		colly.AllowedDomains(s.TargetDomain),
		colly.CacheDir("./cache"),
		colly.MaxDepth(2),
		colly.UserAgent("Sample-Scraper"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  s.TargetDomain,
		Delay:       time.Second,
		RandomDelay: time.Second,
	})

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

	c.OnHTML("a.Button_secondary__cM38g.Button_baseStyle__Vhn6Y.Button_medium__STW9Z.Button_shadow__3xqQY", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatalln("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	url := fmt.Sprintf("%s/topics/%s?order=%s",
		s.TargetURL,
		s.TargetCategory,
		s.Order,
	)
	c.Visit(url)
}
