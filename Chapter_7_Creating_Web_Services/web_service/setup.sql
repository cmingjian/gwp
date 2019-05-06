drop table posts;

create table posts (
  id      serial primary key auto_increment,
  content text,
  author  varchar(255)
);