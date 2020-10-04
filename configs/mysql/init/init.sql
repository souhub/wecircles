DROP TABLE IF EXISTS users;
CREATE TABLE users
(
  `id`
(10) int NOT NULL AUTO_INCREMENT,
  `email` varchar
(64) NOT NULL DEFAULT '',
  `nickname` varchar
(32) NOT NULL DEFAULT '',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `modified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON
UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY
(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='users';