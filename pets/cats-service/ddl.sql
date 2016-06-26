create database cats_db;

use cats_db;

create table cats ( 
  cat_id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  cat_name VARCHAR(32) NOT NULL,
  cat_age INT(10) UNSIGNED NOT NULL,
  cat_type VARCHAR(32) NOT NULL, 
  PRIMARY KEY (cat_id),
  UNIQUE INDEX cat_name (cat_name)
);
