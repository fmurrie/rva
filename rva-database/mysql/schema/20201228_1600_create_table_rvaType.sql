create table rvaType
(
	idRvaType int auto_increment,
    idRvaEntity int not null,
    typeName varchar(100) not null,
    description varchar(1000),
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
    constraint pk_rvaType_idRvaType primary key(idRvaType),
    constraint uk_rvaType_idRvaEntity_typeName unique key(idRvaEntity,typeName),
    constraint fk_rvaType_idRvaEntity foreign key(idRvaEntity) references rvaEntity(idRvaEntity) on update cascade
);