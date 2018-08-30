package main

import (
	"encoding/json"
	"fmt"

	"github.com/xbee/provingwork"
)

func main() {
	sw := provingwork.NewStrongWork(
		[]byte("Just some test data in the string"),
		&provingwork.WorkOptions{BitStrength: 22},
	)
	sw.FindProof()
	fmt.Printf("%v\n", sw.String())

	json, _ := json.Marshal(sw)
	fmt.Println(string(json))
	fmt.Printf("%x\n", sw.ContentHash())
}
