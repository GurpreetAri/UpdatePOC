CREATE TABLE Transaction (
     Account_ID STRING(36) NOT NULL,
     Transaction_ID STRING(36) NOT NULL,
     Last_Update_Time TIMESTAMP NOT NULL OPTIONS (
         allow_commit_timestamp = true
         ),
) PRIMARY KEY(Account_ID, Transaction_ID);

CREATE TABLE TransactionCategoryNew (
    Account_ID STRING(36) NOT NULL,
    Transaction_ID STRING(36) NOT NULL,
    Category_Type STRING(100) NOT NULL,
    Original_Category_ID STRING(36) NOT NULL,
    Recategorised_Category_ID STRING(36),
    Primary_Category_ID STRING(36) AS (IF(Recategorised_Category_ID IS NOT NULL, Recategorised_Category_ID, Original_Category_ID)) STORED,
    Last_Update_Time TIMESTAMP NOT NULL OPTIONS (
        allow_commit_timestamp = true
    )
) PRIMARY KEY(Account_ID, Transaction_ID, Category_Type);

INSERT INTO
    TRANSACTION (Account_ID,
                 Transaction_ID,
                 Last_Update_Time)
VALUES
    ('account_id1',
        'transaction_id1',
        '2021-08-20T00:00:00Z'
    );

INSERT INTO
    TransactionCategoryNew (Account_ID,
                            Transaction_ID,
                            Last_Update_Time,
                            Category_Type,
                            Original_Category_ID,
                            Recategorised_Category_ID)
VALUES
    ('account_id1',
     'transaction_id1',
    PENDING_COMMIT_TIMESTAMP(),
     'category_type1',
     'original_category1',
     'recat_category1'
    );
