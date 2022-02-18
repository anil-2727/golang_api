package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = []byte("gosecretkey")

// var client *mongo.Client

// var coll = client.Database("admin_panel").Collection("users")
// var user = users{ email : "anil@gmail.com", password : "12345@abc"}
// _,err := coll.InsertOne(context.TODO(), user)

// coll := client.Database("admin_panel").Collection("dashboard")
// dashboard1 := dashboard{ companyName : "abc", tinNumber : "123456", numberOfEmployees : 200, freeTrailPlan : "7 days", subscriptionPlan : "1 month", address : "madhapur", contactDetails : "9156843266" }
// _, err = coll.InsertOne(context.TODO(), dashboard1)

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		fmt.Println(hash)
	}
	return string(hash)

}

func GenerateJWT(user map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	claims := token.Claims.(jwt.MapClaims)
	claims["_id"] = user["_id"].(primitive.ObjectID)
	claims["email"] = user["email"].(string)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	fmt.Println(bearToken)
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}
