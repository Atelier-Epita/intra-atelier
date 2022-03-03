CREATE TABLE users_group (
    -- many to one join table
    -- one user have many groups
    id INT NOT NULL auto_increment primary key,
    user_id INT NOT NULL,
    group_id INT NOT NULL
);
