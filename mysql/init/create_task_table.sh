#!/bin/sh

PREFIX_MYSQL="mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}"

$PREFIX_MYSQL -e "CREATE TABLE task (
    id INT(10) AUTO_INCREMENT NOT NULL primary key,
    title VARCHAR(50) NOT NULL,
    userId INT(10)"
    );"

$PREFIX_MYSQL -e "INSERT INTO task VALUES ('task1', '14');"
