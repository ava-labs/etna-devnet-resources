import { useState, useEffect } from "react";
import { useSubnetViewerStore } from "./subnetViewStore";
import { info, pvm } from "@avalabs/avalanchejs";


async function JSONRPC(url: string, method: string, params: any) {
    const response = await fetch(url, {
        method: "POST",
        headers: {
            "content-type": "application/json"
        },
        body: JSON.stringify({ jsonrpc: "2.0", id: 1, method, params })
    });
    const resp = (await response.json())
    if (resp.error) {
        throw new Error(resp.error.message);
    }
    return resp.result;
}

function Header({ title }: { title: string }) {
    return (
        <div className="text-center p-8">
            <h1 className="text-3xl font-bold">{title}</h1>
        </div>
    );
}


export default function SubnetView() {
    const store = useSubnetViewerStore();

    return (
        <div className="container mx-auto px-4 py-8">
            <SelectSubnet />
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8 auto-rows-auto">
                <div className="md:col-span-2">
                    <Header title="Node Info" />
                </div>

                {["P", "X"].map(chainName => (
                    <APIInfoCard
                        key={chainName}
                        subnetId={store.subnetId}
                        dataFetcher={async (subnetId) => {
                            return JSONRPC("http://localhost:9650/ext/info", "info.isBootstrapped", { chain: chainName });
                        }}
                        title={`info.isBootstrapped('${chainName}')`}
                        note="await infoApi.isBootstrapped('P') returns peers for some reason, so we do fetch instead"
                    />
                ))}
                <APIInfoCard
                    subnetId={store.subnetId}
                    dataFetcher={async (subnetId) => {
                        return JSONRPC("http://localhost:9650/ext/info", "info.getNodeIP", {});
                    }}
                    title={`info.getNodeIP()`}
                />
                <APIInfoCard
                    subnetId={store.subnetId}
                    dataFetcher={async (subnetId) => {
                        return JSONRPC("http://localhost:9650/ext/info", "info.getNodeID", {});
                    }}
                    title={`info.getNodeID()`}
                />
                <APIInfoCard
                    subnetId={store.subnetId}
                    dataFetcher={async (subnetId) => {
                        return JSONRPC("http://localhost:9650/ext/info", "info.getNodeVersion", {});
                    }}
                    title={`info.getNodeVersion()`}
                />
                <APIInfoCard
                    subnetId={store.subnetId}
                    dataFetcher={async (subnetId) => {
                        return JSONRPC("http://localhost:9650/ext/info", "info.getVMs", {});
                    }}
                    title={`info.getVMs()`}
                />
                <APIInfoCard
                    subnetId={store.subnetId}
                    dataFetcher={async (subnetId) => {
                        const result = await JSONRPC("http://localhost:9650/ext/info", "info.peers", {});
                        result.peers = "Removed for brevity";
                        return result;
                    }}
                    title={`info.peers()`}
                />
                <div className="md:col-span-2">
                    <Header title="P-Chain" />
                </div>
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
            </div>
        </div>
    );
}

function APIInfoCard({
    subnetId, dataFetcher, title, note
}: {
    subnetId: string,
    dataFetcher: (subnetId: string) => Promise<any>,
    title: string,
    note?: string
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

    const renderContent = () => {
        if (error) {
            return <div className="bg-red-100 text-red-700 p-4 rounded overflow-auto">{error}</div>;
        }
        return <pre className="bg-gray-100 p-4 rounded overflow-auto">
            {JSON.stringify(data, (_, v) => typeof v === 'bigint' ? v.toString() : v, 2)}
        </pre>;
    };

    return (
        <div className="bg-white rounded-lg shadow-md p-6 h-fit">
            <h1 className="text-xl font-bold mb-4">{title}</h1>
            <div className="min-h-[100px]">
                {renderContent()}
            </div>
            {note && <p className="text-sm text-gray-500 mt-2">{note}</p>}
        </div>
    );
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
        <div className="bg-white rounded-lg shadow-md p-6 mb-4">
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
                <div className="flex items-center gap-2">
                    <label htmlFor="nodeRPCUrl" className="font-medium">Node RPC URL:</label>
                    <input
                        id="nodeRPCUrl"
                        type="text"
                        value={store.nodeRPCUrl}
                        onChange={(e) => store.setNodeRPCUrl(e.target.value)}
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
