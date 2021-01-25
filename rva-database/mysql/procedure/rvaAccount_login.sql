create procedure rvaAccount_login
(
	in accountName varchar(50),
    in email varchar(320),
	in accountPassword varchar(1000)	
)
begin

select 
rvaAccount.idRvaAccount,
rvaAccount.accountName,
rvaAccount.firstName,
rvaAccount.lastName,
rvaAccount.email,
rvaAccount.phoneNumber
from rvaAccount 
where 
(rvaAccount.accountName=accountName or rvaAccount.email=email)
and
rvaAccount.accountPassword=sha1(accountPassword)
and
rvaAccount.logicDelete=false;

end