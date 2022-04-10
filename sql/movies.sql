use movies;

create table movies (
    id int(10) not null auto_increment,
    title varchar(50) default null,
    area varchar(10) default null,
    which_type varchar(10) default null, 
    date varchar(10) default null,
    language varchar(10) default null,
    director varchar(25) default null,
    actor varchar(50) default null,
    image_url varchar(100) default null,
    introduction varchar(500) default null,
    rate varchar(5) default null,
    primary key(id)
) 