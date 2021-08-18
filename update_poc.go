package main

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/prometheus/common/log"
	"google.golang.org/api/iterator"
)

type CheckIfRecordExistsCount struct {
	Count spanner.NullInt64
}

func main() {
	ctx := context.Background()

	client, err := spanner.NewClient(ctx, "projects/anz-x-fabric-np-641432/instances/test-instance/databases/example-db")
	if err != nil {
		log.Info(ctx, "Error on creating client")
		return
	}
	defer client.Close()
	log.Info(ctx, "client created...reading data now..")

	//Read from table
	stmt := spanner.Statement{
		SQL: `SELECT count(*) as count from TransactionCategoryNew 
			  WHERE Account_ID = @accountId 
			  AND Transaction_ID = @transactionId;`,
		Params: map[string]interface{}{
			"accountId":     "f6f3c93f-9fd0-5aeb-8599-238a88c5f906",
			"transactionId": "00000000-0000-0000-0000-000000000015",
		},
	}

	iter := client.Single().Query(ctx, stmt)
	row, err := iter.Next()
	if err != nil && err != iterator.Done {
		log.Info(ctx, "iterator failed on reading response: ", err)
	}

	if row != nil {
		var existingRecord CheckIfRecordExistsCount
		if err = row.ToStruct(&existingRecord); err != nil {
			log.Info(ctx, "error: unable to parse row from check if records exists query: ", err)
			return
		}

		if existingRecord.Count.Int64 == 0 {
			log.Info(ctx, "record does not exist: ", err)
			return
		}
	}

	//Update the table
	stmt = spanner.Statement{
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
		log.Info(ctx, "error: failed to update table: ", err)
	}

	log.Info(ctx, "No of records updated:: ", count)
}
