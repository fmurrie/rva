create table rvaAccountRole
(
	idRvaAccount int not null,
	idRvaAccountGroup int not null,
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
	constraint pk_rvaAccountRole_idRvaAccount_idRvaAccountGroup primary key(idRvaAccount,idRvaAccountGroup),
    constraint fk_rvaAccountRole_idRvaAccount foreign key(idRvaAccount) references rvaAccount(idRvaAccount) on update cascade,
    constraint fk_rvaAccountRole_idRvaAccountGroup foreign key(idRvaAccountGroup) references rvaAccountGroup(idRvaAccountGroup) on update cascade
);