package core

import (
	"net/url"
)

func CreateUrl(original url.URL, host string) string {

	if host[len(host)-1:] == "/" {
		host = host[0 : len(host)-1]
	}

	newURL := host + original.RequestURI()
	return newURL

}
