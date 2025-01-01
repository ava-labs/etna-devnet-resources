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


import { secp256k1, UnsignedTx, utils } from '@avalabs/avalanchejs';
import { bytesToHex, hexToBytes } from '@noble/hashes/utils';
import { addr } from 'micro-eth-signer';


export function newPrivateKey(): string {
    return bytesToHex(secp256k1.randomPrivateKey());
}

export function getAddresses(privateKeyHex: string): { C: string, P: string } {
    const publicKey = secp256k1.getPublicKey(hexToBytes(privateKeyHex));

    const pChainAddress = `P-${utils.formatBech32(
        "fuji",
        secp256k1.publicKeyBytesToAddress(publicKey)
    )}`;

    const cChainAddress = addr.fromPublicKey(publicKey);

    return {
        C: cChainAddress,
        P: pChainAddress
    }
}
