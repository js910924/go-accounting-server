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
    LogId int not null auto_increment primary key,
    UserId int not null,
    ActionType int not null,
    DetailType int not null,
    Money int not null,
    Description varchar(20),
    CreateTime timestamp default current_timestamp
);

INSERT INTO Action (ActionType, DetailType, DetailName)
VALUES
(1, 1, "Breakfast"),
(1, 2, "Lunch"),
(1, 3, "Dinner"),
(1, 4, "Midnight Snack"),
(1, 5, "Beverage"),
(1, 6, "Snack"),
(1, 7, "Clothes"),
(1, 8, "Transportation"),
(1, 9, "Books"),
(1, 10, "Entertainment"),
(2, 1, "Salary"),
(2, 2, "Present"),
(2, 3, "Other");

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
