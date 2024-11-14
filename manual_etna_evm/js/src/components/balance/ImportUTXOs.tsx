import Button from "../../lib/Button";
import { Utxo } from "avalanchejs-bleeding-edge";
import { importUTXOs } from "./utxo";
import { useAsync } from "../../lib/hooks";
import { useWalletStore } from "./walletStore";

export function ImportUTXOs({ UTXOs }: { UTXOs: Utxo[] }) {
    const wallet = useWalletStore(state => state.wallet);
    const reloadUTXOs = useWalletStore(state => state.reloadUTXOs);
    const reloadBalances = useWalletStore(state => state.reloadBalances);

    const importPromise = useAsync(async () => {
        await importUTXOs(wallet!, UTXOs)
        await Promise.all([reloadUTXOs(), reloadBalances()])
    });

    if (importPromise.loading) {
        return <div>Importing UTXOs...</div>;
    }

    if (importPromise.error) {
        return <div>Error importing UTXOs: {importPromise.error}</div>;
    }

    return <div>
        <div className=" mb-4">You have {UTXOs.length} UTXOs</div>
        <Button onClick={importPromise.execute} >Import UTXOs to P-Chain</Button>
    </div>;
}

