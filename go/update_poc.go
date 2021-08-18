package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

type Transaction struct {
	TransactionID string `spanner:"Transaction_ID"`
}

func main() {
	run()
	log.Printf("finished %v", time.Now())
}

func printTransactions(ctx context.Context, client *spanner.Client) {
	log.Printf("Reading now %v", time.Now())
	stmt := spanner.Statement{
		SQL: `SELECT Transaction_ID FROM TransactionCategoryNew;`,
	}
	//Reading 1
	iter := client.Single().Query(ctx, stmt)
	row, err := iter.Next()
	if err != nil && err != iterator.Done {
		log.Println("iterator failed on reading response: ", err)
	}

	if row != nil {
		var existingRecord Transaction
		if err = row.ToStruct(&existingRecord); err != nil {
			log.Println("error: unable to parse row from check if records exists query: ", err)
			return
		}

		log.Printf("Transaction ID: %s", existingRecord.TransactionID)
	}
	log.Printf("Finished reading %v", time.Now())
}

func run() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := spanner.NewClient(ctx, "projects/anz-x-fabric-np-641432/instances/test-instance/databases/example-db")
	if err != nil {
		log.Println("Error on creating client")
		return
	}
	defer client.Close()
	log.Println("client created...")

	//Read from table
	printTransactions(ctx, client)

	//Update the table
	log.Printf("Updating the category id  %v", time.Now())
	stmt := spanner.Statement{
		SQL: `UPDATE TransactionCategoryNew
			  SET Recategorised_Category_ID = @newCategoryId,
			  Last_Update_Time = PENDING_COMMIT_TIMESTAMP()
			  WHERE Account_ID = @accountId
			  AND Transaction_ID = @transactionId;`,
		Params: map[string]interface{}{
			"newCategoryId": "055a6b93-cf12-5637-9c65-d44f61615e97",
			"accountId":     "f6f3c93f-9fd0-5aeb-8599-238a88c5f906",
			"transactionId": "00000000-0000-0000-0000-000000000015",
		},
	}
	count, err := client.PartitionedUpdate(ctx, stmt)
	if err != nil {
		log.Println("error: failed to update table: ", err)
	}

	log.Println("No of records updated:: ", count)

	// Reading 2
	printTransactions(ctx, client)
}
