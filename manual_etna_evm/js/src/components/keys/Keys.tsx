import { AbstractWallet } from "../../lib/wallet";
import { useAsync } from "../../lib/hooks";
import { useEffect } from 'react';

export function Keys({ wallet }: { wallet: AbstractWallet }) {
    const address = useAsync(() => wallet.getAddress());

    useEffect(() => {
        address.execute();
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
        <div>C-Chain: {address.loading ? 'Loading...' : address.data?.C}</div>
        <div>P-Chain: {address.loading ? 'Loading...' : address.data?.P}</div>
    </>
}
