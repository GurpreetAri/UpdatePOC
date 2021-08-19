### Context

We are using the `cloud.google.com/go/spanner` library to integrate with spanner in our application and this is a POC 
to demonstrate an issue we have been facing with it. 

### Problem Statement
We have a simple DB use case where we perform 2 operations:
1. A read operation and 
2. An update operation using partitionedDML

At the time of update, it appears as though the spanner connection is not released in time, because of which the 
use case fails with the `Error: context deadline exceeded` error.

__Please note that this behavior is observed only when using the `cloud.google.com/go/spanner` library.__
No issues have been observed when using Java or Python libraries or when executing the same partitioned DML statement on
spanner using `gcloud`.


### Steps to replicate the issue

1. Create a test database in spanner using this DDL stmt:

```sql
    CREATE TABLE TransactionCategoryNew (
      Account_ID STRING(36) NOT NULL,
      Transaction_ID STRING(36) NOT NULL,
      Category_Type STRING(100) NOT NULL,
      Original_Category_ID STRING(36) NOT NULL,
      Recategorised_Category_ID STRING(36),
      Last_Update_Time TIMESTAMP NOT NULL OPTIONS (
        allow_commit_timestamp = true
      ),
      Primary_Category_ID STRING(36) AS (IF(Recategorised_Category_ID IS NOT NULL, Recategorised_Category_ID, Original_Category_ID)) STORED,
    ) PRIMARY KEY(Account_ID, Transaction_ID, Category_Type)
```

2. Insert a sample record using this DML stmt:

```sql
    INSERT INTO TransactionCategoryNew (Account_ID, Transaction_ID, Category_Type, Last_Update_Time, Original_Category_ID)
    VALUES ('f6f3c93f-9fd0-5aeb-8599-238a88c5f906','00000000-0000-0000-0000-000000000015','Type1',PENDING_COMMIT_TIMESTAMP(),'0cbb1f59-01d2-55df-bbbe-6398004ae29a');
```

3. Navigate to the go and python stubs to execute them sequentially. Instructions to run them are enclosed within the 
   individual directories.