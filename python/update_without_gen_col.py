# Imports the Google Cloud Client Library.
from google.cloud import spanner

# Replace these with your own parameters
project_id="anz-x-fabric-np-641432"
instance_id = "test-instance"
database_id = "example-db"

# Instantiate a client.
spanner_client = spanner.Client(project=project_id)

# Get a Cloud Spanner instance by ID.
instance = spanner_client.instance(instance_id)

# Get a Cloud Spanner database by ID.
database = instance.database(database_id)

def printTransactions(database):
    select_query = "SELECT Transaction_ID FROM TransactionCategory"
    with database.snapshot() as snapshot:
        results = snapshot.execute_sql(select_query)
        for row in results:
            print(u"Transaction_ID: {}".format(*row))


dml = '''UPDATE TransactionCategory
     SET Recategorised_Category_ID = @newCategoryId,
     Last_Update_Time = PENDING_COMMIT_TIMESTAMP()
     WHERE Transaction_ID = @transactionId'''

params = {
    "newCategoryId": "055a6b93-cf12-5637-9c65-d44f61615e97",
    "transactionId": "00000000-0000-0000-0000-000000000015",
}
param_types = {
    "newCategoryId": spanner.param_types.STRING,
    "transactionId": spanner.param_types.STRING,
}

row_ct = database.execute_partitioned_dml(
    dml=dml, params=params, param_types=param_types)

print("{} records updated.".format(row_ct))

printTransactions(database)
printTransactions(database)