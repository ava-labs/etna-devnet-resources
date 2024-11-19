import { Context, pvm, utils } from "@avalabs/avalanchejs";
import { useSubnetStore } from "./subnetStore";
import { AbstractWallet } from "../../lib/wallet";
import { useAsync } from "../../lib/hooks";
import { useWalletStore } from "../balance/walletStore";
import Button from "../../lib/Button";
import { useEffect } from "react";

export default function CreateSubnet() {
    const subnetId = useSubnetStore((state) => state.subnetId);
    const wallet = useWalletStore((state) => state.wallet);
    const setSubnetId = useSubnetStore((state) => state.setSubnetId);
    const reloadBalances = useWalletStore((state) => state.reloadBalances);

    const createSubnetPromise = useAsync(async () => {
        const subnetId = await createSubnet(wallet!);
        setSubnetId(subnetId);
        await reloadBalances();
    });

    const getSubnetPromise = useAsync(async () => {
        const pvmApi = new pvm.PVMApi(wallet!.getAPIEndpoint());
        const { subnets } = await pvmApi.getSubnets({ ids: [subnetId] });
        console.log("Subnets: ", subnets)
        return subnets
    })

    useEffect(() => {
        getSubnetPromise.execute()
    }, [subnetId])


    if (createSubnetPromise.loading) {
        return <div>Creating subnet...</div>
    }

    if (getSubnetPromise.loading) return <div>Getting subnet...</div>

    if (createSubnetPromise.error) return <div>Error creating subnet: {createSubnetPromise.error}</div>
    if (getSubnetPromise.error) return <div>Error getting subnet: {getSubnetPromise.error}</div>

    if (subnetId) {
        return <>
            <div className="mb-4">✅ Subnet {subnetId} created!</div>
            <pre className="bg-gray-100 p-4 rounded-lg">{JSON.stringify(getSubnetPromise.data, null, 2)}</pre>
        </>
    } else {
        return <>
            <Button onClick={createSubnetPromise.execute}>Create Subnet</Button>
        </>
    }


}

async function createSubnet(wallet: AbstractWallet): Promise<string> {
    if (!wallet) {
        throw new Error("Wallet not connected");
    }

    const { P: pAddress } = await wallet.getAddress();

    const uri = wallet.getAPIEndpoint();
    const pvmApi = new pvm.PVMApi(uri);
    const feeState = await pvmApi.getFeeState();
    const context = await Context.getContextFromURI(uri);

    const { utxos } = await pvmApi.getUTXOs({ addresses: [pAddress] });

    const testPAddr = utils.bech32ToBytes(pAddress);

    const tx = pvm.e.newCreateSubnetTx(
        {
            feeState,
            fromAddressesBytes: [testPAddr],
            utxos,
            subnetOwners: [testPAddr],
        },
        context,
    );

    await wallet.signRawTx(tx);

    const response = await pvmApi.issueSignedTx(tx.getSignedTx());
    return response.txID;//FIXME: I am not sure if txID is the subnetID
}
