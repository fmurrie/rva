create table rvaEntity
(
    idRvaEntity int auto_increment,
    idRvaModule int not null,
    entityName varchar(100) not null,
    description varchar(1000),
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
    constraint pk_rvaEntity_idValueEntity primary key(idRvaEntity),
    constraint uk_rvaEntity_idRvaModule_entityName unique key(idRvaModule,entityName),
    constraint fk_rvaEntity_idRvaModule foreign key(idRvaModule) references rvaModule(idRvaModule) on update cascade
);