package db

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConfig struct {
	Driver       string
	Dbname       string
	Username     string
	Password     string
	Host         string
	Port         string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

var ConnDB *gorm.DB

func Migrate(db *gorm.DB) {
	//db.AutoMigrate(&models.User{}, &models.AccountFintech{}, &models.Transaction{}, &models.Log{})
	//if db.Error != nil {
	//	log.Fatal(db.Error)
	//}
}

func (conf DatabaseConfig) ConnectionDataBaseMain() {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Dbname,
	)
	var err error

	ConnDB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatal(err)
	}
	Migrate(ConnDB)

	sqlDB, err := ConnDB.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan instance *sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)

}

var Client *mongo.Client
var DBMongo *mongo.Database

func (conf DatabaseConfig) ConnectMongoDB() {
	//uri := fmt.Sprintf("%s://%s:%s", conf.Driver, conf.Host, conf.Port)
	//	uri := fmt.Sprintf("mongodb://shagya:shagyamongo09@103.139.192.137:27017/?authSource=admin")
	uri := fmt.Sprintf("mongodb://appuser:apppassword123@103.139.192.137:27017/shagya-tech?authSource=shagya-tech")

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("connection error mongodb:", err.Error())
	}

	fmt.Println("Connected to MongoDB Local!")
	db := client.Database(conf.Dbname)
	fmt.Println("Using database:", db.Name())
	Client = client
	DBMongo = db

}

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func ConnectRabbitMQ() {
	var err error
	RabbitConn, err = amqp.Dial("amqp://KwVUjT3EVRwRHNwB:H.MObcmUWz2IzL2eeumcOC.kb-hhaAld@switchyard.proxy.rlwy.net:16482")
	//RabbitConn, err = amqp.Dial("amqp://shagya:rabbitmqshagya09@rabbitmq:5672/") //rabbitmq 5672
	if err != nil {
		log.Println("Gagal menghubungkan ke RabbitMQ:", err)
	}

	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		log.Println("Gagal membuka channel RabbitMQ:", err)
	}

	_, err = RabbitChannel.QueueDeclare(
		"payment",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Gagal mendeklarasikan queue:", err)
	}

	log.Println("Terhubung ke RabbitMQ Payment")
}
