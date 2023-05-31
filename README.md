#chainFinder/sniffer

Execution process: Enter the cmd file and enter chainFinder. Use './chainFinder' to execute the chainFinder. Then enter the sniffer file in the cmd file and use './sniffer' to execute the sniffer.

If you have executed tyche, the entire program has been executed successfully. If not, please download the tyche project from Git and follow the instructions provided in the tyche project to execute it.



---------------------
mysql:
USE chainfindata;
CREATE TABLE `event` (
   `id` int(11) NOT NULL AUTO_INCREMENT,
   `address` varchar(42) NOT NULL,
   `chainID` bigint(20) NOT NULL,
   `blockHash` varchar(66) NOT NULL,
   `blockNumber` varchar(255) DEFAULT NULL,
   `txHash` varchar(66) DEFAULT NULL,
   `txIndex` varchar(255) DEFAULT NULL,
   `gas` bigint(20) DEFAULT NULL,
   `gasPrice` bigint(20) DEFAULT NULL,
   `gasTipCap` bigint(20) DEFAULT NULL,
   `gasFeeCap` bigint(20) DEFAULT NULL,
   `value` varchar(255) DEFAULT NULL,
   `nonce` bigint(20) DEFAULT NULL,
   `toAddress` varchar(255) DEFAULT NULL,
   `status` tinyint(1) NOT NULL,
   `timestamp` bigint(20) NOT NULL,
   `newAddress` varchar(255) NOT NULL,  
   `newToAddress` varchar(255) NOT NULL,  
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
ALTER TABLE `event` ADD INDEX `idx_address` (`address`);
ALTER TABLE `event` ADD INDEX `idx_toAddress` (`toAddress`);
ALTER TABLE `event` ADD INDEX `idx_blockNumber` (`blockNumber`);

CREATE TABLE `block` (
   `id` int(11) NOT NULL AUTO_INCREMENT,
   `blockHash` varchar(66) NOT NULL,
   `blockNumber` varchar(255) DEFAULT NULL,   
   `blockReward` varchar(255) NOT NULL, 
   `minerAddress` varchar(42) DEFAULT NULL,
   `size` varchar(255) NOT NULL, 
   `timestamp` bigint(20) NOT NULL,
   `GasLimit` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
ALTER TABLE `block` ADD INDEX `idx_blockNumber` (`blockNumber`);


CREATE TABLE `ercevent` (
   `id` int(11) NOT NULL AUTO_INCREMENT,
   `contractName` varchar(255) NOT NULL,
   `eventName` varchar(255) NOT NULL,
   `data` json NOT NULL,   
   `name` varchar(255) NOT NULL,
   `txHash` varchar(66) DEFAULT NULL,
   `toAddress` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
ALTER TABLE `ercevent` ADD INDEX `idx_txHash` (`txHash`);

CREATE TABLE `ercTop` (
   `id` int(11) NOT NULL AUTO_INCREMENT,
   `contracaddress` varchar(42) NOT NULL,
   `name` varchar(255) NOT NULL,
   `value` varchar(255) DEFAULT NULL,
   `newContracaddress` varchar(255) NOT NULL,
   `contractTxCount` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `addressTop` (
   `id`         int(11) NOT NULL AUTO_INCREMENT,
   `address`    varchar(42)  NOT NULL,
   `Balance`    varchar(255) NOT NULL,
   `Count`      varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
