create table rvaProcedure
(
    idRvaProcedure int auto_increment,
    procedureName varchar(64) not null,
	procedureQuery varchar(1000) not null,
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
    constraint pk_rvaProcedure_idRvaProcedure primary key(idRvaProcedure),
    constraint uk_rvaProcedure_procedureName unique key(procedureName)
);