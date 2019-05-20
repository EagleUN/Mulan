CREATE TABLE `shares` (
  `uuid` VARCHAR(36) NOT NULL,
  `userId` VARCHAR(45) NULL,
  `postId` VARCHAR(45) NULL,
  `sharedAt` DATETIME NULL,
  PRIMARY KEY (`uuid`));