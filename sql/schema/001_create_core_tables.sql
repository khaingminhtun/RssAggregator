-- +goose Up
-- ----------------------------
-- 1. Setup UUID extension
-- ----------------------------
create extension if not exists "uuid-ossp";

-- ----------------------------
-- 2. users table ( Authentication  )
-- ----------------------------
create table users (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    email text not null unique,
    password_hash text not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now() 
);

-- -------------------------------
-- 3. feeds table ( Globla feed metadat )
-- -------------------------------
CREATE TABLE feeds (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    title text NOT NULL,
    url text NOT NULL UNIQUE,
    last_fetched_at timestamp with time zone NOT NULL DEFAULT '2000-01-01 00:00:00Z'
);

-- --------------------------------
-- 4. user_feeds table ( Subscriptions Mapping)
-- --------------------------------
create table user_feeds (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    feed_id uuid not null references feeds(id) on delete cascade,
    created_at timestamp with time zone not null default now(),
    --Constraints : A user cannot subscribe to the same feed multiple times
    constraint user_feed_uniqueness 
    unique(user_id, feed_id)
);

-- ----------------------------
-- 5. posts table ( Aggregated Content)
-- ----------------------------
CREATE TABLE posts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    feed_id uuid NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    title text NOT NULL,
    url text NOT NULL UNIQUE,
    description text,
    published_at timestamp with time zone NOT NULL,
    guid text NOT NULL UNIQUE, -- For robust deduplication
    created_at timestamp with time zone NOT NULL DEFAULT now()
);

-- +goose Down
-- Drop tables in reverse order to respect foregin key dependencies
drop table posts;
drop table user_feeds;
drop table feeds;
drop table users;
drop extension "uuid-ossp";