create procedure rvaEndpoint_getAll
(
)
begin

select 
	rvaEndpoint.idRvaEndpoint,
	rvaEndpoint.path,
	rvaType.typeName as httpVerb
from rvaEndpoint
inner join rvaType 
	on rvaEndpoint.idHTTPVerb=rvaType.idRvaType
where 
rvaEndpoint.logicDelete=false
and
rvaType.logicDelete=false;

end