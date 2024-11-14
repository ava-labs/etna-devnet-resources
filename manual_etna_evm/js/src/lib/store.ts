import { create } from "zustand";
import { AbstractWallet } from "./wallet";
import { Utxo } from "@avalabs/avalanchejs";
import { getUTXOS } from "../components/balance/utxo";
import { getCChainBalance } from "../components/balance/balances";
import { getPChainbalance } from "../components/balance/balances";

interface WalletState {
    pAddress: string;
    cAddress: string;

    pBalance: bigint;
    setPBalance: (pBalance: bigint) => void;

    cBalance: bigint;
    setCBalance: (cBalance: bigint) => void;

    wallet: AbstractWallet | null;
    setWallet: (wallet: AbstractWallet) => Promise<void>;

    utxos: Utxo[];
    setUTXOs: (utxos: Utxo[]) => void;

    reloadUTXOs: () => Promise<void>;
    reloadBalances: () => Promise<void>;
}

export const useWalletStore = create<WalletState>((set, get) => ({
    pAddress: "",
    cAddress: "",

    pBalance: BigInt(0),
    setPBalance: pBalance => set({ pBalance }),

    cBalance: BigInt(0),
    setCBalance: cBalance => set({ cBalance }),

    wallet: null,
    setWallet: async (wallet) => {
        const { C, P } = await wallet.getAddress();
        set({ wallet, pAddress: P, cAddress: C });
    },

    utxos: [],
    setUTXOs: utxos => set({ utxos }),

    reloadUTXOs: async () => {
        const { wallet, setUTXOs } = get();
        if (!wallet) {
            throw new Error("Wallet not set");
        }
        const utxos = await getUTXOS(wallet);
        setUTXOs(utxos);
    },

    reloadBalances: async () => {
        const { setPBalance, setCBalance, pAddress, cAddress } = get();
        const [pBalance, cBalance] = await Promise.all([
            getPChainbalance(pAddress),
            getCChainBalance(cAddress)
        ]);
        setPBalance(pBalance);
        setCBalance(cBalance);
    },
}))
