CLI commands:

## Namespace
ns list
ns create <name:string> --decs <string[optional]> --memory <int[optional]> --route <int[optional]>
ns update <name:string> --desc <string[optional]> --memory <int[optional]> --route <int[optional]>
ns inspect <name:string>
ns remove <name:string>

## Service
ns <name:string> service list
ns <name:string> service create <image:string> --name <string[optional]> --desc <string[optional]> --replicas <int[optional]> --memory <int[optional]>
ns <name:string> service update <name:string> --desc <string[optional]> --replicas <int[optional]> --memory <int[optional]>
ns <name:string> service inspect <name:string>  
ns <name:string> service remove <name:string> 
 
## Secret
ns <name:string> secret list 
ns <name:string> secret create <name:string> <data:string> 
ns <name:string> secret update <name:string> <data:string> 
ns <name:string> secret remove <name:string> 

## Route
ns <name:string> route list 
ns <name:string> route create <endpoint:string> <port:int> 
ns <name:string> route update <name:string> <endpoint:string> <port:int> 
ns <name:string> route inspect <name:string> 
ns <name:string> route remove <name:string> 