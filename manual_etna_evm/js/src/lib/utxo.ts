import { pvm } from "avalanchejs-bleeding-edge";

import { JsonRpcProvider } from 'ethers';
import { Context, evm, utils } from 'avalanchejs-bleeding-edge'
import { AbstractWallet } from './wallet';

export async function exportUTXO(wallet: AbstractWallet, amount: number) {
    const provider = new JsonRpcProvider(wallet.getAPIEndpoint() + '/ext/bc/C/rpc');
    const evmapi = new evm.EVMApi(wallet.getAPIEndpoint());

    const context = await Context.getContextFromURI(wallet.getAPIEndpoint());
    const address = await wallet.getAddress();
    const txCount = await provider.getTransactionCount(address.C);
    const baseFee = await evmapi.getBaseFee();
    const addressBytes = utils.bech32ToBytes(address.P);

    const tx = evm.newExportTxFromBaseFee(
        context,
        baseFee / BigInt(1e9),
        BigInt(amount * 1e9),
        context.pBlockchainID,
        utils.hexToBuffer(address.C),
        [addressBytes],
        BigInt(txCount),
    );

    await wallet.addTxSignatures(tx);
    await evmapi.issueSignedTx(tx.getSignedTx());
}

export async function getUTXOS(wallet: AbstractWallet) {
    //FIXME: ignores pagination
    const pvmApi = new pvm.PVMApi(wallet.getAPIEndpoint());

    const { utxos } = await pvmApi.getUTXOs({
        sourceChain: 'C',
        addresses: [(await wallet.getAddress()).P],
    });

    return utxos;
}
