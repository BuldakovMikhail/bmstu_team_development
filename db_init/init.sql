CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users
(
    id_user       UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    user_name     varchar(50) UNIQUE NOT NULL,
    password_hash varchar(256)       NOT NULL
);

CREATE TABLE task
(
    id_task     UUID PRIMARY KEY      DEFAULT (gen_random_uuid()),
    user_id     UUID         NOT NULL,
    title       varchar(128) NOT NULL,
    description varchar(1000),
    is_done     boolean      NOT NULL DEFAULT false
);

CREATE TABLE category
(
    id_category UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
    user_id     UUID        NOT NULL,
    name        varchar(50) NOT NULL
);

CREATE TABLE task_category
(
    task_id     uuid NOT NULL,
    category_id UUID NOT NULL,
    PRIMARY KEY (task_id, category_id)
);

CREATE INDEX ON task (user_id);
CREATE INDEX ON category (user_id);
CREATE UNIQUE INDEX ON category (user_id, name);

ALTER TABLE category
    ADD FOREIGN KEY (user_id) REFERENCES users (id_user) ON DELETE CASCADE;

ALTER TABLE task_category
    ADD FOREIGN KEY (task_id) REFERENCES task (id_task) ON DELETE CASCADE,
    ADD FOREIGN KEY (category_id) REFERENCES category (id_category) ON DELETE CASCADE;

ALTER TABLE task
    ADD FOREIGN KEY (user_id) REFERENCES users (id_user) ON DELETE CASCADE;