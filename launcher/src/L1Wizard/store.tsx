export const stepList = {
    "genesis": {
        title: "Create genesis",
        description: "Allocations and precompiles",
        icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
            <path stroke="currentColor" strokeLinecap="round" strokeWidth="2" d="M20 6H10m0 0a2 2 0 1 0-4 0m4 0a2 2 0 1 1-4 0m0 0H4m16 6h-2m0 0a2 2 0 1 0-4 0m4 0a2 2 0 1 1-4 0m0 0H4m16 6H10m0 0a2 2 0 1 0-4 0m4 0a2 2 0 1 1-4 0m0 0H4" />
        </svg>

    },
    "generate-keys": {
        title: "Generate keys",
        description: "Generate keys for your nodes",
        icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
            <path fill="currentColor" d="M6.94318 11h-.85227l.96023-2.90909h1.07954L9.09091 11h-.85227l-.63637-2.10795h-.02272L6.94318 11Zm-.15909-1.14773h1.60227v.59093H6.78409v-.59093ZM9.37109 11V8.09091h1.25571c.2159 0 .4048.04261.5667.12784.162.08523.2879.20502.3779.35937.0899.15436.1349.33476.1349.5412 0 .20833-.0464.38873-.1392.54119-.0918.15246-.2211.26989-.3878.35229-.1657.0824-.3593.1236-.5809.1236h-.75003v-.61367h.59093c.0928 0 .1719-.0161.2372-.0483.0663-.03314.1169-.08002.152-.14062.036-.06061.054-.13211.054-.21449 0-.08334-.018-.15436-.054-.21307-.0351-.05966-.0857-.10511-.152-.13636-.0653-.0322-.1444-.0483-.2372-.0483h-.2784V11h-.78981Zm3.41481-2.90909V11h-.7898V8.09091h.7898Z" />
            <path stroke="currentColor" stroke-linejoin="round" stroke-width="2" d="M8.31818 2c-.55228 0-1 .44772-1 1v.72878c-.06079.0236-.12113.04809-.18098.07346l-.55228-.53789c-.38828-.37817-1.00715-.37817-1.39543 0L3.30923 5.09564c-.19327.18824-.30229.44659-.30229.71638 0 .26979.10902.52813.30229.71637l.52844.51468c-.01982.04526-.03911.0908-.05785.13662H3c-.55228 0-1 .44771-1 1v2.58981c0 .5523.44772 1 1 1h.77982c.01873.0458.03802.0914.05783.1366l-.52847.5147c-.19327.1883-.30228.4466-.30228.7164 0 .2698.10901.5281.30228.7164l1.88026 1.8313c.38828.3781 1.00715.3781 1.39544 0l.55228-.5379c.05987.0253.12021.0498.18102.0734v.7288c0 .5523.44772 1 1 1h2.65912c.5523 0 1-.4477 1-1v-.7288c.1316-.0511.2612-.1064.3883-.1657l.5435.2614v.4339c0 .5523.4477 1 1 1H14v.0625c0 .5523.4477 1 1 1h.0909v.0625c0 .5523.4477 1 1 1h.6844l.4952.4823c1.1648 1.1345 3.0214 1.1345 4.1863 0l.2409-.2347c.1961-.191.3053-.454.3022-.7277-.0031-.2737-.1183-.5342-.3187-.7207l-6.2162-5.7847c.0173-.0398.0342-.0798.0506-.12h.7799c.5522 0 1-.4477 1-1V8.17969c0-.55229-.4478-1-1-1h-.7799c-.0187-.04583-.038-.09139-.0578-.13666l.5284-.51464c.1933-.18824.3023-.44659.3023-.71638 0-.26979-.109-.52813-.3023-.71637l-1.8803-1.8313c-.3883-.37816-1.0071-.37816-1.3954 0l-.5523.53788c-.0598-.02536-.1201-.04985-.1809-.07344V3c0-.55228-.4477-1-1-1H8.31818Z" />
        </svg>

    },
    "paste-keys": {
        title: "Paste keys",
        description: "Paste your keys into the nodes",
        icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20H5a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h2.429M7 8h3M8 8V4h4v2m4 0V5h-4m3 4v3a1 1 0 0 1-1 1h-3m9-3v9a1 1 0 0 1-1 1h-7a1 1 0 0 1-1-1v-6.397a1 1 0 0 1 .27-.683l2.434-2.603a1 1 0 0 1 .73-.317H19a1 1 0 0 1 1 1Z" />
        </svg>

    },
    "create-l1": {
        title: "Create an L1",
        description: "Create a subnet, chain and convert to L1",
        icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h3a3 3 0 0 0 0-6h-.025a5.56 5.56 0 0 0 .025-.5A5.5 5.5 0 0 0 7.207 9.021C7.137 9.017 7.071 9 7 9a4 4 0 1 0 0 8h2.167M12 19v-9m0 0-2 2m2-2 2 2" />
        </svg>
    },
    "launch-validators": {
        title: "Launch validators",
        description: "Launch validators on your infra",
        icon: <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12a1 1 0 0 0-1 1v4a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-4a1 1 0 0 0-1-1M5 12h14M5 12a1 1 0 0 1-1-1V7a1 1 0 0 1 1-1h14a1 1 0 0 1 1 1v4a1 1 0 0 1-1 1m-2 3h.01M14 15h.01M17 9h.01M14 9h.01" />
        </svg>
    },
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
    evmChainId: number;
    setEvmChainId: (chainId: number) => void;
    genesisString: string;
    regenerateGenesis: () => Promise<void>;
    nodePopJsons: string[];
    setNodePopJsons: (nodePopJsons: string[]) => void;
    l1Name: string;
    setL1Name: (l1Name: string) => void;
    chainId: string;
    setChainId: (chainId: string) => void;
    subnetId: string;
    setSubnetId: (subnetId: string) => void;
    conversionId: string;
    setConversionId: (conversionId: string) => void;
}


