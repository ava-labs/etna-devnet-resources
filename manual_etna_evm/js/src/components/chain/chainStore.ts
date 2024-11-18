import { create } from "zustand";
import { persist, createJSONStorage } from 'zustand/middleware';
import { AbstractWallet } from "../../lib/wallet";
import { utils } from "avalanchejs-bleeding-edge";
import { Context } from "avalanchejs-bleeding-edge";
import { pvm } from "avalanchejs-bleeding-edge";

const SUBNET_EVM_ID = 'srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy'

interface CreateChainParams {
    chainName: string;
    wallet: AbstractWallet;
    subnetId: string;
    genesisData: string;
}

interface ChainState {
    chainId: string;
    chainName: string;
    createChain: (params: CreateChainParams) => Promise<void>;
}

export const useChainStore = create<ChainState>()(persist<ChainState>((set) => ({
    chainId: '',
    chainName: '',
    createChain: async (params: CreateChainParams): Promise<void> => {
        const { P: pAddress } = await params.wallet.getAddress();

        const uri = params.wallet.getAPIEndpoint();
        const pvmApi = new pvm.PVMApi(uri);
        const feeState = await pvmApi.getFeeState();
        const context = await Context.getContextFromURI(uri);

        const { utxos } = await pvmApi.getUTXOs({ addresses: [pAddress] });

        const testPAddr = utils.bech32ToBytes(pAddress);


        const tx = pvm.e.newCreateChainTx(
            {
                feeState,
                fromAddressesBytes: [testPAddr],
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

        await params.wallet.signRawTx(tx);

        console.log(tx.getSignedTx());

        const { txID } = await pvmApi.issueSignedTx(tx.getSignedTx());
        const chainId = txID;
        set({ chainId, chainName: params.chainName });
    },
}), {
    name: "chain-storage", // Key for localStorage
    storage: createJSONStorage(() => localStorage),
}));
