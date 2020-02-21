package quik

import (
  //"context"
  "errors"

  "go.mongodb.org/mongo-driver/bson/primitive"
  "golang.org/x/crypto/bcrypt"
)

type User struct {
  Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  Email        string             `json:"email,omitempty" bson:"email,omitempty"`
  Password     string             `json:"password,omitempty" bson:"password,omitempty"`
  FirstName    string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
  LastName     string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
  JobSearch    int                `json:"job_search,omitempty" bson:"job_search,omitempty"` // 0 - not looking, 1 - open to opportunities, 2 - active
  Profile      Profile            `json:"profile,omitempty" bson:"profile,omitempty"`
  Applications []string           `json:"applications,omitempty" bson:"applications,omitempty"` // holds the id's of applications
  SavedJobs    []string           `json:"saved_jobs,omitempty" bson:"saved_jobs,omitempty"`     // holds the id's of jobs that they want to come back to later
  Secret       string             // this is the string they will post to their github profile to confirm that the account is theirs?
}

// need context when using SQL DB
type UserService interface {
  CreateUser(user *User) error
  UpsertUser(user *User) error
  GetByID(id string) (*User, error)
  GetByEmail(email string) (*User, error)
  GetByJobStatus(status int) (*User, error)
  // AddApplication
  // AddListing
  // RemoveApplication
  // RemoveListing
}

type Profile struct {
  Skills      []Skill   `json:"skills,omitempty" bson:"skills,omitempty"`
  WorkHistory []Work    `json:"work_history,omitempty" bson:"work_history,omitempty"`
  Projects    []Project `json:"projects,omitempty" bson:"projects,omitempty"`
}

type Skill struct {
  SkillName       string `json:"skill_name,omitempty" bson:"skill_name,omitempty"`
  ExperienceLevel string `json:"experience_level,omitempty" bson:"experience_level,omitempty"` // determined by how many projects you have used it on
}

type Work struct {
  CompanyName  string   `json:"company_name,omitempty" bson:"company_name,omitempty"`
  BulletPoints []string `json:"bullet_points,omitempty" bson:"bullet_points,omitempty"`
  Length       int      `json:"length,omitempty" bson:"length,omitempty"` // in months
}

type Project struct {
  ProjectName string   `json:"project_name,omitempty" bson:"project_name,omitempty"`
  Details     string   `json:"details,omitempty" bson:"details,omitempty"`
  Tech        []string `json:"tech,omitempty" bson:"tech,omitempty"`
  LiveLink    string   `json:"live_link,omitempty" bson:"live_link,omitempty"`
  GitLink     string   `json:"git_link,omitempty" bson:"git_link,omitempty"`
}

func (u *User) HashPassword(pass string) (string, error) {
  if len(pass) == 0 {
    return "", errors.New("password empty")
  }
  hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
  return string(hash), err
}

func (u *User) CheckPassword(pass string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
  return err == nil
}
