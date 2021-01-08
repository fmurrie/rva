create table rvaEndpointStep
(
	idRvaEndpointStep int auto_increment,
	idRvaEndpoint int not null,
    idRvaProcedure int not null,
    stepOrder int not null,
    description varchar(1000),
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
    constraint pk_rvaEndpointStep_idRvaEndpointStep primary key(idRvaEndpointStep),
    constraint uk_rvaEndpointStep_idRvaEndpoint_stepOrder unique key(idRvaEndpoint,stepOrder),
	constraint fk_rvaEndpointStep_idRvaEndpoint foreign key(idRvaEndpoint) references rvaEndPoint(idRvaEndpoint) on update cascade,
    constraint fk_rvaEndpointStep_idRvaProcedure foreign key(idRvaProcedure) references rvaProcedure(idRvaProcedure) on update cascade
);