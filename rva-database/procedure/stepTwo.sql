create procedure stepTwo
(
in stepName varchar(100)
)
begin

select concat(stepName,'-stepTwo') as stepName;

end