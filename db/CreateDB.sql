Create Database if not exists account;

use account;

Create Table if not exists User (
    UId int not null auto_increment primary key,
    Name varchar(20) not null,
    Account varchar(10) not null unique,
    Password varchar(64) not null unique,
    CreateTime timestamp default current_timestamp
);

Create Table if not exists Action (
    ActionType int not null,
    DetailType int not null,
    DetailName varchar(20) not null unique,
    PRIMARY KEY (ActionType, DetailType)
);


Create Table if not exists Log (
    UserId int not null,
    ActionType int not null,
    DetailType int not null,
    Money int not null,
    Description varchar(20),
    CreateTime timestamp default current_timestamp
);

-- Create Table if not exists ActionType (
--     Type int not null primary key,
--     TypeName varchar(20) not null unique
-- )

-- Create Table if not exists OutlayType (
--     Type int not null auto_increment primary key,
--     TypeName varchar(20) not null unique
-- );

-- Create Table if not exists IncomeType (
--     Type int not null auto_increment primary key,
--     TypeName varchar(20) not null unique
-- );