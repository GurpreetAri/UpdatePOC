### Setup the project

Make sure you have go installed. (Tested using `go1.16.5` and `cloud.google.com/go/spanner 1.24.1`)

1. Create a test instance and create tables in it using DDL from README.md
2. Replace `dbConnString` in the .go clients to point to your config.
3. Run the individual clients for testing against table with/without generated column.

    ```bash
    go run update_with_gen_column.go
    ```
    
    ```bash
    go run update_without_gen_column.go
    ```