create
    extension if not exists "uuid-ossp";

create table if not exists thread
(
    id uuid primary key default uuid_generate_v4()
);

create table if not exists string
(
    id uuid primary key default uuid_generate_v4()
);

create table if not exists versioned_thread
(
    id           uuid primary key                  default uuid_generate_v4(),
    name         varchar                  not null,
    version      int                      not null,
    thread_id    uuid                     not null,
    archived     bool                     not null default false,
    deleted      bool                     not null default false,
    date_created timestamp with time zone not null default current_timestamp,

    constraint fk_versioned_thread_thread_id foreign key (thread_id) references thread (id),
    constraint unique_versioned_thread_thread_id_version unique (thread_id, version)
);

create table if not exists versioned_string
(
    id           uuid primary key                  default uuid_generate_v4(),
    name         varchar                  not null,
    version      int                      not null,
    string_id    uuid                     not null,
    thread_id    uuid                     not null,
    "order"      int                      not null,
    active       bool                     not null default false,
    archived     bool                     not null default false,
    deleted      bool                     not null default false,
    date_created timestamp with time zone not null default current_timestamp,

    constraint fk_versioned_string_thread_id foreign key (thread_id) references thread (id),
    constraint fk_versioned_string_string_id foreign key (string_id) references string (id),
    constraint unique_versioned_string_thread_id_version unique (string_id, version)
);

CREATE INDEX ON versioned_thread (id, thread_id, version, date_created);
CREATE INDEX ON versioned_string (id, string_id, thread_id, version, date_created);