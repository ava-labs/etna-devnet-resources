import { create } from "zustand";
import { persist, createJSONStorage } from 'zustand/middleware';

interface SubnetViewerState {
    subnetId: string;
    setSubnetId: (subnetId: string) => void;
}

export const useSubnetViewerStore = create<SubnetViewerState>()(
    persist<SubnetViewerState>(
        (set) => ({
            subnetId: "",
            setSubnetId: (subnetId: string) => set({ subnetId })
        }),
        {
            name: "subnetViewer",
            storage: createJSONStorage(() => localStorage),
        }
    )
);
