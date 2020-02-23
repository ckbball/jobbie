package http

import (
  //"github.com/labstack/echo/v4"

  "github.com/ckbball/quik"
  "github.com/ckbball/quik/utils"
)

type userResponse struct {
  User struct {
    Email     string       `json:"email,omitempty"`
    Token     string       `json:"token,omitempty"`
    FirstName string       `json:"first_name,omitempty"`
    LastName  string       `json:"last_name,omitempty"`
    JobSearch int          `json:"job_search,omitempty"`
    Profile   quik.Profile `json:"profile,omitempty"`
  } `json:"user"`
}

type loginResponse struct {
  Email string `json:"email"`
  Token string `json:"token"`
}

//
func newUserResponse(u *quik.User) *userResponse {
  var r userResponse
  r.User.Email = u.Email
  r.User.FirstName = u.FirstName
  r.User.LastName = u.LastName
  r.User.JobSearch = u.JobSearch
  r.User.Profile = u.Profile
  token, err := utils.Encode(u)
  if err != nil {
    return &r
  }
  r.User.Token = token
  return &r
}

func newLoginResponse(email, token string) *loginResponse {
  return &loginResponse{
    Email: email,
    Token: token,
  }
}
