# Openclimate

The main repo for the openclimate idea incubated at the YaleOpenLab.

### Prerequisites:
1. Go
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
postgres=> CREATE DATABASE openclimate;
CREATE DATABASE
postgres=> \list
                              List of databases
   Name    |   Owner   | Encoding | Collate | Ctype |    Access privileges
-----------+-----------+----------+---------+-------+-------------------------
 openclimate       | ghost     | UTF8     | C       | C     |

postgres=> \c openclimate
You are now connected to database "openclimate" as user "ghost".
openclimate=> CREATE TABLE users (
openclimate(> ID SERIAL PRIMARY KEY,
openclimate(>   name VARCHAR(30),
openclimate(>   email VARCHAR(30)
openclimate(> );
CREATE TABLE
openclimate=> INSERT INTO users (name, email)
openclimate->   VALUES ('Jerry', 'jerry@example.com'), ('George', 'george@example.com');
INSERT 0 2
openclimate=> SELECT * FROM USERS;
 id |  name  |       email
----+--------+--------------------
  1 | Jerry  | jerry@example.com
  2 | George | george@example.com
(2 rows)
```
