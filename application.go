package quik

import (
  "context"
  "errors"

  "github.com/jinzhu/gorm"
)

type Application struct {
  gorm.Model
  Status       string
  JobID        string
  UserID       string
  AbilitiesMet []Ability
}

type ApplicationService interface {
  CreateApplication(ctx context.Context, app *Application) error
  UpsertApplication(ctx context.Context, app *Application) error
  GetApplicationsByJobID(ctx context.Context, id string) ([]*Application, error)
  GetApplicationsByUserID(ctx context.Context, id string) ([]*Application, error)
}

type Ability struct {
  gorm.Model
  Description string
}
