
docker run --name some-postgres -v `pwd`/datadir:/var/lib/postgresql/data -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres


CREATE ROLE "birdman" WITH PASSWORD 'mycape'  LOGIN INHERIT  VALID UNTIL 'infinity';


CREATE DATABASE birddb
  WITH ENCODING='UTF8'
       OWNER=birdman
       CONNECTION LIMIT=-1;  

CREATE TABLE public.birds
(
  bird_id serial NOT NULL,
  bird_name character varying(50) NOT NULL,
  bird_age integer NOT NULL,
  bird_type character varying(50) NOT NULL,
  CONSTRAINT birds_pkey PRIMARY KEY (bird_id),
  CONSTRAINT birds_bird_name_key UNIQUE (bird_name)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.birds
  OWNER TO birdman;
        
INSERT INTO birds( bird_name, bird_age, bird_type)
    VALUES ('first', 2, 'BLACKBIRD'); 
INSERT INTO birds( bird_name, bird_age, bird_type)
    VALUES ('second', 5, 'NUTHATCHESOWL'); 
INSERT INTO birds( bird_name, bird_age, bird_type)
    VALUES ('third', 2, 'DOVESDUCK');   