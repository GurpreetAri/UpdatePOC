### Context

This doc outlines an issue with the spanner emulator when it is used to update an interleaved table that also has a generated column.

### Problem Statement
We have a simple DB use case where we perform 2 operations:
1. A read operation and 
2. An update operation using partitionedDML on a table with a generated column. This table is interleaved under a parent table.
 

Emulator is able to handle the update successfully when:
1. The table being updated is not interleaved OR
2. The table being updated does not have a generated column.

However, when the table is interleaved and has a generated column, the update fails. 
The spanner connection is not released in time, because of which the use case fails with the `Error: context deadline exceeded` error.


### Tools required

1. [Spanner emulator](https://cloud.google.com/spanner/docs/emulator) - You can pull the latest image from [here](https://console.cloud.google.com/gcr/images/cloud-spanner-emulator/GLOBAL/emulator).
1. [Spanner-cli](https://github.com/cloudspannerecosystem/spanner-cli)

### Steps to replicate the issue

1. Create a new database in spanner-emulator and create test tables using these DDL stmts:

```sql
    CREATE TABLE ParentTable (
                                Parent_ID STRING(36) NOT NULL,
                                User_ID STRING(36) NOT NULL,
                                Last_Update_Time TIMESTAMP NOT NULL OPTIONS (
                                   allow_commit_timestamp = true
                                   ),
    ) PRIMARY KEY(Parent_ID, User_ID);
```

Child table with generated column and without interleaving:
```sql
   CREATE TABLE ChildTable (
               Parent_ID STRING(36) NOT NULL,
               User_ID STRING(36) NOT NULL,
               Level STRING(36) NOT NULL,
               Original_Child_ID STRING(36) NOT NULL,
               New_Child_ID STRING(36),
               Primary_Child_ID STRING(36) AS (IF(New_Child_ID IS NOT NULL, New_Child_ID, Original_Child_ID)) STORED,
               Last_Update_Time TIMESTAMP NOT NULL OPTIONS (
                  allow_commit_timestamp = true
                  )
   ) PRIMARY KEY(Parent_ID, User_ID, Level);
```

Child table with generated column and with interleaving:
```sql
   CREATE TABLE ChildTableInterleaved(
             Parent_ID STRING(36) NOT NULL,
             User_ID STRING(36) NOT NULL,
             Level STRING(36) NOT NULL,
             Original_Child_ID STRING(36) NOT NULL,
             New_Child_ID STRING(36),
             Primary_Child_ID STRING(36) AS (IF(New_Child_ID IS NOT NULL, New_Child_ID, Original_Child_ID)) STORED,
             Last_Update_Time TIMESTAMP NOT NULL OPTIONS (
                allow_commit_timestamp = true
                )
   ) PRIMARY KEY(Parent_ID, User_ID, Level),
   INTERLEAVE IN PARENT ParentTable ON DELETE NO ACTION;
```

2. Insert a sample record in each table using these DML stmts:

```sql
   INSERT INTO ParentTable (Parent_ID, User_ID, Last_Update_Time)
   VALUES('parent_id', 'user_id', '2021-08-20T00:00:00Z');
   
   INSERT INTO ChildTable (Parent_ID, User_ID, Level, Last_Update_Time, Original_Child_ID)
   VALUES ('parent_id', 'user_id', 'Level', PENDING_COMMIT_TIMESTAMP(), 'original_child_id');
   
   INSERT INTO ChildTableInterleaved(Parent_ID, User_ID, Level, Last_Update_Time, Original_Child_ID)
   VALUES ('parent_id', 'user_id', 'Level', PENDING_COMMIT_TIMESTAMP(), 'original_child_id');
```

3. Navigate to the go and python stubs to execute them sequentially. Instructions to run them are enclosed within the 
   individual directories.
   
