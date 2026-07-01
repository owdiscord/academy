-- =============================================================
-- Seed data for modmail schema
-- AI slop generated off my real schema
-- =============================================================

-- Clear and reset

DROP TABLE IF EXISTS blocked_users, cases, case_notes, messages, moderator_role_overrides, notes, snippets, schema_migrations, threads, thread_messages, updates;

--
-- Create the schemas 
--

create table blocked_users
(
    user_id    varchar(20)  not null
        primary key,
    user_name  varchar(128) not null,
    blocked_by varchar(20)  null,
    blocked_at datetime     not null,
    expires_at datetime     null
);

create table cases
(
    id             int unsigned auto_increment
        primary key,
    guild_id       bigint unsigned                            not null,
    case_number    int unsigned                               not null,
    user_id        bigint unsigned                            not null,
    user_name      varchar(128)                               not null,
    mod_id         bigint unsigned                            null,
    mod_name       varchar(128)                               null,
    type           int unsigned                               not null,
    audit_log_id   bigint                                     null,
    created_at     datetime         default CURRENT_TIMESTAMP not null,
    is_hidden      tinyint unsigned default '0'               not null,
    pp_id          bigint                                     null,
    pp_name        varchar(128)                               null,
    log_message_id varchar(64)                                null,
    constraint mod_actions_audit_log_id_unique
        unique (audit_log_id),
    constraint mod_actions_guild_id_case_number_unique
        unique (guild_id, case_number)
)
    collate = utf8mb4_general_ci;

create table case_notes
(
    id         int unsigned auto_increment
        primary key,
    case_id    int unsigned                       not null,
    mod_id     bigint unsigned                    null,
    mod_name   varchar(128)                       null,
    body       text                               not null,
    created_at datetime default CURRENT_TIMESTAMP not null,
    constraint case_notes_case_id_fk
        foreign key (case_id) references cases (id)
            on update cascade on delete cascade
)
    collate = utf8mb4_general_ci;

create index mod_action_notes_created_at_index
    on case_notes (created_at);

create index mod_action_notes_mod_action_id_index
    on case_notes (case_id);

create index mod_action_notes_mod_id_index
    on case_notes (mod_id);

create index IDX_9103dcd7dac7ddd60068296cec
    on cases (is_hidden);

create index mod_actions_created_at_index
    on cases (created_at);

create index mod_actions_mod_id_index
    on cases (mod_id);

create index mod_actions_user_id_index
    on cases (user_id);

create table messages
(
    id         bigint auto_increment
        primary key,
    user_id    bigint   default 0                 not null,
    channel_id bigint   default 0                 not null,
    is_bot     smallint default 0                 null,
    posted_at  datetime default CURRENT_TIMESTAMP not null
);

create table moderator_role_overrides
(
    id           int unsigned auto_increment
        primary key,
    moderator_id varchar(20) not null,
    thread_id    varchar(36) null,
    role_id      varchar(20) not null,
    constraint moderator_role_overrides_moderator_id_thread_id_unique
        unique (moderator_id, thread_id)
);

create table notes
(
    id         int unsigned auto_increment
        primary key,
    user_id    varchar(20) null,
    author_id  varchar(20) null,
    body       mediumtext  null,
    created_at datetime    null
);

create index notes_author_id_index
    on notes (author_id);

create index notes_user_id_index
    on notes (user_id);

create table registered_users
(
    discord_id      varchar(25)                        not null
        primary key,
    registered_name varchar(50)                        null,
    created_at      datetime default CURRENT_TIMESTAMP not null,
    updated_at      datetime                           null
);

create table schema_migrations
(
    id             int unsigned auto_increment
        primary key,
    name           varchar(255) null,
    batch          int          null,
    migration_time timestamp    null
);

create table snippets
(
    `trigger`  varchar(32) not null
        primary key,
    body       text        not null,
    created_by varchar(20) null,
    created_at datetime    not null
);

create table threads
(
    id                     varchar(36)                         not null
        primary key,
    status                 int unsigned                        not null,
    is_legacy              int unsigned                        not null,
    user_id                varchar(20)                         not null,
    user_name              varchar(128)                        not null,
    channel_id             varchar(20)                         null,
    scheduled_suspend_name varchar(128)                        null,
    scheduled_suspend_id   varchar(20)                         null,
    scheduled_suspend_at   datetime                            null,
    scheduled_close_name   varchar(128)                        null,
    scheduled_close_silent int                                 null,
    scheduled_close_id     varchar(20)                         null,
    scheduled_close_at     datetime                            null,
    created_at             datetime                            not null,
    next_message_number    int       default 1                 null,
    alert_ids              text                                null,
    log_storage_type       varchar(255)                        null,
    log_storage_data       text                                null,
    metadata               text                                null,
    thread_number          int                                 null,
    closed_by_id           varchar(20)                         null,
    closed_at              datetime                            null,
    roles                  varchar(512)                        null,
    server_join            datetime                            null,
    updated_at             timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint threads_channel_id_unique
        unique (channel_id),
    constraint threads_thread_number_unique
        unique (thread_number)
);

