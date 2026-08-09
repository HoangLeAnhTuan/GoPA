export interface Vehicle { id:string; make:string; model:string; year:number; nickname:string|null; current_mileage:number; }
export interface VehicleInput { make:string; model:string; year:number; nickname:string|null; currentMileage:number; }
export interface VehicleLog { id:string; log_type:"MAINTENANCE"|"FUEL"|"UPGRADE"; mileage:number; description:string; cost:string; log_date:string; }
export interface NetworkNode { id:string; device_name:string; mac_address:string|null; static_ip:string|null; vlan_tag:number|null; connection_type:"LAN"|"WIFI"; parent_node_id:string|null; location:string|null; notes:string|null; }
export interface NetworkNodeInput { deviceName:string; macAddress:string|null; staticIP:string|null; vlan:number|null; connectionType:"LAN"|"WIFI"; parentNodeID:string|null; location:string|null; notes:string|null; }
