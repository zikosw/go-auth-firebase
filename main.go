package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	opt := option.WithCredentialsFile("./credentials/firebaseadmin.json")
	fba, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(fmt.Errorf("error initializing app: %v", err))
	}
	authCl, err := fba.Auth(context.Background())
	if err != nil {
		panic(err)
	}

	if false {
		if err := createUser(authCl, context.Background()); err != nil {
			panic(fmt.Errorf("create user - %w", err))
		}
		listUsers(authCl, context.Background())
	}
	createAuthToken(authCl, user1UID, context.Background())
	// verifyToken(authCl, user2WebToken, context.Background())

	fmt.Println("done")
}

const user1UID = "wXII387hwZQRm6UvxCCpGBgGsoq2"
const user2UID = "l4aLXdXXLMSyFysajImfPVLLUdl1"

func createAuthToken(client *auth.Client, uid string, ctx context.Context) {

	token, err := client.CustomToken(ctx, uid)
	if err != nil {
		log.Fatalf("error minting custom token: %v\n", err)
	}

	log.Printf("Got custom token: %v\n", token)
}
func verifyToken(client *auth.Client, idToken string, ctx context.Context) {
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	log.Printf("Verified ID token: %v\n", token)
}

func createUser(authClient *auth.Client, ctx context.Context) error {
	params := (&auth.UserToCreate{}).
		Email("user3@zik.ooo").
		EmailVerified(false).
		Password("password2").
		DisplayName("User Two").
		PhotoURL("https://images.dog.ceo/breeds/mountain-bernese/n02107683_7183.jpg"). //"https://images.dog.ceo/breeds/beagle/n02088364_12973.jpg").
		Disabled(false)
	u, err := authClient.CreateUser(ctx, params)
	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
		return err
	}
	log.Printf("Successfully created user: %v\n", u)
	spew.Dump(u)
	return nil
}

func getUser(client *auth.Client, ctx context.Context, uid string) {

	u, err := client.GetUser(ctx, uid)
	if err != nil {
		log.Fatalf("error getting user %s: %v\n", uid, err)
	}
	log.Printf("Successfully fetched user data: %v\n", u)
}

func findUsers(client *auth.Client, ctx context.Context) {
	getUsersResult, err := client.GetUsers(ctx, []auth.UserIdentifier{
		// auth.UIDIdentifier{UID: "uid1"},
		// auth.EmailIdentifier{Email: "user@example.com"},
		// auth.PhoneIdentifier{PhoneNumber: "+15555551234"},
		auth.ProviderIdentifier{ProviderID: "password", ProviderUID: "user3@zik.ooo"},
	})
	if err != nil {
		log.Fatalf("error retriving multiple users: %v\n", err)
	}

	log.Printf("Successfully fetched user data:")
	for _, u := range getUsersResult.Users {
		log.Printf("%v", u)
	}

	log.Printf("Unable to find users corresponding to these identifiers:")
	for _, id := range getUsersResult.NotFound {
		log.Printf("%v", id)
	}

}

func listUsers(client *auth.Client, ctx context.Context) {
	// // Note, behind the scenes, the Users() iterator will retrive 1000 Users at a time through the API
	// iter := client.Users(ctx, "")
	// for {
	// 	user, err := iter.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("error listing users: %s\n", err)
	// 	}
	// 	log.Printf("read user user: %v\n", user)
	// }

	// Iterating by pages 100 users at a time.
	// Note that using both the Next() function on an iterator and the NextPage()
	// on a Pager wrapping that same iterator will result in an error.
	pager := iterator.NewPager(client.Users(ctx, ""), 100, "")
	for {
		var users []*auth.ExportedUserRecord
		nextPageToken, err := pager.NextPage(&users)
		if err != nil {
			log.Fatalf("paging error %v\n", err)
		}
		for _, u := range users {
			log.Printf("read user user: %v - %+v\n", u, u.UserInfo)
		}
		if nextPageToken == "" {
			break
		}
	}
}
