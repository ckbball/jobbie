package http

import (
  //"context"
  "net/http"

  "github.com/ckbball/quik"
  "github.com/go-chi/chi"
  "github.com/mholt/binding"
)

type userHandler struct {
  router chi.Router

  // Services
  userService quik.UserService
}

func newUserHandler() *userHandler {
  h := &userHandler{router: chi.NewRouter()}
  h.router.Post("/api/v1/signup", h.handleNewUser)
  return h
}

// ServeHTTP implements http.Handler
func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  h.router.ServeHTTP(w, r)
}

func (h *userHandler) handleNewUser(w http.ResponseWriter, r *http.Request) {
  // initialize variables
  var user quik.User
  req := &userRegisterRequest{}

  // bind request to request struct
  if err := binding.Bind(req); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    //log.Infof("Error in reading request body. line 27. createTeamHandler(). \nbody: %v", r.Body)
    return
  }

  // assign user the User from the request
  user = req.User

  // call the datastore method to create a new user in the datastore
  if err := h.userService.CreateUser(&u); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    //log.Infof("Error in reading request body. line 27. createTeamHandler(). \nbody: %v", r.Body)
    return
  }

  // marshall a successful response
  marshalledResp, err := json.Marshal(newUserResponse(&u))
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
