-- MySQL dump 10.13  Distrib 8.0.27, for Linux (x86_64)
--
-- Host: localhost    Database: suit
-- ------------------------------------------------------
-- Server version	8.0.27

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `cart_items`
--

DROP TABLE IF EXISTS `cart_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cart_items` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `category_id` bigint DEFAULT NULL,
  `weight` bigint NOT NULL,
  `checkout_id` bigint DEFAULT NULL,
  `cart_user_id` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_cart_items_category` (`category_id`),
  KEY `fk_cart_items_checkout` (`checkout_id`),
  KEY `fk_cart_items_cart` (`cart_user_id`),
  CONSTRAINT `fk_cart_items_cart` FOREIGN KEY (`cart_user_id`) REFERENCES `carts` (`user_id`),
  CONSTRAINT `fk_cart_items_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`),
  CONSTRAINT `fk_cart_items_checkout` FOREIGN KEY (`checkout_id`) REFERENCES `checkouts` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cart_items`
--

LOCK TABLES `cart_items` WRITE;
/*!40000 ALTER TABLE `cart_items` DISABLE KEYS */;
INSERT INTO `cart_items` (`id`, `category_id`, `weight`, `checkout_id`, `cart_user_id`, `created_at`) VALUES (2,1,4,1,1,'2021-11-07 22:14:34.397'),(3,2,3,2,1,'2021-11-07 22:14:42.267');
/*!40000 ALTER TABLE `cart_items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `carts`
--

DROP TABLE IF EXISTS `carts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `carts` (
  `user_id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`user_id`),
  CONSTRAINT `fk_carts_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `carts`
--

LOCK TABLES `carts` WRITE;
/*!40000 ALTER TABLE `carts` DISABLE KEYS */;
INSERT INTO `carts` (`user_id`, `created_at`) VALUES (1,'2021-11-07 20:56:56.960');
/*!40000 ALTER TABLE `carts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `categories` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `point` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` (`id`, `name`, `point`, `created_at`, `updated_at`) VALUES (1,'kertas',2,'2021-11-07 22:04:04.443','2021-11-07 22:07:07.288'),(2,'plastik',4,'2021-11-07 22:04:24.311','2021-11-07 22:04:24.311'),(3,'kaca',4,'2021-11-07 22:04:46.982','2021-11-07 22:04:46.982'),(4,'kaleng',5,'2021-11-07 22:05:06.754','2021-11-07 22:05:06.754');
/*!40000 ALTER TABLE `categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `checkouts`
--

DROP TABLE IF EXISTS `checkouts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `checkouts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `checkouts`
--

LOCK TABLES `checkouts` WRITE;
/*!40000 ALTER TABLE `checkouts` DISABLE KEYS */;
INSERT INTO `checkouts` (`id`, `created_at`) VALUES (1,'2021-11-07 22:16:15.594'),(2,'2021-11-07 22:17:28.104');
/*!40000 ALTER TABLE `checkouts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `drop_points`
--

DROP TABLE IF EXISTS `drop_points`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `drop_points` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `address` longtext NOT NULL,
  `longitude` longtext,
  `latitude` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `drop_points`
--

LOCK TABLES `drop_points` WRITE;
/*!40000 ALTER TABLE `drop_points` DISABLE KEYS */;
INSERT INTO `drop_points` (`id`, `address`, `longitude`, `latitude`, `created_at`, `updated_at`) VALUES (1,'universitas indonesia','106.8270482','-6.3627638','2021-11-07 21:14:08.586','2021-11-07 21:26:35.540'),(2,'universitas padjadjaran','107.7746881','-6.926132099999999','2021-11-07 22:37:51.973','2021-11-07 22:37:51.973'),(3,'Jl TB Simatupang, Cilandak Timur','106.8171689','-6.2926515','2021-11-07 22:39:01.460','2021-11-07 22:39:01.460'),(4,'Jl. Matraman Raya No.46-48, RT.12/RW.2, Kb. Manggis, Kec. Matraman, Kota Jakarta Timur, Daerah Khusus Ibukota Jakarta 13150','106.8561887','-6.202977799999999','2021-11-07 22:40:18.969','2021-11-07 22:40:18.969'),(5,'Jl. RE Martadinata No.5, Cipayung, Kec. Ciputat, Kota Tangerang Selatan, Banten 15411','106.748685','-6.3221015','2021-11-07 22:41:37.131','2021-11-07 22:41:37.131');
/*!40000 ALTER TABLE `drop_points` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `logins`
--

DROP TABLE IF EXISTS `logins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `logins` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `email` varchar(55) DEFAULT NULL,
  `username` varchar(55) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `role` enum('staff','user') DEFAULT NULL,
  `user_id` bigint DEFAULT NULL,
  `staff_id` bigint DEFAULT NULL,
  `token` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `username` (`username`),
  KEY `fk_logins_staff` (`staff_id`),
  KEY `fk_logins_user` (`user_id`),
  CONSTRAINT `fk_logins_staff` FOREIGN KEY (`staff_id`) REFERENCES `staffs` (`id`),
  CONSTRAINT `fk_logins_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `logins`
--

LOCK TABLES `logins` WRITE;
/*!40000 ALTER TABLE `logins` DISABLE KEYS */;
INSERT INTO `logins` (`id`, `email`, `username`, `password`, `role`, `user_id`, `staff_id`, `token`, `created_at`, `updated_at`) VALUES (1,'darienkentanu@gmail.com','darienkentanu','$2a$14$gmxstHyq6u8ashaTGvuYOOZNAwVmic5MtP2JH998nrmX.Kl7nnbOq','user',1,NULL,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MzYzMDE0OTMsImlkIjoxLCJsb2dpbklEIjoxLCJyb2xlIjoidXNlciJ9.ifd7NLzbKzHhC9vAzss0XVQWrIlXkWRWDP0JCWmQav0','2021-11-07 20:56:54.758','2021-11-07 22:11:33.037'),(2,'staff1@gmail.com','staff1','$2a$14$c5qRR0TnF1MYpD.6pBrT2.UsvJzz5XXDwMcPoHkUdstS.4S5CQybu','staff',NULL,1,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MzYzMDI5NzAsImlkIjoxLCJsb2dpbklEIjoyLCJyb2xlIjoic3RhZmYifQ.I1gM7JtEFvmeiMvahXS7yoJhGuFHXRcmxnWjYf6uuBo','2021-11-07 21:17:12.658','2021-11-07 22:36:10.382'),(3,'staff2@gmail.com','staff2','$2a$14$068kkSg0w8grG3S3zzKVDeonqME/wIT2/2TOdjybt8SmucsyJEUnS','staff',NULL,2,NULL,'2021-11-07 22:55:01.330','2021-11-07 22:55:01.330'),(4,'staff3@gmail.com','staff3','$2a$14$8/AcA.eBiJtDBnbsuSoBr.oHqsPJ/mRR3E5Bu0VOzf93PEhkxIiKu','staff',NULL,3,NULL,'2021-11-07 22:55:28.460','2021-11-07 22:55:28.460'),(5,'staff4@gmail.com','staff4','$2a$14$I5tQwVzkdm2QRs.5QN5zqO63llk0ip8VKqyyedK2142HzjHvGjXO.','staff',NULL,4,NULL,'2021-11-07 22:55:45.565','2021-11-07 22:55:45.565'),(6,'staff5@gmail.com','staff5','$2a$14$kKESAwjzgGXWc7Jz0Ifi4OBPN.NsdxA/ovxDv4RtYyX8r2e60oOM2','staff',NULL,5,NULL,'2021-11-07 22:56:07.617','2021-11-07 22:56:07.617');
/*!40000 ALTER TABLE `logins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `staffs`
--

DROP TABLE IF EXISTS `staffs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `staffs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `fullname` varchar(100) NOT NULL,
  `phone_number` varchar(15) NOT NULL,
  `drop_point_id` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `phone_number` (`phone_number`),
  KEY `fk_staffs_drop_point` (`drop_point_id`),
  CONSTRAINT `fk_staffs_drop_point` FOREIGN KEY (`drop_point_id`) REFERENCES `drop_points` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `staffs`
--

LOCK TABLES `staffs` WRITE;
/*!40000 ALTER TABLE `staffs` DISABLE KEYS */;
INSERT INTO `staffs` (`id`, `fullname`, `phone_number`, `drop_point_id`, `created_at`, `updated_at`) VALUES (1,'staff1','081258963218',1,'2021-11-07 21:17:12.413','2021-11-07 21:41:28.045'),(2,'staff2','081258963214',2,'2021-11-07 22:55:01.140','2021-11-07 22:55:01.140'),(3,'staff3','08125896321',3,'2021-11-07 22:55:27.938','2021-11-07 22:55:27.938'),(4,'staff4','081258963213',4,'2021-11-07 22:55:45.296','2021-11-07 22:55:45.296'),(5,'staff5','081258963215',5,'2021-11-07 22:56:07.359','2021-11-07 22:56:07.359');
/*!40000 ALTER TABLE `staffs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `transactions`
--

DROP TABLE IF EXISTS `transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `transactions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `status` tinyint NOT NULL DEFAULT '0',
  `point` bigint NOT NULL,
  `method` enum('dropoff','pickup') DEFAULT NULL,
  `drop_point_id` bigint NOT NULL,
  `checkout_id` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_transactions_user` (`user_id`),
  KEY `fk_transactions_drop_point` (`drop_point_id`),
  KEY `fk_transactions_checkout` (`checkout_id`),
  CONSTRAINT `fk_transactions_checkout` FOREIGN KEY (`checkout_id`) REFERENCES `checkouts` (`id`),
  CONSTRAINT `fk_transactions_drop_point` FOREIGN KEY (`drop_point_id`) REFERENCES `drop_points` (`id`),
  CONSTRAINT `fk_transactions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `transactions`
--

LOCK TABLES `transactions` WRITE;
/*!40000 ALTER TABLE `transactions` DISABLE KEYS */;
INSERT INTO `transactions` (`id`, `user_id`, `status`, `point`, `method`, `drop_point_id`, `checkout_id`, `created_at`, `updated_at`) VALUES (1,1,1,8,'pickup',1,1,'2021-11-07 22:16:17.695','2021-11-07 22:20:40.866'),(2,1,1,12,'dropoff',1,2,'2021-11-07 22:17:29.145','2021-11-07 22:23:08.882');
/*!40000 ALTER TABLE `transactions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_vouchers`
--

DROP TABLE IF EXISTS `user_vouchers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_vouchers` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `voucher_id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `status` enum('used','available') DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_user_vouchers_voucher` (`voucher_id`),
  KEY `fk_user_vouchers_user` (`user_id`),
  CONSTRAINT `fk_user_vouchers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_user_vouchers_voucher` FOREIGN KEY (`voucher_id`) REFERENCES `vouchers` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_vouchers`
--

LOCK TABLES `user_vouchers` WRITE;
/*!40000 ALTER TABLE `user_vouchers` DISABLE KEYS */;
INSERT INTO `user_vouchers` (`id`, `voucher_id`, `user_id`, `status`, `created_at`, `updated_at`) VALUES (1,1,1,'used','2021-11-07 22:24:57.370','2021-11-07 22:26:49.091');
/*!40000 ALTER TABLE `user_vouchers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `fullname` varchar(100) NOT NULL,
  `phone_number` varchar(15) NOT NULL,
  `gender` enum('male','female') NOT NULL,
  `address` longtext NOT NULL,
  `point` bigint NOT NULL DEFAULT '0',
  `longitude` longtext,
  `latitude` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `phone_number` (`phone_number`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` (`id`, `fullname`, `phone_number`, `gender`, `address`, `point`, `longitude`, `latitude`, `created_at`, `updated_at`) VALUES (1,'darien kentanu','0895379933379','male','jalan kesehatan 3 ciputat',10,'106.7543732','-6.3400566','2021-11-07 20:56:54.194','2021-11-07 22:24:57.586');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vouchers`
--

DROP TABLE IF EXISTS `vouchers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vouchers` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `point` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vouchers`
--

LOCK TABLES `vouchers` WRITE;
/*!40000 ALTER TABLE `vouchers` DISABLE KEYS */;
INSERT INTO `vouchers` (`id`, `name`, `point`, `created_at`, `updated_at`) VALUES (1,'voucher pulsa 10rb',10,'2021-11-07 21:45:29.307','2021-11-07 21:59:21.327'),(3,'voucher pulsa 20rb',20,'2021-11-07 22:45:31.689','2021-11-07 22:45:31.689'),(4,'voucher ayam mcd',30,'2021-11-07 22:46:22.627','2021-11-07 22:46:22.627'),(5,'voucher uang tunai 100rb',100,'2021-11-07 22:46:51.615','2021-11-07 22:46:51.615'),(6,'voucher indomaret 100rb',100,'2021-11-07 22:47:23.042','2021-11-07 22:47:23.042');
/*!40000 ALTER TABLE `vouchers` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-11-07 22:58:23
