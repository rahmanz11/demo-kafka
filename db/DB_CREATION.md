DB Creation for microservice
=============================

Following steps will be done as "postgres" system user

1. Login to the postgres console using following command

psql -U postgres

2. Create microservice db user using following command

create user micros with encrypted password 'micro$';

3. Create OrderDB with following command

create database order_db;

4. Grant permission of order_db to 'micros' user

grant all privileges on database order_db to micros;

5. Create MatchOrderDB with following command

create database match_order_db;

6. Grant permission of match_order_db to 'micros' user

grant all privileges on database order_db to micros;

7. To show the list of database created execute the following command

\l

In the output there is a database list and you'll see the two above created microservice databases there