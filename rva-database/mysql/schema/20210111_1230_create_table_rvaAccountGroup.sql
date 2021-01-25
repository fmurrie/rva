create table rvaAccountGroup
(
	idRvaAccountGroup int auto_increment,
    idRvaModule int not null,
    groupName varchar(100) not null,
    description varchar(1000),
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
	constraint pk_rvaAccountGroup_idRvaAccountGroup primary key(idRvaAccountGroup),
    constraint uk_rvaAccountGroup_idRvaModule_groupName unique key(idRvaModule,groupName),
    constraint fk_rvaAccountGroup_idRvaModule foreign key(idRvaModule) references rvaModule(idRvaModule) on update cascade
);