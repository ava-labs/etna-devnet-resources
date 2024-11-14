import { AbstractWallet } from "../../lib/wallet";
import Button from "../../lib/Button";
import { Utxo } from "@avalabs/avalanchejs";
import { importUTXOs } from "./utxo";
import { useAsync } from "../../lib/hooks";
import { useWalletStore } from "../../lib/store";
import { getCChainBalance, getPChainbalance } from "./balances";

export function ImportUTXOs({ wallet, UTXOs }: { wallet: AbstractWallet, UTXOs: Utxo[] }) {
    const address = useWalletStore(state => state.address);

    const setCBalance = useWalletStore(state => state.setCBalance);
    const setPBalance = useWalletStore(state => state.setPBalance);

    const importPromise = useAsync(async () => {
        await importUTXOs(wallet, UTXOs)
        const balance = await getCChainBalance(address.C)
        setCBalance(balance)
        const pBalance = await getPChainbalance(address.P)
        setPBalance(pBalance)
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

