# DBZ 

## Overview
dbz is a cli program that is used to quickly and easily create databases from the command line. The default is to create the databases as docker containers. dbz is easily installed with a curl command. once installed you can use the command `dbz` with specified options. dbz can be used to create production and dev databases. 

## features
- install with curl or from source
    - should have a curl command users can run that will globally install dbz. 
    - can install from source as well (obviously)
- cli usage
    - should be used primarily with the `dbz` command and options
- create and delete databases
    - ability to specify version of database
    - can create from sql files
    - can create dev databases (testcontainers?)
- able to use config files or options to modify databases (privileges, bin log, etc. )
    - advanced user stuff
- ability to use seeders to seed test (faker) data into database
- ability to run migrations to modify database 

## Supported database
### OLTP 
- MySQL 
- PostgreSQL
- MariaDB
- SQLite
### OLAP
- DuckDB
- ClickHouse

## Considerations
- databases should be created as docker containers
- seeders and migrations should use standardized format, most likely sql files. if there is a good package that can be used for seeders / migrations that may be worth looking into.
    - seeders may benefit from using faker 
- look into packages that can be used that would simplify the implementation. For example testcontainers could be used for dev databases.

## Example usage 
- `dbz create postgres`
    - creates postgres database using latest image
- `dbz create mysql@8.4` 
    - creates mysql database using specified version 
- `dbz create postgres db_file.sql`
    - create postgres database and run sql file
- `dbz seed [database] [table]`
    - seed data 
- `dbz migrate [file]`
    - run migrations