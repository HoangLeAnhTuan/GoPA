import { CarFront, Network, Plus } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import { GlassPanel } from "../../../components/design-system/glass-panel";
import { PageHeader } from "../../../components/design-system/page-header";
import { useCreateNode, useCreateVehicle, useNetworkNodes, useVehicles } from "../hooks/use-assets";

export function AssetsPage() {
	const { t } = useTranslation("assets");
	const vehicles = useVehicles();
	const nodes = useNetworkNodes();
	const createVehicle = useCreateVehicle();
	const createNode = useCreateNode();
	const [make, setMake] = useState("Mazda");
	const [model, setModel] = useState("CX-5");
	const [node, setNode] = useState("");
	if (vehicles.isPending || nodes.isPending) return <LoadingState />;
	if (vehicles.isError || nodes.isError) return <ErrorState />;
	return <div><PageHeader description={t("description")} title={t("title")} /><section className="mt-6 grid gap-5 lg:grid-cols-2"><GlassPanel className="p-5"><div className="flex items-center gap-2"><CarFront className="text-primary" /><h2 className="font-semibold">{t("vehicles.title")}</h2></div><form className="mt-4 flex flex-wrap gap-2" onSubmit={(event) => { event.preventDefault(); createVehicle.mutate({ make, model, year: 2024, nickname: null, currentMileage: 0 }); }}><input className="h-10 min-w-24 flex-1 rounded-xl border bg-white/70 px-3 dark:bg-slate-900/70" onChange={(event) => setMake(event.target.value)} value={make} /><input className="h-10 min-w-24 flex-1 rounded-xl border bg-white/70 px-3 dark:bg-slate-900/70" onChange={(event) => setModel(event.target.value)} value={model} /><button className="inline-flex h-10 items-center gap-1 rounded-xl bg-primary px-3 text-sm font-semibold text-primary-foreground" type="submit"><Plus size={16} />{t("add")}</button></form><div className="mt-5 space-y-3">{vehicles.data.map((vehicle) => <article className="rounded-xl border bg-white/70 p-3 dark:bg-slate-900/70" key={vehicle.id}><p className="font-medium">{vehicle.nickname ?? `${vehicle.make} ${vehicle.model}`}</p><p className="text-sm text-muted-foreground">{vehicle.year} · {vehicle.current_mileage.toLocaleString()} km</p></article>)}{vehicles.data.length === 0 ? <p className="text-sm text-muted-foreground">{t("vehicles.empty")}</p> : null}</div></GlassPanel><GlassPanel className="p-5"><div className="flex items-center gap-2"><Network className="text-primary" /><h2 className="font-semibold">{t("network.title")}</h2></div><form className="mt-4 flex gap-2" onSubmit={(event) => { event.preventDefault(); if (node.trim() !== "") createNode.mutate({ deviceName: node, macAddress: null, staticIP: null, vlan: null, connectionType: "WIFI", parentNodeID: null, location: null, notes: null }); setNode(""); }}><input className="h-10 flex-1 rounded-xl border bg-white/70 px-3 dark:bg-slate-900/70" onChange={(event) => setNode(event.target.value)} placeholder={t("network.placeholder")} value={node} /><button className="inline-flex h-10 items-center gap-1 rounded-xl bg-primary px-3 text-sm font-semibold text-primary-foreground" type="submit"><Plus size={16} />{t("add")}</button></form><div className="mt-5 space-y-3">{nodes.data.map((networkNode) => <article className="rounded-xl border bg-white/70 p-3 dark:bg-slate-900/70" key={networkNode.id}><p className="font-medium">{networkNode.device_name}</p><p className="text-sm text-muted-foreground">{networkNode.connection_type} · {networkNode.static_ip ?? t("network.noIP")}</p></article>)}{nodes.data.length === 0 ? <p className="text-sm text-muted-foreground">{t("network.empty")}</p> : null}</div></GlassPanel></section></div>;
}
