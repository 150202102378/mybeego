package utils

import (
	"fmt"
	"testing"
)

func TestGetProjectPath(t *testing.T) {
	path := GetProjectPath("utxo")
	fmt.Println("path: ", path)
	t.Log(path)
}
