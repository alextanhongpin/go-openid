package oidc_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/SermoDigital/jose/jwt"
	"github.com/stretchr/testify/assert"

	"github.com/alextanhongpin/go-openid"
)

// const (
//         privKeyPath = "keys/app.rsa"     // openssl genrsa -out keys/app.rsa 1024
//         pubKeyPath  = "keys/app.rsa.pub" // openssl rsa -in keys/app.rsa -pubout > keys/app.rsa.pub
// )
//
// var signKey []byte
//
// func init() {
//         var err error
//         signKey, err = ioutil.ReadFile(privKeyPath)
//         if err != nil {
//                 log.Fatal("Error reading private key")
//                 os.Exit(1)
//         }
// }
//
// func createJWT() string {
//         claims := jws.Claims{}
//         // claims.Set("AccessToken", "level1")
//         signMethod := jws.GetSigningMethod("HS512")
//         token := jws.NewJWT(claims, signMethod)
//         byteToken, err := token.Serialize(signKey)
//         if err != nil {
//                 log.Fatal("Error signing the key. ", err)
//                 os.Exit(1)
//         }
//
//         return string(byteToken)
// }

func TestJWT(t *testing.T) {
	assert := assert.New(t)

	token := &oidc.IDToken{
		Issuer:   "server",
		Audience: "domr",
	}
	b, err := json.Marshal(token)
	assert.Nil(err)
	log.Println(string(b))

	var j map[string]interface{}
	err = json.Unmarshal(b, &j)
	jj := jwt.Claims(j)
	log.Println(err, jj)
	log.Println(jj.Issuer())

	// ttoken := createJWT()
	// fmt.Println("Created token", ttoken)
}
