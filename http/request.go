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

func (u *userUpsertRequest) FieldMap(req *http.Request) binding.FieldMap {
  return binding.FieldMap{
    &u.User.Email:        "email",
    &u.User.Password:     "password",
    &u.User.FirstName:    "first_name",
    &u.User.LastName:     "last_name",
    &u.User.JobSearch:    "job_search",
    &u.User.Profile:      "profile",
    &u.User.Applications: "applications",
    &u.User.SavedJobs:    "saved_jobs",
  }
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

func (r *userUpsertRequest) bind(c echo.Context, u *quik.User) error {
  if err := c.Bind(r); err != nil {
    return err
  }
  if err := c.Validate(r); err != nil {
    return err
  }
  u.Email = r.User.Email
  u.FirstName = r.User.FirstName
  u.LastName = r.User.LastName
  u.JobSearch = r.User.JobSearch
  u.Profile = r.User.Profile
  if r.User.Password != u.Password {
    h, err := u.HashPassword(r.User.Password)
    if err != nil {
      return err
    }

    u.Password = h
  }
  return nil
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

func (r *userRegisterRequest) bind(c echo.Context, u *quik.User) error {
  if err := c.Bind(r); err != nil {
    return err
  }
  if err := c.Validate(r); err != nil {
    return err
  }
  u.FirstName = r.User.FirstName
  u.LastName = r.User.LastName
  u.Email = r.User.Email
  h, err := u.HashPassword(r.User.Password)
  if err != nil {
    return err
  }
  u.Password = h
  return nil
}

type userLoginRequest struct {
  User struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
  } `json:"user"`
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
