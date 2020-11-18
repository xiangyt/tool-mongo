package input

type User struct {
	Id      string               `bson:"id"`
	Name    string               `bson:"user_name"`
	Age     int32                `bson:"age"`
	Money   float64              `bson:"money"`
	Friends []*People            `bson:"friends"`
	Houses  []string             `bson:"houses"`
	Map1    map[string]int32     `bson:"map1"`
	Map2    map[int32]*People    `bson:"map2"`
	Map3    map[string][]*People `bson:"map3"`
}

func (u User) TableName() string {
	return "atest"
}

type People struct {
	Name string `bson:"name"`
	Age  int32  `bson:"age"`
}
