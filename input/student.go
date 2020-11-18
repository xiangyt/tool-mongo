package input

type Student struct {
	Name      string            `bson:"_id"`
	Age       int32             `bson:"age"`
	PeopleMap map[int32]*People `bson:"people_map"`
}

func (s *Student) TableName() string {
	return "student"
}
