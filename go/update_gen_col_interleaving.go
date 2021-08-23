package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

// Replace this with your own connection string
const dbConnString = "projects/test-project/instances/test-instance/databases/test-database"

type ChildTableInterleaved struct {
	UserID string `spanner:"User_ID"`
}

func main() {
	log.Printf("client execution started at: %v", time.Now())
	runGenColWithInterleaving()
	log.Printf("client execution finished at : %v", time.Now())
}

func readAndPrintChildTableInterleaved(ctx context.Context) {
	stmt := spanner.Statement{
		SQL: `SELECT User_ID FROM ChildTableInterleaved;`,
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
		var existingRecord ChildTableInterleaved
		if err = row.ToStruct(&existingRecord); err != nil {
			log.Println("error: unable to parse response row: ", err)
			return
		}
		log.Printf("read user id: %s", existingRecord.UserID)
	}
	log.Printf("finished reading at: %v", time.Now())
}

func runGenColWithInterleaving() {
	ctx := context.Background()

	//Read 1 from table
	readAndPrintChildTableInterleaved(ctx)

	client, err := spanner.NewClient(ctx, dbConnString)
	if err != nil {
		log.Println("error: client could not created for update..")
		return
	}
	defer client.Close()
	log.Printf("client created for update, updating record now: %v", time.Now())

	//Update table
	stmt := spanner.Statement{
		SQL: `UPDATE ChildTableInterleaved
			  SET New_Child_ID = @newChildId,
			  Last_Update_Time = PENDING_COMMIT_TIMESTAMP()
			  WHERE User_ID = @userId;`,
		Params: map[string]interface{}{
			"newChildId": "new_child_id",
			"userId":     "user_id",
		},
	}
	count, err := client.PartitionedUpdate(ctx, stmt)
	if err != nil {
		log.Println("error: failed to update ChildTableInterleaved table: ", err)
		return
	}

	log.Printf("finished updating at : %v", time.Now())
	log.Print("no of records updated:", count)

	//Read 2 from table
	readAndPrintChildTableInterleaved(ctx)
}
