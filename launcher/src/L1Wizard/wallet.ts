import { createWalletClient, custom } from 'viem'

declare global {
    interface Window {
        ethereum?: any;
    }
}

export async function getWalletAddress() {
    const walletClient = createWalletClient({
        transport: custom(window.ethereum!)
    })
    const [account] = await walletClient.requestAddresses()
    if (!account) {
        throw new Error('No account found')
    }
    return account
}
