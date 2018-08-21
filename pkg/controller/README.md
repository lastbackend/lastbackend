# Last.Backend OpenSource

## API server

## Controller

### Meta:
#### TODO Describe runtime workflow

#### Workflow:
Controller workflow:

#### Cluster observer

- [ ] TODO: Cluster observer runtime:
  [ ] TODO: Cluster observer struct with methods:
      - [ ] TODO: Restore cluster state:
            - TODO: [x] start cluster IPAM
            - TODO: [x] get cluster information
            - TODO: [x] get all nodes in cluster
            - TODO: [ ] for each node - run node observer
            - TODO: [ ] sync current cluster state
            - TODO: [ ] get network information
            - TODO: [ ] get all endpoints
            - TODO: [ ] get all cluster manifests
      - [ ] Runtime:
            - [ ]  Watch node stats changes
            - [ ]  Watch requests for node lease
            - [ ]  Watch requests for node release
      - [ ] Lease node:
            - [ ]  Find node matches scheduling type
            - [ ]  Calculate new cluster state
            - [ ]  Return node
      - [ ] Release node:
            - [ ]  Release node state
            - [ ]  Calculate new cluster state
      - [ ] Endpoint:
            - [ ]  Provision: check endpoint spec and network spec
            - [ ]  create new endpoint and save it to state storage
            - [ ]  update endpoint spec if previous condition is true
            - [ ]  remove endpoint if network spec is empty


 [ ] Node observer:
  [ ] Node observer struct
  [ ] Restore: restore current node state:
        - [ ]  get all manifests for node
        - [ ]  calculate memory usage
        - [ ]  calculate storage usage
        - [ ]  check network is ready
        - [ ]  calculate node global state:
        Node state conditions:
          - CNI state: [Ready, Error] - [ Warning ] state
          - CPI state: [Ready, Error] - [ Warning ] state
          - CRI state: [Ready, Error] - [ Not Ready ] state
          - CSI state: [Ready, Error] - [ Warning ] state
  [ ] Update: update current node state
  [ ] Remove: clean up node and remove
  [ ] Manifest: add manifest to node
  [ ] Manifest: set manifest to node
  [ ] Manifest: del manifest from node

[ ] Namespaces observer

[ ] Namespace observer
  [ ] Namespace observer struct
  [ ] Restore namespaces observer:
     - [ ] Get all namespaces
     - [ ] Loop over namespaces and get services for each
     - [ ] Loop over services and create service controller
     - [ ] Start service controller loop
  [ ] Create watcher for service
     - [ ] if service state is initialized -> pass service for update
     - [ ] if service is not initialized -> initialize
  [ ] Create watcher for pods
     - [ ] if service for this pod is initialized > pass pod for service update
     - [ ] if service is not initialized > skip
  [ ] Runtime loop:
     - [ ] Start service watching changes

 [ ] Services observer
     Service observer should watch service state and manage service runtime
  Restore: service controller state
     - [x] get all pods
     - [x] get all deployments
     - [x] get endpoint for service
     - [x] loop over pods and pass pod to service update
     - [x] loop over deployment and pass deployment to service update
     - [x] check current service spec and provision service if needed
  Handle changes:
  Runtime loop:
  Scale:
  Update:
  Remove:


 [ ] Pod package:
 Pod status states:
 created > provision > ready > destroy > destroyed
                     > error 

  [x] Create: Create new pod
  [x] Destroy: Mark Pod for destroy
  [x] Provision: Provision pod
  [x] Remove: Remove pod if previous state was destroy

 [x] Deployment package:
 Deployment status states:
 crated > provision > ready > destroy > destroyed
                    > error

  [x] Create: Create new deployment based on service spec
  [x] Scale: Scale deployment for replicas, it should work as increment pods and decrement
  [x] Destroy: Mark all pods for destroy
  [x] Remove: Remove deployment if service is marked as removed
  [x] Cancel: Cancel deployment if not need anymore

 [x] Service package:
  [x] Remove: Remove service
  [x] Sync: Service sync state
  [x] Provision: Provision service
  [x] Destroy: destroy service

 [x] Endpoint package:
  [x] Create: create new endpoint
  [x] Update: update endpoint spec
  [x] Destroy: destroy endpoint

