#chainFinder/sniffer

Execution process: Enter the cmd file and enter chainFinder. Use './chainFinder' to execute the chainFinder. Then enter the sniffer file in the cmd file and use './sniffer' to execute the sniffer.

If you have executed tyche, the entire program has been executed successfully. If not, please download the tyche project from Git and follow the instructions provided in the tyche project to execute it.



---------------------
mysql:
USE chainfindata;
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

-- --------------------------------------------------------

--
-- 琛ㄧ殑缁撴瀯 `addressTop`
--

CREATE TABLE `addressTop` (
  `id` bigint NOT NULL,
  `address` varchar(42) NOT NULL,
  `Balance` decimal(28,0) NOT NULL,
  `Count` int UNSIGNED DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- 琛ㄧ殑缁撴瀯 `block`
--

CREATE TABLE `block` (
  `id` bigint NOT NULL,
  `blockHash` varchar(66) NOT NULL,
  `blockNumber` bigint NOT NULL DEFAULT '0',
  `blockReward` bigint NOT NULL DEFAULT '0',
  `minerAddress` varchar(42) DEFAULT NULL,
  `size` int UNSIGNED NOT NULL DEFAULT '0',
  `timestamp` int UNSIGNED NOT NULL,
  `GasLimit` int UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- 琛ㄧ殑缁撴瀯 `ercevent`
--

CREATE TABLE `ercevent` (
  `id` bigint NOT NULL,
  `contractName` varchar(255) NOT NULL,
  `eventName` varchar(255) NOT NULL,
  `data` json NOT NULL,
  `name` varchar(255) NOT NULL,
  `txHash` varchar(66) DEFAULT NULL,
  `toAddress` varchar(44) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '0x0000000000000000000000000000000000000000'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- 琛ㄧ殑缁撴瀯 `ercTop`
--

CREATE TABLE `ercTop` (
  `id` bigint NOT NULL,
  `contracaddress` varchar(42) NOT NULL,
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `value` decimal(64,0) DEFAULT '0',
  `newContracaddress` varchar(44) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `contractTxCount` int UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- 琛ㄧ殑缁撴瀯 `event`
--

CREATE TABLE `event` (
  `id` bigint NOT NULL,
  `address` varchar(42) NOT NULL,
  `chainID` int UNSIGNED NOT NULL,
  `blockHash` varchar(66) NOT NULL,
  `blockNumber` bigint DEFAULT '0',
  `txHash` varchar(66) DEFAULT NULL,
  `txIndex` int DEFAULT '0',
  `gas` bigint DEFAULT '0',
  `gasPrice` bigint DEFAULT '0',
  `gasTipCap` bigint DEFAULT '0',
  `gasFeeCap` bigint DEFAULT '0',
  `value` decimal(27,0) NOT NULL DEFAULT '0',
  `nonce` bigint DEFAULT '0',
  `toAddress` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '0x0000000000000000000000000000000000000000',
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `timestamp` int UNSIGNED NOT NULL DEFAULT '0',
  `newAddress` varchar(44) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '0x0000000000000000000000000000000000000000',
  `newToAddress` varchar(44) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '0x0000000000000000000000000000000000000000'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- 杞偍琛ㄧ殑绱㈠紩
--

--
-- 琛ㄧ殑绱㈠紩 `addressTop`
--
ALTER TABLE `addressTop`
  ADD PRIMARY KEY (`id`),
  ADD KEY `address` (`address`);

--
-- 琛ㄧ殑绱㈠紩 `block`
--
ALTER TABLE `block`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `blockHash` (`blockHash`),
  ADD KEY `idx_blockNumber` (`blockNumber`);

--
-- 琛ㄧ殑绱㈠紩 `ercevent`
--
ALTER TABLE `ercevent`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_txHash` (`txHash`);

--
-- 琛ㄧ殑绱㈠紩 `ercTop`
--
ALTER TABLE `ercTop`
  ADD PRIMARY KEY (`id`),
  ADD KEY `contracaddress` (`contracaddress`);

--
-- 琛ㄧ殑绱㈠紩 `event`
--
ALTER TABLE `event`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_address` (`address`),
  ADD KEY `idx_toAddress` (`toAddress`),
  ADD KEY `idx_blockNumber` (`blockNumber`);

--
-- 鍦ㄥ鍑虹殑琛ㄤ娇鐢ˋUTO_INCREMENT
--

--
-- 浣跨敤琛ˋUTO_INCREMENT `addressTop`
--
ALTER TABLE `addressTop`
  MODIFY `id` bigint NOT NULL AUTO_INCREMENT;

--
-- 浣跨敤琛ˋUTO_INCREMENT `block`
--
ALTER TABLE `block`
  MODIFY `id` bigint NOT NULL AUTO_INCREMENT;

--
-- 浣跨敤琛ˋUTO_INCREMENT `ercevent`
--
ALTER TABLE `ercevent`
  MODIFY `id` bigint NOT NULL AUTO_INCREMENT;

--
-- 浣跨敤琛ˋUTO_INCREMENT `ercTop`
--
ALTER TABLE `ercTop`
  MODIFY `id` bigint NOT NULL AUTO_INCREMENT;

--
-- 浣跨敤琛ˋUTO_INCREMENT `event`
--
ALTER TABLE `event`
  MODIFY `id` bigint NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
