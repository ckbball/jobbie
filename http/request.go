package http

import (
  "net/http"

  "github.com/labstack/echo/v4"
  "github.com/mholt/binding"

  "github.com/ckbball/quik"
)

type userUpsertRequest struct {
  User struct {
    Email        string       `json:"email,omitempty"`
    Password     string       `json:"password,omitempty"`
    FirstName    string       `json:"first_name,omitempty"`
    LastName     string       `json:"last_name,omitempty"`
    JobSearch    int          `json:"job_search,omitempty"`
    Profile      quik.Profile `json:"profile,omitempty"`
    Applications []string     `json:"applications"`
    SavedJobs    []string     `json:"saved_jobs"`
  } `json:"user"`
}

func newUserUpsertRequest() *userUpsertRequest {
  return &userUpsertRequest{}
}

func (r *userUpsertRequest) populate(u *quik.User) {
  r.User.Email = u.Email
  r.User.Password = u.Password
  r.User.FirstName = u.FirstName
  r.User.LastName = u.LastName
  r.User.JobSearch = u.JobSearch
  r.User.Profile = u.Profile
}

type userRegisterRequest struct {
  User struct {
    FirstName string `json:"first_name" validate:"required"`
    LastName  string `json:"last_name" validate:"required"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required"`
  } `json:"user"`
}

func (u *userRegisterRequest) FieldMap(req *http.Request) binding.FieldMap {
  return binding.FieldMap{
    &u.User.Email:     "email",
    &u.User.Password:  "password",
    &u.User.FirstName: "first_name",
    &u.User.LastName:  "last_name",
  }
}

type userLoginRequest struct {
  Email    string `json:"email" validate:"required,email"`
  Password string `json:"password" validate:"required"`
}

func (r *userLoginRequest) bind(c echo.Context) error {
  if err := c.Bind(r); err != nil {
    return err
  }
  if err := c.Validate(r); err != nil {
    return err
  }
  return nil
}
