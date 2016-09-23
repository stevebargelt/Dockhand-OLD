package jenkins

import (
	"errors"
	"io/ioutil"
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
	r.SetBasicAuth(username, password)

	resp, err := client.Do(r)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err := errors.New("ERROR: Response code: " + string(resp.StatusCode) + " from " + url)
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
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
	r.SetBasicAuth(username, password)

	resp, err := client.Do(r)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err := errors.New("ERROR: Response code: " + string(resp.StatusCode) + " from " + url)
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if strings.Contains(string(body), "false") {
		//TODO: log
		return false, nil
	}

	return true, nil
}
