package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// Create a new Elasticsearch client
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false))
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	ctx := context.Background()

	user := User{
		ID:   "1",
		Name: "John Doe",
		Age:  30,
	}

	// Index a document
	indexDoc(ctx, client, user)
	// Getting a document
	getDoc(ctx, client, "1")

	// Updating a document
	updateDoc(ctx, client, "1", 31)
	// Getting a document
	getDoc(ctx, client, "1")

	// Deleting a document
	deleteDoc(ctx, client, "1")
	// Getting a document
	getDoc(ctx, client, "1")

}

func indexDoc(ctx context.Context, client *elastic.Client, user User) {
	_, err := client.Index().
		Index("users").
		Id(user.ID).
		BodyJson(user).
		Do(ctx)
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}
	fmt.Printf("Indexed document %s \n", user.ID)
}

func getDoc(ctx context.Context, client *elastic.Client, id string) {
	getResult, err := client.Get().
		Index("users").
		Id(id).
		Do(ctx)
	if err != nil && !elastic.IsNotFound(err) {
		log.Fatalf("Error getting document: %s \n", err)
	}

	if getResult != nil && getResult.Found {
		var user User
		err := json.Unmarshal(getResult.Source, &user)
		if err != nil {
			log.Fatalf("Error unmarshalling document: %s", err)
		}
		fmt.Printf("Got document: %+v\n", user)
	} else {
		fmt.Println("Document not found")
	}
}

func updateDoc(ctx context.Context, client *elastic.Client, id string, newAge int) {
	_, err := client.Update().
		Index("users").
		Id(id).
		Doc(map[string]interface{}{
			"age": newAge,
		}).
		Do(ctx)
	if err != nil {
		log.Fatalf("Error updating document: %s", err)
	}
	fmt.Printf("Updated document %s \n", id)
}

func deleteDoc(ctx context.Context, client *elastic.Client, id string) {
	_, err := client.Delete().
		Index("users").
		Id(id).
		Do(ctx)
	if err != nil {
		log.Fatalf("Error deleting document: %s", err)
	}

	fmt.Printf("Deleted document %s \n", id)
}
