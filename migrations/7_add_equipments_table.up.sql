CREATE TABLE equipements (
    id INT NOT NULL auto_increment primary key,
    name VARCHAR(255) NOT NULL,
    -- permission INT NOT NULL
    files INT -- one to many
);
