package mongo

import (
  "context"
  "fmt"
  "os"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  //"go.mongodb.org/mongo-driver/mongo"
  //"go.mongodb.org/mongo-driver/mongo/options"

  "github.com/ckbball/quik"
)

var _ quik.UserService = &UserService{}

type UserService struct {
  db *DB //
}

func NewUserService(client *DB) *UserService {
  return &UserService{
    db: client,
  }
}

func (s *UserService) GetByID(id string) (*quik.User, error) {
  primitiveId, _ := primitive.ObjectIDFromHex(id)

  var user quik.User
  err := s.db.ds.FindOne(context.TODO(), quik.User{Id: primitiveId}).Decode(&user)
  if err != nil {
    return nil, err
  }

  return &user, nil
}

func (s *UserService) GetByJobStatus(status int) (*quik.User, error) {

  var user quik.User
  err := s.db.ds.FindOne(context.TODO(), quik.User{JobSearch: status}).Decode(&user)
  if err != nil {
    return nil, err
  }

  return &user, nil
}

func (s *UserService) GetByEmail(email string) (*quik.User, error) {

  fmt.Fprintf(os.Stderr, "email after db: %s\n", email)

  var user quik.User
  err := s.db.ds.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
  if err != nil {
    return nil, err
  }

  return &user, nil
}

func (service *UserService) CreateUser(user *quik.User) error {
  // add a duplicate email and a duplicate username check

  insertUser := bson.D{
    {"email", user.Email},
    {"password", user.Password},
    {"first_name", user.FirstName},
    {"last_name", user.LastName},
  }

  /*result*/
  _, err := service.db.ds.InsertOne(context.TODO(), insertUser)

  if err != nil {
    return err
  }
  /*
     id := result.InsertedID
     w, _ := id.(primitive.ObjectID)

     out := w.Hex()
  */
  return err

}

func (service *UserService) UpsertUser(user *quik.User, id string) (int64, int64, error) {
  // add a duplicate email and a duplicate username check

  primitiveId, _ := primitive.ObjectIDFromHex(id)

  insertUser := bson.D{
    {"email", user.Email},
    {"password", user.Password},
    {"first_name", user.FirstName},
    {"last_name", user.LastName},
    {"job_search", user.JobSearch},
    {"profile", user.Profile},
    // in the future add other fieldb
  }

  result, err := service.db.ds.UpdateOne(context.TODO(),
    bson.D{
      {"_id", primitiveId},
    },
    bson.D{
      {"$set", insertUser},
    },
  )

  if err != nil {
    return -1, -1, err
  }

  return result.MatchedCount, result.ModifiedCount, nil
}

/*
func (service *UserService) Delete(id string) (int64, error) {
  primitiveId, _ := primitive.ObjectIDFromHex(id)
  filter := bson.D{{"_id", primitiveId}}

  result, err := Service.db.DeleteOne(context.TODO(), filter)
  if err != nil {
    return -1, err
  }
  return result.DeletedCount, nil
}
*/

/*
func (s *UserService) FilterUsers(req *v1.FindRequest) ([]*User, error) {

  findOptions := options.Find()
  findOptions.SetLimit(int64(req.Limit))
  findOptions.SetSort(bson.D{{"_id", -1}})
  findOptions.SetSkip(int64(req.Page))

  var users []*User
  cur, err := s.db.Find(context.TODO(),
    bson.D{
      {"experience", req.Experience},
    },
    findOptions)

  if err != nil {
    return nil, err
  }
  defer cur.Close(context.TODO())

  for cur.Next(context.TODO()) {
    var elem *User
    err := cur.Decode(&elem)
    if err != nil {
      return nil, err
    }

    users = append(users, elem)
  }

  if err := cur.Err(); err != nil {
    return users, err
  }

  return users, nil
}
*/
