import { create } from "zustand";
import { persist, createJSONStorage } from 'zustand/middleware';
import defaultgenesis from "./defaultGenesis.json";

interface GenerateGenesisParams {
    userAddress: string;
}

interface GenesisState {
    genesis: string;
    generateGenesis: (params: GenerateGenesisParams) => void;
    clearGenesis: () => void;
}

export const useGenesisStore = create<GenesisState>()(persist<GenesisState>((set, get) => ({
    genesis: "",
    hasGeneratedGenesis: false,
    generateGenesis: ({ userAddress }: GenerateGenesisParams) => {
        let genesisString = JSON.stringify(defaultgenesis, null, 2)

        if (userAddress.startsWith("0x")) {
            userAddress = userAddress.slice(2);
        }

        genesisString = genesisString.split("%REPLACE_ME%").join(userAddress);
        set({ genesis: genesisString });
    },
    clearGenesis: () => {
        set({ genesis: "" });
    }
}), {
    name: "genesis-storage", // Key for localStorage
    storage: createJSONStorage(() => localStorage),
}));
