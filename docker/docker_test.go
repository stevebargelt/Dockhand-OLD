package docker

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	flag.Parse()
	exitCode := m.Run()

	// Exit
	os.Exit(exitCode)
}

func TestGetDockerImage(t *testing.T) {

	fmt.Println("test")

}

func ExampleGetDockerImage() {
	//numbers := []int{5, 5, 5}
	//fmt.Println(Sum(numbers))
	// Output:
	// 15
}
