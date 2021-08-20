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
4. Run the individual clients for testing against table with/without generated column.
   
   ```bash
   python update_with_gen_col.py
   ```

 ```bash
   python update_without_gen_col.py
   ```