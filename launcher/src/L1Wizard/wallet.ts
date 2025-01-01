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

export function getAddresses(privateKeyHex: string): { C: `0x${string}`, P: string } {
    const publicKey = secp256k1.getPublicKey(hexToBytes(privateKeyHex));

    const pChainAddress = `P-${utils.formatBech32(
        "fuji",
        secp256k1.publicKeyBytesToAddress(publicKey)
    )}`

    const cChainAddress = addr.fromPublicKey(publicKey) as `0x${string}`

    return {
        C: cChainAddress,
        P: pChainAddress
    }
}

export async function addSignature(tx: UnsignedTx, privateKeyHex: string) {
    const privateKey = hexToBytes(privateKeyHex);
    const unsignedBytes = tx.toBytes();
    const publicKey = secp256k1.getPublicKey(privateKey);

    if (tx.hasPubkey(publicKey)) {
        const signature = await secp256k1.sign(unsignedBytes, privateKey);
        tx.addSignature(signature);
    }
}
