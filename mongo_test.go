package tool_mongo

import (
	"context"
	"fmt"
	"log"
	"testing"
	"tool-mongo/input"
	"tool-mongo/output"
)

func TestConnect(t *testing.T) {
	s := output.NewStore()
	if err := s.Connect(); err != nil {
		log.Fatal(err)
	}

	_ = s.Close()
}

func TestUser(t *testing.T) {
	s := output.NewStore()
	if err := s.Connect(); err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	u := output.NewUser(&input.User{
		Id:      "4",
		Name:    "xyt",
		Age:     26,
		Money:   1000,
		Friends: []*input.People{{"a", 20}},

		Houses: []string{"A", "B"},
		Map1:   map[string]int32{"C": 21},
		Map2:   map[int32]*input.People{23: {"b", 22}},
		Map3:   map[string][]*input.People{"map3": {{"c", 24}, {"d", 25}}},
	})
	err := u.Insert(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", u)

	//err = u.SetId("2").UpsertById(context.TODO())
	//u := output.NewUser(&input.User{
	//	Id:      "2",
	//	Name:    "xyt",
	//	Age:     26,
	//	Money:   1000,
	//	Friends: []*input.People{{"a", 20}},
	//
	//	Houses: []string{"A", "B"},
	//	Map1:   map[string]int32{"C": 21},
	//	Map2:   map[int32]*input.People{23: {"b", 22}},
	//	Map3:   map[string][]*input.People{"map3": {{"c", 24}, {"d", 25}}},
	//})
	//err = u.UpsertById(context.TODO())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Found a single document: %+v\n", u)
}

func TestStudent(t *testing.T) {
	store := output.NewStore()
	if err := store.Connect(); err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	s := output.NewStudent(&input.Student{
		Name:      "xyt",
		Age:       26,
		PeopleMap: map[int32]*input.People{23: {"b", 22}},
	})
	err := s.Insert(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", s)
}
func TestCreate(t *testing.T) {
	if err := Create([]interface{}{
		&input.User{},
		&input.Student{},
	}, "output"); err != nil {
		log.Fatal(err)
	}
}
