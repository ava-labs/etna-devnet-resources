import { AbstractWallet } from '../../lib/wallet';
import { exportUTXO, getUTXOS } from '../../lib/utxo';
import { useAsync } from '../../lib/hooks';

export default function Transfer({ wallet, amount }: { wallet: AbstractWallet, amount: number }) {
    const exportTx = useAsync(() => exportUTXO(wallet, amount));
    const getUtxos = useAsync(() => getUTXOS(wallet));

    return <>
        <div className={`${exportTx.loading || getUtxos.loading ? 'opacity-50' : ''}`}>
            <h2 className="text-2xl font-bold">Transfer</h2>
            <button
                className="px-4 py-2 bg-blue-500 text-white rounded mr-2"
                onClick={exportTx.execute}
                disabled={exportTx.loading || getUtxos.loading}
            >
                1. Export Tx
            </button>
            <button
                className="px-4 py-2 bg-blue-500 text-white rounded"
                onClick={getUtxos.execute}
                disabled={exportTx.loading || getUtxos.loading}
            >
                2. Get My UTXOS
            </button>
            {(exportTx.loading || getUtxos.loading) && (
                <span className="ml-2">Loading...</span>
            )}

            {(exportTx.error || getUtxos.error) && (
                <div className="mt-4 p-4 bg-red-100 text-red-700 rounded">
                    {exportTx.error || getUtxos.error}
                </div>
            )}

            {getUtxos.data && getUtxos.data.length > 0 && (
                <div className="mt-4">
                    <h3 className="text-xl font-semibold">UTXOs:</h3>
                    <pre className="mt-2 p-4 bg-gray-100 rounded overflow-auto">
                        {JSON.stringify(getUtxos.data, null, 2)}
                    </pre>
                </div>
            )}
        </div>
    </>;
}
