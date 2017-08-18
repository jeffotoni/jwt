# Jwt example

The purpose of this project is to make a simplified version using JWT for educational purposes.

We will create a simplistic version for understanding, we will generate a token and the expiration date and send in json format to the requester.

The requester would have to have a way to register the user before using our authenticated API.

This part of the registry we will leave last, first we will simulate with keys that we create in the server for tests purposes.

We will Generate the Token and with the same token we will validate our handlers to see if they can be executed.

# Packages

go get -u github.com/dgrijalva/jwt-go

# Generate the keys

```sh

$ openssl genrsa -out private.rsa 1024

$ openssl rsa -in private.rsa -pubout > public.rsa.pub

```

# Install

$ go build jwt.go

$ sudo cp jwt /usr/bin

# Main function

```go



```