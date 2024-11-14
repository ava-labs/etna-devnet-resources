import { useEffect } from "react";
import { UTXOs } from "./UTXOs";
import { useAsync } from '../../lib/hooks';
import { getCChainBalance, getPChainbalance, MINIMUM_P_CHAIN_BALANCE, P_CHAIN_DIVISOR, C_CHAIN_DIVISOR } from "./balances";
import { useWalletStore } from "../../lib/store";


const MINIMUM_BALANCE = 5.1;

export default function Balance() {
    const wallet = useWalletStore(state => state.wallet);
    const setCBalance = useWalletStore(state => state.setCBalance);
    const setPBalance = useWalletStore(state => state.setPBalance);
    const cAddress = useWalletStore(state => state.cAddress);
    const pAddress = useWalletStore(state => state.pAddress);

    const pBalance = useWalletStore(state => state.pBalance);
    const cBalance = useWalletStore(state => state.cBalance);

    const cBalancePromise = useAsync(() => getCChainBalance(cAddress));
    const pBalancePromise = useAsync(() => getPChainbalance(pAddress));

    useEffect(() => {
        cBalancePromise.execute();
        pBalancePromise.execute();
    }, [cAddress, pAddress]);

    useEffect(() => {
        if (cBalancePromise.data !== null) {
            setCBalance(cBalancePromise.data);
        }
        if (pBalancePromise.data !== null) {
            setPBalance(pBalancePromise.data);
        }
    }, [cBalancePromise.data, pBalancePromise.data]);

    if (cBalancePromise.error || pBalancePromise.error) {
        return <div className="error">Error: {cBalancePromise.error || pBalancePromise.error}</div>;
    }
    if (cBalancePromise.loading || pBalancePromise.loading) {
        return <div>Loading...</div>;
    }

    return <>
        <div>C-Chain balance: {Number(cBalance) / C_CHAIN_DIVISOR} AVAX</div>
        <div>P-Chain balance: {Number(pBalance) / P_CHAIN_DIVISOR} AVAX</div>
        {pBalance <= MINIMUM_P_CHAIN_BALANCE && <div className="mt-4">
            <UTXOs minAmount={MINIMUM_BALANCE} />
        </div>}
    </>;
}
