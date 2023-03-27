
# Database Queries

* Steps
  * Create the database (Note: I am using postgresql so the commands might differ)
  
    ```sql
    CREATE DATABASE todo_app;
    ```

    * Select the database

    ```sql
    // \c <database_name>

    // \c todo_app

    /* 
    You will get a message like "You are now connected to database "todo_app" as user "postgres".`
    */

  * Create the tables

    ```sql
    // This is for the groups table, updated_at has to be linked with a trigger
    CREATE table groups(
        id SERIAL PRIMARY KEY, 
        title TEXT NOT NULL, 
        description TEXT, 
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW(),
        is_public BOOLEAN DEFAULT true,
        is_active BOOLEAN DEFAULT true
    );

    // This is for the items table, updated_at has to be linked with a trigger

    CREATE table items(
    id SERIAL PRIMARY KEY,
        content TEXT NOT NULL,
        is_active BOOLEAN DEFAULT true,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW(),
        remind_at TIMESTAMP NULL
    );


    // This is for the GroupedItems table

        CREATE table grouped_items(
            id SERIAL PRIMARY KEY,
            group_id INTEGER REFERENCES groups(id),
            item_id INTEGER REFERENCES items(id),
        created_at TIMESTAMP DEFAULT NOW(),
            is_active BOOLEAN DEFAULT true
        );

    // This is for the Contents table

        CREATE TABLE contents(
            id SERIAL PRIMARY KEY,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
        )

    // This is for the ItemContents table

        CREATE TABLE item_contents(
            id SERIAL PRIMARY KEY,
            item_id INTEGER REFERENCES items(id),
            content_id INTEGER REFERENCES contents(id),
            created_at TIMESTAMP DEFAULT NOW(),
            is_active BOOLEAN DEFAULT true
        )
    ```

  * Add the triggers

      ```sql
      // Trigger and Function to update updated at for groups table

      // create function
      CREATE OR REPLACE FUNCTION update_groups_updated_at()
      RETURNS TRIGGER AS $$
      BEGIN
          NEW.updated_at = NOW();
          RETURN NEW;
      END;
      $$ LANGUAGE plpgsql;

      // create trigger

      CREATE TRIGGER update_groups_trigger 
      BEFORE UPDATE ON groups
      FOR EACH ROW
      EXECUTE FUNCTION update_groups_updated_at();


      // Trigger and Function to update updated at for items table

      // create function
      CREATE OR REPLACE FUNCTION update_items_updated_at()
      RETURNS TRIGGER AS $$
      BEGIN
          NEW.updated_at = NOW();
          RETURN NEW;
      END;
      $$ LANGUAGE plpgsql;

      // create trigger

      CREATE TRIGGER update_items_trigger 
      BEFORE UPDATE ON items
      FOR EACH ROW
      EXECUTE FUNCTION update_items_updated_at();

      // Trigger and Function to update updated at for contents table

      // create function
      
      CREATE OR REPLACE FUNCTION update_contents_updated_at()
      RETURNS TRIGGER AS $$
      BEGIN
          NEW.updated_at = NOW();
          RETURN NEW;
      END;
      $$ LANGUAGE plpgsql;

      // create trigger

      CREATE TRIGGER update_contents_trigger 
      BEFORE UPDATE ON contents
      FOR EACH ROW
      EXECUTE FUNCTION update_contents_updated_at();```
