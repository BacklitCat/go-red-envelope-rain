CREATE TABLE `rain` (
                        `user_account` varchar(255) NOT NULL DEFAULT '' COMMENT '账号',
                        `status` bool NOT NULL DEFAULT '1' COMMENT '是否可以抽取（黑名单）',
                        `remaining` int NOT NULL DEFAULT '0' COMMENT '剩余次数',
                        `balance` int NOT NULL DEFAULT '0' COMMENT '余额',
                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                        `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        PRIMARY KEY (`user_account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

