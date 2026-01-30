-- ========================================
-- Homebox 数据库初始化脚本
-- 创建日期: 2025-01-14
-- 说明: 包含所有表的创建语句
-- ========================================

-- 1. 用户表
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `name` varchar(50) NOT NULL COMMENT '用户名',
    `email` varchar(255) NOT NULL COMMENT '邮箱',
    `password` varchar(255) NOT NULL COMMENT '密码(加密)',
    `avatar` varchar(255) DEFAULT '' COMMENT '头像URL',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1-正常，0-禁用',
    `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 2. 空间表
DROP TABLE IF EXISTS `spaces`;
CREATE TABLE `spaces` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '空间ID',
    `user_id` bigint unsigned NOT NULL COMMENT '拥有者用户ID',
    `name` varchar(100) NOT NULL COMMENT '空间名称',
    `description` varchar(500) DEFAULT '' COMMENT '空间描述',
    `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_spaces_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='空间表';

-- 3. 标签分组表
DROP TABLE IF EXISTS `tag_groups`;
CREATE TABLE `tag_groups` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '分组ID',
    `space_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '空间ID',
    `name` varchar(100) NOT NULL COMMENT '分组名称',
    `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tag_groups_space_id` (`space_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签分组表';

-- 4. 标签表
DROP TABLE IF EXISTS `tags`;
CREATE TABLE `tags` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '标签ID',
    `space_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '空间ID',
    `content` text COMMENT '标签内容',
    `group_id` bigint unsigned NOT NULL COMMENT '所属分组ID',
    `color` varchar(20) NOT NULL COMMENT '标签颜色',
    `order_no` int NOT NULL DEFAULT '0' COMMENT '排序序号',
    `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_tags_space_id` (`space_id`),
    KEY `idx_tags_group_id` (`group_id`),
    CONSTRAINT `fk_tags_group` FOREIGN KEY (`group_id`) REFERENCES `tag_groups` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签表';

-- 5. 空间成员表
DROP TABLE IF EXISTS `space_members`;
CREATE TABLE `space_members` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '成员ID',
    `space_id` bigint unsigned NOT NULL COMMENT '空间ID',
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `is_owner` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为所有者：1-是，0-否',
    `joined_at` datetime(3) DEFAULT NULL COMMENT '加入时间',
    `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_space_members_space_id` (`space_id`),
    KEY `idx_space_members_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='空间成员表';

-- 6. 空间邀请表
DROP TABLE IF EXISTS `space_invites`;
CREATE TABLE `space_invites` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '邀请ID',
    `space_id` bigint unsigned NOT NULL COMMENT '空间ID',
    `invitee_id` bigint unsigned NOT NULL COMMENT '被邀请人用户ID',
    `inviter_id` bigint unsigned NOT NULL COMMENT '邀请人用户ID',
    `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending-待处理，accepted-已接受，rejected-已拒绝',
    `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_space_invites_space_id` (`space_id`),
    KEY `idx_space_invites_invitee_id` (`invitee_id`),
    KEY `idx_space_invites_inviter_id` (`inviter_id`),
    KEY `idx_space_invites_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='空间邀请表';

DROP TABLE IF EXISTS `notes`;
CREATE TABLE `notes` (
                         `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '笔记ID',
                         `space_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '空间ID',
                         `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
                         `content` text COMMENT '笔记内容（JSON格式，Tiptap格式）',
                         `created_at` datetime(3) DEFAULT NULL,
                         `updated_at` datetime(3) DEFAULT NULL,
                         `deleted_at` datetime(3) DEFAULT NULL,
                         PRIMARY KEY (`id`),
                         KEY `idx_notes_space_id` (`space_id`),
                         KEY `idx_notes_user_id` (`user_id`),
                         KEY `idx_notes_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='笔记表';
