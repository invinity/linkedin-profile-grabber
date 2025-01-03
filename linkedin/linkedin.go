package linkedin

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
	. "github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

func NewBrowser(browser *rod.Browser) *LinkedInBrowser {
	return &LinkedInBrowser{browser: browser}
}

type LinkedInBrowser struct {
	browser *rod.Browser
}

type LinkedInProfile struct {
	GeneratedAt    time.Time
	Name           string
	Headline       string
	Summary        string
	Experience     []*LinkedInExperience
	Education      []*LinkedInEducation
	Projects       []*LinkedInProject
	Certifications []*LinkedInCertification
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
	StartDate   string
	EndDate     string
	Description string
}

type LinkedInCertification struct {
	Title       string
	Institution string
	ImgUrl      *string
	IssuedOn    string
	ExpiresOn   string
}

func (r *LinkedInBrowser) RetrieveProfileViaSearch(firstName string, lastName string, profileAlias string) (*LinkedInProfile, error) {
	page, err := r.navigateToProfilePageViaSearch(firstName, lastName, profileAlias)
	if err != nil {
		return nil, err
	}
	return r.extractProfileData(page)
}

func (r *LinkedInBrowser) RetrieveProfileViaLogin(email string, password string) (*LinkedInProfile, error) {
	page, err := r.navigateToProfilePageViaLogin(email, password)
	if err != nil {
		return nil, err
	}
	return r.extractProfileData(page)
}

func (r *LinkedInBrowser) RetrieveProfileViaGoogleLogin(email string, password string) (*LinkedInProfile, error) {
	page, err := r.navigateToProfilePageViaGoogleLogin(email, password)
	if err != nil {
		return nil, err
	}
	return r.extractProfileData(page)
}

func (r *LinkedInBrowser) navigateToProfilePageViaLogin(email string, password string) (*rod.Page, error) {
	page, err := r.performLinkedInLogin(email, password)
	if err != nil {
		return nil, err
	}
	err = page.Navigate("https://www.linkedin.com/public-profile/settings?trk=d_flagship3_profile_self_view_public_profile")
	if err != nil {
		return nil, err
	}
	err = page.WaitDOMStable(time.Second*3, .2)
	if err != nil {
		return nil, err
	}
	title := page.MustInfo().Title
	log.Println("Got page ", title)
	return page, nil
}

func (r *LinkedInBrowser) performLinkedInLogin(email string, password string) (*rod.Page, error) {
	attempt := 0
	var page *rod.Page
	var err error
	for (page == nil || !strings.Contains(page.MustInfo().Title, "Feed")) && attempt <= 10 {
		attempt += 1
		log.Println("Not yet at a logged in page, performing login attempt", attempt)
		page, err = r.attemptLinkedInLogin(email, password)
		if err != nil {
			return nil, err
		}
	}

	if page == nil {
		return nil, fmt.Errorf("unable to login into Linkedin after %d attempts", attempt)
	}

	return page, nil
}

func (r *LinkedInBrowser) attemptLinkedInLogin(email string, password string) (*rod.Page, error) {
	page, err := r.browser.Page(proto.TargetCreateTarget{URL: "https://www.linkedin.com/login"})
	if err != nil {
		return nil, err
	}
	typeDur := 200 * time.Millisecond
	waitDur := 2 * time.Second
	err = page.WaitDOMStable(waitDur, .2)
	if err != nil {
		return nil, err
	}
	log.Println("Got page ", page.MustInfo().Title)
	usernameInput, err := page.Element("input[id=username]")
	if err != nil {
		return nil, err
	}
	for _, v := range email {
		usernameInput.MustType(Key(v))
		time.Sleep(typeDur)
	}
	passwordInput, err := page.Element("input[id=password]")
	if err != nil {
		return nil, err
	}
	for _, v := range password {
		passwordInput.MustType(Key(v))
		time.Sleep(typeDur)
	}
	passwordInput.MustType(Enter)
	err = page.WaitDOMStable(waitDur, .2)
	if err != nil {
		return nil, err
	}
	title := page.MustInfo().Title
	log.Println("Arrived at page after login: ", title)
	if strings.Contains(title, "Security Verification") {
		log.Println("Got Security Verification page, sleeping for a bit")
		time.Sleep(10 * time.Second)
		err = page.WaitDOMStable(waitDur, .2)
		if err != nil {
			return nil, err
		}
	}
	title = page.MustInfo().Title
	if !strings.Contains(title, "Feed") {
		return nil, errors.New("Expected to get feed page after login, but was " + title)
	}
	return page, nil
}

