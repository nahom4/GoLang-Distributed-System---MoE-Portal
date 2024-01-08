package main

import (
	database "authServer2/config"
	"authServer2/controller"
	models "authServer2/model"
	"fmt"
	"net"
	"net/rpc"
)

type AuthServer int


func init(){
	database.ConnectDB()
}
func (c *AuthServer) RegisterUser(newUser *controller.NewUser, result *bool) error {
	fmt.Println("Registration is called from the auth-server")

	*result = controller.RegisterUser(*newUser)
	return nil
}
func (c *AuthServer) ValidateToken(token *string, result *bool) error {
	*result = controller.ValidateToken(*token)
	fmt.Println(controller.ValidateToken(*token))
	return nil
}
func (c *AuthServer) AuthenticateUser(user *controller.User, result *controller.LoginResult) error {
	loginResult, err := controller.Login(*user)
	if err != nil {
		return err
	}

	*result = loginResult
	return nil
}
func (c *AuthServer) RefreshToken(token *string, result *controller.LoginResult) error {
	loginResult, err := controller.Refresh(*token)
	if err != nil {
		return err
	}

	*result = loginResult
	return nil
}

func main() {
	database.DB.AutoMigrate(&models.User{})
	authServer := new(AuthServer)
	rpc.Register(authServer)

	listener, err := net.Listen("tcp", ":8002")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	fmt.Println("Server is listening on port 8002...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}