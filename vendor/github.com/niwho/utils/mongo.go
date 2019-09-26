package utils

import (
	"fmt"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Imgo interface {
	SaveRecord(collection, rd string) error
	UpdateRecord(collection string, id int64) error
	DelRecord(collection string, id int64) error
	GetRecord(collection string, id int64) (bson.M, error)
	LoadReord(collection string, from, to int64) ([]bson.M, error)
}

type MongoDb struct {
	session *mgo.Session
	db      string
}

func InitMongo(user, pwd, dbname string, servers []string) (*MongoDb, error) {
	mdb := &MongoDb{db: dbname}
	var err error
	connStr := fmt.Sprintf("mongodb://%s:%s@%s/%s",
		user, pwd, strings.Join(servers, ","), dbname)
	//connStr = fmt.Sprintf("mongodb://%s/%s", strings.Join(servers, ","), dbname)
	mdb.session, err = mgo.Dial(connStr)
	if err == nil {
		mdb.session.SetMode(mgo.Monotonic, true)
	}
	return mdb, err
}

func (mdb *MongoDb) Insert(col string, record interface{}) error {
	cloned := mdb.session.Clone()
	defer func() {
		cloned.Close()
	}()
	return cloned.DB(mdb.db).C(col).Insert(record)
}

func (mdb *MongoDb) Update(col string, selector interface{}, update interface{}) error {
	cloned := mdb.session.Clone()
	defer func() {
		cloned.Close()
	}()
	return cloned.DB(mdb.db).C(col).Update(selector, update)
}

func (mdb *MongoDb) Upsert(col string, selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	cloned := mdb.session.Clone()
	defer func() {
		cloned.Close()
	}()
	info, err = cloned.DB(mdb.db).C(col).Upsert(selector, update)
	return
}

func (mdb *MongoDb) Find(col string, query interface{}, result interface{}, sortFileds ...string) error {
	cloned := mdb.session.Clone()
	defer func() {
		cloned.Close()
	}()
	return cloned.DB(mdb.db).C(col).Find(query).Sort(sortFileds...).Limit(99).All(result)
}

func (mdb *MongoDb) FindX(col string, query interface{}, result interface{}, count int, sortFileds ...string) error {
	cloned := mdb.session.Clone()
	defer func() {
		cloned.Close()
	}()
	return cloned.DB(mdb.db).C(col).Find(query).Sort(sortFileds...).Limit(count).All(result)
}

func (mdb *MongoDb) FindAll(col string, query interface{}, result interface{}, sortFileds ...string) error {
	cloned := mdb.session.Clone()
	defer func() {
		cloned.Close()
	}()
	return cloned.DB(mdb.db).C(col).Find(query).Sort(sortFileds...).All(result)
}

func (mdb *MongoDb) FindOne(col string, query interface{}, result interface{}) error {
	cloned := mdb.session.Clone()
	defer func() {
		cloned.Close()
	}()
	return cloned.DB(mdb.db).C(col).Find(query).One(result)
}
func (mdb *MongoDb) Remove(col string, selector interface{}) error {
	cloned := mdb.session.Clone()
	defer func() {
		cloned.Close()
	}()

	return cloned.DB(mdb.db).C(col).Remove(selector)

}
