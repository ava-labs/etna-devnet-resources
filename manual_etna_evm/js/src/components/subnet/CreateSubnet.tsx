import { Context, pvm, utils } from "@avalabs/avalanchejs";
import { useSubnetStore } from "./subnetStore";
import { AbstractWallet } from "../../lib/wallet";
import { useAsync } from "../../lib/hooks";
import { useWalletStore } from "../balance/walletStore";
import Button from "../../lib/Button";

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

    if (subnetId) {
        return <div>✅ Subnet {subnetId} created!</div>
    }

    if (createSubnetPromise.loading) {
        return <div>Creating subnet...</div>
    }

    if (createSubnetPromise.error) {
        return <div>Error creating subnet: {createSubnetPromise.error}</div>
    }

    return <>
        <Button onClick={createSubnetPromise.execute}>Create Subnet</Button>
    </>
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
    return response.txID;
}
