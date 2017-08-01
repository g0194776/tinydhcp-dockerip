/*
SQLyog Ultimate v11.24 (32 bit)
MySQL - 5.6.21 : Database - tiny_dhcp
*********************************************************************
*/


/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`tiny_dhcp` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_bin */;

USE `tiny_dhcp`;

/*Table structure for table `ip_pool` */

DROP TABLE IF EXISTS `ip_pool`;

CREATE TABLE `ip_pool` (
  `ip` char(18) COLLATE utf8_bin NOT NULL,
  `source_ip` char(15) COLLATE utf8_bin DEFAULT NULL,
  `env_id` smallint(6) NOT NULL,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `desc` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `owner` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`ip`,`env_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

/*Data for the table `ip_pool` */

/*Table structure for table `ref_env` */

DROP TABLE IF EXISTS `ref_env`;

CREATE TABLE `ref_env` (
  `env_id` smallint(6) NOT NULL,
  `env_name` varchar(128) COLLATE utf8_bin NOT NULL,
  `is_enabled` bit(1) NOT NULL DEFAULT b'1'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

/*Data for the table `ref_env` */

insert  into `ref_env`(`env_id`,`env_name`,`is_enabled`) values (1,'prod-bj-yz',''),(2,'test-sh-h5',''),(255,'un-assigned','');

/* Procedure structure for procedure `USP_CheckinIP` */

/*!50003 DROP PROCEDURE IF EXISTS  `USP_CheckinIP` */;

DELIMITER $$

/*!50003 CREATE DEFINER=`root`@`%` PROCEDURE `USP_CheckinIP`(
`node_ip` varchar(18),
`env_id` smallint,
`owner` varchar(32),
`desc` varchar(64)
)
BEGIN
	DECLARE orgIp char(18) DEFAULT NULL;
	SELECT ip INTO orgIp FROM tiny_dhcp.ip_pool t WHERE t.source_ip = `node_ip` and t.env_id = `env_id`;
	IF (orgIp IS NOT NULL) THEN
	   SELECT orgIp;
	ELSE
     SELECT ip INTO orgIp FROM tiny_dhcp.ip_pool t WHERE t.env_id = `env_id` AND (t.source_ip IS NULL or t.source_ip = '') LIMIT 1;
		 IF (orgIp IS NOT NULL) THEN
			  UPDATE tiny_dhcp.ip_pool t SET t.source_ip = `node_ip`, t.`desc` = `desc`, t.`owner` = `owner`, t.`update_time` = NOW() WHERE t.ip = orgIp;
		 END IF;
	   SELECT orgIp;
	END IF;
END */$$
DELIMITER ;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
