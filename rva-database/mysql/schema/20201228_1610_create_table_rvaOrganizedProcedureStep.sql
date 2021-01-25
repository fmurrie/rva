create table rvaOrganizedProcedureStep
(
	idRvaOrganizedProcedureStep int auto_increment,
	idRvaOrganizedProcedure int not null,
    idRvaProcedure int not null,
    stepOrder int not null,
	creatorAccount varchar(100) not null,
    updaterAccount varchar(100) not null,
	createdDate datetime default(now()),
	updatedDate datetime default(now()),
    logicDelete boolean default(false),
    constraint pk_rvaOrganizedProcedureStep_idRvaOrganizedProcedureStep primary key(idRvaOrganizedProcedureStep),
    constraint uk_rvaOrganizedProcedureStep_idRvaOrganizedProcedure_stepOrder unique key(idRvaOrganizedProcedure,stepOrder),
	constraint fk_rvaOrganizedProcedureStep_idRvaOrganizedProcedure foreign key(idRvaOrganizedProcedure) references rvaOrganizedProcedure(idRvaOrganizedProcedure) on update cascade,
    constraint fk_rvaOrganizedProcedureStep_idRvaProcedure foreign key(idRvaProcedure) references rvaProcedure(idRvaProcedure) on update cascade
);