CREATE DATABASE wecircles;
USE wecircles;

DROP TABLE IF EXISTS users;

CREATE TABLE users
(
  id           INT(10),
  name     VARCHAR(40)
);

INSERT users (id, name) VALUES (1, "Nagaoka");
INSERT users (id, name) VALUES (2, "Tanaka");
INSERT users (id, name) VALUES (3, "SOICHI");
INSERT users (id, name) VALUES (4, "DEMURA");
