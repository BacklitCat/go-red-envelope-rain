CREATE TABLE `user` (
                        `id` bigint NOT NULL AUTO_INCREMENT,
                        `user_account` varchar(255) NOT NULL DEFAULT '' COMMENT '账号',
                        `user_name` varchar(255)  NOT NULL DEFAULT '' COMMENT '用户名',
                        `user_password` varchar(255)  NOT NULL DEFAULT '' COMMENT '用户密码',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `user_account_unique` (`user_account`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 ;