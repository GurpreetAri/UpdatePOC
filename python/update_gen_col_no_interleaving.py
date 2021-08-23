# Imports the Google Cloud Client Library.
from google.cloud import spanner

# Replace these with your own parameters
project_id="test-project"
instance_id = "test-instance"
database_id = "test-database"

# Instantiate a client.
spanner_client = spanner.Client(project=project_id)

# Get a Cloud Spanner instance by ID.
instance = spanner_client.instance(instance_id)

# Get a Cloud Spanner database by ID.
database = instance.database(database_id)

def printTransactions(database):
    select_query = "SELECT User_ID FROM ChildTableInterleaved"
    with database.snapshot() as snapshot:
        results = snapshot.execute_sql(select_query)
        for row in results:
            print(u"User_ID: {}".format(*row))


dml = '''UPDATE ChildTable
     SET New_Child_ID = @newChildId,
     Last_Update_Time = PENDING_COMMIT_TIMESTAMP()
     WHERE User_ID = @userId'''

params = {
    "userId": "user_id",
    "newChildId": "new_child_id",
}
param_types = {
    "userId": spanner.param_types.STRING,
    "newChildId": spanner.param_types.STRING,
}

row_ct = database.execute_partitioned_dml(
    dml=dml, params=params, param_types=param_types)

print("{} records updated.".format(row_ct))

printTransactions(database)
printTransactions(database)