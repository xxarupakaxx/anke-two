# anke-two
[![codecov](https://codecov.io/gh/xxarupakaxx/anke-two/branch/main/graph/badge.svg?token=1IMEND5RZR)](https://codecov.io/gh/xxarupakaxx/anke-two)


## TODO
usecaseからinterfacesに伝えるLoggerが必要

## DB

```MariaDB [anke-two]> show tables;
+--------------------+
| Tables_in_anke-two |
+--------------------+
| administrators     |
| options            |
| question           |
| question_types     |
| questionnaires     |
| res_shared_tos     |
| respondents        |
| response           |
| scale_labels       |
| targets            |
| validations        |
+--------------------+
11 rows in set (0.000 sec)

MariaDB [anke-two]> SHOW columns from administrators;
+------------------+-------------+------+-----+---------+-------+
| Field            | Type        | Null | Key | Default | Extra |
+------------------+-------------+------+-----+---------+-------+
| questionnaire_id | int(11)     | NO   | PRI | NULL    |       |
| user_traqid      | varchar(32) | NO   | PRI | NULL    |       |
+------------------+-------------+------+-----+---------+-------+
2 rows in set (0.001 sec)

MariaDB [anke-two]> SHOW columns from options;                       
+-------------+---------+------+-----+---------+----------------+
| Field       | Type    | Null | Key | Default | Extra          |
+-------------+---------+------+-----+---------+----------------+
| id          | int(11) | NO   | PRI | NULL    | auto_increment |
| question_id | int(11) | NO   | MUL | NULL    |                |
| option_num  | int(11) | NO   |     | NULL    |                |
| body        | text    | YES  |     | NULL    |                |
+-------------+---------+------+-----+---------+----------------+
4 rows in set (0.000 sec)

MariaDB [anke-two]> SHOW columns from question;                             
+------------------+------------+------+-----+---------------------+---------------
-+
| Field            | Type       | Null | Key | Default             | Extra          |
+------------------+------------+------+-----+---------------------+---------------
-+
| id               | int(11)    | NO   | PRI | NULL                | auto_increment |
| questionnaire_id | int(11)    | NO   | MUL | NULL                |                |
| page_num         | int(11)    | NO   |     | NULL                |                |
| question_num     | int(11)    | NO   |     | NULL                |                |
| type             | int(11)    | NO   |     | NULL                |                |
| body             | text       | YES  |     | NULL                |                |
| is_required      | tinyint(4) | NO   |     | 0                   |                |
| deleted_at       | datetime   | YES  |     | NULL                |                |
| created_at       | datetime   | NO   |     | current_timestamp() |                |
+------------------+------------+------+-----+---------------------+---------------
-+
9 rows in set (0.000 sec)

MariaDB [anke-two]> SHOW columns from question_types ;
+--------+-------------+------+-----+---------+----------------+
| Field  | Type        | Null | Key | Default | Extra          |
+--------+-------------+------+-----+---------+----------------+
| id     | int(11)     | NO   | PRI | NULL    | auto_increment |
| name   | varchar(30) | NO   |     | NULL    |
| active | tinyint(1)  | NO   |     | 1       |                |
+--------+-------------+------+-----+---------+----------------+
3 rows in set (0.001 sec)

MariaDB [anke-two]> SHOW columns from questionnaires; 
+----------------+-------------+------+-----+----------------------+---------------
-+
| Field          | Type        | Null | Key | Default              | Extra          |
+----------------+-------------+------+-----+----------------------+---------------
-+
| id             | int(11)     | NO   | PRI | NULL                 | auto_increment |
| title          | char(50)    | NO   |     | NULL                 |                |
| description    | text        | NO   |     | NULL                 |                |
| res_time_limit | datetime    | YES  |     | NULL                 |                |
| deleted_at     | datetime    | YES  |     | NULL                 |                |
| res_shared_to  | int(11)     | NO   |     | 0                    |                |
| created_at     | datetime(3) | NO   |     | current_timestamp(3) |                |
| modified_at    | datetime(3) | NO   |     | current_timestamp(3) |                |
+----------------+-------------+------+-----+----------------------+---------------
-+
8 rows in set (0.001 sec)

MariaDB [anke-two]> SHOW columns from  res_shared_tos ;              
+--------+-------------+------+-----+---------+----------------+
| Field  | Type        | Null | Key | Default | Extra          |
+--------+-------------+------+-----+---------+----------------+
| id     | int(11)     | NO   | PRI | NULL    | auto_increment |
| name   | varchar(30) | NO   |     | NULL    |                |
| active | tinyint(1)  | NO   |     | 1       |                |
+--------+-------------+------+-----+---------+----------------+
3 rows in set (0.000 sec)

MariaDB [anke-two]> SHOW columns from   respondents;                
+------------------+-------------+------+-----+---------------------+--------------
--+
| Field            | Type        | Null | Key | Default             | Extra          |
+------------------+-------------+------+-----+---------------------+--------------
--+
| response_id      | int(11)     | NO   | PRI | NULL                | auto_increment |
| questionnaire_id | int(11)     | NO   | MUL | NULL                |                |
| user_traqid      | varchar(32) | YES  |     | NULL                |                |
| updated_at       | datetime    | NO   |     | current_timestamp() |                |
| submitted_at     | datetime    | YES  |     | NULL                |                |
| deleted_at       | datetime    | YES  |     | NULL                |                |
+------------------+-------------+------+-----+---------------------+--------------
--+
6 rows in set (0.000 sec)

MariaDB [anke-two]> SHOW columns from  response;                       
+-------------+----------+------+-----+---------------------+-------+
| Field       | Type     | Null | Key | Default             | Extra |
+-------------+----------+------+-----+---------------------+-------+
| response_id | int(11)  | NO   | PRI | NULL                |       |
| question_id | int(11)  | NO   | PRI | NULL                |       |
| body        | text     | YES  |     | NULL                |       |
| updated_at  | datetime | NO   |     | current_timestamp() |       |
| deleted_at  | datetime | YES  |     | NULL                |       |
+-------------+----------+------+-----+---------------------+-------+
5 rows in set (0.001 sec)

MariaDB [anke-two]> SHOW columns from  targets;                            
+------------------+-------------+------+-----+---------+----------------+
| Field            | Type        | Null | Key | Default | Extra          |
+------------------+-------------+------+-----+---------+----------------+
| questionnaire_id | int(11)     | NO   | PRI | NULL    | auto_increment |
| user_traqid      | varchar(32) | NO   |     | NULL    |                |
+------------------+-------------+------+-----+---------+----------------+
2 rows in set (0.000 sec)

MariaDB [anke-two]> SHOW columns from  scale_labels ;                       
+-------------------+-------------+------+-----+---------+----------------+
| Field             | Type        | Null | Key | Default | Extra          |
+-------------------+-------------+------+-----+---------+----------------+
| question_id       | int(11)     | NO   | PRI | NULL    | auto_increment |
| scale_label_right | varchar(50) | YES  |     | NULL    |                |
| scale_label_left  | varchar(50) | YES  |     | NULL    |                |
| scale_min         | int(11)     | YES  |     | NULL    |                |
| scale_max         | int(11)     | YES  |     | NULL    |                |
+-------------------+-------------+------+-----+---------+----------------+
5 rows in set (0.000 sec)

MariaDB [anke-two]> SHOW columns from   validations ;                
+---------------+---------+------+-----+---------+-------+
| Field         | Type    | Null | Key | Default | Extra |
+---------------+---------+------+-----+---------+-------+
| question_id   | int(11) | NO   | PRI | NULL    |       |
| regex_pattern | text    | YES  |     | NULL    |       |
| min_bound     | text    | YES  |     | NULL    |       |
| max_bound     | text    | YES  |     | NULL    |       |
+---------------+---------+------+-----+---------+-------+
4 rows in set (0.001 sec)
```