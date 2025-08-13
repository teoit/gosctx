package mongodb

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/teoit/gosctx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrDBNamEmpty = errors.New("database name not empty")
)

type MongoComponent interface {
	GetMongoDB() *mongo.Client
	GetDatabase() *mongo.Database
}

type MongoOpt struct {
	dsn                    string
	dbName                 string
	serverSelectionTimeout int
	maxConnIdleTime        int
}

type mongoDB struct {
	id     string
	prefix string
	logger gosctx.Logger
	client *mongo.Client
	db     *mongo.Database
	*MongoOpt
}

func NewMongoDB(id, prefix string) *mongoDB {
	return &mongoDB{
		id:       id,
		prefix:   prefix,
		MongoOpt: new(MongoOpt),
	}
}

func (mdb *mongoDB) ID() string {
	return mdb.id
}

func (mdb *mongoDB) InitFlags() {
	prefix := mdb.prefix
	if mdb.prefix != "" {
		prefix += "-"
	}
	flag.StringVar(&mdb.dsn, fmt.Sprintf("%smongo-dsn", prefix), "", "Database dsn mongo")
	flag.StringVar(&mdb.dbName, fmt.Sprintf("%smongo-database", prefix), "", "Database Name mongo")
	flag.IntVar(&mdb.serverSelectionTimeout, fmt.Sprintf("%smongo-timeout", prefix), 60, "timeout connections to the mongodb - Default 60")
	flag.IntVar(&mdb.maxConnIdleTime, fmt.Sprintf("%smongo-maxconnidletime", prefix), 10, "maxconnidletime to the mongodb - Default 10")
}

func (mdb *mongoDB) Activate(_ gosctx.ServiceContext) error {
	mdb.logger = gosctx.GlobalLogger().GetLogger(mdb.id)
	if mdb.dbName == "" {
		return ErrDBNamEmpty
	}

	opt := options.Client().ApplyURI(mdb.dsn)
	opt.SetServerSelectionTimeout(time.Duration(mdb.serverSelectionTimeout) * time.Second)
	opt.SetMaxConnIdleTime(time.Duration(mdb.maxConnIdleTime) * time.Second)

	client, err := mongo.Connect(context.Background(), opt)
	if err != nil {
		return err
	}

	mdb.logger.Info("Monogodb is connected")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		mdb.logger.Error("err mongodn connect ping: ", err)
		return err
	}

	mdb.client = client
	mdb.db = client.Database(mdb.dbName)

	return nil
}

func (mdb *mongoDB) Stop() error {
	if mdb.client != nil {
		err := mdb.client.Disconnect(context.Background())
		if err != nil {
			mdb.logger.Error("err mongodb disconnect: ", err)
			return err
		}
	}
	return nil
}

func (mdb *mongoDB) GetMongoDB() *mongo.Client {
	return mdb.client
}

func (mdb *mongoDB) GetDatabase() *mongo.Database {
	return mdb.db
}
