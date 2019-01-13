DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS blog_article;
DROP TABLE IF EXISTS blog_auth;
create table blog_article (
    id serial primary key,
    tag_id integer DEFAULT 0,
    title varchar(100) DEFAULT '',
    descri varchar(255) DEFAULT '',
    content text,
    cover_image_url varchar(255) DEFAULT '',
    created_on integer DEFAULT 0,
    created_by integer DEFAULT 0,
    modified_on integer DEFAULT 0,
    modified_by varchar(255) DEFAULT '',
    deleted_on integer DEFAULT 0,
    stateCode smallint DEFAULT 1
);

create table blog_auth (
    id serial primary key,
    username varchar(50) DEFAULT '',
    password varchar(50) DEFAULT ''
);

insert into blog_auth(id, username, password) values (1, 'vincentlin113', 'test123456')

