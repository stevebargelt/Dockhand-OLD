package jenkins

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"
)

//CheckLabelIsUnique checks jenkins to see if a label already exists
//it will return false if the label DOES exists and true if it does not exist
func CheckLabelIsUnique(jenkinsURL string, cloudName string, label string, username string, password string) (bool, error) {

	//TODO: LOG

	client := &http.Client{}
	url := jenkinsURL + "/scriptler/run/getLabels.groovy?cloudName=" + cloudName
	r, err := http.NewRequest("GET", url, nil)
	r.Header.Add("Accept-Encoding", "gzip")
	r.SetBasicAuth(username, password)

	response, err := client.Do(r)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		defer reader.Close()
	default:
		reader = response.Body
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	body := buf.String()

	if response.StatusCode != 200 {
		err := errors.New("ERROR: Response code: " + string(response.StatusCode) + " from " + url)
		return false, err
	}

	if err != nil {
		return false, err
	}

	if strings.Contains(string(body), label) {
		//TODO: log
		//fmt.Println("Label:", label, "already exists. will return false.")
		return false, nil
	}
	return true, nil
}

//CreateDockerTemplate calls a script on the jenkins instance to create a slave template
//given the cloundname, label (must be unique) and dockerImage to use
func CreateDockerTemplate(jenkinsURL string, cloudName string, label string, dockerImage string, username string, password string) (bool, error) {

	client := &http.Client{}

	url := jenkinsURL + "/scriptler/run/createDockerTemplate.groovy?cloudName=" + cloudName + "&label=" + label + "&image=" + dockerImage

	r, err := http.NewRequest("GET", url, nil)
	r.Header.Add("Accept-Encoding", "gzip")
	r.SetBasicAuth(username, password)

	response, err := client.Do(r)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		defer reader.Close()
	default:
		reader = response.Body
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	body := buf.String()

	if response.StatusCode != 200 {
		err := errors.New("ERROR: Response code: " + string(response.StatusCode) + " from " + url)
		return false, err
	}

	if strings.Contains(string(body), "false") {
		//TODO: log
		return false, nil
	}

	return true, nil
}
