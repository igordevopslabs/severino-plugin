# severino-plugin
Um plugin simples de autorização e autenticação JWT para Kong API Gateway


## Ciclo de vida de um plugin para o Kong Api Gateway
* O Kong intercepta as requisições antes de chegar ao seu upstream (serviço de destino).
* Em cada fase (por exemplo, access, header_filter, response, etc.), o Kong chama o método correspondente do seu plugin para executar alguma lógica.
* A validação JWT, será realizada na fase de Access, pois será necessário inspecionar a request para verificar a existência de um header específco.

