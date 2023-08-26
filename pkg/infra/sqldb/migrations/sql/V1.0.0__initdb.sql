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

create table if not exists thread_versioned
(
    id           uuid primary key                  default uuid_generate_v4(),
    name         varchar                  not null,
    version      int                      not null,
    thread_id    uuid                     not null,
    archived     bool                     not null default false,
    deleted      bool                     not null default false,
    date_created timestamp with time zone not null default current_timestamp,

    constraint fk_thread_versioned_thread_id foreign key (thread_id) references thread (id),
    constraint unique_thread_versioned_thread_id_version unique (thread_id, version),
    constraint unique_thread_versioned_thread_id_name unique (thread_id, name)
);

create table if not exists string_versioned
(
    id           uuid primary key                  default uuid_generate_v4(),
    name         varchar                  not null,
    version      int                      not null,
    string_id    uuid                     not null,
    thread_id    uuid                     not null,
    "order"      int                      not null,
    active       bool,
    archived     bool                     not null default false,
    deleted      bool                     not null default false,
    date_created timestamp with time zone not null default current_timestamp,

    constraint fk_string_versioned_thread_id foreign key (thread_id) references thread (id),
    constraint fk_string_versioned_string_id foreign key (string_id) references string (id),
    constraint unique_string_versioned_thread_id_version unique (string_id, version)
);
