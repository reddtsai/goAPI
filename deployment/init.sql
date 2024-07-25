CREATE TABLE `user` (
  `id` bigint NOT NULL,
  `account` varchar(45) NOT NULL,
  `secret` varchar(255) NOT NULL,
  `name` varchar(100) NOT NULL,
  `description` varchar(200) NOT NULL,
  `state` smallint NOT NULL,
  `created_at` bigint NOT NULL,
  `creator` bigint NOT NULL,
  `updated_at` bigint NOT NULL,
  `updater` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_merchant_user_account` (`account`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;