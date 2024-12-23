declare global {
    interface Window {
        ethereum?: any;
    }
}

export async function getWalletAddress() {
    if (!window.ethereum) {
        throw new Error('No wallet detected');
    }

    const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
    if (!accounts || accounts.length === 0) {
        throw new Error('No account found');
    }

    return accounts[0]; // Return the first account
}
