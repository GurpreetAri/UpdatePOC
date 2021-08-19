package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

const dbConnString = "projects/anz-x-fabric-np-641432/instances/test-instance/databases/example-db"

type TransactionCategory struct {
	TransactionID string `spanner:"Transaction_ID"`
}

func main() {
	log.Printf("stub execution started at: %v", time.Now())
	runWithoutGenCol()
	log.Printf("stub execution finished at : %v", time.Now())
}

func readAndPrintTransactionCategory(ctx context.Context) {
	stmt := spanner.Statement{
		SQL: `SELECT Transaction_ID FROM TransactionCategory;`,
	}

	client, err := spanner.NewClient(ctx, dbConnString)
	if err != nil {
		log.Println("error: client could not created..")
		return
	}
	defer client.Close()
	log.Printf("client created for read. Reading record now: %v", time.Now())

	iter := client.Single().Query(ctx, stmt)
	row, err := iter.Next()
	if err != nil && err != iterator.Done {
		log.Println("error: failed while iterating over response: ", err)
		return
	}

	if row != nil {
		var existingRecord TransactionCategory
		if err = row.ToStruct(&existingRecord); err != nil {
			log.Println("error: unable to parse response row: ", err)
			return
		}
		log.Printf("read Transaction ID successfully: %s", existingRecord.TransactionID)
	}
	log.Printf("finished reading at: %v", time.Now())
}

func runWithoutGenCol() {
	ctx := context.Background()

	//Read 1 from table
	readAndPrintTransactionCategory(ctx)

	client, err := spanner.NewClient(ctx, dbConnString)
	if err != nil {
		log.Println("error: client could not created..")
		return
	}
	defer client.Close()
	log.Printf("client created for update. Updating record now: %v", time.Now())

	//Update the table
	stmt := spanner.Statement{
		SQL: `UPDATE TransactionCategory
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
		log.Println("error: failed to update TransactionCategory table: ", err)
		return
	}

	log.Println("no of records updated:", count)

	//Read 2 from table
	readAndPrintTransactionCategory(ctx)
}
