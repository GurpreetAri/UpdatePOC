### Context

This doc outlines an issue with the spanner emulator when it is used to update an interleaved/non-interleaved table that also has a generated column.

### Problem Statement
We have a simple DB use case where we perform 2 operations:
1. A read operation and 
2. An update operation using `partitionedDML` on a table with a generated column. This table is interleaved under a parent table.
On updating the table, it appears that the spanner connection is not released in time, causing a `context deadline exceeded` error.

### Tools required

1. [Spanner emulator](https://cloud.google.com/spanner/docs/emulator) - You can pull the latest image from [here](https://console.cloud.google.com/gcr/images/cloud-spanner-emulator/GLOBAL/emulator).
2. [Spanner-cli](https://github.com/cloudspannerecosystem/spanner-cli)

### Steps to replicate the issue

Make sure you have go installed. (Tested using `go1.16.5` and `cloud.google.com/go/spanner 1.24.1`)

1. Create a test instance and db in spanner-emulator using these DDL statements:
   
```sql
    CREATE TABLE ParentTable (
                 Parent_ID STRING(36) NOT NULL,
                 User_ID STRING(36) NOT NULL,
                 Last_Update_Time TIMESTAMP NOT NULL OPTIONS ( allow_commit_timestamp = true ),
    ) PRIMARY KEY(Parent_ID, User_ID);
```

Table with generated column and NO interleaving:
```sql
   CREATE TABLE ChildTable (
               Parent_ID STRING(36) NOT NULL,
               User_ID STRING(36) NOT NULL,
               Level STRING(36) NOT NULL,
               Original_Child_ID STRING(36) NOT NULL,
               New_Child_ID STRING(36),
               Primary_Child_ID STRING(36) AS (IF(New_Child_ID IS NOT NULL, New_Child_ID, Original_Child_ID)) STORED,
               Last_Update_Time TIMESTAMP NOT NULL OPTIONS ( allow_commit_timestamp = true )
   ) PRIMARY KEY(Parent_ID, User_ID, Level);
```

Table with generated column AND interleaving:
```sql
   CREATE TABLE ChildTableInterleaved(
             Parent_ID STRING(36) NOT NULL,
             User_ID STRING(36) NOT NULL,
             Level STRING(36) NOT NULL,
             Original_Child_ID STRING(36) NOT NULL,
             New_Child_ID STRING(36),
             Primary_Child_ID STRING(36) AS (IF(New_Child_ID IS NOT NULL, New_Child_ID, Original_Child_ID)) STORED,
             Last_Update_Time TIMESTAMP NOT NULL OPTIONS ( allow_commit_timestamp = true )
   ) PRIMARY KEY(Parent_ID, User_ID, Level),
   INTERLEAVE IN PARENT ParentTable ON DELETE NO ACTION;
```
3. Insert a sample record in each table using these DML stmts:

```sql
   INSERT INTO ParentTable (Parent_ID, User_ID, Last_Update_Time)
   VALUES('parent_id', 'user_id', '2021-08-20T00:00:00Z');
   
   INSERT INTO ChildTable (Parent_ID, User_ID, Level, Last_Update_Time, Original_Child_ID)
   VALUES ('parent_id', 'user_id', 'Level', PENDING_COMMIT_TIMESTAMP(), 'original_child_id');
   
   INSERT INTO ChildTableInterleaved(Parent_ID, User_ID, Level, Last_Update_Time, Original_Child_ID)
   VALUES ('parent_id', 'user_id', 'Level', PENDING_COMMIT_TIMESTAMP(), 'original_child_id');
```

4. Modify `dbConnString` in the client files to point to your test db config.
5. Run individual go clients for testing against table with generated column and with or without interleaving to simulate the issue:

 ```bash
  go run update_gen_col_interleaving.go
```

```bash
  go run update_gen_col_no_interleaving.go
```

6. You will see that both clients wait on the update operation to complete, eventually failing with the `context deadline exceeded` error.

Here are some logs:

go % `go run update_gen_col_no_interleaving.go`
2021/08/24 09:30:22 client execution started at: 2021-08-24 09:30:22.937477 +1000 AEST m=+0.004029451
2021/08/24 09:30:22 client created for read, reading record now: 2021-08-24 09:30:22.938289 +1000 AEST m=+0.004841948
2021/08/24 09:30:23 read user id: user_id
2021/08/24 09:30:23 finished reading at: 2021-08-24 09:30:23.001076 +1000 AEST m=+0.067626458
2021/08/24 09:30:23 client created for update, updating record now: 2021-08-24 09:30:23.038443 +1000 AEST m=+0.104992965
2021/08/24 09:30:48 error: failed to update ChildTable table:  context deadline exceeded
2021/08/24 09:30:48 client execution finished at : 2021-08-24 09:30:48.100405 +1000 AEST m=+25.071050551

go % `go run update_gen_col_interleaving.go`
2021/08/24 09:32:09 client execution started at: 2021-08-24 09:32:09.573085 +1000 AEST m=+0.004799988
2021/08/24 09:32:09 client created for read, reading record now: 2021-08-24 09:32:09.574115 +1000 AEST m=+0.005830078
2021/08/24 09:32:09 read user id: user_id
2021/08/24 09:32:09 finished reading at: 2021-08-24 09:32:09.639838 +1000 AEST m=+0.071551068
2021/08/24 09:32:09 client created for update, updating record now: 2021-08-24 09:32:09.688079 +1000 AEST m=+0.119790658
2021/08/24 09:32:34 error: failed to update ChildTableInterleaved table:  context deadline exceeded
2021/08/24 09:32:34 client execution finished at : 2021-08-24 09:32:34.630369 +1000 AEST m=+25.061332033

__The issue is intermittent but is reproducible most times. You may have to re-try a few times to simulate it.__