CLI commands:

## Namespace
namespace ls
namespace create [NAME] --decs <string[optional]> --memory <int[optional]> --route <int[optional]>
namespace update [NAME] --decs <string[optional]> --memory <int[optional]> --route <int[optional]>
namespace remove [NAME]
namespace inspect [NAME]
namespace watch [NAME]

## Service
service ls [NAMESPACE] 
service create [NAMESPACE] [IMAGE] --name <string[optional]> --desc <string[optional]> --replicas <int[optional]> --memory <int[optional]>
service update [NAMESPACE] [NAME] --desc <string[optional]> --replicas <int[optional]> --memory <int[optional]>
service remove [NAMESPACE] [NAME] 
service inspect [NAMESPACE] [NAME]  
service watch [NAMESPACE] [NAME] 
service logs [NAMESPACE] [NAME] --pod=<string> --container=<string> 

## Route
route ls [NAMESPACE] 
route create [NAMESPACE] [ENDPOINT] [PORT] 
route update [NAMESPACE] [NAME] [ENDPOINT] [PORT] 
route remove [NAMESPACE] [NAME] 
route inspect [NAMESPACE] [NAME]

## Secret
secret ls [NAMESPACE] 
secret create [NAMESPACE] [NAME] [DATA] 
secret update [NAMESPACE] [NAME] [DATA] 
secret remove [NAMESPACE] [NAME]  
