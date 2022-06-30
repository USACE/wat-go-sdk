package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("wat-go-sdk")
	fmt.Println(uuid.NodeID())
	//uuid.SetNodeID([]byte{1, 2, 3, 4, 5, 6})
	//uuid.Variant
	//fmt.Println(uuid.NodeID())
	manifestID := uuid.New()
	fmt.Println(manifestID)
	file1 := uuid.NewSHA1(manifestID, []byte("Input-Muncie.tmp.hdf"))
	file2 := uuid.NewSHA1(manifestID, []byte("Output-Muncie.tmp.hdf"))
	fmt.Println(file1)
	fmt.Println(file2)
	parsed1 := uuid.MustParse(file1.String())
	parsed2 := uuid.MustParse(file2.String())
	fmt.Println(parsed1.Value())
	fmt.Println(parsed2.Domain().String())
	/*for i := 1; i < 40; i++ {
		fmt.Println(fmt.Sprintf("%v,%v", i, uuid.New()))
	}
	*/
}
