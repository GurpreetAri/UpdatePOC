### Set up the project.

Make sure you have python 3 installed.

1. Setup environment
    ```bash
    python3 -m venv venv
    source venv/bin/activate
    pip install -r requirements.txt
    ```
2. Create a test instance and create tables in it using DDL from README.md
3. Replace project_id, instance-id and database_id values in the .py clients to point to your config.
3. Run the individual clients for testing against table with generated column and with or without interleaving.
   
   ```bash
   python update_gen_col_interleaving.py
   ```

 ```bash
   python update_gen_col_no_interleaving.py
   ```