import { useEffect, useState } from "react";
import { useChainStore } from "./chainStore";
import { useAsync } from "../../lib/hooks";
import { useWalletStore } from "../balance/walletStore";
import { useSubnetStore } from "../subnet/subnetStore";
import { useGenesisStore } from "../genesis/genesisStore";
import Button from "../../lib/Button";
import { pvm } from "@avalabs/avalanchejs";

export default function CreateChain() {
    const chainStore = useChainStore();
    const [chainName, setChainName] = useState("MyL1");
    const { subnetId } = useSubnetStore()
    const { genesis } = useGenesisStore()
    const { wallet } = useWalletStore()

    useEffect(() => {
        if (chainStore.chainName) {
            setChainName(chainStore.chainName)
        }
    }, [chainStore.chainName])

    const createChainPromise = useAsync(async () => {
        console.log("Creating chain with params: ", {
            chainName: chainName,
            wallet: wallet!,
            subnetId: subnetId,
            genesisData: genesis,
        })

        await chainStore.createChain({
            chainName: chainName,
            wallet: wallet!,
            subnetId: subnetId,
            genesisData: genesis,
        })
    })

    if (createChainPromise.loading) {
        return <div>Loading...</div>
    }

    if (createChainPromise.error) {
        return <div>Error: {createChainPromise.error}</div>
    }

    const isChainCreated = chainStore.chainId.length > 0
    return (
        <>
            {isChainCreated ? (
                <div>✅ Chain {chainStore.chainId} created!</div>
            ) : (
                <>
                    <label htmlFor="chainName" className="block text-gray-700 mb-2">
                        Chain Name:
                    </label>
                    <input
                        type="text"
                        id="chainName"
                        value={chainName}
                        onChange={(e) => setChainName(e.target.value)}
                        className="w-full p-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 mb-4"
                    />
                    <Button onClick={createChainPromise.execute}>
                        Create Chain
                    </Button>
                </>
            )}
        </>
    );
}
