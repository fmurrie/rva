create procedure rvaEndpointStep_getByIdRvaEndpoint
(
in idRvaEndpoint int
)
begin

select 
	rvaEndpointStep.idRvaEndpointStep,
	rvaProcedure.procedureQuery
from rvaEndpointStep 
inner join rvaProcedure
	on rvaProcedure.idRvaProcedure=rvaEndpointStep.idRvaProcedure
where 
rvaEndpointStep.idRvaEndpoint=idRvaEndpoint 
and 
rvaEndpointStep.logicDelete=false 
and 
rvaProcedure.logicDelete=false 
order by rvaEndpointStep.stepOrder;

end