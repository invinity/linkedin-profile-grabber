package linkedin

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
	. "github.com/go-rod/rod/lib/input"
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
	Summary    string
	Experience []*LinkedInExperience
	Education  []*LinkedInEducation
	Projects   []*LinkedInProject
}

type LinkedInPosition struct {
	Title       string
	StartDate   string
	EndDate     string
	Location    string
	Description string
}

type LinkedInExperience struct {
	Company      string
	CompanyImage *string
	Positions    []*LinkedInPosition
}

type LinkedInEducation struct {
	Title       string
	Subtitle    string
	StartDate   string
	EndDate     string
	Description string
}

type LinkedInProject struct {
	Title       string
	Dates       string
	Description string
}

func (r *LinkedIn) getPage(firstName string, lastName string, profileAlias string) (*rod.Page, error) {
	page, err := r.browser.Page(proto.TargetCreateTarget{URL: "https://www.linkedin.com"})
	if err != nil {
		return nil, err
	}
	waitDur, _ := time.ParseDuration("10s")
	err = page.WaitDOMStable(waitDur, .5)
	if err != nil {
		return nil, err
	}
	log.Println("Got page ", page.MustInfo().Title)
	peopleLink, err := page.Element("a[data-tracking-control-name='guest_homepage-basic_guest_nav_menu_people']")
	if err != nil {
		return nil, err
	}
	peopleLink.MustType(Enter)
	err = page.WaitDOMStable(waitDur, .5)
	if err != nil {
		return nil, err
	}
	log.Println("Got page ", page.MustInfo().Title)
	firstNameElm, err := page.Element("input[name='firstName']")
	if err != nil {
		return nil, err
	}
	lastNameElm, err := page.Element("input[name='lastName']")
	if err != nil {
		return nil, err
	}
	for _, v := range firstName {
		firstNameElm.MustType(Key(v))
	}
	for _, v := range lastName {
		lastNameElm.MustType(Key(v))
	}
	lastNameElm.MustType(Enter)
	err = page.WaitDOMStable(waitDur, .5)
	if err != nil {
		return nil, err
	}
	log.Println("Got page ", page.MustInfo().Title)
	profileUrl := fmt.Sprintf("https://www.linkedin.com/in/%s?trk=people-guest_people_search-card", profileAlias)
	profileLink, err := page.Element(fmt.Sprintf("a[href='%s']", profileUrl))
	if err != nil {
		return nil, err
	}
	profileLink.MustType(Enter)
	err = page.WaitDOMStable(waitDur, .5)
	if err != nil {
		return nil, err
	}
	log.Println("Got page ", page.MustInfo().Title)
	return page, nil
}

func (r *LinkedIn) RetrieveProfile(firstName string, lastName string, profileAlias string) (*LinkedInProfile, error) {
	page, err := r.getPage(firstName, lastName, profileAlias)
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
	summaryElm, err := page.Element("section[data-section=summary] p")
	if err != nil {
		return nil, err
	}
	summary, err := summaryElm.Text()
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
	projects, err := ExtractProjectList(page)
	if err != nil {
		return nil, err
	}
	return &LinkedInProfile{Name: name, Headline: headline, Summary: summary, Experience: experience, Education: education, Projects: projects}, nil
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
	start, end, err := ExtractStartEndDates(element)
	if err != nil {
		return nil, err
	}
	desc, err := ExtractDescription(element)
	if err != nil {
		return nil, err
	}
	return &LinkedInPosition{Title: title, Description: desc, Location: location, StartDate: start, EndDate: end}, nil
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
	subtitleElms, err := element.Elements("h4 > span")
	if err != nil {
		return nil, err
	}
	subtitle := ""
	for _, v := range subtitleElms {
		if subtitle != "" {
			subtitle += " "
		}
		subtitle += v.MustText()
	}
	desc, err := ExtractDescription(element)
	if err != nil {
		return nil, err
	}
	start, end, err := ExtractStartEndDates(element)
	if err != nil {
		return nil, err
	}
	return &LinkedInEducation{Title: title, Subtitle: subtitle, Description: desc, StartDate: start, EndDate: end}, nil
}

func ExtractProjectList(page *rod.Page) ([]*LinkedInProject, error) {
	items, err := page.Elements("ul.projects__list > li")
	if err != nil {
		return nil, err
	}
	return MapElements(items, ExtractProject)
}

func ExtractProject(element *rod.Element) (*LinkedInProject, error) {
	titleE, err := element.Element("div > h3")
	if err != nil {
		return nil, err
	}
	title, err := titleE.Text()
	if err != nil {
		return nil, err
	}
	datesE, err := element.Element("div > h4 > span.date-range")
	if err != nil {
		return nil, err
	}
	dates, err := datesE.Text()
	if err != nil {
		return nil, err
	}
	desc, err := ExtractDescription(element)
	if err != nil {
		return nil, err
	}
	return &LinkedInProject{Title: title, Dates: dates, Description: desc}, nil
}

func ExtractStartEndDates(element *rod.Element) (string, string, error) {
	start, end := "", ""
	dates, err := element.Elements("span.date-range > time")
	if err != nil {
		return "", "", err
	}
	if len(dates) > 0 {
		start, err = dates[0].Text()
		if err != nil {
			return "", "", err
		}
	}
	if len(dates) > 1 {
		end, err = dates[1].Text()
		if err != nil {
			return "", "", err
		}
	}
	return start, end, nil
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
