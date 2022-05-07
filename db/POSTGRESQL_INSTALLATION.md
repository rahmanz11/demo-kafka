Steps
=====
1. Login as root

2. Install dependencies

apt-get update -y
apt-get install libncurses5-dev libncursesw5-dev
apt-get install -y libreadline-dev
apt install build-essential
apt install linux-headers-$(uname -r)
apt-get install zlib1g-dev
apt install zlib1g

3. Add postgres system user

adduser postgres
Password: psql

4. Create following directories to install postgres binary

mkdir -p /p01/pgsql/14.2/dbs/data
mkdir -p /p01/pgsql/diag/trace

5. Give permission to postgres system user

chown -R postgres:postgres /p01/*

6. Login to postgres user using following command:

su - postgres

7. Download postgres binary: wget https://ftp.postgresql.org/pub/source/v14.2/postgresql-14.2.tar.gz

8. Open bash_profile file 

vi ~/.bash_profile

9. Set following environment variables in bash_profile

export LD_LIBRARY_PATH=/p01/pgsql/12.4/lib
export PATH=/p01/pgsql/14.2/bin:$PATH
export MANPATH=/p01/pgsql/14.2/share/man:$MANPATH
export PGDATA=/p01/pgsql/14.2/dbs/data
export PGLOGS=/p01/pgsql/diag/trace
export PGHOME=/p01/pgsql/14.2/bin
export PGPORT=5432

10. Then initialize the environment variables

source ~/.bash_profile

11. Extract postgres binary

tar -vxzf postgresql-14.2.tar.gz

12. Enter into the extracted directory 

cd postgresql-14.2/

13. Start building postgres from source using following commands step by step

./configure --prefix=/p01/pgsql/14.2

make -j 8

make install

14. Initialize database database 

initdb -D /p01/pgsql/14.2/dbs/data/

In the end, you'll see a message with "Success" following the database start command. Keep that command copied somewhereelse

15. Browse to the directory to check configuration files

cd /p01/pgsql/14.2/dbs/data/

16. Open postgresql configuration file

vi postgresql.conf

17. Modify few properties as follows: 
- Set max_connections property to 1000
- Uncomment authentication_timeout property
- Uncomment password_encryption property and set it to md5
- Set shared_buffers property to 4096MB
- Uncomment huge_pages property
- Uncomment work_mem property and set it to 128MB
- Uncomment max_files_per_process and set it to 65565
- Uncomment logging_collector and set it to on
- Uncomment log_directory and set it to /p01/pgsql/diag/trace
- Uncomment log_rotation_age, log_rotation_size, log_min_messages, log_min_error_statement and log_checkpoints properties

18. Start database with the command copied in step 14

pg_ctl -D /p01/pgsql/14.2/dbs/data/ -l logfile start

You'll see a message ended with "server started"

19. Enter following command to open postgresql console

psql

20. Execute following sql to the console

ALTER USER postgres PASSWORD 'pgROOTpasswd';

21. Then execute following command to exit from console

exit

22. Now stop database using following command

pg_ctl -D /p01/pgsql/14.2/dbs/data stop

You'll see a message ended with "server stopped"

23. Go to the postgresql configuration directory

cd /p01/pgsql/14.2/dbs/data/

24. Open pg_hba.conf file

vi pg_hba.conf

25. In the below there is a table with following headers:

--TYPE  DATABASE        USER            ADDRESS                 METHOD

Go there and set all "trust" to "md5" under METHOD column

26. Now start database with following command

pg_ctl -D /p01/pgsql/14.2/dbs/data start

27. To check whether database is running use following command:

pg_ctl -D /p01/pgsql/14.2/dbs/data status