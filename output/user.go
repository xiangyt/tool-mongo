package output

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tool-mongo/input"
)

const (
	UserFieldId      = "id"
	UserFieldName    = "user_name"
	UserFieldAge     = "age"
	UserFieldMoney   = "money"
	UserFieldFriends = "friends"
	UserFieldHouses  = "houses"
	UserFieldMap1    = "map1"
	UserFieldMap2    = "map2"
	UserFieldMap3    = "map3"
)

// User 数据表内存模型
type UserModel struct {
	MongoId    primitive.ObjectID `bson:"_id,omitempty"`
	input.User `bson:",inline"`
	setter     map[string]interface{} `bson:"-"`
}

// 获取完整数据
func (m *UserModel) Object() *input.User {
	return &(m.User)
}

// 获取 id 的值
func (m *UserModel) GetId() string {
	return m.User.Id
}

// 获取 user_name 的值
func (m *UserModel) GetName() string {
	return m.User.Name
}

// 获取 age 的值
func (m *UserModel) GetAge() int32 {
	return m.User.Age
}

// 获取 money 的值
func (m *UserModel) GetMoney() float64 {
	return m.User.Money
}

// 获取 friends 的值
func (m *UserModel) GetFriends() []*input.People {
	return m.User.Friends
}

// 获取 houses 的值
func (m *UserModel) GetHouses() []string {
	return m.User.Houses
}

// 获取 map1 的值
func (m *UserModel) GetMap1() map[string]int32 {
	return m.User.Map1
}

// 获取 map2 的值
func (m *UserModel) GetMap2() map[int32]*input.People {
	return m.User.Map2
}

// 获取 map3 的值
func (m *UserModel) GetMap3() map[string][]*input.People {
	return m.User.Map3
}

// 设置 id 的值
func (m *UserModel) SetId(id string) *UserModel {
	m.setter[UserFieldId] = id
	m.Id = id
	return m
}

// 设置 user_name 的值
func (m *UserModel) SetName(userName string) *UserModel {
	m.setter[UserFieldName] = userName
	m.Name = userName
	return m
}

// 设置 age 的值
func (m *UserModel) SetAge(age int32) *UserModel {
	m.setter[UserFieldAge] = age
	m.Age = age
	return m
}

// 设置 money 的值
func (m *UserModel) SetMoney(money float64) *UserModel {
	m.setter[UserFieldMoney] = money
	m.Money = money
	return m
}

// 设置 friends 的值
func (m *UserModel) SetFriends(friends []*input.People) *UserModel {
	m.setter[UserFieldFriends] = friends
	m.Friends = friends
	return m
}

// 设置 houses 的值
func (m *UserModel) SetHouses(houses []string) *UserModel {
	m.setter[UserFieldHouses] = houses
	m.Houses = houses
	return m
}

// 设置 map1 的值
func (m *UserModel) SetMap1(map1 map[string]int32) *UserModel {
	m.setter[UserFieldMap1] = map1
	m.Map1 = map1
	return m
}

// 设置 map2 的值
func (m *UserModel) SetMap2(map2 map[int32]*input.People) *UserModel {
	m.setter[UserFieldMap2] = map2
	m.Map2 = map2
	return m
}

// 设置 map3 的值
func (m *UserModel) SetMap3(map3 map[string][]*input.People) *UserModel {
	m.setter[UserFieldMap3] = map3
	m.Map3 = map3
	return m
}

// 获取表名
func (m *UserModel) TableName() string {
	if v, ok := interface{}(m.User).(Table); ok {
		return v.TableName()
	}
	return "user"
}

func NewUser(v *input.User) *UserModel {
	return &UserModel{
		//MongoId: primitive.ObjectID{},
		User:   *v,
		setter: map[string]interface{}{},
	}
}

// 插入新数据
func (m *UserModel) Insert(ctx context.Context) error {
	result, err := GetCollection(m.TableName()).InsertOne(ctx, m.Object())
	if err != nil {
		return err
	}
	//fmt.Printf("UserModel Insert result:%+v\n", result)
	m.MongoId = result.InsertedID.(primitive.ObjectID)
	return nil
}

// 删除数据
func (m *UserModel) Delete(ctx context.Context) error {
	_, err := GetCollection(m.TableName()).DeleteOne(ctx, bson.D{{"_id", m.MongoId}})
	if err != nil {
		return err
	}
	//fmt.Printf("UserModel Delete result:%+v\n", result)
	return nil
}

// 更新数据
func (m *UserModel) Update(ctx context.Context) error {
	var update bson.D
	for key, value := range m.setter {
		update = append(update, bson.E{
			Key:   key,
			Value: value,
		})
	}
	_, err := GetCollection(m.TableName()).UpdateOne(ctx, bson.D{{"_id", m.MongoId}}, bson.D{{"$set", update}})
	if err != nil {
		return err
	}
	//fmt.Printf("UserModel Update result:%+v\n", result)
	return nil
}

// 更新或新增
func (m *UserModel) UpsertById(ctx context.Context) error {
	result, err := GetCollection(m.TableName()).
		UpdateOne(ctx, bson.D{{"id", m.Id}}, bson.D{{"$set", m.Object()}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	//fmt.Printf("UserModel Update result:%+v\n", result)
	if result.MatchedCount != 0 {
		//fmt.Println("matched and replaced an existing document")
		return nil
	}
	if result.UpsertedCount != 0 {
		//fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
		m.MongoId = result.UpsertedID.(primitive.ObjectID)
		return nil
	}
	return nil
}

// 根据 Id 查询
func (m *UserModel) GetUserById(ctx context.Context, id string) error {
	return GetCollection(m.TableName()).FindOne(ctx, bson.D{{"id", id}}).Decode(m)
}
