create table rvaAccount
(
	idRvaAccount int auto_increment,
    accountName varchar(50) not null,
	accountPassword varchar(1000) not null,
	firstName varchar(100) not null,
    lastName varchar(100) not null,
    email varchar(320) not null,
    phoneNumber varchar(320),
    ipAddress varchar(100) not null,
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
	constraint pk_rvaAccount_idRvaAccount primary key(idRvaAccount),
    constraint uk_rvaAccount_accountName unique key(accountName),
    constraint uk_rvaAccount_email unique key(email)
);