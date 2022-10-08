package main

import (
	"fmt"
	"go-designpattern/iterator/iterator"
)

func main() {
	user1 := &iterator.User{
		Name: "a",
		Age:  30,
	}
	user2 := &iterator.User{
		Name: "b",
		Age:  20,
	}

	userCollection := &iterator.UserCollection{
		Users: []*iterator.User{user1, user2},
	}

	i := userCollection.CreateIterator()
	for i.HasNext() {
		user := i.GetNext()
		fmt.Printf("User is %+v\n", user)
	}
}
