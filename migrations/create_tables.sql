CREATE TABLE user (
    id int  not null auto_increment primary key,
    login VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    promotion smallint NOT NULL
);

CREATE TABLE account (
    id int not null auto_increment primary key,
    user_id int not null, 
    foreign key (user_id) references user(id)
);