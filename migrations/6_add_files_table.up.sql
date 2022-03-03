CREATE TABLE files (
    id INT NOT NULL auto_increment primary key,
    type ENUM('public', 'private') NOT NULL,
    owner_id INT NOT NULL,
    group_id INT NOT NULL,

    -- file path is ./files/$(file_hash)$(file_name)
    file_name VARCHAR(255) NOT NULL,
    file_hash BINARY(32) NOT NULL -- sha256 hash is 256 bits = 32 bytes
);
