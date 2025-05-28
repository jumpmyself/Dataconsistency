CREATE TABLE IF NO EXISTS `info`
(
                        `id` int(11) NOT NULL,
                        `name` varchar(255) DEFAULT NULL,
                        `created_at` datetime DEFAULT NULL,
                        `updated_at` datetime DEFAULT NULL,
                        PRIMARY KEY (`id`)
) ENGINE=InnoDB
    DEFAULT CHARSET=utf8mb4;