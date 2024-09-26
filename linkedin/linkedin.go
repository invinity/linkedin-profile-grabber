package linkedin

import (
	"log"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func New() *LinkedIn {
	return &LinkedIn{}
}

type LinkedIn struct {
}

type LinkedInProfile struct {
	Experience []*LinkedInExperience
}

type LinkedInPosition struct {
	Title       string
	Start       string
	End         string
	Location    string
	Description string
}

type LinkedInExperience struct {
	Company   string
	Positions []*LinkedInPosition
}

func (r *LinkedIn) RetrieveProfile() *LinkedInProfile {
	// Launch a new browser with default options, and connect to it.
	browser := rod.New().ControlURL(launcher.New().Leakless(false).MustLaunch()).Trace(true).MustConnect()

	// Even you forget to close, rod will close it after main process ends.
	defer browser.MustClose()

	// Create a new page
	log.Println("Retrieving https://linkedin.com/in/mattpitts")
	page := browser.MustPage("https://linkedin.com/in/mattpitts")
	waitDur, _ := time.ParseDuration("5s")
	page.WaitDOMStable(waitDur, 0.7)
	log.Println("Page Stable")
	expList := page.MustElements("ul.experience__list > li")
	log.Println("Got experience list")

	return &LinkedInProfile{Experience: MapElements(expList, ExtractExperience)}
}

func ExtractExperience(element *rod.Element) *LinkedInExperience {
	title := element.MustElement(".experience-item__title")
	positionElements := element.MustElements("li")
	return &LinkedInExperience{Company: title, Positions: MapElements(positionElements, ExtractPosition)}
}

func ExtractPosition(element *rod.Element) *LinkedInPosition {
	title := element.MustElement(".experience-item__title").MustText()
	metaElements := element.MustElements(".experience-item__meta-item")
	desc := element.MustElement("p.show-more-less-text__text--more").MustText()
	return &LinkedInPosition{Title: title, Description: desc, Location: metaElements[0].MustText()}
}

func MapElements[V any](ts rod.Elements, fn func(*rod.Element) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
