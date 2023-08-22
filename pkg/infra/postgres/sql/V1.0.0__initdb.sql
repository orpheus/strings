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

create table if not exists thread_version
(
    id           uuid primary key                  default uuid_generate_v4(),
    name         varchar                  not null,
    version      int                      not null unique,
    thread_id    uuid                     not null,
    archived     bool                     not null default false,
    deleted      bool                     not null default false,
    date_created timestamp with time zone not null default current_timestamp,

    constraint fk_thread_version_thread_id foreign key (thread_id) references thread (id)
);

create table if not exists string_version
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

    constraint fk_string_version_thread_id foreign key (thread_id) references thread (id),
    constraint fk_string_version_string_id foreign key (string_id) references string (id)
);
