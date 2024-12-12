import { create } from "zustand";
import { persist, createJSONStorage } from 'zustand/middleware';

export interface SubnetViewerState {
    subnetId: string;
    setSubnetId: (subnetId: string) => void;
    nodeRPCUrl: string;
    setNodeRPCUrl: (nodeRPCUrl: string) => void;
}

export const useSubnetViewerStore = create<SubnetViewerState>()(
    persist<SubnetViewerState>(
        (set) => ({
            subnetId: "",
            setSubnetId: (subnetId: string) => set({ subnetId }),
            nodeRPCUrl: "http://localhost:9650",
            setNodeRPCUrl: (nodeRPCUrl: string) => set({ nodeRPCUrl })
        }),
        {
            name: "subnetViewer",
            storage: createJSONStorage(() => localStorage),
        }
    )
);
