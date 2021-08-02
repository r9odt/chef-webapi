DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) DEFAULT "",
  `password` varchar(255) DEFAULT "",
  `fullName` varchar(255) DEFAULT "",
  `avatar` text DEFAULT "",
  `isAdmin` boolean DEFAULT 0,
  `isBlocked` boolean DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE IF NOT EXISTS `sessions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uuid` varchar(255) DEFAULT "",
  `username` varchar(255) DEFAULT "",
  `expire` int(11) DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `tasks`;
CREATE TABLE IF NOT EXISTS `tasks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `resource` varchar(255) DEFAULT "",
  `onlyResource` boolean DEFAULT 0,
  `selectedResource` boolean DEFAULT 0,
  `resources` text DEFAULT "",
  `name` varchar(255) DEFAULT "",
  `status` varchar(255) DEFAULT "",
  `initiatorID` int(11) DEFAULT -1,
  `timestamp` int(11) DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
