insert into rvaOrganizedProcedure
(procedureName,creatorAccount,updaterAccount)
values
('signin','System','System');

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
	(select idRvaOrganizedProcedure from rvaOrganizedProcedure where procedureName='signin'),
    (select idRvaProcedure from rvaProcedure where procedureName='rvaAccount_login'),
    1,
    'System',
    'System'
);

insert into rvaEndpoint
(
	idRvaOrganizedProcedure,
	idHTTPVerb,
	path,
    createAuth,
	creatorAccount,
	updaterAccount
)
values
(
	(select idRvaOrganizedProcedure from rvaOrganizedProcedure where procedureName='signin'),
	(select idRvaType from rvaType where typeName='POST' and idRvaEntity=(select idRvaEntity from rvaEntity where entityName='HTTPVerb' and idRvaModule=(select idRvaModule from rvaModule where moduleName='RVA'))),
	'/rva/signin',
    true,
	'System',
	'System'
);