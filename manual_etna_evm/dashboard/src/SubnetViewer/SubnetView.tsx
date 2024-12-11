import { useState, useEffect } from "react";
import { useSubnetViewerStore } from "./subnetViewStore";
import { info, pvm } from "@avalabs/avalanchejs";

import fastJsonStringify from "fast-json-stringify"

export default function SubnetView() {
    const store = useSubnetViewerStore();

    return (
        <div>
            <SelectSubnet />
            <APIInfoCard
                subnetId={store.subnetId}
                dataFetcher={async (subnetId) => {
                    const pvmapi = new pvm.PVMApi("https://api.avax-test.network");
                    return pvmapi.getSubnet({ subnetID: subnetId });
                }}
                title="GetSubnet"
            />
            <APIInfoCard
                subnetId={store.subnetId}
                dataFetcher={async (subnetId) => {
                    const pvmapi = new pvm.PVMApi("https://api.avax-test.network");
                    return pvmapi.getValidatorsAt({ subnetID: subnetId, height: "proposed" });
                }}
                title="GetValidatorsAt('proposed')"
            />
            <APIInfoCard
                subnetId={store.subnetId}
                dataFetcher={async (subnetId) => {
                    const pvmapi = new pvm.PVMApi("https://api.avax-test.network");
                    const subnet = await pvmapi.getSubnet({ subnetID: subnetId });
                    return pvmapi.getBalance({ addresses: subnet.controlKeys });
                }}
                title="GetBalance(subnet.controlKeys)"
            />
            <APIInfoCard
                subnetId={store.subnetId}
                dataFetcher={async (subnetId) => {
                    const pvmapi = new pvm.PVMApi("http://localhost:9650");
                    const subnet = await pvmapi.getSubnet({ subnetID: subnetId });
                    return pvmapi.getBlockchainStatus(subnet.managerChainID);
                }}
                title="GetBlockchainStatus(subnet.managerChainID)"
            />
            <APIInfoCard
                subnetId={store.subnetId}
                dataFetcher={async (subnetId) => {
                    const infoApi = new info.InfoApi("http://localhost:9650");
                    const nodeId = await infoApi.getNodeId();
                    const pvmapi = new pvm.PVMApi("http://localhost:9650");
                    return pvmapi.getCurrentValidators({ nodeIDs: [nodeId.nodeID] });
                }}
                title="GetCurrentValidators(nodeIDs: [nodeId])"
            />
        </div>
    );
}

function APIInfoCard({
    subnetId, dataFetcher, title
}: {
    subnetId: string,
    dataFetcher: (subnetId: string) => Promise<any>,
    title: string
}) {
    const [data, setData] = useState<any>(null);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        dataFetcher(subnetId)
            .then((data) => {
                setData(data);
                setError(null);
            })
            .catch((err) => {
                setError(String(err.message));
                setData(null);
            });
    }, [subnetId]);

    return (
        <div className="bg-white rounded-lg shadow-md p-6 m-4">
            <h1 className="text-xl font-bold mb-4">{title}</h1>
            <JSONOrError data={data} error={error} />
        </div>
    );
}

function JSONOrError({ data, error }: { data: any; error: string | null }) {
    if (error) {
        return <div className="bg-red-100 text-red-700 p-4 rounded overflow-auto">{error}</div>;
    }
    return <pre className="bg-gray-100 p-4 rounded overflow-auto">{JSON.stringify(data, (_, v) => typeof v === 'bigint' ? v.toString() : v, 2)}</pre>;
}

function SelectSubnet() {
    const store = useSubnetViewerStore();
    const [inputValue, setInputValue] = useState("");

    useEffect(() => {
        if (store) {
            setInputValue(store.subnetId);
        }
    }, [store?.subnetId]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (store) {
            store.setSubnetId(inputValue);
        }
    };

    return (
        <div className="bg-white rounded-lg shadow-md p-6 m-4">
            <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                <div className="flex items-center gap-2">
                    <label htmlFor="subnetId" className="font-medium">Subnet ID:</label>
                    <input
                        id="subnetId"
                        type="text"
                        value={inputValue}
                        onChange={(e) => setInputValue(e.target.value)}
                        className="border rounded px-2 py-1 flex-1"
                    />
                </div>
                <button
                    type="submit"
                    className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition-colors"
                >
                    Save
                </button>
            </form>
        </div>
    );
}
