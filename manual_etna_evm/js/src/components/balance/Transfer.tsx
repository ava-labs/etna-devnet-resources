import { pvm, Utxo } from '@avalabs/avalanchejs';
import { AbstractWallet } from '../../lib/wallet';
import { useState } from 'react';
import { exportTx } from '../../lib/exportUTXO';

async function getMyUTXOS(wallet: AbstractWallet) {
    const pvmApi = new pvm.PVMApi("https://etna.avax-dev.network");

    const { utxos } = await pvmApi.getUTXOs({
        sourceChain: 'X',
        addresses: [(await wallet.getAddress()).P],
    });

    return utxos;
}

export default function Transfer({ wallet, amount }: { wallet: AbstractWallet, amount: number }) {
    const [error, setError] = useState<string | null>(null);
    const [utxos, setUtxos] = useState<Utxo[]>([]);
    const [loading, setLoading] = useState(0);

    const handleExportTx = async () => {
        try {
            setLoading(l => l + 1);
            setError(null);
            await exportTx(wallet, amount);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred during export');
        } finally {
            setLoading(l => l - 1);
        }
    };

    const handleGetUTXOS = async () => {
        try {
            setLoading(l => l + 1);
            setError(null);
            const result = await getMyUTXOS(wallet);
            setUtxos(result);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred fetching UTXOs');
        } finally {
            setLoading(l => l - 1);
        }
    };

    return <>
        <h2 className="text-2xl font-bold">Transfer</h2>
        <button
            className="px-4 py-2 bg-blue-500 text-white rounded mr-2"
            onClick={handleExportTx}
            disabled={loading > 0}
        >
            1. Export Tx
        </button>
        <button
            className="px-4 py-2 bg-blue-500 text-white rounded"
            onClick={handleGetUTXOS}
            disabled={loading > 0}
        >
            2. Get My UTXOS
        </button>
        {loading > 0 && (
            <span className="ml-2">Loading...</span>
        )}

        {error && (
            <div className="mt-4 p-4 bg-red-100 text-red-700 rounded">
                {error}
            </div>
        )}

        {utxos.length > 0 && (
            <div className="mt-4">
                <h3 className="text-xl font-semibold">UTXOs:</h3>
                <pre className="mt-2 p-4 bg-gray-100 rounded overflow-auto">
                    {JSON.stringify(utxos, null, 2)}
                </pre>
            </div>
        )}
    </>;
}
