import { JsonRpcProvider } from 'ethers';
import { addTxSignatures, Context, evm, utils } from 'avalanchejs-bleeding-edge'

const C_CHAIN_ADDRESS = "0x700046af78cD0E55D5d25025bCEaE992B87A5686"
const X_CHAIN_ADDRESS = "X-custom125uzjyu824gwlz6e4jqthvks9er6xwvx7ve80v"
const TEST_PRIVATE_KEY = "0963764d0da7b2f6a4327eb49042bf9e05e88534f22acb9fbf469eb44a8ed3f8"

const main = async (publicUrl: string) => {
    const provider = new JsonRpcProvider(publicUrl + '/ext/bc/C/rpc');

    const evmapi = new evm.EVMApi(publicUrl);

    const context = await Context.getContextFromURI(publicUrl);
    const txCount = await provider.getTransactionCount(C_CHAIN_ADDRESS);
    const baseFee = await evmapi.getBaseFee();
    const addressBytes = utils.bech32ToBytes(X_CHAIN_ADDRESS);

    const tx = evm.newExportTxFromBaseFee(
        context,
        baseFee / BigInt(1e9),
        BigInt(0.1 * 1e9),
        context.pBlockchainID,
        utils.hexToBuffer(C_CHAIN_ADDRESS),
        [addressBytes],
        BigInt(txCount),
    );

    await addTxSignatures({
        unsignedTx: tx,
        privateKeys: [utils.hexToBuffer(TEST_PRIVATE_KEY)],
    });

    await evmapi.issueSignedTx(tx.getSignedTx());

    console.log("Done");
};


console.log(`\n\n\nExecuting on Etna:\n\n\n`)
const ETNA_PUBLIC_URL = "https://etna.avax-dev.network"
await main(ETNA_PUBLIC_URL)
