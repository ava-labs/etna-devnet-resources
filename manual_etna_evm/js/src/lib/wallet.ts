import { secp256k1, UnsignedTx, utils } from 'avalanchejs-bleeding-edge';
import { bytesToHex, hexToBytes } from '@noble/hashes/utils';
import { addr } from 'micro-eth-signer';


export interface AbstractWallet {
    getAddress: () => Promise<{
        C: string;
        P: string;
    }>;
    addTxSignatures: (unsignedTx: UnsignedTx) => Promise<void>;
    getAPIEndpoint: () => string;
}


export function getPrivateKeyWallet(privateKey: Uint8Array, apiEndpoint: string): AbstractWallet {
    const publicKey = secp256k1.getPublicKey(privateKey);

    const pChainAddress = `P-${utils.formatBech32(
        "custom",
        secp256k1.publicKeyBytesToAddress(publicKey)
    )}`;

    const cChainAddress = addr.fromPublicKey(publicKey);

    return {
        getAddress: async () => {
            return { C: cChainAddress, P: pChainAddress };
        },
        addTxSignatures: async (unsignedTx: UnsignedTx) => {
            const unsignedBytes = unsignedTx.toBytes();
            const publicKey = secp256k1.getPublicKey(privateKey);

            if (unsignedTx.hasPubkey(publicKey)) {
                const signature = await secp256k1.sign(unsignedBytes, privateKey);
                unsignedTx.addSignature(signature);
            }
        },
        getAPIEndpoint: () => apiEndpoint
    }
}

export function getLocalStorageWallet(apiEndpoint: string): AbstractWallet {
    let privateKeyHex = localStorage.getItem("privateKey");
    if (privateKeyHex === null) {
        privateKeyHex = bytesToHex(secp256k1.randomPrivateKey());
        localStorage.setItem("privateKey", privateKeyHex as string);
    }

    const privateKey = hexToBytes(privateKeyHex as string);
    return getPrivateKeyWallet(privateKey, apiEndpoint);
}
