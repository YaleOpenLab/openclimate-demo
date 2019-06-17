# Openclimate

The main repo for the openclimate idea incubated at the YaleOpenLab.

### Prerequisites:
1. Node
2. Postgresql

Getting postgresql setup:
```
brew install postgresql
brew services start postgresql
```

```
psql postgres
postgres-# \conninfo
You are connected to database "postgres" as user "" via socket in "/tmp" at port "5432".
postgres=# CREATE ROLE ghost WITH LOGIN PASSWORD 'password';
CREATE ROLE
postgres=# ALTER ROLE ghost CREATEDB;
ALTER ROLE
postgres=# \du
                                   List of roles
 Role name |                         Attributes                         | Member of
-----------+------------------------------------------------------------+-----------
 ghost     | Create DB                                                  | {}

postgres=# \q

psql -d postgres -u ghost
postgres=> CREATE DATABASE api;
CREATE DATABASE
postgres=> \list
                              List of databases
   Name    |   Owner   | Encoding | Collate | Ctype |    Access privileges
-----------+-----------+----------+---------+-------+-------------------------
 api       | ghost     | UTF8     | C       | C     |

postgres=> \c api
You are now connected to database "api" as user "ghost".
api=> CREATE TABLE users (
api(> ID SERIAL PRIMARY KEY,
api(>   name VARCHAR(30),
api(>   email VARCHAR(30)
api(> );
CREATE TABLE
api=> INSERT INTO users (name, email)
api->   VALUES ('Jerry', 'jerry@example.com'), ('George', 'george@example.com');
INSERT 0 2
api=> SELECT * FROM USERS;
 id |  name  |       email
----+--------+--------------------
  1 | Jerry  | jerry@example.com
  2 | George | george@example.com
(2 rows)
```

After getting the relevant dbs for postgres up, run `npm install` and `npm start` to get the backend running
```
