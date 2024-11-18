
export const P_CHAIN_DIVISOR = 1e9;
export const C_CHAIN_DIVISOR = 1e18;
export const MINIMUM_P_CHAIN_BALANCE_AVAX = 6.0

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

export async function getCChainBalance(address: string): Promise<bigint> {
    const response = await makeGRPCRequest<string>(
        "https://etna.avax-dev.network/ext/bc/C/rpc",
        "eth_getBalance",
        [address, "latest"]
    );

    return BigInt(response);
}

export async function getPChainbalance(address: string): Promise<bigint> {
    const response = await makeGRPCRequest(
        "https://etna.avax-dev.network/ext/bc/P",
        "platform.getBalance",
        { addresses: [address] }
    ) as { balance: string };

    return BigInt(response.balance);
}
