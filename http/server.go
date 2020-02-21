package http

import (
  "io"
  "io/ioutil"
  "net"
  "net/http"

  "github.com/ckbball/quik"
  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"
  "golang.org/x/crypto/acme/autocert"
)

// Server represents an HTTP server.
type Server struct {
  ln net.Listener

  // Services
  UserService quik.UserService

  // Server options.
  Addr        string // bind address
  Host        string // external hostname
  Autocert    bool   // ACME autocert
  Recoverable bool   // panic recovery

  LogOutput io.Writer
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
  return &Server{
    Recoverable: true,
    LogOutput:   ioutil.Discard,
  }
}

// Open opens the server.
func (s *Server) Open() error {
  // Open listener on specified bind address.
  // Use HTTPS port if autocert is enabled.
  if s.Autocert {
    s.ln = autocert.NewListener(s.Host)
  } else {
    ln, err := net.Listen("tcp", s.Addr)
    if err != nil {
      return err
    }
    s.ln = ln
  }

  // Start HTTP server.
  go http.Serve(s.ln, s.Router())

  return nil
}

// Close closes the socket.
func (s *Server) Close() error {
  if s.ln != nil {
    s.ln.Close()
  }
  return nil
}

func (s *Server) Router() http.Handler {
  r := chi.NewRouter()

  // Attach router middleware.
  r.Use(middleware.RealIP)
  r.Use(middleware.Logger)
  if s.Recoverable {
    r.Use(middleware.Recoverer)
  }
  // r.Mount("/debug", middleware.Profiler())

  //r.Use(s.attachLogOutputToContext)
  //r.Use(s.detectAccept)

  // Create API routes.
  r.Route("/api/v1/", func(r chi.Router) {
    //r.Use(middleware.DefaultCompress)
    //r.Get("/ping", s.handlePing)
    r.Mount("/users", s.userHandler())
  })

  return r
}

func (s *Server) userHandler() *userHandler {
  h := newUserHandler()
  h.userService = s.UserService
  return h
}
