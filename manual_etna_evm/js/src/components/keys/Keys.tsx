import { AbstractWallet } from "../../lib/wallet";
import { useAsync } from "../../lib/hooks";
import { useEffect } from 'react';
import { useWalletStore } from "../../lib/store";

export function Keys({ wallet }: { wallet: AbstractWallet }) {
    const addressPromise = useAsync(() => wallet.getAddress());
    const cAddress = useWalletStore(state => state.cAddress);
    const pAddress = useWalletStore(state => state.pAddress);

    useEffect(() => {
        addressPromise.execute();
    }, [wallet]);

    useEffect(() => {
        if (addressPromise.data && addressPromise.data.C !== cAddress) {
            useWalletStore.getState().setCAddress(addressPromise.data.C);
        }
        if (addressPromise.data && addressPromise.data.P !== pAddress) {
            useWalletStore.getState().setPAddress(addressPromise.data.P);
        }
    }, [addressPromise.data]);

    if (addressPromise.error) {
        return <div>{addressPromise.error}</div>
    }

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
