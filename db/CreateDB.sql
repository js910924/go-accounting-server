Create Database Account;

Create Table User (
    UId int not null auto_increment primary key,
    Name varchar(20) not null,
    Account varchar(10) not null unique,
    Password varchar(10) not null unique,
    CreateTime timestamp default current_timestamp
)
