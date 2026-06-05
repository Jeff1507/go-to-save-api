-- TODO: Criar indexes para consultas mais rápidas --

CREATE TYPE bag_status AS ENUM ('ACTIVE', 'CLOSED');

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    verified_at TIMESTAMPTZ
);

CREATE TABLE stores (
    id UUID PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    cnpj VARCHAR(20) UNIQUE NOT NULL,
    picture_key TEXT,
    description TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
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

CREATE TABLE pickup_templates (
    id UUID PRIMARY KEY NOT NULL,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,

    UNIQUE(store_id, name)
);

CREATE TABLE bags (
    id UUID PRIMARY KEY NOT NULL,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
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

CREATE TABLE bag_items (
    bag_id UUID NOT NULL REFERENCES bags(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,

    PRIMARY KEY(bag_id, product_id)
);
