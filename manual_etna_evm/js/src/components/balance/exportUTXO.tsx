import Button from "../../lib/Button";
import { useAsync } from "../../lib/hooks";
import { useWalletStore } from "../../lib/store";
import { exportUTXO } from "./utxo";

export default function ExportUTXO({ minAmount }: { minAmount: number }) {
    const pBalance = useWalletStore(state => state.pBalance);
    const wallet = useWalletStore(state => state.wallet);
    const reloadUTXOs = useWalletStore(state => state.reloadUTXOs);
    const reloadBalances = useWalletStore(state => state.reloadBalances);

    const toExport = BigInt(minAmount * 1e9) - BigInt(pBalance)

    const exportPromise = useAsync(async () => {
        await exportUTXO(wallet!, minAmount)

        await Promise.all([reloadUTXOs(), reloadBalances()])
    });

    if (exportPromise.loading) {
        return <div>Exporting UTXO...</div>;
    }

    if (exportPromise.error) {
        return <div>Error exporting UTXO: {exportPromise.error}</div>;
    }

    return <div>
        <p className="mb-4">You'll need {minAmount} AVAX on P chain. You have to export {Number(toExport) / 1e9} AVAX from C chain to P chain.</p>
        <Button onClick={exportPromise.execute} >Export a UTXO for {Number(toExport) / 1e9} AVAX from C to P</Button>
    </div>;
}
