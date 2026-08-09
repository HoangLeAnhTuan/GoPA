import { apiClient } from "../../../lib/api-client";
import { unwrap, unwrapCollection, type ApiEnvelope } from "../../../lib/api-envelope";
import type { NetworkNode, NetworkNodeInput, Vehicle, VehicleInput, VehicleLog } from "../types";
export async function listVehicles(){return unwrapCollection(await apiClient.get<ApiEnvelope<Vehicle[] | null>>("/vehicles"));}
export async function createVehicle(input:VehicleInput){return unwrap(await apiClient.post<ApiEnvelope<Vehicle>>("/vehicles",{make:input.make,model:input.model,year:input.year,nickname:input.nickname,current_mileage:input.currentMileage}));}
export async function listVehicleLogs(id:string){return unwrapCollection(await apiClient.get<ApiEnvelope<VehicleLog[] | null>>(`/vehicles/${id}/logs`));}
export async function createVehicleLog(id:string,input:{logType:"MAINTENANCE"|"FUEL"|"UPGRADE";mileage:number;description:string;cost:string;logDate:string}){return unwrap(await apiClient.post<ApiEnvelope<VehicleLog>>(`/vehicles/${id}/logs`,{log_type:input.logType,mileage:input.mileage,description:input.description,cost:input.cost,log_date:`${input.logDate}T00:00:00Z`}));}
export async function maintenance(id:string){return unwrap(await apiClient.get<ApiEnvelope<{status:string;remaining_distance:number;next_mileage:number}>>(`/vehicles/${id}/maintenance-status`));}
export async function listNodes(){return unwrapCollection(await apiClient.get<ApiEnvelope<NetworkNode[] | null>>("/network-nodes"));}
export async function createNode(input:NetworkNodeInput){return unwrap(await apiClient.post<ApiEnvelope<NetworkNode>>("/network-nodes",{device_name:input.deviceName,mac_address:input.macAddress,static_ip:input.staticIP,vlan_tag:input.vlan,connection_type:input.connectionType,parent_node_id:input.parentNodeID,location:input.location,notes:input.notes}));}
