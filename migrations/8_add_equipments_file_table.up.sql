CREATE TABLE equipements_file (
    -- many to one join table
    -- one equipement have many files
    id INT NOT NULL auto_increment primary key,
    equipements_id INT NOT NULL,
    file_id INT NOT NULL
);
