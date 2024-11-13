import { AbstractWallet } from "../../lib/wallet";
import { useState, useEffect } from "react";
import Transfer from "./Transfer";

async function makeGRPCRequest<T>(url: string, method: string, params: string[] | Record<string, string[]>): Promise<T> {
    const response = await fetch(url, {
        headers: {
            "content-type": "application/json",
            "accept": "*/*"
        },
        method: "POST",
        body: JSON.stringify({
            jsonrpc: "2.0",
            id: 0,
            method: method,
            params: params
        })
    });


    const data = await response.json();

    if (data.error?.message) {
        throw new Error(data.error.message);
    }

    return data.result as T;
}

async function getCChainBalance(address: string): Promise<bigint> {
    const response = await makeGRPCRequest<string>(
        "https://etna.avax-dev.network/ext/bc/C/rpc",
        "eth_getBalance",
        [address, "latest"]
    );

    return BigInt(response);
}

async function getPChainbalance(address: string): Promise<bigint> {
    const response = await makeGRPCRequest(
        "https://etna.avax-dev.network/ext/bc/P",
        "platform.getBalance",
        { addresses: [address] }
    ) as { balance: string };

    return BigInt(response.balance);
}

const P_CHAIN_DECIMALS = 9;
const C_CHAIN_DECIMALS = 18;

export default function Balance({ wallet }: { wallet: AbstractWallet }) {
    const [address, setAddress] = useState<{ C: string, P: string } | null>(null);
    const [balances, setBalances] = useState<{ C: bigint, P: bigint } | null>(null);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        wallet.getAddress()
            .then(setAddress)
            .catch(err => setError(err.message));
    }, [wallet]);

    useEffect(() => {
        if (!address) return;

        Promise.all([
            getCChainBalance(address.C),
            getPChainbalance(address.P)
        ])
            .then(([cBalance, pBalance]) => {
                setBalances({
                    C: cBalance,
                    P: pBalance
                });
                setError(null); // Clear any previous errors
            })
            .catch(err => setError(err.message));
    }, [address]);

    if (error) return <div className="error">Error: {error}</div>;
    if (!address || !balances) return <div>Loading...</div>;

    return <>
        <div>C-Chain address: {address.C}</div>
        <div>C-Chain balance: {Number(balances.C) / 10 ** C_CHAIN_DECIMALS} AVAX</div>
        <div>P-Chain balance: {Number(balances.P) / 10 ** P_CHAIN_DECIMALS} AVAX</div>
        <Transfer wallet={wallet} amount={0.2} />
    </>;
}
