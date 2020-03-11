CREATE TABLE `rule` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `threshold` int(11) NOT NULL,
    `coefficient` decimal(3,2) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `threshold` (`threshold`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;