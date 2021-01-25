insert into rvaEntity
(
	idRvaModule,
	entityName,
	description,
	creatorAccount,
	updaterAccount
)
values
(
	(select idRvaModule from rvaModule where moduleName='RVA'),
	'HTTPVerb',
	'HTTPVerb Entity',
	'System',
	'System'
);