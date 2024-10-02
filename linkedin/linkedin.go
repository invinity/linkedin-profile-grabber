package linkedin

import (
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func New(browser *rod.Browser) *LinkedIn {
	return &LinkedIn{browser: browser}
}

type LinkedIn struct {
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

func (r *LinkedIn) getPage() (*rod.Page, error) {
	page, err := r.browser.Page(proto.TargetCreateTarget{URL: "https://linkedin.com/in/mattpitts"})
	if err != nil {
		return nil, err
	}
	waitDur, _ := time.ParseDuration("5s")
	err = page.WaitDOMStable(waitDur, .2)
	if err != nil {
		return nil, err
	}
	log.Println("Page Stable")
	return page, nil
}

func (r *LinkedIn) RetrieveProfile() (*LinkedInProfile, error) {
	page, err := r.getPage()
	if err != nil {
		return nil, err
	}
	_, err = page.Element("body")
	if err != nil {
		return nil, err
	}
	nameElm, err := page.Element(".top-card-layout h1.top-card-layout__title")
	if err != nil {
		return nil, err
	}
	name, err := nameElm.Text()
	if err != nil {
		return nil, err
	}
	headlineElm, err := page.Element(".top-card-layout h2.top-card-layout__headline")
	if err != nil {
		return nil, err
	}
	headline, err := headlineElm.Text()
	if err != nil {
		return nil, err
	}
	experience, err := ExtractExperienceList(page)
	if err != nil {
		return nil, err
	}
	education, err := ExtractEducationList(page)
	if err != nil {
		return nil, err
	}
	return &LinkedInProfile{Name: name, Headline: headline, Experience: experience, Education: education}, nil
}

func ExtractExperienceList(page *rod.Page) ([]*LinkedInExperience, error) {
	expList, err := page.Elements("ul.experience__list > li")
	if err != nil {
		return nil, err
	}
	return MapElements(expList, ExtractExperience)
}

func ExtractExperience(element *rod.Element) (*LinkedInExperience, error) {
	titleE, err := element.Element(".experience-item__subtitle")
	if err != nil {
		return nil, err
	}
	title, err := titleE.Text()
	if err != nil {
		return nil, err
	}
	companyImageE, err := element.Element("img.profile-section-card__image")
	if err != nil {
		return nil, err
	}
	companyImage, err := companyImageE.Attribute("src")
	if err != nil {
		return nil, err
	}
	positions, err := ExtractPositions(element)
	if err != nil {
		return nil, err
	}
	return &LinkedInExperience{Company: title, CompanyImage: companyImage, Positions: positions}, nil
}

func ExtractPositions(element *rod.Element) ([]*LinkedInPosition, error) {
	class, err := element.Attribute("class")
	if err != nil {
		return nil, err
	}
	isGroup := strings.Contains(*class, "experience-group")
	if isGroup {
		elements, err := element.Elements("li")
		if err != nil {
			return nil, err
		}
		return MapElements(elements, ExtractPosition)
	} else {
		singlePosition, err := ExtractPosition(element)
		if err != nil {
			return nil, err
		}
		return []*LinkedInPosition{singlePosition}, nil
	}
}

func ExtractPosition(element *rod.Element) (*LinkedInPosition, error) {
	titleE, err := element.Element(".experience-item__title")
	if err != nil {
		return nil, err
	}
	title, err := titleE.Text()
	if err != nil {
		return nil, err
	}
	metaElements, err := element.Elements(".experience-item__meta-item")
	if err != nil {
		return nil, err
	}
	location, err := metaElements[1].Text()
	if err != nil {
		return nil, err
	}
	desc, err := ExtractDescription(element)
	if err != nil {
		return nil, err
	}
	return &LinkedInPosition{Title: title, Description: desc, Location: location}, nil
}

func ExtractDescription(element *rod.Element) (string, error) {
	var desc string
	moreText, err := element.Element("p.show-more-less-text__text--more")
	if err != nil {
		lessText, err := element.Element("p.show-more-less-text__text--less")
		if err != nil {
			return "", err
		}
		desc, err = lessText.Text()
		if err != nil {
			return "", err
		}
	} else {
		desc, err = moreText.Text()
		if err != nil {
			return "", err
		}
	}
	return desc, nil
}

func ExtractEducationList(page *rod.Page) ([]*LinkedInEducation, error) {
	items, err := page.Elements("ul.education__list > li")
	if err != nil {
		return nil, err
	}
	return MapElements(items, ExtractEducation)
}

func ExtractEducation(element *rod.Element) (*LinkedInEducation, error) {
	titleE, err := element.Element("h3 > a")
	if err != nil {
		return nil, err
	}
	title, err := titleE.Text()
	if err != nil {
		return nil, err
	}
	return &LinkedInEducation{Title: title}, nil
}

func MapElements[V any](ts []*rod.Element, fn func(*rod.Element) (V, error)) ([]V, error) {
	result := make([]V, len(ts))
	for i, t := range ts {
		elm, err := fn(t)
		if err != nil {
			return nil, err
		}
		result[i] = elm
	}
	return result, nil
}
