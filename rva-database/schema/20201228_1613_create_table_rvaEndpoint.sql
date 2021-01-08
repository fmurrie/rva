create table rvaEndpoint
(
	idRvaEndpoint int auto_increment,
    idHTTPVerb int not null,
    path varchar(500) not null,
    description varchar(1000),
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
    constraint pk_rvaEndpoint_idRvaEndpoint primary key(idRvaEndpoint),
    constraint uk_rvaEndpoint_idHTTPVerb_path unique key(idHTTPVerb,path),
    constraint fk_rvaEndpoint_idHTTPVerb foreign key(idHTTPVerb) references rvaType(idRvaType) on update cascade
);