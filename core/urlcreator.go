package core

import (
	"net/url"
	"strings"
)

// CreateUrl replaces the host part with the given host part
func CreateUrl(original url.URL, host string) string {

	if host[len(host)-1:] == "/" {
		host = host[0 : len(host)-1]
	}

	newURL := host + original.RequestURI()
	return newURL

}

// ExtractFile extracts the file part of a URL (i.e index.html)
func ExtractFile(url url.URL) string {

	path := url.Path

	var element = "index.html"
	indexOfLastTrailing := strings.LastIndex(path, "/")

	if indexOfLastTrailing != -1 {
		file := strings.Trim(path[indexOfLastTrailing+1:], " ")

		if len(file) > 0 && specifiesAFile(file) {
			element = file
		}

	}

	return element
}

func specifiesAFile(file string) bool {
	return strings.Contains(file, ".")
}
