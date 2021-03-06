/*
*
* Project printServer, an api rest responsible for printing to a Zebra thermal printer.
* The printServer will receive a cryptographic POST containing
* the Zpl content so that the printer can print.
*
* @package     auth
* @author      @jeffotoni
* @size        11/08/2017
*
 */

package models

//
// Claim structure, where we will use
// to validate our token with jwt
//
import "github.com/dgrijalva/jwt-go"

//
// jwt
//
type Claim struct {

	//
	//
	//
	User string `json:"user"`

	//
	//
	//
	jwt.StandardClaims
}