create table thread_messages
(
    id                int unsigned auto_increment
        primary key,
    thread_id         varchar(36)                         not null,
    message_type      int unsigned                        not null,
    user_id           varchar(20)                         null,
    user_name         varchar(128)                        not null,
    is_anonymous      int unsigned                        not null,
    dm_message_id     varchar(20)                         null,
    created_at        datetime                            not null,
    message_number    int unsigned                        null,
    inbox_message_id  varchar(20)                         null,
    dm_channel_id     varchar(20)                         null,
    role_name         varchar(255)                        null,
    attachments       text                                null,
    small_attachments text                                null,
    use_legacy_format tinyint(1)                          null,
    metadata          text                                null,
    body              text                                null,
    updated_at        timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint thread_messages_dm_message_id_unique
        unique (dm_message_id),
    constraint thread_messages_inbox_message_id_unique
        unique (inbox_message_id),
    constraint thread_messages_thread_id_foreign
        foreign key (thread_id) references threads (id)
            on delete cascade
);

create index thread_messages_created_at_index
    on thread_messages (created_at);

create index thread_messages_thread_id_index
    on thread_messages (thread_id);

create index closed_by_id_idx
    on threads (closed_by_id);

create index threads_created_at_index
    on threads (created_at);

create index threads_scheduled_close_at_index
    on threads (scheduled_close_at);

create index threads_scheduled_suspend_at_index
    on threads (scheduled_suspend_at);

create index threads_status_index
    on threads (status);

create index threads_user_id_index
    on threads (user_id);

create table updates
(
    available_version varchar(16) null,
    last_checked      datetime    null
);

-- -------------------------------------------------------------
-- blocked_users
-- -------------------------------------------------------------
INSERT INTO blocked_users (user_id, user_name, blocked_by, blocked_at, expires_at) VALUES
('111000000000000001', 'spammer', '999000000000000001', '2026-01-10 09:00:00', NULL),
('111000000000000002', 'troll',   '999000000000000001', '2026-02-14 12:00:00', '2026-08-14 12:00:00');

-- -------------------------------------------------------------
-- cases
-- -------------------------------------------------------------
INSERT INTO cases (guild_id, case_number, user_id, user_name, mod_id, mod_name, type, created_at, is_hidden) VALUES
(123456789000000001, 1, 111000000000000003, 'user_one',   999000000000000001, 'mod_alice', 1, '2026-01-15 10:00:00', 0),
(123456789000000001, 2, 111000000000000004, 'user_two',   999000000000000002, 'mod_bob',   2, '2026-02-20 14:30:00', 0),
(123456789000000001, 3, 111000000000000001, 'spammer',    999000000000000001, 'mod_alice', 4, '2026-03-05 08:45:00', 0),
(123456789000000001, 4, 111000000000000005, 'user_three', 999000000000000002, 'mod_bob',   1, '2026-04-01 16:00:00', 1);

-- -------------------------------------------------------------
-- case_notes
-- -------------------------------------------------------------
INSERT INTO case_notes (case_id, mod_id, mod_name, body, created_at) VALUES
(1, 999000000000000001, 'mod_alice', 'User was warned for posting spam links in #general.', '2026-01-15 10:05:00'),
(1, 999000000000000002, 'mod_bob',   'Follow-up: user acknowledged the warning.', '2026-01-16 09:00:00'),
(2, 999000000000000002, 'mod_bob',   'Muted for 1 hour after repeated off-topic posts.', '2026-02-20 14:35:00'),
(3, 999000000000000001, 'mod_alice', 'Permanent ban issued after third spam offence.', '2026-03-05 08:50:00');

-- -------------------------------------------------------------
-- registered_users
-- -------------------------------------------------------------
INSERT INTO registered_users (discord_id, registered_name, created_at, updated_at) VALUES
('111000000000000003', 'user_one',   '2025-06-01 12:00:00', NULL),
('111000000000000004', 'user_two',   '2025-07-15 08:30:00', '2026-01-10 10:00:00'),
('111000000000000005', 'user_three', '2025-09-20 17:45:00', NULL);

-- -------------------------------------------------------------
-- snippets
-- -------------------------------------------------------------
INSERT INTO snippets (`trigger`, body, created_by, created_at) VALUES
('rules',   'Please read our server rules in #rules before continuing.',        '999000000000000001', '2025-01-01 00:00:00'),
('appeal',  'To appeal a moderation action, fill out the form at example.com/appeal.', '999000000000000001', '2025-01-01 00:00:00'),
('faq',     'Check out our FAQ channel at #faq for common questions.',          '999000000000000002', '2025-03-10 11:00:00');

