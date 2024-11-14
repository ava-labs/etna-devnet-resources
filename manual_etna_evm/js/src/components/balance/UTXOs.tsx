import { AbstractWallet } from '../../lib/wallet';
import { exportUTXO, getUTXOS } from './utxo';
import { useAsync } from '../../lib/hooks';
import { ImportUTXOs } from './ImportUTXOs';
import { useEffect } from 'react';

export function UTXOs({ wallet, minAmount }: { wallet: AbstractWallet, minAmount: number }) {
    const loadUTXOs = useAsync(() => getUTXOS(wallet));

    useEffect(() => {
        loadUTXOs.execute();
    }, [wallet]);


    if (loadUTXOs.error) return <div className="error">Error: {loadUTXOs.error}</div>;
    if (loadUTXOs.loading) return <div>Loading...</div>;

    if (loadUTXOs.data && loadUTXOs.data.length > 0) {
        return <>
            <ImportUTXOs wallet={wallet} UTXOs={loadUTXOs.data} />
        </>;
    }

    return <>
        <div className="text-lg font-bold">
            TODO: not implemented
        </div>
    </>;
}
