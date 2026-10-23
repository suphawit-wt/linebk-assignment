INSERT INTO `users` (`user_id`, `name`, `password_hash`, `passcode_hash`, `dummy_col_1`) VALUES
('000027b0e9a211ef95a30242ac180002', 'Clare', '$argon2id$v=19$m=65536,t=1,p=4$T5H3iErbpg65G7Bk43TQcA$csMp4p8LkYScFLqk7bShc+HeGvYtGECGV2YRmc+9cLs', '$argon2id$v=19$m=65536,t=1,p=4$/z8gVrh4ixJkFeUACubmug$E5S6SMW+zlOdFXMkeVPWtx94/sO8fdoJ83FOVvWBltw', 'dummy_value_1');

INSERT INTO `user_greetings` (`user_id`, `greeting`, `dummy_col_2`) VALUES
('000027b0e9a211ef95a30242ac180002', 'Have a nice day Clare', 'dummy_value_2');


INSERT INTO `accounts` (`account_id`, `user_id`, `type`, `currency`, `account_number`, `issuer`, `dummy_col_3`) VALUES
('00000acc01a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'saving-account', 'THB', '568-2-81740-9', 'TestLab', 'dummy_value_3'),
('00000acc02a211ef95a30242ac180002',   '000027b0e9a211ef95a30242ac180002', 'credit-loan',    'THB', '568-2-81740-9', 'TestLab', 'dummy_value_3'),
('00000acc03a211ef95a30242ac180002',   '000027b0e9a211ef95a30242ac180002', 'goal-saving-account', 'THB', '568-2-81740-9', 'TestLab', 'dummy_value_3'),
('00000acc04a211ef95a30242ac180002',   '000027b0e9a211ef95a30242ac180002', 'credit-loan',    'THB', '568-2-81740-9', 'TestLab', 'dummy_value_3');

INSERT INTO `account_balances` (`account_id`, `user_id`, `amount`, `dummy_col_4`) VALUES
('00000acc01a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 62000.00, 'dummy_value_4'),
('00000acc02a211ef95a30242ac180002',   '000027b0e9a211ef95a30242ac180002',   300.10, 'dummy_value_4'),
('00000acc03a211ef95a30242ac180002',   '000027b0e9a211ef95a30242ac180002', 30000.00, 'dummy_value_4'),
('00000acc04a211ef95a30242ac180002',   '000027b0e9a211ef95a30242ac180002', 30000.00, 'dummy_value_4');

INSERT INTO `account_details` (`account_id`, `user_id`, `color`, `is_main_account`, `progress`, `dummy_col_5`) VALUES
('00000acc01a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', '#24c875', 1, NULL, 'dummy_value_5'),
('00000acc02a211ef95a30242ac180002',   '000027b0e9a211ef95a30242ac180002', '#9366ed', 0, NULL, 'dummy_value_5'),
('00000acc03a211ef95a30242ac180002',   '000027b0e9a211ef95a30242ac180002', '#00a1e2', 0, 24,   'dummy_value_5'),
('00000acc04a211ef95a30242ac180002',   '000027b0e9a211ef95a30242ac180002', '#15bbc7', 0, NULL, 'dummy_value_5');

INSERT INTO `account_flags` (`account_id`, `user_id`, `flag_type`, `flag_value`) VALUES
('00000acc04a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'STATUS', 'Disbursement'),
('00000acc04a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'ALERT', 'Overdue');


INSERT INTO `transactions` (`transaction_id`, `user_id`, `name`, `image`, `isBank`, `dummy_col_6`, `created_at`) VALUES
('0000tx01e1a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'Emily', 'https://dummyimage.com/54x54/999/fff', 0, 'dummy_value_6', NOW()),
('0000tx02e1a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'Jone Kiersten', 'https://dummyimage.com/54x54/999/fff', 0, 'dummy_value_6', NOW());


INSERT INTO `debit_cards` (`card_id`, `user_id`, `name`, `dummy_col_7`) VALUES
('0000card01a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'My Salary', 'dummy_value_8'),
('0000card02a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'For My Dream', 'dummy_value_8'),
('0000card03a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'My Debit card', 'dummy_value_8');

INSERT INTO `debit_card_status` (`card_id`, `user_id`, `status`, `dummy_col_8`) VALUES
('0000card01a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'in-progress', 'dummy_value_8'),
('0000card02a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'in-progress', 'dummy_value_8'),
('0000card03a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'Active', 'dummy_value_8');

INSERT INTO `debit_card_design` (`card_id`, `user_id`, `color`, `border_color`, `dummy_col_9`) VALUES
('0000card01a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', '#00a1e2', NULL, 'dummy_value_9'),
('0000card02a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', '#ff8300', NULL, 'dummy_value_9'),
('0000card03a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', '#ffffff', '#f2f3f7', 'dummy_value_9');

INSERT INTO `debit_card_details` (`card_id`, `user_id`, `issuer`, `number`, `dummy_col_10`) VALUES
('0000card01a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'TestLab', NULL, 'dummy_value_10'),
('0000card02a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'TestLab', NULL, 'dummy_value_10'),
('0000card03a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'TestLab', '9440 7841 2222 3115', 'dummy_value_10');


INSERT INTO `banners` (`banner_id`, `user_id`, `title`, `description`, `image`, `dummy_col_11`) VALUES
('0000bn01e1a211ef95a30242ac180002', '000027b0e9a211ef95a30242ac180002', 'Want some money?', 'You can start apply ''Clare''', 'https://dummyimage.com/54x54/999/fff', 'dummy_value_11');