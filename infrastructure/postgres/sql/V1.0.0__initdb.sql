--
-- Initial setup
--
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--
-- Thread
--
CREATE TABLE IF NOT EXISTS thread
(
    id            UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    name          VARCHAR UNIQUE           NOT NULL,
    description   VARCHAR,
    date_created  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date_modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--
-- String
--
CREATE TABLE IF NOT EXISTS string
(
    id            UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    name          VARCHAR UNIQUE           NOT NULL,
    "order"       INT                      NOT NULL,
    thread        UUID                     NOT NULL,
    description   VARCHAR,
    date_created  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date_modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_thread_id FOREIGN KEY (thread) REFERENCES thread (id)
);

-- --
-- -- Skill
-- --
-- CREATE TABLE IF NOT EXISTS skill
-- (
--     id                UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
--     skill_id          CHAR(3)                  NOT NULL,
--     user_id           UUID,
--     exp               INTEGER                  NOT NULL DEFAULT 0,
--     txp               INTEGER                  NOT NULL DEFAULT 0,
--     level             SMALLINT                          DEFAULT 1,
--     date_created      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     date_modified     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     date_last_txp_add TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
--
--     CONSTRAINT fk_skill_id
--         FOREIGN KEY (skill_id)
--             REFERENCES skill_config (id),
--     CONSTRAINT fk_user_id
--         FOREIGN KEY (user_id)
--             REFERENCES user_account (id)
-- );
--