-- -------------------------------------------------------------
-- notes
-- -------------------------------------------------------------
INSERT INTO notes (user_id, author_id, body, created_at) VALUES
('111000000000000003', '999000000000000001', 'Previously warned on old account.',        '2025-12-01 09:00:00'),
('111000000000000004', '999000000000000002', 'Known to evade mutes — watch closely.',    '2026-01-05 14:00:00'),
('111000000000000005', '999000000000000001', 'Friendly user, no issues so far.',         '2026-02-28 10:30:00');

-- -------------------------------------------------------------
-- threads
-- -------------------------------------------------------------
INSERT INTO threads (
    id, status, is_legacy, user_id, user_name, channel_id,
    created_at, next_message_number, thread_number,
    closed_by_id, closed_at
) VALUES
(
    'a1b2c3d4-0001-0001-0001-000000000001', 1, 0,
    '111000000000000003', 'user_one', '222000000000000001',
    '2026-03-01 10:00:00', 3, 1,
    NULL, NULL
),
(
    'a1b2c3d4-0002-0002-0002-000000000002', 2, 0,
    '111000000000000004', 'user_two', '222000000000000002',
    '2026-03-15 14:00:00', 5, 2,
    '999000000000000002', '2026-03-15 15:30:00'
),
(
    'a1b2c3d4-0003-0003-0003-000000000003', 2, 0,
    '111000000000000005', 'user_three', '222000000000000003',
    '2026-04-10 09:00:00', 2, 3,
    '999000000000000001', '2026-04-10 09:45:00'
);

-- -------------------------------------------------------------
-- thread_messages
-- -------------------------------------------------------------
INSERT INTO thread_messages (
    thread_id, message_type, user_id, user_name, is_anonymous,
    dm_message_id, inbox_message_id, dm_channel_id,
    body, created_at, message_number
) VALUES
-- Thread 1: open thread, two messages
(
    'a1b2c3d4-0001-0001-0001-000000000001', 1,
    '111000000000000003', 'user_one', 0,
    '333000000000000001', '444000000000000001', '555000000000000001',
    'Hi, I need help with my mute appeal.', '2026-03-01 10:01:00', 1
),
(
    'a1b2c3d4-0001-0001-0001-000000000001', 2,
    '999000000000000001', 'mod_alice', 0,
    '333000000000000002', '444000000000000002', '555000000000000001',
    'Hi! We will look into this for you shortly.', '2026-03-01 10:05:00', 2
),
-- Thread 2: closed thread
(
    'a1b2c3d4-0002-0002-0002-000000000002', 1,
    '111000000000000004', 'user_two', 0,
    '333000000000000003', '444000000000000003', '555000000000000002',
    'Why was I muted?', '2026-03-15 14:01:00', 1
),
(
    'a1b2c3d4-0002-0002-0002-000000000002', 2,
    '999000000000000002', 'mod_bob', 0,
    '333000000000000004', '444000000000000004', '555000000000000002',
    'You were muted for repeated off-topic posts. Please review #rules.', '2026-03-15 14:10:00', 2
),
-- Thread 3: closed thread, quick close
(
    'a1b2c3d4-0003-0003-0003-000000000003', 1,
    '111000000000000005', 'user_three', 0,
    '333000000000000005', '444000000000000005', '555000000000000003',
    'Hello, just checking in about something.', '2026-04-10 09:01:00', 1
),
(
    'a1b2c3d4-0003-0003-0003-000000000003', 2,
    '999000000000000001', 'mod_alice', 0,
    '333000000000000006', '444000000000000006', '555000000000000003',
    'Resolved! Let us know if you need anything else.', '2026-04-10 09:44:00', 2
);

-- -------------------------------------------------------------
-- messages (activity log)
-- -------------------------------------------------------------
INSERT INTO messages (user_id, channel_id, is_bot, posted_at) VALUES
(111000000000000003, 666000000000000001, 0, '2026-05-01 10:00:00'),
(111000000000000003, 666000000000000001, 0, '2026-05-01 10:02:00'),
(111000000000000004, 666000000000000002, 0, '2026-05-01 11:00:00'),
(111000000000000005, 666000000000000001, 0, '2026-05-02 09:30:00'),
(999000000000000001, 666000000000000001, 0, '2026-05-02 09:35:00'),
(777000000000000001, 666000000000000003, 1, '2026-05-02 10:00:00');

-- -------------------------------------------------------------
-- updates
-- -------------------------------------------------------------
INSERT INTO updates (available_version, last_checked) VALUES
('3.4.1', '2026-06-25 00:00:00');
