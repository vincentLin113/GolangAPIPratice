DROP TABLE IF EXISTS blog_tag;
create table blog_tag (
    id serial primary key,
    name varchar(100) DEFAULT '', 
    created_on integer DEFAULT '0',
    created_by varchar(100) DEFAULT '',
    modified_on integer DEFAULT '',
    modified_by varchar(100) DEFAULT '',
    state smallint DEFAULT '1',
    deleted_on integer DEFAULT '0'
);