import generateName from 'boring-name-generator'

const wizardStoreFunc: StateCreator<WizardState> = (set, get) => ({
    ownerEthAddress: "",
    setOwnerEthAddress: (address: string) => set(() => ({
        ownerEthAddress: address,
        genesisString: "",
    })),

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

    nodesCount: 1,
    setNodesCount: (count: number) => set(() => ({ nodesCount: count })),

    genesisString: "",
    regenerateGenesis: async () => {
        const params = new URLSearchParams({
            ownerEthAddressString: get().ownerEthAddress,
            evmChainId: get().evmChainId.toString()
        });

        const response = await fetch(`/api/genesis?${params}`);
        if (!response.ok) {
            let errorMessage = response.statusText;
            try {
                errorMessage = await response.text()
            } catch (e) {
                console.error(e)
            }
            throw new Error('Failed to generate genesis: ' + errorMessage);
        }
        const genesis = await response.text();
        set({ genesisString: genesis });
    },

    nodePopJsons: ["", "", "", "", "", "", "", "", "", ""],
    setNodePopJsons: (nodePopJsons: string[]) => set(() => ({ nodePopJsons: nodePopJsons })),

    l1Name: (generateName().spaced.split('-').join(' ').split(' ').map((word: string) => word.charAt(0).toUpperCase() + word.slice(1)).join(' ') + " L1"),
    setL1Name: (l1Name: string) => set(() => ({ l1Name: l1Name })),

    chainId: "",
    setChainId: (chainId: string) => set(() => ({ chainId: chainId })),

    subnetId: "",
    setSubnetId: (subnetId: string) => set(() => ({ subnetId: subnetId })),

    conversionId: "",
    setConversionId: (conversionId: string) => set(() => ({ conversionId: conversionId })),

    evmChainId: Math.floor(Math.random() * 1000000) + 1,
    setEvmChainId: (chainId: number) => set(() => ({
        evmChainId: chainId,
        genesisString: ""
    })),
})



export const useWizardStore = window.location.origin.startsWith("http://localhost:") || window.location.origin.startsWith("http://tokyo:")
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
