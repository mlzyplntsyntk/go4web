// server should be wrapped up with a custom application package
// Start and AddHandler functions should be used explicitly, otherwise
// this server will not be able to handle any requests at all
// (it can do some redirections if it's multilingual)

package core 

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"errors"
	"strconv"
)

// TODO: These variables should not be visible to other packages. We should consider
// encapsulating the needed properties into functions, later... Last thing to do :)

var (
	SiteConfig 		Config

	// this parameter is important because it decides if the language prefix
	// should be used at the requested urls or not.

	isMultiLanguage bool		= false
	Handlers []HandlerInterface = make([]HandlerInterface, 0)
)

type HandlerInterface interface {
	RunBeforeHandled() bool
	RunAfterHandled() bool
	Pattern() string
	Run(http.ResponseWriter, *http.Request, *Page) (bool, error)
}

// gets the current language from the requested url's first path
// if the requested path is not empty and does not match any of the languages we
// provide, it also returns an error.
// note that this error is ommited if our application is single language.
// this function always returns a language. If it matchs it returns the matched
// language, otherwise it will return the default language.

func getCurrentLanguage(path string) (Language, error) {
	var defaultLang Language

	for i:=0; i<len(SiteConfig.Languages); i++ {

		if SiteConfig.Languages[i].IsDefault {
			defaultLang = SiteConfig.Languages[i]
		}

		if SiteConfig.Languages[i].Prefix == path {
			return SiteConfig.Languages[i], nil
		}
	}

	if len(path) == 0 {
		return defaultLang, nil
	} else {
		return defaultLang, errors.New("Language %s is not supported" + path)
	}
}

// returns the path defined in config.json file. See mainHandler function
// to know more about the parameters.

func getCurrentRoute(path string, langPrefix string) (Page, error) {
	
	for i:=0; i<len(SiteConfig.Routes); i++ {

		if SiteConfig.Routes[i].LanguagePrefix != langPrefix {
			continue
		}

		for j:=0; j<len(SiteConfig.Routes[i].Pages); j++ {
			if len(path)==0 && SiteConfig.Routes[i].Pages[j].IsDefault {
				return SiteConfig.Routes[i].Pages[j], nil
			}

			if SiteConfig.Routes[i].Pages[j].Pattern == path {
				return SiteConfig.Routes[i].Pages[j], nil
			}
		}
	}
	var matchingPage Page
	return matchingPage, errors.New("No Pages Available")
}

// this function manages all the pipeline of the framework
// it delegates requests to appropriate handlers. 
// assuming that our domain is http://example.com/ handler first
// searches for the default path for our application. This path
// should be defined explicitly in our config.json file, in form
// of Page type, and the definition should be as it's described 
// in checkSiteConfig function.

func MainHandler(config Config) http.HandlerFunc {

	// SiteConfig setting moved here, for testing purposes.
	// TODO : Find a better way to implement this.

	SiteConfig = config

	// here we decide if there are more than one languages

	isMultiLanguage = len(SiteConfig.Languages) > 1

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// paths variable is an array of strings. 
		// it should always have a minimum size of 1. Which is 
		// an empty string at worse.

		paths := strings.Split(r.URL.Path[1:], "/")

		// pathToHandle means so much to this context. it is either
		// empty, or a language prefix (if the app is multilanguage)
		// orrrr a path for single language apps.

		pathToHandle := paths[0]

		// see getCurrentLanguage Comments 

		language, langStatus := getCurrentLanguage(pathToHandle)

		// here we will check if the app is multi language, and if so,
		// we need to know that the requested url's first path should 
		// either match a language prefix or empty

		if isMultiLanguage {

			// langStatus is an error, if the path does not match any
			// provided language. So we will redirect it to the default
			// language for now.

			if langStatus != nil {
				http.Redirect(w, r, "/"+language.Prefix, 301)
			}

			// it should be noted that, if the application is multi-language
			// we need to shift the path to the next index (if available) 
			// so we have the real path to deal with.

			if len(paths) > 1 {
				pathToHandle = paths[1]
			} else {
				pathToHandle = ""
			}
		}

		// pathToHandle variable was defined at the very beginning of our function
		// so it can change if our app is multi language or not. Either way, we will
		// have the current language for current context and the pathToHandle should be 
		// executed by the given language.

		if route, err := getCurrentRoute(pathToHandle, language.Prefix); err != nil {

			// TODO : Implement either a native 404 page or expect the developer to do so
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "not found")

		} else {

			var currentHandler HandlerInterface

			for item := range Handlers {

				//here we have an exact match

				if Handlers[item].Pattern() == route.Pattern {
					currentHandler = Handlers[item]
				}

				if Handlers[item].Pattern() == "*" && currentHandler == nil {
					currentHandler = Handlers[item]
				}

				if Handlers[item].RunBeforeHandled() {
					Handlers[item].Run(w,r,&route)
				}
			}

			currentHandler.Run(w,r,&route)

			for item := range Handlers {
				if Handlers[item].RunAfterHandled() {
					Handlers[item].Run(w,r,&route)
				}
			}

		}
	})
}

// check the site configuration file once the application starts.
// we should give here directions to correctly setup the config
// file

func checkSiteConfig() {

}

// adding our custom handlers to respond to requests.
// fn => the function signature to respond to request
// pattern => the pattern to match against this handler
// before => if this handler will run before the main Handler.
// after => if this handler will run after the main handler 

func AddHandler(hd HandlerInterface) {
	Handlers = append(Handlers, hd)
}

//Gets the config file from given path. 

func GetConfigFromJSON(configFile string) Config {
	var config Config
	contents, _ := ioutil.ReadFile(configFile)
	json.Unmarshal(contents, &config)
	return config
}

// Server starts here. it needs to parse config.json file first.
// checkSiteConfig function should explain everything that is needed 
// by the framework, if the Start fails due to configuration settings

func Start(portNumber int, configFile string) {
	config := GetConfigFromJSON(configFile)

	checkSiteConfig()

	http.HandleFunc("/", MainHandler(config))

	http.ListenAndServe(":" + strconv.Itoa(portNumber), nil)
}
