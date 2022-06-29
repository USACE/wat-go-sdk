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

	for i := 1; i < 40; i++ {
		fmt.Println(fmt.Sprintf("%v,%v", i, uuid.New()))
	}
}
