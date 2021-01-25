create table rvaEndpoint
(
	idRvaEndpoint int auto_increment,
	idRvaOrganizedProcedure int not null,
    idHTTPVerb int not null,
    path varchar(500) not null,
    createAuth boolean default(false),
    validAuth boolean default(false),
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
    constraint pk_rvaEndpoint_idRvaEndpoint primary key(idRvaEndpoint),
    constraint uk_rvaEndpoint_idHTTPVerb_path unique key(idHTTPVerb,path),
    constraint fk_rvaEndpoint_idRvaOrganizedProcedure foreign key(idRvaOrganizedProcedure) references rvaOrganizedProcedure(idRvaOrganizedProcedure) on update cascade,
    constraint fk_rvaEndpoint_idHTTPVerb foreign key(idHTTPVerb) references rvaType(idRvaType) on update cascade
);