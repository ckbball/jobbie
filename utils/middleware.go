package main

import (
  "context"
  "fmt"
  "net/http"
  "strings"
  "time"

  "github.com/dgrijalva/jwt-go"
  "github.com/dgrijalva/jwt-go/request"
  "github.com/golang/glog"
  "github.com/google/uuid"
  "github.com/sirupsen/logrus"
)

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func AllowCORS(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if origin := r.Header.Get("Origin"); origin != "" {
      w.Header().Set("Access-Control-Allow-Origin", origin)
      if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
        preflightHandler(w, r)
        return
      }
    }
    h.ServeHTTP(w, r)
  })
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func PreflightHandler(w http.ResponseWriter, r *http.Request) {
  headers := []string{"Content-Type", "Accept", "Authorization"}
  w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
  methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
  w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
  glog.Infof("preflight request for %s", r.URL.Path)
}

// add auth middleware here

func UserAuth(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    token, err := request.ParseFromRequest(r, MyAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
      b := GetKey()
      return b, nil
    }, request.WithClaims(&CustomClaims{}))

    if err != nil {
      http.Error(w, http.StatusText(401), 401)
      return
    }

    claims, err := Decode(token.Raw)
    if err != nil {
      fmt.Println(err)
      http.Error(w, http.StatusText(401), 401)
      return
    }

    user := claims.User

    ctx := context.WithValue(r.Context(), "user", user)
    h.ServeHTTP(w, r.WithContext(ctx))
  })
}

func stripBearerPrefixFromTokenString(tok string) (string, error) {
  if len(tok) > 5 && strings.ToUpper(tok[0:6]) == "TOKEN " {
    return tok[6:], nil
  }
  return tok, nil
}

// Extract  token from Authorization header
// Uses PostExtractionFilter to strip "TOKEN " prefix from header
var AuthorizationHeaderExtractor = &request.PostExtractionFilter{
  request.HeaderExtractor{"Authorization"},
  stripBearerPrefixFromTokenString,
}

// Extractor for OAuth2 access tokens.  Looks in 'Authorization'
// header then 'access_token' argument for a token.
var MyAuth2Extractor = &request.MultiExtractor{
  AuthorizationHeaderExtractor,
  request.ArgumentExtractor{"access_token"},
}

type ctxKeyLog struct{}
type ctxKeyRequestID struct{}

type logHandler struct {
  log  *logrus.Logger
  next http.Handler
}

type responseRecorder struct {
  b      int
  status int
  w      http.ResponseWriter
}

func (r *responseRecorder) Header() http.Header { return r.w.Header() }

func (r *responseRecorder) Write(p []byte) (int, error) {
  if r.status == 0 {
    r.status = http.StatusOK
  }
  n, err := r.w.Write(p)
  r.b += n
  return n, err
}

func (r *responseRecorder) WriteHeader(statusCode int) {
  r.status = statusCode
  r.w.WriteHeader(statusCode)
}

func (lh *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  ctx := r.Context()
  requestID, _ := uuid.NewRandom()
  ctx = context.WithValue(ctx, ctxKeyRequestID{}, requestID.String())

  start := time.Now()
  rr := &responseRecorder{w: w}
  log := lh.log.WithFields(logrus.Fields{
    "http.req.path":   r.URL.Path,
    "http.req.method": r.Method,
    "http.req.id":     requestID.String(),
  })
  /*
     if v, ok := r.Context().Value(ctxKeySessionID{}).(string); ok {
       log = log.WithField("session", v)
     }*/
  log.Debug("request started")
  defer func() {
    log.WithFields(logrus.Fields{
      "http.resp.took_ms": int64(time.Since(start) / time.Millisecond),
      "http.resp.status":  rr.status,
      "http.resp.bytes":   rr.b}).Debugf("request complete")
  }()

  ctx = context.WithValue(ctx, ctxKeyLog{}, log)
  r = r.WithContext(ctx)
  lh.next.ServeHTTP(rr, r)
}

/*
func ensureSessionID(next http.Handler) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    var sessionID string
    c, err := r.Cookie(cookieSessionID)
    if err == http.ErrNoCookie {
      u, _ := uuid.NewRandom()
      sessionID = u.String()
      http.SetCookie(w, &http.Cookie{
        Name:   cookieSessionID,
        Value:  sessionID,
        MaxAge: cookieMaxAge,
      })
    } else if err != nil {
      return
    } else {
      sessionID = c.Value
    }
    ctx := context.WithValue(r.Context(), ctxKeySessionID{}, sessionID)
    r = r.WithContext(ctx)
    next.ServeHTTP(w, r)
  }
}
*/
