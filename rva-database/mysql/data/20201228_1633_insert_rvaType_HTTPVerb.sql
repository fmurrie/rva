insert into rvaType
(idRvaEntity,typeName,description,creatorAccount,updaterAccount)
values
((select idRvaEntity from rvaEntity where idRvaModule=(select idRvaModule from rvaModule where moduleName='RVA') and entityName='HTTPVerb'),'GET','HTTPVerb GET','System','System'),
((select idRvaEntity from rvaEntity where idRvaModule=(select idRvaModule from rvaModule where moduleName='RVA') and entityName='HTTPVerb'),'HEAD','HTTPVerb HEAD','System','System'),
((select idRvaEntity from rvaEntity where idRvaModule=(select idRvaModule from rvaModule where moduleName='RVA') and entityName='HTTPVerb'),'POST','HTTPVerb POST','System','System'),
((select idRvaEntity from rvaEntity where idRvaModule=(select idRvaModule from rvaModule where moduleName='RVA') and entityName='HTTPVerb'),'PUT','HTTPVerb PUT','System','System'),
((select idRvaEntity from rvaEntity where idRvaModule=(select idRvaModule from rvaModule where moduleName='RVA') and entityName='HTTPVerb'),'DELETE','HTTPVerb DELETE','System','System'),
((select idRvaEntity from rvaEntity where idRvaModule=(select idRvaModule from rvaModule where moduleName='RVA') and entityName='HTTPVerb'),'CONNECT','HTTPVerb CONNECT','System','System'),
((select idRvaEntity from rvaEntity where idRvaModule=(select idRvaModule from rvaModule where moduleName='RVA') and entityName='HTTPVerb'),'OPTIONS','HTTPVerb OPTIONS','System','System'),
((select idRvaEntity from rvaEntity where idRvaModule=(select idRvaModule from rvaModule where moduleName='RVA') and entityName='HTTPVerb'),'TRACE','HTTPVerb TRACE','System','System'),
((select idRvaEntity from rvaEntity where idRvaModule=(select idRvaModule from rvaModule where moduleName='RVA') and entityName='HTTPVerb'),'PATCH','HTTPVerb PATCH','System','System');