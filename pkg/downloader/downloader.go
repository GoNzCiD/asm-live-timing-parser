// https://www.zenrows.com/blog/selenium-golang

package downloader

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type Downloader struct {
	chromeDriverPath string
	seleniumUrl      string
	debug            bool
	user             string
	pass             string
	url              string
	serverNo         int
}

func NewDownloader(chromeDriverPath string, seleniumUrl string, user string, pass string, url string, serverNo int, debug bool) *Downloader {
	return &Downloader{
		chromeDriverPath: chromeDriverPath,
		seleniumUrl:      seleniumUrl,
		debug:            debug,
		user:             user,
		pass:             pass,
		url:              url,
		serverNo:         serverNo,
	}
}

func (d *Downloader) Download() (string, error) {
	// Using a fixed selenium docker image/version because of the JSON & W3C capability
	// https://stackoverflow.com/questions/76574419/golang-docker-selenium-chrome
	// selenium/standalone-chrome:4.8.3

	//// initialize a Chrome browser instance on port 4444 [local mode]
	//service, err := selenium.NewChromeDriverService(d.chromeDriverPath, 4444)
	//if err != nil {
	//	return "", fmt.Errorf("error initializing selenium: %v", err)
	//}
	//defer service.Stop()

	// configure the browser options
	caps := selenium.Capabilities{}
	if !d.debug {
		caps.AddChrome(chrome.Capabilities{
			Args: []string{
				"--headless", // comment out this line for testing
				//"--disable-dev-shm-usage",
			},
		})
	}

	// create a new remote client with the specified options
	driver, err := selenium.NewRemote(caps, d.seleniumUrl)
	//driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		return "", fmt.Errorf("error initializing remote client: %v", err)
	}
	defer driver.Quit()

	//// maximize the current window to avoid responsive rendering
	//err = driver.MaximizeWindow("")
	//if err != nil {
	//	log.Fatal("Error:", err)
	//}

	// build login URL
	u, err := url.Parse(d.url)
	if err != nil {
		return "", fmt.Errorf("error parsing the url: %v", err)
	}
	u.Path = "login"
	q := u.Query()
	q.Add("server", strconv.Itoa(d.serverNo))
	u.RawQuery = q.Encode()

	err = driver.Get(u.String())
	if err != nil {
		return "", fmt.Errorf("error getting login: %v", err)
	}

	// select the login form
	formElement, err := driver.FindElement(selenium.ByXPATH, "//form[@action='/login']")
	if err != nil {
		return "", fmt.Errorf("error finding the login form: %v", err)
	}
	// fill in the login form fields
	user, err := formElement.FindElement(selenium.ByName, "Username")
	if err != nil {
		return "", fmt.Errorf("error finding the 'Username' input: %v", err)
	}
	err = user.SendKeys(d.user)
	if err != nil {
		return "", fmt.Errorf("error setting the 'Username' value: %v", err)
	}

	pass, err := formElement.FindElement(selenium.ByName, "Password")
	if err != nil {
		return "", fmt.Errorf("error finding the 'Password' input: %v", err)
	}
	err = pass.SendKeys(d.pass)
	if err != nil {
		return "", fmt.Errorf("error setting the 'password' value: %v", err)
	}

	// submit the form
	err = formElement.Submit()
	if err != nil {
		return "", fmt.Errorf("error submitting the form: %v", err)
	}

	// build the leaderboard.json url
	u, err = url.Parse(d.url)
	if err != nil {
		return "", fmt.Errorf("error parsing the url: %v", err)
	}
	u.Path = "api/live-timings/leaderboard.json"

	// get the leaderboard.json URL
	err = driver.Get(u.String())
	if err != nil {
		return "", fmt.Errorf("error getting leaderboard: %v", err)
	}

	// obtain just the json content inside <pre> tag
	html, err := driver.FindElement(selenium.ByTagName, "pre")
	if err != nil {
		return "", fmt.Errorf("error finding 'pre' tag: %v", err)
	}

	json, err := html.Text()
	if err != nil {
		return "", fmt.Errorf("error retrieving json text: %v", err)
	}

	return json, nil
}
