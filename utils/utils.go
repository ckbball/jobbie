package utils

import (
  "github.com/dgrijalva/jwt-go"

  "github.com/ckbball/quik"
)

const key = "hehehoohoo9wf"

// CustomClaims is our custom metadata, which will be hashed
// and sent as the second segment in our JWT
type CustomClaims struct {
  User *quik.User
  jwt.StandardClaims
}

// Decode a token string into a token object
func Decode(tokenString string) (*CustomClaims, error) {

  // Parse the token
  token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
    return key, nil
  })

  // Validate the token and return the custom claims
  if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
    return claims, nil
  } else {
    return nil, err
  }
}

// Encode a claim into a JWT
func Encode(user *quik.User) (string, error) {

  expireToken := time.Now().Add(time.Hour * 24).Unix()

  // Create the Claims
  claims := CustomClaims{
    user,
    jwt.StandardClaims{
      ExpiresAt: expireToken,
      Issuer:    "dev.user",
    },
  }

  // Create token
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

  // Sign token and return
  return token.SignedString(key)
}
