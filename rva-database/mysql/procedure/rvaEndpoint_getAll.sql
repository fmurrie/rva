create procedure rvaEndpoint_getAll
(
)
begin

select 
    rvaEndpoint.idRvaOrganizedProcedure,
	rvaEndpoint.path,
	rvaType.typeName as httpVerb,
    rvaEndpoint.createAuth,
    rvaEndpoint.validAuth
from rvaEndpoint
inner join rvaType 
	on rvaEndpoint.idHTTPVerb=rvaType.idRvaType
where 
rvaEndpoint.logicDelete=false
and
rvaType.logicDelete=false;

end