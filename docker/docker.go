package docker

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

//  InitClient
// is this the comment you are looking for
func InitClient(dockerHost string, certFile string, keyFile string, caFile string) (*client.Client, error) {

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpCli := &http.Client{Transport: transport}

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(dockerHost, "v1.24", httpCli, defaultHeaders)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func GetDockerImage(cli *client.Client, imageName string, registryUsername string, registryPassword string, registryURL string) (*types.ImageInspect, error) {

	fmt.Println("Looking for image", imageName, "...")
	image, _, imageErr := cli.ImageInspectWithRaw(context.TODO(), imageName)

	newImage, err := pullImage(cli, imageName, registryUsername, registryPassword, registryURL)
	if err != nil {
		if imageErr == nil {
			fmt.Println("Cannot pull the latest version of image", imageName, ":", err)
			fmt.Println("Locally found image will be used instead.")
			return &image, nil
		}
		return nil, err
	}
	return newImage, nil
}

func pullImage(cli *client.Client, imageName string, registryUsername string, registryPassword string, registryURL string) (*types.ImageInspect, error) {

	fmt.Println("Pulling docker image", imageName, "...")

	ref := imageName
	// Add :latest to limit the download results
	if !strings.ContainsAny(ref, ":@") {
		ref += ":latest"
	}

	authConfig := types.AuthConfig{Username: registryUsername, Password: registryPassword, ServerAddress: registryURL}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(authConfig); err != nil {
		return nil, err
	}
	encodedAuth := base64.URLEncoding.EncodeToString(buf.Bytes())

	options := types.ImagePullOptions{RegistryAuth: encodedAuth}
	readCloser, err := cli.ImagePull(context.Background(), ref, options)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	if _, err := io.Copy(ioutil.Discard, readCloser); err != nil {
		return nil, fmt.Errorf("Failed to pull image: %s: %s", ref, err)
	}

	image, _, err := cli.ImageInspectWithRaw(context.Background(), imageName)
	return &image, err
}

// Stuff I used at the start to test the connection to the Docker Host - might come in handy some day.
// TODO: delete after adding to git
// func randomTestStuff() {

// 	fmt.Print("\n\n********************\nList Images\n********************\n")
// 	imageOptions := types.ImageListOptions{All: true}
// 	images, err := cli.ImageList(context.Background(), imageOptions)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, i := range images {
// 		fmt.Println(i.ID, i.Labels)
// 	}

// 	fmt.Print("\n\n********************\nList Containers\n********************\n")
// 	options := types.ContainerListOptions{All: true}
// 	containers, err := cli.ContainerList(context.Background(), options)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, c := range containers {
// 		fmt.Println(c.ID, c.Names)
// 	}

// 	fmt.Print("\n\n********************\nOther Info???\n********************\n")
// 	fmt.Println("cli type=", reflect.TypeOf(cli))

// }
