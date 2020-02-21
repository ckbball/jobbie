package http

import (
  //"context"
  "net/http"

  "github.com/ckbball/quik"
  "github.com/go-chi/chi"
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
  ctx := r.Context()
  var user quik.User
  req := &userRegisterRequest{}
  if err := req.bind(c, &u); err != nil {
    return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
  }
  if err := h.userService.CreateUser(&u); err != nil {
    return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
  }
  return c.JSON(http.StatusCreated, newUserResponse(&u))
}
