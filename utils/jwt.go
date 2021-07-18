package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	infraCache "github.com/iamrz1/ab-auth/infra/cache"
	"log"
	"os"
	"strconv"
	"time"
)

type claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	UserType string `json:"user_type"`
	jwt.StandardClaims
}

func GenerateTokens(username, role, usertype string) (string, string) {
	t := time.Now()
	accessTokenValidityStr := os.Getenv("ACCESS_TOKEN_EXPIRATION_MINUTES")
	//refreshTokenValidityStr := os.Getenv("REFRESH_TOKEN_VALIDITY")
	accessTokenValidity, err := strconv.Atoi(accessTokenValidityStr)
	if err != nil {
		log.Println("ACCESS_TOKEN_EXPIRATION_MINUTES variable is not found in env")
		//return token, err
		accessTokenValidity = 30
	}

	refreshTokenValidityStr := os.Getenv("REFRESH_TOKEN_EXPIRATION_MINUTES")
	//refreshTokenValidityStr := os.Getenv("REFRESH_TOKEN_VALIDITY")
	refreshTokenValidity, err := strconv.Atoi(refreshTokenValidityStr)
	if err != nil {
		log.Println("REFRESH_TOKEN_EXPIRATION_MINUTES variable is not found in env")
		//return token, err
		refreshTokenValidity = 7 * 24 * 60
	}

	now := time.Now().UTC()
	accessExpTime := now.Add(time.Minute * time.Duration(accessTokenValidity))
	refreshExpTime := now.Add(time.Minute * time.Duration(refreshTokenValidity))
	// Create the JWT accessClaims, which includes the username and expiry time
	accessClaims := &claims{
		Username: username,
		Role:     role,
		UserType: usertype,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: accessExpTime.Unix(),
		},
	}

	refreshClaims := &claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: refreshExpTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the accessClaims
	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	jwtRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	// Create the JWT string
	accessTokenString, _ := jwtAccessToken.SignedString([]byte(accessTokenKey))
	refreshTokenString, _ := jwtRefreshToken.SignedString([]byte(refreshTokenKey))

	log.Printf("CreateTokens method took: %s \n", time.Now().Sub(t).String())

	return accessTokenString, refreshTokenString

}

func VerifyToken(token string, isRefresh bool) (*claims, error) {
	thisClaims := &claims{}
	secretKey := accessTokenKey
	if isRefresh {
		secretKey = refreshTokenKey
	}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, thisClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, fmt.Errorf("%s", "Invalid token")
	}

	return thisClaims, nil
}

func GetLastResetAt(username string) (int64, error) {
	scmd := infraCache.Client().Get(fmt.Sprintf("%s_%s", username, LastResetEventAtKey))
	err := scmd.Err()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	lastResetAt, err := scmd.Int64()
	if err != nil {
		log.Println("Unsupported value as last_reset_at")
		return 0, fmt.Errorf("%s", "Unsupported value as last_reset_at")
	}

	return lastResetAt, nil
}

func SetLastResetAt(username string, in int64) {
	infraCache.Client().Set(fmt.Sprintf("%s_%s", username, LastResetEventAtKey), in, 0)
}
