-- TODO: Criar indexes para consultas mais rápidas --
-- TODO: Adicionar query para anonimizar users, stores e customers por causa da LGPD --

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE bag_status AS ENUM ('ACTIVE', 'CLOSED');

CREATE TYPE order_status AS ENUM ('PENDING', 'COMPLETED', 'CANCELLED', 'NO_SHOW');

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    verified_at TIMESTAMPTZ
);

CREATE TABLE customers (
    id UUID PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    avatar_key TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE stores (
    id UUID PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    cnpj VARCHAR(14) UNIQUE NOT NULL,
    avatar_key TEXT,
    description TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE addresses (
    store_id UUID PRIMARY KEY NOT NULL REFERENCES stores(id) ON DELETE RESTRICT,
    postal_code VARCHAR(20) NOT NULL,
    street VARCHAR(150) NOT NULL,
    number VARCHAR(20) NOT NULL,
    neighborhood VARCHAR(100) NOT NULL,
    complement TEXT,
    city VARCHAR(100) NOT NULL,
    state CHAR(2) NOT NULL,

    latitude NUMERIC(10, 7) NOT NULL,
    longitude NUMERIC(10, 7) NOT NULL,
    
    location GEOGRAPHY(POINT, 4326) GENERATED ALWAYS AS (
        ST_SetSRID(
            ST_MakePoint(longitude, latitude),
            4326
        )::geography
    ) STORED
);

CREATE TABLE products (
    id UUID PRIMARY KEY NOT NULL,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    image_key TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    UNIQUE(store_id, name)
);

-- Preset para cadastrar sacolas de forma mais rápida --
CREATE TABLE pickup_presets (
    id UUID PRIMARY KEY NOT NULL,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,

    UNIQUE(store_id, name)
);

CREATE TABLE bags (
    id UUID PRIMARY KEY NOT NULL,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE RESTRICT,
    title VARCHAR(50) NOT NULL,
    description TEXT,
    stock INTEGER NOT NULL,
    original_price BIGINT NOT NULL,
    final_price BIGINT NOT NULL,
    status bag_status DEFAULT 'ACTIVE' NOT NULL,
    pickup_start TIMESTAMPTZ NOT NULL,
    pickup_end TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Snapshot de product --
CREATE TABLE bag_items (
    id UUID PRIMARY KEY NOT NULL,
    bag_id UUID NOT NULL REFERENCES bags(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    image_key TEXT,
    quantity INTEGER NOT NULL,

    UNIQUE(bag_id, name)
);

CREATE TABLE orders (
    id UUID PRIMARY KEY NOT NULL,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE RESTRICT,
    status order_status DEFAULT 'PENDING' NOT NULL,
    total_price BIGINT NOT NULL,
    pickup_start TIMESTAMPTZ NOT NULL,
    pickup_end TIMESTAMPTZ NOT NULL,
    confirmation_code CHAR(4) NOT NULL,
    completed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    UNIQUE(store_id, confirmation_code)
);

CREATE TABLE order_items (
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    bag_id UUID NOT NULL REFERENCES bags(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL,
    unit_price BIGINT NOT NULL,

    PRIMARY KEY(order_id, bag_id)
);

CREATE TABLE reviews (
    id UUID PRIMARY KEY NOT NULL,
    order_id UUID UNIQUE NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    rating SMALLINT NOT NULL,
    comment TEXT,

    reviewed_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE stores_rating (
    store_id UUID PRIMARY KEY NOT NULL REFERENCES stores(id) ON DELETE RESTRICT,
    sum INTEGER NOT NULL,
    count INTEGER NOT NULL
);
