CREATE TABLE participants (
    -- many to one join table
    -- one user have many events
    id INT NOT NULL auto_increment primary key,
    event_id INT NOT NULL,
    user_id INT NOT NULL
);
