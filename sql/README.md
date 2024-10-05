# SQL Notes

## Install SQLite

For Debian

```
sudo apt-get install sqlite3
```

For Windows

https://dev.to/dendihandian/installing-sqlite3-in-windows-44eb

Check Version

```
sqlite3 --version
```


## Run a SQL file

```
sqlite3 ./data/sql.db < ./sql/fill.sql
sqlite3 ./data/sql.db < ./sql/tables-drop.sql
```

> If sql.db doesn't exist sqlite will make it

## Test SQL Files

