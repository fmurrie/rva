create procedure rvaAccount_insert
(
	in accountName varchar(50),
	in accountPassword varchar(1000),
	in firstName varchar(100),
	in lastName varchar(100),
	in email varchar(320),
	in phoneNumber varchar(320),
	in ipAddress varchar(100)
)
begin

insert into rvaAccount
(
	rvaAccount.accountName,
	rvaAccount.accountPassword,
	rvaAccount.firstName,
	rvaAccount.lastName,
	rvaAccount.email,
	rvaAccount.phoneNumber,
	rvaAccount.ipAddress,
	rvaAccount.creatorAccount,
    rvaAccount.updaterAccount
)
values
(
	accountName,
	sha1(accountPassword),
	firstName,
	lastName,
	email,
	phoneNumber,
	ipAddress,
    'System',
    'System'
);

end