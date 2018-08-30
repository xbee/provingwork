package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/sstelfox/provingwork"
)

func main() {
	sw := provingwork.NewStrongWork(
		[]byte("Just some test data in the string"),
		&provingwork.WorkOptions{BitStrength: 20},
	)
	sw.FindProof()
	fmt.Printf("%v\n", sw.String())

	json, _ := json.Marshal(sw)
	fmt.Println(string(json))
	fmt.Printf("%x\n", sha256.Sum256(sw.ContentHash()))
}
