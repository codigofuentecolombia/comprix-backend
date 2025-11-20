-- ============================================
-- Comprix Database Schema
-- ============================================

-- Crear la base de datos (opcional, comentar si ya existe)
-- CREATE DATABASE IF NOT EXISTS comprix_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- USE comprix_db;

-- ============================================
-- Tabla: roles
-- ============================================
CREATE TABLE IF NOT EXISTS `roles` (
  `id` TINYINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(250) NOT NULL UNIQUE,
  `display_name` VARCHAR(250) NOT NULL,
  `description` TEXT NOT NULL DEFAULT '',
  `created_at` TIMESTAMP NULL DEFAULT NULL,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: users
-- ============================================
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `role_id` TINYINT UNSIGNED NOT NULL,
  `email` VARCHAR(250) NOT NULL UNIQUE,
  `picture` VARCHAR(250) NULL,
  `username` VARCHAR(250) NOT NULL UNIQUE,
  `password` VARCHAR(250) NOT NULL,
  `first_name` VARCHAR(250) NOT NULL,
  `last_name` VARCHAR(250) NOT NULL,
  `phone_number` VARCHAR(20) NOT NULL,
  `is_verified` TINYINT(1) NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP NULL DEFAULT NULL,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_users_deleted_at` (`deleted_at`),
  FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: verification_codes
-- ============================================
CREATE TABLE IF NOT EXISTS `verification_codes` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `code` VARCHAR(255) NOT NULL,
  `expiration` TIMESTAMP NOT NULL,
  `claimed` TINYINT(1) NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_verification_codes_user_id` (`user_id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: user_shipping_addresses
-- ============================================
CREATE TABLE IF NOT EXISTS `user_shipping_addresses` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `date` TIMESTAMP NOT NULL,
  `time` VARCHAR(50) NOT NULL,
  `alias` VARCHAR(255) NOT NULL DEFAULT 'Principal',
  `street` VARCHAR(255) NOT NULL,
  `colony` VARCHAR(255) NOT NULL,
  `state` VARCHAR(255) NOT NULL,
  `city` VARCHAR(255) NOT NULL,
  `postal_code` VARCHAR(20) NOT NULL,
  `phone_number` VARCHAR(20) NOT NULL,
  `reference` TEXT NULL,
  `external_number` VARCHAR(50) NOT NULL,
  `internal_number` VARCHAR(50) NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_user_shipping_addresses_user_id` (`user_id`),
  INDEX `idx_user_shipping_addresses_deleted_at` (`deleted_at`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: brands
-- ============================================
CREATE TABLE IF NOT EXISTS `brands` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(250) NOT NULL UNIQUE,
  `created_at` TIMESTAMP NULL DEFAULT NULL,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_brands_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: categories
-- ============================================
CREATE TABLE IF NOT EXISTS `categories` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(250) NOT NULL UNIQUE,
  `parent_id` INT UNSIGNED NULL,
  `created_at` TIMESTAMP NULL DEFAULT NULL,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_categories_deleted_at` (`deleted_at`),
  INDEX `idx_categories_parent_id` (`parent_id`),
  FOREIGN KEY (`parent_id`) REFERENCES `categories`(`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: products
-- ============================================
CREATE TABLE IF NOT EXISTS `products` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `sku` VARCHAR(100) NOT NULL UNIQUE,
  `name` VARCHAR(255) NOT NULL,
  `brand_id` INT NULL,
  `category_id` INT UNSIGNED NULL,
  `is_disabled` TINYINT(1) NOT NULL DEFAULT 0,
  `description` TEXT NULL,
  `is_in_discount` TINYINT(1) NOT NULL DEFAULT 0,
  `is_recommended` TINYINT(1) NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP NULL DEFAULT NULL,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_products_deleted_at` (`deleted_at`),
  INDEX `idx_products_brand_id` (`brand_id`),
  INDEX `idx_products_category_id` (`category_id`),
  FOREIGN KEY (`brand_id`) REFERENCES `brands`(`id`) ON DELETE SET NULL,
  FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: product_groups
-- ============================================
CREATE TABLE IF NOT EXISTS `product_groups` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `reference_product_id` BIGINT UNSIGNED NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_product_groups_reference_product_id` (`reference_product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: product_group_members
-- ============================================
CREATE TABLE IF NOT EXISTS `product_group_members` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `group_id` BIGINT UNSIGNED NOT NULL,
  `product_id` BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_product_group_members_group_id` (`group_id`),
  INDEX `idx_product_group_members_product_id` (`product_id`),
  FOREIGN KEY (`group_id`) REFERENCES `product_groups`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT,
  FOREIGN KEY (`product_id`) REFERENCES `products`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: pages
-- ============================================
CREATE TABLE IF NOT EXISTS `pages` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `url` VARCHAR(255) NOT NULL UNIQUE,
  `logo` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT NULL,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_pages_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: page_products
-- ============================================
CREATE TABLE IF NOT EXISTS `page_products` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `url` VARCHAR(500) NOT NULL,
  `images` JSON NULL,
  `page_id` BIGINT UNSIGNED NOT NULL,
  `product_id` BIGINT UNSIGNED NOT NULL,
  `main_product_id` BIGINT UNSIGNED NOT NULL,
  `price` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `discount_price` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `original_price` VARCHAR(100) NULL,
  `original_discount_price` VARCHAR(100) NULL,
  `min_quantity_to_apply_discount` INT(11) NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP NULL DEFAULT NULL,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `deleted_at` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `idx_page_products_page_id` (`page_id`),
  INDEX `idx_page_products_product_id` (`product_id`),
  INDEX `idx_page_products_main_product_id` (`main_product_id`),
  INDEX `idx_page_products_deleted_at` (`deleted_at`),
  FOREIGN KEY (`page_id`) REFERENCES `pages`(`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  FOREIGN KEY (`product_id`) REFERENCES `products`(`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  FOREIGN KEY (`main_product_id`) REFERENCES `products`(`id`) ON UPDATE CASCADE ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: orders
-- ============================================
CREATE TABLE IF NOT EXISTS `orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `status` ENUM('pending','processing','completed','cancelled') NOT NULL DEFAULT 'pending',
  `total` DECIMAL(10,2) NOT NULL,
  `subtotal` DECIMAL(10,2) NOT NULL,
  `total_discount` DECIMAL(10,2) NOT NULL,
  `shipping_cost` DECIMAL(10,2) NOT NULL,
  `payment_method` VARCHAR(50) NOT NULL,
  `user_shipping_address_id` BIGINT UNSIGNED NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_orders_user_id` (`user_id`),
  INDEX `idx_orders_user_shipping_address_id` (`user_shipping_address_id`),
  INDEX `idx_orders_deleted_at` (`deleted_at`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`user_shipping_address_id`) REFERENCES `user_shipping_addresses`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: order_products
-- ============================================
CREATE TABLE IF NOT EXISTS `order_products` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `quantity` INT UNSIGNED NOT NULL,
  `order_id` BIGINT UNSIGNED NOT NULL,
  `page_product_id` BIGINT UNSIGNED NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_order_products_quantity` (`quantity`),
  INDEX `idx_order_products_order_id` (`order_id`),
  INDEX `idx_order_products_page_product_id` (`page_product_id`),
  INDEX `idx_order_products_deleted_at` (`deleted_at`),
  FOREIGN KEY (`order_id`) REFERENCES `orders`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`page_product_id`) REFERENCES `page_products`(`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Tabla: calendar
-- ============================================
CREATE TABLE IF NOT EXISTS `calendar` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `start_time` TIME NOT NULL,
  `end_time` TIME NOT NULL,
  `price` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `monday` TINYINT(1) NOT NULL DEFAULT 0,
  `tuesday` TINYINT(1) NOT NULL DEFAULT 0,
  `wednesday` TINYINT(1) NOT NULL DEFAULT 0,
  `thursday` TINYINT(1) NOT NULL DEFAULT 0,
  `friday` TINYINT(1) NOT NULL DEFAULT 0,
  `saturday` TINYINT(1) NOT NULL DEFAULT 0,
  `sunday` TINYINT(1) NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_calendar_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Datos iniciales (opcional)
-- ============================================

-- Insertar roles por defecto
INSERT IGNORE INTO `roles` (`id`, `name`, `display_name`, `description`, `created_at`) VALUES
(1, 'admin', 'Administrador', 'Usuario con acceso completo al sistema', NOW()),
(2, 'user', 'Usuario', 'Usuario regular del sistema', NOW());

-- ============================================
-- FIN DEL SCRIPT
-- ============================================
