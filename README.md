# prototypes
1  git status
    2  sudo apt-get update
    3  sudo apt-get install mysql-server
    4  sudo service mysql start
    5  sudo mysql_secure_installation
    6  sudo mysql -u root -p
    7  mysql -u root
    8  sudo service mysql stop
    9  sudo service mysql start
   10  sudo mysql -u root
   11  sudo mysqladmin shutdown
   12  sudo service mysql start
   13  cd Connection_Pool/
   14  go connection.go 
   15  go run connection.go 
   16  sudo mysql -u root -p
   17  sudo netstat -tuln | grep 3306
   18  go run connection.go 
   19  sudo mysql -u root -p
   20  sudo service mysql restart
   21  sudo mysql -u root -p
   22  go run connection.go 
   23  history

   sudo mysql -u root -p
Enter password: 
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 14
Server version: 8.0.41-0ubuntu0.20.04.1 (Ubuntu)

Copyright (c) 2000, 2025, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> USE mysql;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY 'new_password';
Query OK, 0 rows affected (0.00 sec)

mysql> ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'Mysql@1234';
Query OK, 0 rows affected (0.00 sec)

mysql> FLUSH PRIVILEGES;
Query OK, 0 rows affected (0.00 sec)

mysql> exit
Bye
@PragnyaAitha ➜ /workspaces/prototypes/Connection_Pool (main) $ sudo service mysql restart
 * Stopping MySQL database server mysqld                                                                                [ OK ] 
 * Starting MySQL database server mysqld                                                                                       su: warning: cannot change directory to /nonexistent: No such file or directory
/opt/conda/bin/xz