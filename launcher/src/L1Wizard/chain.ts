import { utils } from "@avalabs/avalanchejs";
import { Context } from "@avalabs/avalanchejs";
import { pvm } from "@avalabs/avalanchejs";
import { RPC_ENDPOINT } from "./utxo";
import { addSignature, getAddresses } from "./wallet";

export async function createSubnet(privateKeyHex: string): Promise<string> {
    if (!privateKeyHex) {
        throw new Error("Private key required");
    }

    const pvmApi = new pvm.PVMApi(RPC_ENDPOINT);
    const feeState = await pvmApi.getFeeState();
    const context = await Context.getContextFromURI(RPC_ENDPOINT);

    const { P: pAddress } = await getAddresses(privateKeyHex);
    const addressBytes = utils.bech32ToBytes(pAddress);

    const { utxos } = await pvmApi.getUTXOs({
        addresses: [pAddress]
    });

    const tx = pvm.e.newCreateSubnetTx(
        {
            feeState,
            fromAddressesBytes: [addressBytes],
            utxos,
            subnetOwners: [addressBytes],
        },
        context,
    );

    await addSignature(tx, privateKeyHex);

    const response = await pvmApi.issueSignedTx(tx.getSignedTx());
    return response.txID;
}

export interface CreateChainParams {
    privateKeyHex: string;
    chainName: string;
    subnetId: string;
    genesisData: string;
}

export const SUBNET_EVM_ID = "srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy";


export async function createChain(params: CreateChainParams): Promise<string> {
    if (!params.privateKeyHex) {
        throw new Error("Private key required");
    }

    const pvmApi = new pvm.PVMApi(RPC_ENDPOINT);
    const feeState = await pvmApi.getFeeState();
    const context = await Context.getContextFromURI(RPC_ENDPOINT);

    const { P: pAddress } = await getAddresses(params.privateKeyHex);
    const addressBytes = utils.bech32ToBytes(pAddress);

    const { utxos } = await pvmApi.getUTXOs({
        addresses: [pAddress]
    });

    const tx = pvm.e.newCreateChainTx(
        {
            feeState,
            fromAddressesBytes: [addressBytes],
            utxos,
            chainName: params.chainName,
            subnetAuth: [0],
            subnetId: params.subnetId,
            vmId: SUBNET_EVM_ID,
            fxIds: [],
            genesisData: JSON.parse(params.genesisData),
        },
        context,
    );

    await addSignature(tx, params.privateKeyHex);

    const response = await pvmApi.issueSignedTx(tx.getSignedTx());
    return response.txID;
}
