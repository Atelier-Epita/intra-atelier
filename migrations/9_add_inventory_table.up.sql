CREATE TABLE inventory (
    id INT NOT NULL auto_increment primary key,
    group_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount INT NOT NULL
);
