package linkedin

import (
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/kofalt/go-memoize"
)

func New(browser *rod.Browser) *LinkedIn {
	return &LinkedIn{cache: memoize.NewMemoizer(5*time.Minute, 10*time.Minute), browser: browser}
}

type LinkedIn struct {
	cache   *memoize.Memoizer
	browser *rod.Browser
}

type LinkedInProfile struct {
	Name       string
	Headline   string
	Experience []*LinkedInExperience
	Education  []*LinkedInEducation
}

type LinkedInPosition struct {
	Title       string
	Start       string
	End         string
	Location    string
	Description string
}

type LinkedInExperience struct {
	Company      string
	CompanyImage *string
	Positions    []*LinkedInPosition
}

type LinkedInEducation struct {
	Title string
}

func (r *LinkedIn) retrievePage() (*rod.Page, error) {
	page := r.browser.MustPage("https://linkedin.com/in/mattpitts")
	waitDur, _ := time.ParseDuration("10s")
	err := page.WaitDOMStable(waitDur, .2)
	log.Println("Page Stable")
	return page, err
}

func (r *LinkedIn) getPage() *rod.Page {
	result, err, cached := memoize.Call(r.cache, "page", r.retrievePage)
	if err != nil {
		log.Fatal(err)
	}
	if cached {
		log.Println("Using cached Page")
	}
	return result
}

func (r *LinkedIn) RetrieveProfile() *LinkedInProfile {
	nameElm, err := r.getPage().Element(".top-card-layout h1.top-card-layout__title")
	if err != nil {
		html := r.getPage().MustElement("body")
		log.Fatal("Unable to get initial profile name element from DOM", html, err)
	}
	name := nameElm.MustText()
	headline := r.getPage().MustElement(".top-card-layout h2.top-card-layout__headline").MustText()
	return &LinkedInProfile{Name: name, Headline: headline, Experience: ExtractExperienceList(r.getPage()), Education: ExtractEducationList(r.getPage())}
}

func ExtractExperienceList(page *rod.Page) []*LinkedInExperience {
	expList := page.MustElements("ul.experience__list > li")
	return MapElements(expList, ExtractExperience)
}

func ExtractExperience(element *rod.Element) *LinkedInExperience {
	title := element.MustElement(".experience-item__subtitle").MustText()
	companyImage := element.MustElement("img.profile-section-card__image").MustAttribute("src")
	return &LinkedInExperience{Company: title, CompanyImage: companyImage, Positions: ExtractPositions(element)}
}

func ExtractPositions(element *rod.Element) []*LinkedInPosition {
	class, _ := element.Attribute("class")
	isGroup := strings.Contains(*class, "experience-group")
	if isGroup {
		return MapElements(element.MustElements("li"), ExtractPosition)
	} else {
		return []*LinkedInPosition{ExtractPosition(element)}
	}
}

func ExtractPosition(element *rod.Element) *LinkedInPosition {
	title := element.MustElement(".experience-item__title").MustText()
	metaElements := element.MustElements(".experience-item__meta-item")
	desc := element.MustElement("p.show-more-less-text__text--more").MustText()
	return &LinkedInPosition{Title: title, Description: desc, Location: metaElements[1].MustText()}
}

func ExtractEducationList(page *rod.Page) []*LinkedInEducation {
	items := page.MustElements("ul.education__list > li")
	return MapElements(items, ExtractEducation)
}

func ExtractEducation(element *rod.Element) *LinkedInEducation {
	title := element.MustElement("h3 > a").MustText()
	return &LinkedInEducation{Title: title}
}

func MapElements[V any](ts []*rod.Element, fn func(*rod.Element) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
