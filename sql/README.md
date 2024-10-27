# SQL Notes


## Install SQLite

For Debian

```
sudo apt-get install sqlite3
```

For Windows

https://dev.to/dendihandian/installing-sqlite3-in-windows-44eb

Check if installed

```
sqlite3 --version
```


## Run a SQL file

```
sqlite3 ./data/sql.db < ./sql/init.sql
```

```
sqlite3 ./data/sql.db < ./sql/temp.sql
```

> If the file "sql.db" doesn't exist sqlite3 will create it

## Test SQL Files

