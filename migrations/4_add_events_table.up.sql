CREATE TABLE events (
    id INT NOT NULL auto_increment primary key,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    owner INT NOT NULL, -- user_id
    image INT -- file_id
);
