CREATE TABLE IF NOT EXISTS clients
(
    id          serial          not null unique,
    external_id varchar (64)    not null unique,
    username    varchar(128)    not null,
    phone       varchar(16)     not null unique,
    email       varchar(128)
);

CREATE TABLE IF NOT EXISTS manufacturers
(
    id          serial          not null unique,
    external_id varchar (64)    not null unique,
    name        varchar(128)    not null unique,
    code        varchar(16)     not null unique
);

CREATE TABLE IF NOT EXISTS products
(
    id              serial          not null unique,
    external_id     varchar (64)    not null unique,
    name            varchar(256)    not null,
    expires_at      timestamp       not null,
    manufacturer_id int             not null references manufacturers(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS orders
(
    id              serial          not null unique,
    external_id     varchar (64)    not null unique,
    quantity        smallint        not null,
    manufacturer_id int             not null references manufacturers(id) on delete cascade,
    client_id       int             not null references clients(id) on delete cascade
);