func (r *LinkedInBrowser) navigateToProfilePageViaGoogleLogin(email string, password string) (*rod.Page, error) {
	page, err := r.browser.Page(proto.TargetCreateTarget{URL: "https://www.linkedin.com/login"})
	if err != nil {
		return nil, err
	}
	typeDur := 200 * time.Millisecond
	waitDur := 2 * time.Second
	err = page.WaitDOMStable(waitDur, .2)
	if err != nil {
		return nil, err
	}
	log.Println("Got page ", page.MustInfo().Title)
	spans, err := page.Elements("div[id=container] span")
	if err != nil {
		return nil, err
	}
	var googleLogin *rod.Element
	for _, v := range spans {
		innerText := v.MustText()
		if innerText == "Continue with Google" {
			googleLogin = v
			break
		}
	}
	if googleLogin == nil {
		return nil, errors.New("unable to find google login element")
	}
	googleLogin.Click(proto.InputMouseButtonLeft, 1)
	pages, err := r.browser.Pages()
	if err != nil {
		return nil, err
	}
	googleLoginPage, err := pages.FindByURL("google\\.com")
	if err != nil {
		return nil, err
	}
	err = googleLoginPage.WaitDOMStable(waitDur, .2)
	if err != nil {
		return nil, err
	}
	usernameInput, err := googleLoginPage.Element("input[id=identifierId]")
	if err != nil {
		return nil, err
	}
	for _, v := range email {
		usernameInput.MustType(Key(v))
		time.Sleep(typeDur)
	}
	usernameInput.MustType(Enter)
	err = googleLoginPage.WaitDOMStable(waitDur, .2)
	if err != nil {
		return nil, err
	}
	passwordInput, err := googleLoginPage.Element("input[name=Passwd]")
	if err != nil {
		return nil, err
	}
	for _, v := range password {
		passwordInput.MustType(Key(v))
		time.Sleep(typeDur)
	}
	passwordInput.MustType(Enter)
	err = page.WaitDOMStable(waitDur, .2)
	if err != nil {
		return nil, err
	}
	err = page.Navigate("https://www.linkedin.com/public-profile/settings?trk=d_flagship3_profile_self_view_public_profile")
	if err != nil {
		return nil, err
	}
	err = page.WaitDOMStable(time.Second*3, .2)
	if err != nil {
		return nil, err
	}
	title := page.MustInfo().Title
	log.Println("Got page ", title)
	return page, nil
}

func (r *LinkedInBrowser) navigateToProfilePageViaSearch(firstName string, lastName string, profileAlias string) (*rod.Page, error) {
	typeDelay := 200 * time.Millisecond
	page, err := r.browser.Page(proto.TargetCreateTarget{URL: "https://www.linkedin.com"})
	if err != nil {
		return nil, err
	}
	waitDur, _ := time.ParseDuration("2s")
	err = page.WaitDOMStable(waitDur, .2)
	if err != nil {
		return nil, err
	}
	log.Println("Got page ", page.MustInfo().Title)
	peopleLink, err := page.Element("a[data-tracking-control-name='guest_homepage-basic_guest_nav_menu_people']")
	if err != nil {
		return nil, err
	}
	peopleLink.MustType(Enter)
	err = page.WaitDOMStable(waitDur, .2)
	if err != nil {
		return nil, err
	}
	title := page.MustInfo().Title
	log.Println("Got page ", title)
	if !strings.HasPrefix(title, "Search for people") {
		return nil, errors.New("Unexpected page: " + title)
	}
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
		time.Sleep(typeDelay)
	}
	for _, v := range lastName {
		lastNameElm.MustType(Key(v))
		time.Sleep(typeDelay)
	}
	lastNameElm.MustType(Enter)
	err = page.WaitDOMStable(waitDur, .2)
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
	err = page.WaitDOMStable(waitDur, .2)
	if err != nil {
		return nil, err
	}
	log.Println("Got page ", page.MustInfo().Title)
	return page, nil
}

