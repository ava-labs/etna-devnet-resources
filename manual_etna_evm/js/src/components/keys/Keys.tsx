import { useAsync } from "../../lib/hooks";
import { useEffect } from 'react';
import { useWalletStore } from "../balance/walletStore";

export function Keys() {
    const wallet = useWalletStore(state => state.wallet);
    const cAddress = useWalletStore(state => state.cAddress);
    const pAddress = useWalletStore(state => state.pAddress);

    const addressPromise = useAsync(() => wallet!.getAddress());

    useEffect(() => {
        if (wallet) {
            addressPromise.execute();
        }
    }, [wallet]);

    return <>
        <p className="pb-4">
            Here are your keys. They are randomly generated in your browser and stored in local storage.
            Your private keys are never transmitted over the network.
        </p>
        <p className="pb-4">
            You can close your browser and
            come back later, but if you clear your browser's local storage before completing the subnet
            creation process, any transferred funds will be permanently lost.
        </p>
        <div>Generated addresses:</div>
        <div>C-Chain: {cAddress}</div>
        <div>P-Chain: {pAddress}</div>
    </>
}
