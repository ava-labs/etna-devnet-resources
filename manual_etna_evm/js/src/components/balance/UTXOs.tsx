import { useAsync } from '../../lib/hooks';
import { ImportUTXOs } from './ImportUTXOs';
import { useEffect } from 'react';
import ExportUTXO from './exportUTXO';
import { useWalletStore } from './walletStore';


export function UTXOs({ minAmount }: { minAmount: number }) {
    const reloadUTXOs = useWalletStore(state => state.reloadUTXOs);
    const reloadUTXOsPromise = useAsync(reloadUTXOs);
    const utxos = useWalletStore(state => state.utxos);

    useEffect(() => {
        reloadUTXOsPromise.execute();
    }, []);

    if (reloadUTXOsPromise.error) return <div className="error">Error: {reloadUTXOsPromise.error}</div>;
    if (reloadUTXOsPromise.loading) return <div>Loading...</div>;

    if (utxos.length > 0) {
        return <>
            <ImportUTXOs UTXOs={utxos} />
        </>;
    }

    return <>
        <ExportUTXO minAmount={minAmount} />
    </>;
}
