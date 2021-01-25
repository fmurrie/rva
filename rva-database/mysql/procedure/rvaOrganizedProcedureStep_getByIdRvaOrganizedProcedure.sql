create procedure rvaOrganizedProcedureStep_getByIdRvaOrganizedProcedure
(
in idRvaOrganizedProcedure int
)
begin

select 
	rvaOrganizedProcedureStep.idRvaOrganizedProcedureStep,
	rvaProcedure.procedureQuery
from rvaOrganizedProcedureStep 
inner join rvaProcedure
	on rvaProcedure.idRvaProcedure=rvaOrganizedProcedureStep.idRvaProcedure
where 
rvaOrganizedProcedureStep.idRvaOrganizedProcedure=idRvaOrganizedProcedure 
and 
rvaOrganizedProcedureStep.logicDelete=false 
and 
rvaProcedure.logicDelete=false 
order by rvaOrganizedProcedureStep.stepOrder;

end