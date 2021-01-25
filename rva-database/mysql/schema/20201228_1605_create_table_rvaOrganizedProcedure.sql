create table rvaOrganizedProcedure
(
	idRvaOrganizedProcedure int auto_increment,
    procedureName varchar(250) not null,
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
    constraint pk_rvaOrganizedProcedure_idRvaOrganizedProcedure primary key(idRvaOrganizedProcedure),
    constraint uk_rvaOrganizedProcedure_procedureName unique key(procedureName)
);