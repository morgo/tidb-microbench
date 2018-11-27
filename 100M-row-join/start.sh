#!/bin/bash

mysql test -e "DROP TABLE IF EXISTS t1, t2;
CREATE TABLE t1 (id INT NOT NULL PRIMARY KEY auto_increment, b blob, c blob);
CREATE TABLE t2 (id INT NOT NULL PRIMARY KEY auto_increment, table1id INT NOT NULL, b blob, c blob);";

echo "starting threads";

JOBS=8

for i in `seq 1 8`; do
 echo "starting t1 generator"
 ./t1-generator.sh &
 echo "starting t2 generator"
 ./t2-generator.sh &
done;