func (r *LinkedInBrowser) extractProfileData(page *rod.Page) (*LinkedInProfile, error) {
	_, err := page.Element("body")
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
	certifications, err := ExtractCertificationList(page)
	if err != nil {
		return nil, err
	}
	return &LinkedInProfile{GeneratedAt: time.Now().UTC(), Name: name, Headline: headline, Summary: summary, Experience: experience, Education: education, Projects: projects, Certifications: certifications}, nil
}

func ExtractExperienceList(page *rod.Page) ([]*LinkedInExperience, error) {
	expList, err := page.Elements("ul.experience__list > li")
	if err != nil {
		return nil, err
	}
	return MapElements(expList, ExtractExperience)
}

func ExtractExperience(element *rod.Element) (*LinkedInExperience, error) {
	titleElements, err := element.Elements(".profile-section-card__subtitle,.experience-item__subtitle")
	if err != nil {
		return nil, err
	}
	title, err := titleElements[0].Text()
	if err != nil {
		return nil, err
	}
	companyImageE, err := element.Element("img")
	if err != nil {
		return nil, err
	}
	companyImage, err := companyImageE.Attribute("data-delayed-url")
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
	class, err := element.Attribute("class")
	if err != nil {
		return nil, err
	}
	isGroup := strings.Contains(*class, "experience-group")
	titleElements, err := element.Elements(".profile-section-card__title,.experience-item__title")
	if err != nil {
		return nil, err
	}
	title, err := titleElements[0].Text()
	if err != nil {
		return nil, err
	}
	metaElements, err := element.Elements(".experience-item__meta-item")
	if err != nil {
		return nil, err
	}
	locationXPath := ".experience-item__location"
	if isGroup {
		locationXPath = ".experience-group-position__location"
	}
	locationElements, err := element.Elements(locationXPath)
	if err != nil {
		return nil, err
	}
	location := ""
	if len(locationElements) > 0 {
		location, err = locationElements[0].Text()
		if err != nil {
			return nil, err
		}
	} else {
		location, err = metaElements[1].Text()
		if err != nil {
			return nil, err
		}
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
		desc = strings.Replace(desc, "Show more", "", 1)
	} else {
		desc, err = moreText.Text()
		if err != nil {
			return "", err
		}
		desc = strings.Replace(desc, "Show less", "", 1)
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
	start, end, err := ExtractStartEndDates(element)
	if err != nil {
		return nil, err
	}
	desc, err := ExtractDescription(element)
	if err != nil {
		return nil, err
	}
	return &LinkedInProject{Title: title, StartDate: start, EndDate: end, Description: desc}, nil
}

func ExtractCertificationList(page *rod.Page) ([]*LinkedInCertification, error) {
	items, err := page.Elements("section[data-section=certifications] ul li")
	if err != nil {
		return nil, err
	}
	return MapElements(items, ExtractCertifcation)
}

func ExtractCertifcation(element *rod.Element) (*LinkedInCertification, error) {
	titleE, err := element.Element("h3")
	if err != nil {
		return nil, err
	}
	title, err := titleE.Text()
	if err != nil {
		return nil, err
	}
	institutionE, err := element.Element("h4 > a")
	if err != nil {
		return nil, err
	}
	institution, err := institutionE.Text()
	if err != nil {
		return nil, err
	}
	institutionImgE, err := element.Element("img")
	if err != nil {
		return nil, err
	}
	institutionImg, err := institutionImgE.Attribute("data-delayed-url")
	if err != nil {
		return nil, err
	}
	datesE, err := element.Elements("span > time")
	if err != nil {
		return nil, err
	}
	issued, expires := "", ""
	if len(datesE) >= 1 {
		issued, err = datesE[0].Text()
		if err != nil {
			return nil, err
		}
	}
	if len(datesE) >= 2 {
		expires, err = datesE[1].Text()
		if err != nil {
			return nil, err
		}
	}
	return &LinkedInCertification{Title: title, Institution: institution, ImgUrl: institutionImg, IssuedOn: issued, ExpiresOn: expires}, nil
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
