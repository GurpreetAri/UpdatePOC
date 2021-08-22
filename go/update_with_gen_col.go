package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

const interleavedDBConnString string = "projects/anz-x-fabric-np-641432/instances/test-instance/databases/interleaved1"

type TransactionCategoryNew struct {
	TransactionID string `spanner:"Transaction_ID"`
}

func main() {
	log.Printf("client execution started at : %v", time.Now())
	ctx := context.Background()
	log.Println("client created...")
	client, err := spanner.NewClient(ctx, interleavedDBConnString)
	if err != nil {
		log.Println("error: client could not created..")
		return
	}
	defer client.Close()
	runWithGenCol(ctx, client)
}

func readAndPrintTransactionCategoryNew(ctx context.Context) {
	stmt := spanner.Statement{
		SQL: `SELECT Transaction_ID FROM TransactionCategoryNew;`,
	}

	client, err := spanner.NewClient(ctx, dbConnString)
	if err != nil {
		log.Println("error: client could not created for read..")
		return
	}
	defer client.Close()
	log.Printf("client created for read, reading record now: %v", time.Now())

	iter := client.Single().Query(ctx, stmt)
	row, err := iter.Next()
	if err != nil && err != iterator.Done {
		log.Println("error: failed while iterating over response: ", err)
		return
	}

	if row != nil {
		var existingRecord TransactionCategoryNew
		if err = row.ToStruct(&existingRecord); err != nil {
			log.Println("error: unable to parse response row: ", err)
			return
		}
		log.Printf("read transaction id: %s", existingRecord.TransactionID)
	}
	log.Printf("finished reading at: %v", time.Now())
}

func runWithGenCol(ctx context.Context, client *spanner.Client) {
	//Read 1 from table
	readAndPrintTransactionCategoryNew(ctx)


	log.Printf("client created for update, updating record now: %v", time.Now())
	//Update table
	stmt := spanner.Statement{
		SQL: `UPDATE TransactionCategoryNew
			  SET Recategorised_Category_ID = @newCategoryId,
			  Last_Update_Time = PENDING_COMMIT_TIMESTAMP()
			  WHERE Account_ID = @accountId
			  AND Transaction_ID = @transactionId;`,
		Params: map[string]interface{}{
			"newCategoryId": "new_cat2",
			"transactionId": "transaction_id1",
		},
	}
	count, err := client.PartitionedUpdate(ctx, stmt)
	if err != nil {
		log.Println("error: failed to update TransactionCategoryNew table: ", err)
		return
	}
	log.Printf("finished updating at : %v", time.Now())
	log.Print("no of records updated:", count)

	//Read 2 from table
	readAndPrintTransactionCategoryNew(ctx)
}
