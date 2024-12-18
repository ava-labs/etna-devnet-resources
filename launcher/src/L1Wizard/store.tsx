export const stepList = {
    "genesis": {
        title: "Create genesis",
        description: "Allocations and precompiles",
        icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24" >
            <path fillRule="evenodd" d="M3.559 4.544c.355-.35.834-.544 1.33-.544H19.11c.496 0 .975.194 1.33.544.356.35.559.829.559 1.331v9.25c0 .502-.203.981-.559 1.331-.355.35-.834.544-1.33.544H15.5l-2.7 3.6a1 1 0 0 1-1.6 0L8.5 17H4.889c-.496 0-.975-.194-1.33-.544A1.868 1.868 0 0 1 3 15.125v-9.25c0-.502.203-.981.559-1.331ZM7.556 7.5a1 1 0 1 0 0 2h8a1 1 0 0 0 0-2h-8Zm0 3.5a1 1 0 1 0 0 2H12a1 1 0 1 0 0-2H7.556Z" clipRule="evenodd" />
        </svg >
    },
    "generate-keys": {
        title: "Generate keys",
        description: "Generate keys for your nodes",
        icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24" >
            <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M16 12H4m12 0-4 4m4-4-4-4m3-4h2a3 3 0 0 1 3 3v10a3 3 0 0 1-3 3h-2" />
        </svg >
    },
    "create-l1": {
        title: "Create an L1",
        description: "Create a subnet, chain and convert to L1",
        icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24" >
            <path d="m3 2 1.578 17.834L12 22l7.468-2.165L21 2H3Zm13.3 14.722-4.293 1.204H12l-4.297-1.204-.297-3.167h2.108l.15 1.526 2.335.639 2.34-.64.245-3.05h-7.27l-.187-2.006h7.64l.174-2.006H6.924l-.176-2.006h10.506l-.954 10.71Z" />
        </svg >

        // },
        // "launch-nodes": {
        //     title: "Launch nodes",
        //     description: "Launch nodes on your infra",
        //     icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24" >
        //         <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M3 10h18M6 14h2m3 0h5M3 7v10a1 1 0 0 0 1 1h16a1 1 0 0 0 1-1V7a1 1 0 0 0-1-1H4a1 1 0 0 0-1 1Z" />
        //     </svg >

        // },
        // "add-to-wallet": {
        //     title: "Add to wallet",
        //     description: "Add your L1 to your wallet",
        //     icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24" >
        //         <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 8v8m0-8h8M8 8H6a2 2 0 1 1 2-2v2Zm0 8h8m-8 0H6a2 2 0 1 0 2 2v-2Zm8 0V8m0 8h2a2 2 0 1 1-2 2v-2Zm0-8h2a2 2 0 1 0-2-2v2Z" />
        //     </svg >

        // },
        // "deploy-validator-manager": {
        //     title: "Deploy validator manager",
        //     description: "Deploy contract on your L1",
        //     icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24" >
        //         <path fillRule="evenodd" d="M12 8a1 1 0 0 0-1 1v10H9a1 1 0 1 0 0 2h11a1 1 0 0 0 1-1V9a1 1 0 0 0-1-1h-8Zm4 10a2 2 0 1 1 0-4 2 2 0 0 1 0 4Z" clipRule="evenodd" />
        //         <path fillRule="evenodd" d="M5 3a2 2 0 0 0-2 2v6h6V9a3 3 0 0 1 3-3h8c.35 0 .687.06 1 .17V5a2 2 0 0 0-2-2H5Zm4 10H3v2a2 2 0 0 0 2 2h4v-4Z" clipRule="evenodd" />
        //     </svg >

        // },
        // "initialize-validator-manager": {
        //     title: "Initialize validator manager",
        //     description: "Initialize validator manager on the genesis",
        //     icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24" >
        //         <path fillRule="evenodd" d="M3 6a3 3 0 1 1 4 2.83v6.34a3.001 3.001 0 1 1-2 0V8.83A3.001 3.001 0 0 1 3 6Zm11.207-2.707a1 1 0 0 1 0 1.414L13.914 5H15a4 4 0 0 1 4 4v6.17a3.001 3.001 0 1 1-2 0V9a2 2 0 0 0-2-2h-1.086l.293.293a1 1 0 0 1-1.414 1.414l-2-2a1 1 0 0 1 0-1.414l2-2a1 1 0 0 1 1.414 0Z" clipRule="evenodd" />
        //     </svg >
        // },
    }
}

import { create } from 'zustand'
import { StateCreator } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'

interface WizardState {
    ownerEthAddress: string;
    setOwnerEthAddress: (address: string) => void;
    currentStep: keyof typeof stepList;
    advanceFrom: (givenStep: keyof typeof stepList, direction?: "up" | "down") => void;
    nodesCount: number;
    setNodesCount: (count: number) => void;
    chainId: number;
    setChainId: (chainId: number) => void;
    genesisString: string;
    regenerateGenesis: () => Promise<void>;
    nodePopJsons: string[];
    setNodePopJsons: (nodePopJsons: string[]) => void;
}


const wizardStoreFunc: StateCreator<WizardState> = (set, get) => ({
    ownerEthAddress: "",
    setOwnerEthAddress: (address: string) => set(() => ({ ownerEthAddress: address })),
    currentStep: Object.keys(stepList)[0] as keyof typeof stepList,
    advanceFrom: (givenStep, direction: "up" | "down" = "up") => set((state) => {
        const stepKeys = Object.keys(stepList) as (keyof typeof stepList)[];
        const currentIndex = stepKeys.indexOf(givenStep);
        if (direction === "up" && currentIndex < stepKeys.length - 1) {
            return { currentStep: stepKeys[currentIndex + 1] };
        }
        if (direction === "down" && currentIndex > 0) {
            return { currentStep: stepKeys[currentIndex - 1] };
        }
        return state;
    }),
    nodesCount: 3,
    setNodesCount: (count: number) => set(() => ({ nodesCount: count })),
    chainId: Math.floor(Math.random() * 1000000) + 1,
    setChainId: (chainId: number) => set(() => ({ chainId: chainId })),
    genesisString: "",
    regenerateGenesis: async () => {
        const params = new URLSearchParams({
            ownerEthAddressString: get().ownerEthAddress,
            evmChainId: get().chainId.toString()
        });

        const response = await fetch(`/api/generateGenesis?${params}`);
        if (!response.ok) {
            throw new Error('Failed to generate genesis');
        }
        const genesis = await response.text();
        set({ genesisString: genesis });
    },
    nodePopJsons: ["", "", "", "", "", "", "", "", "", ""],
    setNodePopJsons: (nodePopJsons: string[]) => set(() => ({ nodePopJsons: nodePopJsons })),
})

export const useWizardStore = window.location.origin.startsWith("http://localhost:")
    ? create<WizardState>()(
        persist(
            wizardStoreFunc,
            {
                name: 'wizard-storage',
                storage: createJSONStorage(() => localStorage),
            }
        )
    )
    : create<WizardState>()(wizardStoreFunc);
