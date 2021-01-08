create table rvaModule
(
    idRvaModule varchar(10) not null,
    moduleName varchar(100) not null,
    description varchar(1000),
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
    constraint pk_rvaModule_idRvaModule primary key(idRvaModule),
    constraint uk_rvaModule_moduleName unique key(moduleName)
);