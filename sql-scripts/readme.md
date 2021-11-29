Run script for mysql:

```bash
// from docker container shell
mysql -h hostname -u user database < path/to/test.sql

// from docker container if already running mysql command line (login: mysql -u user -p)
mysql> source \home\user\Desktop\test.sql;
```