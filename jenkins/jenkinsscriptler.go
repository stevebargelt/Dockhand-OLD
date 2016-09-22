package jenkinsscriptler

import (
	"fmt"
    "net/http"
    "io/ioutil"
)


func Test() {
    
    resp, err := http.Get("http://dockerbuild.harebrained-apps.com/scriptler/run/getLabels.groovy?cloudName=AzureJenkins")
    if err != nil {
	 panic(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    fmt.Println(body)
}
