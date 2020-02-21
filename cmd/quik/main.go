package main

import (
  "fmt"
  "os"

  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"

  "github.com/ckbball/quik"
  "github.com/ckbball/quik/http"
  myMongo "github.com/ckbball/quik/mongo"
)

var (
  port = "8000"
)

// Config is configuration for Server
type Config struct {

  // the port to listen for http calls
  HTTPPort string

  // DB Datastore parameters section
  // DatastoreDBHost is host of database
  DatastoreDBHost string
  // DatastoreDBUser is username to connect to database
  DatastoreDBUser string
  // DatastoreDBPassword password to connect to database
  DatastoreDBPassword string
  // DatastoreDBSchema is schema of database
  DatastoreDBSchema string
  // address for single redis node
  RedisAddress string
  // uri for mongodb
  MongoAddress string
  // name for mongodb
  MongoName string
  // collection for mongodb
  MongoCollection string

  // Log parameters section
  // LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
  LogLevel int
  // LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
  LogTimeFormat string
}

func main() {
  if err := run(os.Args, os.Stdout); err != nil {
    fmt.Fprintf(os.Stderr, "%s\n", err)
    os.Exit(exitFail)
  }
}

func run(args []string, stdout io.Writer) error {

  srvPort := port
  if os.Getenv("PORT") != "" {
    srvPort = os.Getenv("PORT")
  }
  addr := os.Getenv("LISTEN_ADDR")
  var cfg Config
  cfg.MongoAddress = os.Getenv("MONGO_URI")
  cfg.MongoName = os.Getenv("MONGO_NAME")

  // grab db info needed here from env vars

  sqlDB := db.NewSql() // pass info here
  db.AutoMigrate(sqlDB)

  // make client
  clientOptions := options.Client().ApplyURI(cfg.MongoAddress)
  client, err := mongo.Connect(context.TODO(), clientOptions)
  if err != nil {
    return err
  }
  // call mongo.NewMongoCollection
  uColl := myMongo.NewMongoCollection(cfg.MongoName, "users", client)

  // pass client to mongo services
  // Instantiate mongo services
  userService := myMongo.NewUserService(uColl)

  // Initialize HTTP server
  httpServer := http.NewServer()
  httpServer.Addr = addr + ":" + srvPort

  httpServer.UserService = userService

  // Open HTTP server
  if err := httpServer.Open(); err != nil {
    return err
  }
  fmt.Fprintf(m.Stdout, "http listening: %s\n", httpServer.Addr)

}

/*
func run(args []string, stdout io.Writer) error {
  r := router.New()
  v1 := r.Group("/api")

  srvPort := port
  if os.Getenv("PORT") != "" {
    srvPort = os.Getenv("PORT")
  }
  addr := os.Getenv("LISTEN_ADDR")

  // grab db info needed here from env vars

  sqlDB := db.NewSql() // pass info here
  db.AutoMigrate(sqlDB)

  // start mongodb here
  datastore := db.NewMongo()

  us := store.NewUserStore(datastore)
  cs := store.NewCompanyStore(datastore)
  js := store.NewJobStore(sqlDB)
  as := store.NewApplicationStore(sqlDB)

  h := handler.NewHandler(us, cs, js, as)
  h.Register(v1)
  r.Logger.Fatal(r.Start(addr + ":" + srvPort))
}
*/
