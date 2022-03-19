CREATE TABLE users_group (
    -- many to one join table
    -- one user have many groups
    id INT NOT NULL auto_increment primary key,
    userID INT NOT NULL,
    groupID INT NOT NULL
);
