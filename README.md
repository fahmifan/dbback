# DB BACKUPPER

A CLI to backups DB currently support MySQL & PostgreSQL. It uses the `mysqldump` & `pg_dump`. So, you need to have them installed in the machine.

## Prerequisite
- install `pg_dump (PostgreSQL) 12.9`
    - ubuntu:
    ```
    sudo apt install postgresql-client
    ```

- install `mysqldump Ver 10.19 Distrib 10.3.34-MariaDB`
    - https://mariadb.com/products/skysql/docs/connect/clients/mariadb-client/#debian-ubuntu

## Usage
Create a `config.json` based on `example.config.json` on the current working directory.

Options
```
backupper [options]

options:
    --help		show the help menu
    --dbaname	the database name
    --driver 	the database driver [mysql, postgres]
    --cron		run as cron job

```

### Example
- `dbback --driver=mysql --dbname=foo_bar`