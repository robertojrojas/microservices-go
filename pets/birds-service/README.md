

CREATE ROLE 'birdman' WITH PASSWORD 'mycape'  LOGIN INHERIT  VALID UNTIL 'infinity';


CREATE DATABASE birddb
  WITH ENCODING='UTF8'
       OWNER=birdman
       CONNECTION LIMIT=-1;  

CREATE TABLE birds(
 bird_id serial PRIMARY KEY,
 bird_name VARCHAR (50) UNIQUE NOT NULL,
 bird_age  integer NOT NULL,
 bird_type VARCHAR (50) NOT NULL 
);

ALTER TABLE birds
  OWNER TO birdman;        
        
INSERT INTO birds( bird_name, bird_age, bird_type)
    VALUES ('first', 2, 'BLACKBIRD'); 
INSERT INTO birds( bird_name, bird_age, bird_type)
    VALUES ('first', 5, 'NUTHATCHESOWL'); 
INSERT INTO birds( bird_name, bird_age, bird_type)
    VALUES ('first', 2, 'DOVESDUCK');   