CREATE TABLE ParentTable (
         Parent_ID STRING(36) NOT NULL,
         User_ID STRING(36) NOT NULL,
         Last_Update_Time TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true))
PRIMARY KEY(Parent_ID, User_ID);

CREATE TABLE ChildTable (
        Parent_ID STRING(36) NOT NULL,
        User_ID STRING(36) NOT NULL,
        Original_Child_ID STRING(36) NOT NULL,
        New_Child_ID STRING(36),
        Primary_Child_ID STRING(36) AS (IF(New_Child_ID IS NOT NULL, New_Child_ID, Original_Child_ID)) STORED,
        Last_Update_Time TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true))
PRIMARY KEY(Parent_ID, User_ID);

CREATE TABLE ChildTableInterleaved(
        Parent_ID STRING(36) NOT NULL,
        User_ID STRING(36) NOT NULL,
        Original_Child_ID STRING(36) NOT NULL,
        New_Child_ID STRING(36),
        Primary_Child_ID STRING(36) AS (IF(New_Child_ID IS NOT NULL, New_Child_ID, Original_Child_ID)) STORED,
        Last_Update_Time TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true))
        PRIMARY KEY(Parent_ID, User_ID);
) PRIMARY KEY(Parent_ID, User_ID)
INTERLEAVE IN PARENT ParentTable ON DELETE NO ACTION;


INSERT INTO ParentTable (Parent_ID, User_ID, Last_Update_Time)
VALUES('parent_id', 'user_id', '2021-08-20T00:00:00Z');

INSERT INTO ChildTable (Parent_ID, User_ID, Last_Update_Time, Original_Child_ID)
VALUES ('parent_id', 'user_id', PENDING_COMMIT_TIMESTAMP(), 'original_child_id');

INSERT INTO ChildTableInterleaved(Parent_ID, User_ID, Last_Update_Time, Original_Child_ID)
VALUES ('parent_id', 'user_id', PENDING_COMMIT_TIMESTAMP(), 'original_child_id');