package dao

import (
	"context"
	"fmt"
	"github.com/junhaideng/IPProxy/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var MongoDB *mongo.Collection

// 初始化数据库
func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	username := viper.GetString("database.mongodb.username")
	password := viper.GetString("database.mongodb.password")
	port := viper.GetString("database.mongodb.port")
	host := viper.GetString("database.mongodb.host")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)))
	if err != nil {
		panic(err)
	}
	database := viper.GetString("database.mongodb.db")
	collection := viper.GetString("database.mongodb.collection")
	fmt.Println(database, collection)
	MongoDB = client.Database(database).Collection(collection)
}

// 插入一个文档
func InsertOne(data interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return MongoDB.InsertOne(ctx, data, opts...)
}

// 同时插入多个文档
func InsertMany(data []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return MongoDB.InsertMany(ctx, data, opts...)
}

// 查找所有的文档
func GetAll() ([]model.IP, error) {
	var ips []model.IP
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := MongoDB.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &ips); err != nil {
		return nil, err
	}
	return ips, nil
}

// 获取到文档
func GetLimit(limit int64, filter interface{}, sort interface{}) ([]model.IP, error) {
	if limit <= 0 {
		return nil, nil
	}
	var ips []model.IP
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt := &options.FindOptions{
		Limit: &limit,
		Sort:  sort,
	}

	cursor, err := MongoDB.Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &ips); err != nil {
		return nil, err
	}
	return ips, nil
}

// 删除一个文档
func Delete(filter interface{}, opts ...*options.DeleteOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := MongoDB.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return err
	}
	return nil
}

// ip地址是否已经存在
func FindByIP(ipaddr string) (*model.IP, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var ip = &model.IP{}
	result, err := MongoDB.Find(ctx, bson.M{"ip": ipaddr})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ip":  ipaddr,
			"err": err,
		}).Error("查找ip地址发生错误")
		return nil, err
	}
	for result.Next(ctx) {
		if err := result.Decode(ip); err != nil {
			logrus.WithFields(logrus.Fields{
				"ip":  ipaddr,
				"err": err,
			}).Error("反序列化失败")
			return nil, err
		}
	}
	return ip, err
}

// ip地址是否存在
func ExistIP(ipaddr string) bool {
	ip, err := FindByIP(ipaddr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
			"ip":  ipaddr,
		})
		return false
	}
	if ip.IP != "" {
		return true
	}
	return false
}

func Update(filter interface{}, update interface{}, option ...*options.UpdateOptions) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := MongoDB.UpdateOne(ctx, filter, update, option...)
	return err
}
