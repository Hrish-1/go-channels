package main

import (
	"fmt"
	"sync"
	"time"
)

var userMap map[int]User = map[int]User{
	1: {1, "John", "john@g.com", "9988773", 1},
	2: {2, "Jim", "jim@g.com", "4433221", 2},
	3: {3, "Mary", "mary@g.com", "4455667", 3},
}

var petMap map[int]Pet = map[int]Pet{
	1: {1, "Max", "dog", User{}},
	2: {2, "Fiona", "cat", User{}},
	3: {3, "Bonnie", "cat", User{}},
}

func lookUpUsers(id int) User {
	time.Sleep(500 * time.Millisecond)
	return userMap[id]
}

func lookupPets(id int) Pet {
	time.Sleep(500 * time.Millisecond)
	return petMap[id]
}

type User struct {
	id    int
	name  string
	email string
	phone string
	petId int
}

type Pet struct {
	id      int
	name    string
	petType string
	user    User
}

func getUsers(ids []int) chan User {
	var wg sync.WaitGroup
	c := make(chan User, len(ids))
	wg.Add(len(ids))
	for _, id := range ids {
		go func(id int) {
			c <- lookUpUsers(id)
			wg.Done()
		}(id)
	}
	wg.Wait()
	close(c)
	return c
}

func getPets(users []User) chan Pet {
	var wg sync.WaitGroup
	c := make(chan Pet, len(users))
	wg.Add(len(users))
	for _, user := range users {
		go func(u User) {
			pet := lookupPets(u.petId)
			pet.user = u
			c <- pet
			wg.Done()
		}(user)
	}
	wg.Wait()
	close(c)
	return c
}

func main() {
	start := time.Now().UnixMilli()
	ids := []int{1, 2, 3}
	userChan := getUsers(ids)
	var users []User
	for user := range userChan {
		users = append(users, user)
	}
	petChan := getPets(users)
	m := make(map[int]string)
	for pet := range petChan {
		str := fmt.Sprintf("| %d | %s | %s | %s | %s |", pet.user.id, pet.user.name, pet.user.email, pet.name, pet.petType)
		m[pet.user.id] = str
	}
	for _, id := range ids {
		fmt.Println(m[id])
	}
	end := time.Now().UnixMilli()
	fmt.Println(end - start)
}
