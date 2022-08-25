drop index IF EXISTS index_threads_on_user_id;
drop table IF EXISTS threads;
drop index IF EXISTS index_sessions_on_user_id;
drop table IF EXISTS sessions;
drop index IF EXISTS index_shadows_on_user_id;
drop table IF EXISTS shadows;
drop index IF EXISTS index_users_on_person_id;
drop table IF EXISTS users;
drop index IF EXISTS index_persons_on_aparta_id;
drop table IF EXISTS persons;
drop index IF EXISTS index_egresos_on_tipo_id;
drop index IF EXISTS index_egresos_on_period_id;
drop table IF EXISTS egresos;
drop index IF EXISTS index_ingresos_on_tipo_id;
drop index IF EXISTS index_ingresos_on_period_id;
drop table IF EXISTS ingresos;
drop index IF EXISTS index_cuotas_on_tipo_id;
drop index IF EXISTS index_cuotas_on_period_id;
drop index IF EXISTS index_cuotas_on_aparta_id;
drop table IF EXISTS cuotas;
drop index IF EXISTS index_balances_on_period_id;
drop table IF EXISTS balances;
drop table IF EXISTS tipos;
drop table IF EXISTS periods;
drop table IF EXISTS apartas;

CREATE TABLE IF NOT EXISTS apartas (
    id          serial primary key,
    codigo      char(2),
    descripcion varchar(48),
    created_at  timestamp not null,   
    updated_at  timestamp not null   
);

CREATE TABLE IF NOT EXISTS periods (
    id          serial primary key,
    inicio      date not null,
    final       date not null,
    created_at  timestamp not null,   
    updated_at  timestamp not null   
) ;

CREATE TABLE IF NOT EXISTS tipos (
    id          serial primary key,
    codigo      char(2),
    aplica      char(2),
    descripcion varchar(48),
    created_at  timestamp not null,   
    updated_at  timestamp not null   
);

CREATE TABLE IF NOT EXISTS balances (
    id          serial primary key,
    period_id   integer references periods(id) NOT NULL,
    amount      bigint,
    cuota       bigint,
    created_at  timestamp not null,   
    updated_at  timestamp not null   
);

CREATE INDEX  index_balances_on_period_id  ON balances(period_id);

CREATE TABLE IF NOT EXISTS cuotas (
    id          serial primary key,
    period_id   integer references periods(id) NOT NULL,
    aparta_id   integer references apartas(id) NOT NULL,
    tipo_id     integer references tipos(id) NOT NULL,
    fecha       date not null,   
    amount      bigint,
    created_at  timestamp not null,   
    updated_at  timestamp not null   
);

CREATE INDEX  index_cuotas_on_aparta_id  ON cuotas(aparta_id);
CREATE INDEX  index_cuotas_on_tipo_id  ON cuotas(tipo_id);
CREATE INDEX  index_cuotas_on_period_id  ON cuotas(period_id);


CREATE TABLE IF NOT EXISTS ingresos (
    id          serial primary key,
    period_id   integer references periods(id) NOT NULL,
    tipo_id     integer references tipos(id) NOT NULL,
    fecha       date not null,   
    amount      bigint,
    description  varchar(48),
    created_at  timestamp not null,   
    updated_at  timestamp not null   
);


CREATE INDEX  index_ingresos_on_tipo_id  ON ingresos(tipo_id);
CREATE INDEX  index_ingresos_on_period_id  ON ingresos(period_id);

CREATE TABLE IF NOT EXISTS egresos (
    id          serial primary key,
    period_id   integer references periods(id) NOT NULL,
    tipo_id     integer references tipos(id) NOT NULL,
    fecha       date not null,   
    amount      bigint,
    description  varchar(48),
    created_at  timestamp not null,   
    updated_at  timestamp not null   
);


CREATE INDEX  index_egresos_on_tipo_id  ON egresos(tipo_id);
CREATE INDEX  index_egresos_on_period_id  ON egresos(period_id);


CREATE TABLE persons (
    id         serial primary key,
    aparta_id  integer references  apartas(id),
    fname      varchar(32),
    lname      varchar(32),
    fecNac     date,
    email      varchar(32),
    address    varchar(64),
    tele       varchar(16),
    mobil      varchar(16),
    tipo       character(1), -- D owner, I inquilino, A admin, R referenciante 
    photo      text,
    created_at timestamp not null,   
    updated_at timestamp not null   
);

CREATE INDEX  index_persons_on_aparta_id  ON persons(aparta_id);

CREATE TABLE users (
    id         serial primary key,
    person_id  integer references persons(id),
    uuid       varchar(64) not null unique,
    cuenta     character(16),
    password   character(64),
    nivel      integer,
    created_at timestamp not null,   
    updated_at timestamp not null   
);

CREATE UNIQUE INDEX index_users_on_person_id ON users (person_id);


create TABLE IF NOT EXISTS shadows (
    id         serial primary key,
    user_id    integer references users(id) NOT NULL,
    uuid       varchar(64) not null unique,
    password   varchar(64),
    created_at timestamp not null,   
    updated_at timestamp not null   
);

CREATE UNIQUE INDEX index_shadows_on_user_id ON shadows (user_id);

create table IF NOT EXISTS sessions (
  id         serial primary key,
  user_id    integer references users(id) NOT NULL,
  uuid       varchar(64) not null unique,
  created_at timestamp not null,   
  updated_at timestamp not null   
);

CREATE UNIQUE INDEX  index_sessions_on_user_id ON  sessions (user_id);

create table IF NOT EXISTS threads (
  id         serial primary key,
  user_id    integer references users(id) NOT NULL,
  uuid       varchar(64) not null unique,
  created_at timestamp not null,       
  updated_at  timestamp not null   
);

CREATE UNIQUE INDEX  index_threads_on_user_id  ON threads (user_id);

