
CREATE TABLE `shares` (
  `userId` VARCHAR(45) NOT NULL,
  `postId` VARCHAR(45) NOT NULL,
  `sharedAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`userId`, `postId`));
