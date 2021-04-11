package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/environment"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/auth"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/conf"
	"github.com/dapperAuteur/dashboard-go-api/internal/platform/database"
	"github.com/dapperAuteur/dashboard-go-api/internal/user"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run() error {

	// ==
	// Configuration

	var cfg struct {
		DB struct {
			// AtlasURI string `conf:"default:"`
			AtlasURI string `conf:"default:"`
		}
		Args conf.Args
	}

	// ==
	// Get Configuration
	// Helpful info in case of error
	if err := conf.Parse(os.Args[1:], "DASHBOARD", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("DASHBOARD", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// This is used for multiple commands below.
	dbConfig := database.Config{
		AtlasURI: environment.MongoDBURI,
	}

	var err error
	switch cfg.Args.Num(0) {
	case "useradd":
		err = useradd(dbConfig, cfg.Args.Num(1), cfg.Args.Num(2))
	case "keygen":
		err = keygen(cfg.Args.Num(1))
	default:
		err = errors.New("Must specify a command from the list: 'adduser', 'keygen'")
	}

	// print config values when app starts
	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	// ==
	// Start Database
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// client, err := database.Open(dbConfig)
	// if err != nil {
	// 	return errors.Wrap(err, "connecting to db")
	// }
	// defer client.Disconnect(ctx)

	return nil

}

func useradd(cfg database.Config, email, password string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	if email == "" || password == "" {
		return errors.New("useradd command must be called with two additional arguments for email and password")
	}

	fmt.Printf("Admin user will be created with email %q and password %q\n", email, password)
	fmt.Print("Continue? (1/0)")

	var confirm bool
	if _, err := fmt.Scanf("%t\n", &confirm); err != nil {
		return errors.Wrap(err, "processing response")
	}

	if !confirm {
		fmt.Println("Canceling")
		return nil
	}

	nu := user.NewUser{
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{auth.RoleAdmin, auth.RoleUser},
	}

	// send the db to the handler and let the router determine which collection to use
	myDatabase := client.Database(("quickstart")) // development database
	// myDatabase := client.Database(("palabras-express-api")) // production database

	userCollection := myDatabase.Collection("users")

	u, err := user.CreateUser(ctx, userCollection, nu, time.Now())
	if err != nil {
		return err
	}

	fmt.Println("User created with _id:", u.ID)
	return nil
}

// keygen creates an x509 private key for signing auth tokens.
func keygen(path string) error {

	if path == "" {
		return errors.New("keygen missing argument for key path")
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return errors.Wrap(err, "generating keys")
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "creating private file")
	}
	defer file.Close()

	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	if err := pem.Encode(file, &block); err != nil {
		return errors.Wrap(err, "encoding to private file")
	}

	if err := file.Close(); err != nil {
		return errors.Wrap(err, "encoding to private file")
	}

	return nil
}
