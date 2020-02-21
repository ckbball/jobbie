package quik

import (
  //"context"
  "errors"

  "go.mongodb.org/mongo-driver/bson/primitive"
  "golang.org/x/crypto/bcrypt"
)

type Company struct {
  Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  Email        string             `json:"email,omitempty" bson:"email,omitempty"`
  Password     string             `json:"password,omitempty" bson:"password,omitempty"`
  CompanyName  string             `json:"company_name,omitempty" bson:"company_name,omitempty"`
  Size         int                `json:"size,omitempty" bson:"size,omitempty"` // 0 - not looking, 1 - open to opportunities, 2 - active
  Profile      CompanyProfile     `json:"profile,omitempty" bson:"profile,omitempty"`
  Applications []string           `json:"applications,omitempty" bson:"applications,omitempty"` // holds the id's of applications
}

type CompanyProfile struct {
  Technologies []string // string list of technology used at company
  Openings     int
}

func (u *Company) HashPassword(pass string) (string, error) {
  if len(pass) == 0 {
    return "", errors.New("password empty")
  }
  hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
  return string(hash), err
}

func (u *Company) CheckPassword(pass string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
  return err == nil
}
