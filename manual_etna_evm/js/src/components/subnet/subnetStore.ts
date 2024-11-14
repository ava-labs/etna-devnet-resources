import { create } from "zustand";
import { persist, createJSONStorage } from 'zustand/middleware'

interface SubnetStore {
    subnetId: string;
    setSubnetId: (subnetId: string) => void;
}

export const useSubnetStore = create<SubnetStore>()(persist<SubnetStore>((set, get) => ({
    subnetId: "",
    setSubnetId: (subnetId: string) => set({ subnetId }),
}), {
    name: "subnet-storage",
    storage: createJSONStorage(() => localStorage),
}))
