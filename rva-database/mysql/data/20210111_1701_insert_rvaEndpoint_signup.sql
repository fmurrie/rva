insert into rvaOrganizedProcedure
(procedureName,creatorAccount,updaterAccount)
values
('signup','System','System');

insert into rvaOrganizedProcedureStep
(
	idRvaOrganizedProcedure,
	idRvaProcedure,
	stepOrder,
	creatorAccount,
	updaterAccount
)
values
(
	(select idRvaOrganizedProcedure from rvaOrganizedProcedure where procedureName='signup'),
    (select idRvaProcedure from rvaProcedure where procedureName='rvaAccount_insert'),
    1,
    'System',
    'System'
);

insert into rvaEndpoint
(
	idRvaOrganizedProcedure,
	idHTTPVerb,
	path,
	creatorAccount,
	updaterAccount
)
values
(
	(select idRvaOrganizedProcedure from rvaOrganizedProcedure where procedureName='signup'),
	(select idRvaType from rvaType where typeName='POST' and idRvaEntity=(select idRvaEntity from rvaEntity where entityName='HTTPVerb' and idRvaModule=(select idRvaModule from rvaModule where moduleName='RVA'))),
	'/rva/signup',
	'System',
	'System'
);