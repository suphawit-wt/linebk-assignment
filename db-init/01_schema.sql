-- -------------------------------------------------------------
-- Database: assignment
-- -------------------------------------------------------------


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

SET GLOBAL event_scheduler = ON;

CREATE TABLE `account_balances` (
  `account_id` varchar(50) NOT NULL,
  `user_id` varchar(50) DEFAULT NULL,
  `amount` decimal(15,2) DEFAULT NULL,
  `dummy_col_4` varchar(255) DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`account_id`),
  KEY `idx_account_balances_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `account_details` (
  `account_id` varchar(50) NOT NULL,
  `user_id` varchar(50) DEFAULT NULL,
  `color` varchar(10) DEFAULT NULL,
  `is_main_account` tinyint(1) DEFAULT NULL,
  `progress` int DEFAULT NULL,
  `dummy_col_5` varchar(255) DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`account_id`),
  KEY `idx_account_details_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `account_flags` (
  `flag_id` int NOT NULL AUTO_INCREMENT,
  `account_id` varchar(50) NOT NULL,
  `user_id` varchar(50) NOT NULL,
  `flag_type` varchar(50) NOT NULL,
  `flag_value` varchar(30) NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`flag_id`),
  KEY `idx_account_flags_account_id` (`account_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6000001 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `accounts` (
  `account_id` varchar(50) NOT NULL,
  `user_id` varchar(50) DEFAULT NULL,
  `type` varchar(50) DEFAULT NULL,
  `currency` varchar(10) DEFAULT NULL,
  `account_number` varchar(20) DEFAULT NULL,
  `issuer` varchar(100) DEFAULT NULL,
  `dummy_col_3` varchar(255) DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`account_id`),
  KEY `idx_accounts_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `banners` (
  `banner_id` varchar(50) NOT NULL,
  `user_id` varchar(50) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `description` text,
  `image` varchar(255) DEFAULT NULL,
  `dummy_col_11` varchar(255) DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`banner_id`),
  KEY `idx_banners_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `debit_card_design` (
  `card_id` varchar(50) NOT NULL,
  `user_id` varchar(50) DEFAULT NULL,
  `color` varchar(10) DEFAULT NULL,
  `border_color` varchar(10) DEFAULT NULL,
  `dummy_col_9` varchar(255) DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`card_id`),
  KEY `idx_card_design_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `debit_card_details` (
  `card_id` varchar(50) NOT NULL,
  `user_id` varchar(50) DEFAULT NULL,
  `issuer` varchar(100) DEFAULT NULL,
  `number` varchar(25) DEFAULT NULL,
  `dummy_col_10` varchar(255) DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`card_id`),
  KEY `idx_card_details_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `debit_card_status` (
  `card_id` varchar(50) NOT NULL,
  `user_id` varchar(50) DEFAULT NULL,
  `status` varchar(20) DEFAULT NULL,
  `dummy_col_8` varchar(255) DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`card_id`),
  KEY `idx_card_status_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `debit_cards` (
  `card_id` varchar(50) NOT NULL,
  `user_id` varchar(50) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `dummy_col_7` varchar(255) DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`card_id`),
  KEY `idx_card_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `transactions` (
  `transaction_id` varchar(50) NOT NULL,
  `user_id` varchar(50) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `image` varchar(255) DEFAULT NULL,
  `isBank` tinyint(1) DEFAULT NULL,
  `dummy_col_6` varchar(255) DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`transaction_id`, `created_at`),
  KEY `idx_transactions_user_id_created_at` (`user_id`, `created_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
PARTITION BY RANGE COLUMNS(`created_at`) (
  PARTITION p_2026_01 VALUES LESS THAN ('2026-02-01'),
  PARTITION p_2026_02 VALUES LESS THAN ('2026-03-01'),
  PARTITION p_2026_03 VALUES LESS THAN ('2026-04-01'),
  PARTITION p_2026_04 VALUES LESS THAN ('2026-05-01'),
  PARTITION p_2026_05 VALUES LESS THAN ('2026-06-01'),
  PARTITION p_2026_06 VALUES LESS THAN ('2026-07-01'),
  PARTITION p_2026_07 VALUES LESS THAN ('2026-08-01'),
  PARTITION p_2026_08 VALUES LESS THAN ('2026-09-01'),
  PARTITION p_2026_09 VALUES LESS THAN ('2026-10-01'),
  PARTITION p_2026_10 VALUES LESS THAN ('2026-11-01'),
  PARTITION p_2026_11 VALUES LESS THAN ('2026-12-01'),
  PARTITION p_2026_12 VALUES LESS THAN ('2027-01-01'),
  PARTITION p_future   VALUES LESS THAN MAXVALUE
);

CREATE TABLE `user_greetings` (
  `user_id` varchar(50) NOT NULL,
  `greeting` text,
  `dummy_col_2` varchar(255) DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `users` (
  `user_id` varchar(50) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `dummy_col_1` varchar(255) DEFAULT NULL,
  `password_hash` varchar(255) DEFAULT NULL,
  `passcode_hash` varchar(255) DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DELIMITER //

DROP PROCEDURE IF EXISTS usp_transactions_manage_partitions //

CREATE PROCEDURE usp_transactions_manage_partitions()
BEGIN
    DECLARE i INT DEFAULT 1;
    DECLARE target_date DATE;
    DECLARE partition_name VARCHAR(32);
    DECLARE next_month_date DATE;
    DECLARE sql_query TEXT;
    
    WHILE i <= 3 DO
        SET target_date = DATE_ADD(DATE_FORMAT(NOW(), '%Y-%m-01'), INTERVAL i MONTH);
        SET partition_name = CONCAT('p_', DATE_FORMAT(target_date, '%Y_%m'));
        SET next_month_date = DATE_ADD(target_date, INTERVAL 1 MONTH);
        
        IF NOT EXISTS (
            SELECT 1 
            FROM INFORMATION_SCHEMA.PARTITIONS 
            WHERE TABLE_SCHEMA = DATABASE() 
              AND TABLE_NAME = 'transactions' 
              AND PARTITION_NAME = partition_name
        ) THEN
            SET sql_query = CONCAT(
                'ALTER TABLE transactions REORGANIZE PARTITION p_future INTO (',
                'PARTITION ', partition_name, ' VALUES LESS THAN (''', DATE_FORMAT(next_month_date, '%Y-%m-%d'), '''), ',
                'PARTITION p_future VALUES LESS THAN MAXVALUE)'
            );
            
            SET @stmt = sql_query;
            PREPARE stmt FROM @stmt;
            EXECUTE stmt;
            DEALLOCATE PREPARE stmt;
        END IF;
        
        SET i = i + 1;
    END WHILE;
END //

DELIMITER ;


DROP EVENT IF EXISTS evt_monthly_transactions_partition_maintenance;

CREATE EVENT evt_monthly_transactions_partition_maintenance
ON SCHEDULE EVERY 1 MONTH
STARTS '2026-11-01 00:00:00'
ON COMPLETION PRESERVE
ENABLE
DO
  CALL usp_transactions_manage_partitions();

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;