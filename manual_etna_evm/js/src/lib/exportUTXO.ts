import { JsonRpcProvider } from 'ethers';
import { Context, evm, utils } from '@avalabs/avalanchejs';
import { AbstractWallet } from './wallet';

export async function exportUTXOFromCChain(wallet: AbstractWallet, amount: number) {
    const provider = new JsonRpcProvider('https://etna.avax-dev.network/ext/bc/C/rpc');

    const evmapi = new evm.EVMApi('https://etna.avax-dev.network');

    const context = await Context.getContextFromURI('https://etna.avax-dev.network');
    console.log("context is ", context);
    const txCount = await provider.getTransactionCount((await wallet.getAddress()).C);
    const baseFee = await evmapi.getBaseFee();
    const pAddressBytes = utils.bech32ToBytes((await wallet.getAddress()).P);

    const tx = evm.newExportTxFromBaseFee(
        context,
        baseFee / BigInt(1e9),
        BigInt(amount * 1e9),
        context.networkID.toString(),
        utils.hexToBuffer((await wallet.getAddress()).C),
        [pAddressBytes],
        BigInt(txCount),
    );

    await wallet.addTxSignatures(tx);
    const signedTx = tx.getSignedTx();

    const issueResponse = await evmapi.issueSignedTx(signedTx);

    return issueResponse;
}
