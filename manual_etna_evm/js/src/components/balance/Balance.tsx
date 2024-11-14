import { useEffect } from "react";
import { UTXOs } from "./UTXOs";
import { useAsync } from '../../lib/hooks';
import { MINIMUM_P_CHAIN_BALANCE_AVAX, P_CHAIN_DIVISOR, C_CHAIN_DIVISOR } from "./balances";
import { useWalletStore } from "./walletStore";


export default function Balance() {
    const reloadBalances = useWalletStore(state => state.reloadBalances);
    const reloadBalancesPromise = useAsync(reloadBalances);

    useEffect(() => {
        reloadBalancesPromise.execute();
    }, []);

    const cBalance = useWalletStore(state => state.cBalance);
    const pBalance = useWalletStore(state => state.pBalance);

    if (reloadBalancesPromise.error) {
        return <div className="error">Error: {reloadBalancesPromise.error}</div>;
    }
    if (reloadBalancesPromise.loading) {
        return <div>Loading...</div>;
    }

    return <>
        <div>C-Chain balance: {Number(cBalance) / C_CHAIN_DIVISOR} AVAX</div>
        <div>P-Chain balance: {Number(pBalance) / P_CHAIN_DIVISOR} AVAX</div>
        {pBalance <= MINIMUM_P_CHAIN_BALANCE_AVAX && <div className="mt-4">
            <UTXOs minAmount={MINIMUM_P_CHAIN_BALANCE_AVAX} />
        </div>}
    </>;
}
