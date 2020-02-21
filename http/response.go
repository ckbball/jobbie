package http

import (
  "github.com/labstack/echo/v4"

  "github.com/ckbball/quik"
  "github.com/ckbball/quik/utils"
)

type userResponse struct {
  User struct {
    Email     string       `json:"email,omitempty"`
    Token  string       `json:"token,omitempty"`
    FirstName string       `json:"first_name,omitempty"`
    LastName  string       `json:"last_name,omitempty"`
    JobSearch int          `json:"job_search,omitempty"`
    Profile   quik.Profile `json:"profile,omitempty"`
  } `json:"user"`
}

func newUserResponse(u *model.User) *userResponse {
  return &userResponse{
    &User{
      Email: u.Email,
      FirstName: u.FirstName,
      LastName: u.LastName,
      JobSearch: u.JobSearch,
      Profile: u.Profile,
      Token: utils.Encode(u)
    }
  }
}
