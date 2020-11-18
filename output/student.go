package output

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tool-mongo/input"
)

const (
	StudentFieldAge       = "age"
	StudentFieldPeopleMap = "people_map"
)

// Student 数据表内存模型
type StudentModel struct {
	input.Student `bson:",inline"`
	setter        map[string]interface{} `bson:"-"`
}

// 获取完整数据
func (m *StudentModel) Object() *input.Student {
	return &(m.Student)
}

// 获取 age 的值
func (m *StudentModel) GetAge() int32 {
	return m.Student.Age
}

// 获取 people_map 的值
func (m *StudentModel) GetPeopleMap() map[int32]*input.People {
	return m.Student.PeopleMap
}

// 设置 age 的值
func (m *StudentModel) SetAge(age int32) *StudentModel {
	m.setter[StudentFieldAge] = age
	m.Age = age
	return m
}

// 设置 people_map 的值
func (m *StudentModel) SetPeopleMap(peopleMap map[int32]*input.People) *StudentModel {
	m.setter[StudentFieldPeopleMap] = peopleMap
	m.PeopleMap = peopleMap
	return m
}

// 获取表名
func (m *StudentModel) TableName() string {
	if v, ok := interface{}(m.Student).(Table); ok {
		return v.TableName()
	}
	return "student"
}

func NewStudent(v *input.Student) *StudentModel {
	return &StudentModel{
		//MongoId: primitive.ObjectID{},
		Student: *v,
		setter:  map[string]interface{}{},
	}
}

// 插入新数据
func (m *StudentModel) Insert(ctx context.Context) error {
	result, err := GetCollection(m.TableName()).InsertOne(ctx, m.Object())
	if err != nil {
		return err
	}
	//fmt.Printf("StudentModel Insert result:%+v\n", result)
	m.Name = result.InsertedID.(string)
	return nil
}

// 删除数据
func (m *StudentModel) Delete(ctx context.Context) error {
	_, err := GetCollection(m.TableName()).DeleteOne(ctx, bson.D{{"_id", m.Name}})
	if err != nil {
		return err
	}
	//fmt.Printf("StudentModel Delete result:%+v\n", result)
	return nil
}

// 更新数据
func (m *StudentModel) Update(ctx context.Context) error {
	var update bson.D
	for key, value := range m.setter {
		update = append(update, bson.E{
			Key:   key,
			Value: value,
		})
	}
	_, err := GetCollection(m.TableName()).UpdateOne(ctx, bson.D{{"_id", m.Name}}, bson.D{{"$set", update}})
	if err != nil {
		return err
	}
	//fmt.Printf("StudentModel Update result:%+v\n", result)
	return nil
}

// 更新或新增
func (m *StudentModel) UpsertByName(ctx context.Context) error {
	result, err := GetCollection(m.TableName()).
		UpdateOne(ctx, bson.D{{"_id", m.Name}}, bson.D{{"$set", m.Object()}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	//fmt.Printf("StudentModel Update result:%+v\n", result)
	if result.MatchedCount != 0 {
		//fmt.Println("matched and replaced an existing document")
		return nil
	}
	if result.UpsertedCount != 0 {
		//fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
		m.Name = result.UpsertedID.(string)
		return nil
	}
	return nil
}

// 根据 Name 查询
func (m *StudentModel) GetStudentByName(ctx context.Context, name string) error {
	return GetCollection(m.TableName()).FindOne(ctx, bson.D{{"_id", name}}).Decode(m)
}
