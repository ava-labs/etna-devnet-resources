import { create } from "zustand";

interface WalletState {
    pAddress: string;
    cAddress: string;
    pBalance: bigint;
    cBalance: bigint;
    setPAddress: (pAddress: string) => void;
    setCAddress: (cAddress: string) => void;
    setPBalance: (pBalance: bigint) => void;
    setCBalance: (cBalance: bigint) => void;
}

export const useWalletStore = create<WalletState>((set) => ({
    pAddress: "",
    cAddress: "",
    pBalance: BigInt(0),
    cBalance: BigInt(0),
    setPAddress: pAddress => set({ pAddress }),
    setCAddress: cAddress => set({ cAddress }),
    setPBalance: pBalance => set({ pBalance }),
    setCBalance: cBalance => set({ cBalance }),
}))
