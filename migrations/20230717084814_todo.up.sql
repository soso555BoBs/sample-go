CREATE TABLE todo (
    id int NOT NULL AUTO_INCREMENT,
    content varchar(255) NOT NULL,
    created_at datetime default current_timestamp,
    PRIMARY KEY (id)
)