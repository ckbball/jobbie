package http

import (
  //"context"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"

  "github.com/go-chi/chi"
  "golang.org/x/crypto/bcrypt"

  "github.com/ckbball/quik"
  "github.com/ckbball/quik/utils"
)

type userHandler struct {
  router chi.Router

  // Services
  userService quik.UserService
}

func newUserHandler() *userHandler {
  h := &userHandler{router: chi.NewRouter()}
  h.router.Post("/signup", h.handleNewUser)
  h.router.Post("/login", h.handleLogin)
  return h
}

// ServeHTTP implements http.Handler
func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  h.router.ServeHTTP(w, r)
}

func (h *userHandler) handleNewUser(w http.ResponseWriter, r *http.Request) {
  // initialize variables
  user := &quik.User{}
  req := &userRegisterRequest{}

  // bind request to request struct
  reqBody, err := ioutil.ReadAll(r.Body)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    //log.Infof("Error in reading request body. line 27. createTeamHandler(). \nbody: %v", r.Body)
    return
  }

  // unmarshal json body into team request struct
  err = json.Unmarshal(reqBody, &req)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    //log.Infof("Error in unmarshalling body. line 35. createTeamHandler(). \nbody: %v", reqBody)
    return
  }

  hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    //log.Infof("Error in reading request body. line 27. createTeamHandler(). \nbody: %v", r.Body)
    return
  }
  user.Email = req.User.Email
  user.Password = string(hashedPass)
  user.FirstName = req.User.FirstName
  user.LastName = req.User.LastName
  fmt.Fprintf(os.Stderr, "user after bind: %s\n", req.User)

  // call the datastore method to create a new user in the datastore

  if err := h.userService.CreateUser(user); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    //log.Infof("Error in reading request body. line 27. createTeamHandler(). \nbody: %v", r.Body)
    return
  }

  // marshall a successful response
  marshalledResp, err := json.Marshal(newUserResponse(user))
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    //log.WithField("error", err).Error("marshall error")
    //log.Infof("Error in marshalling successful response. line 52. createTeamHandler(). \nerr: %v", err.Error())
    return
  }

  // write headers and the response
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  w.Write(marshalledResp)
}

func (h *userHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
  // initialize variables
  user := &quik.User{}
  req := &userLoginRequest{}

  // bind request to request struct
  reqBody, err := ioutil.ReadAll(r.Body)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    //log.Infof("Error in reading request body. line 27. createTeamHandler(). \nbody: %v", r.Body)
    return
  }

  // unmarshal json body into team request struct
  err = json.Unmarshal(reqBody, &req)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    //log.Infof("Error in unmarshalling body. line 35. createTeamHandler(). \nbody: %v", reqBody)
    return
  }

  fmt.Fprintf(os.Stderr, "user after bind: %s\n", req.Email)

  // find user by email
  user, err = h.userService.GetByEmail(req.Email)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    //log.Infof("Error in unmarshalling body. line 35. createTeamHandler(). \nbody: %v", reqBody)
    return
  }

  fmt.Fprintf(os.Stderr, "user after db: %s\n", user)

  // compare user's password from db with password from request, if match generate new token and return it
  // Compare given password to stored hash
  if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    fmt.Fprintf(os.Stderr, "err caught in compare and hash: \n")
    //log.Infof("Error in unmarshalling body. line 35. createTeamHandler(). \nbody: %v", reqBody)
    return
  }

  userModel := &quik.User{
    Id:    user.Id,
    Email: user.Email,
  }

  token, err := utils.Encode(userModel)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    fmt.Fprintf(os.Stderr, "error caught in token encode: \n")
    //log.Infof("Error in unmarshalling body. line 35. createTeamHandler(). \nbody: %v", reqBody)
    return
  }

  // marshall a successful response
  marshalledResp, err := json.Marshal(newLoginResponse(user.Email, token))
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    fmt.Fprintf(os.Stderr, "error caught in response marshall: \n")
    //log.WithField("error", err).Error("marshall error")
    //log.Infof("Error in marshalling successful response. line 52. createTeamHandler(). \nerr: %v", err.Error())
    return
  }

  // write headers and the response
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(marshalledResp)

}
