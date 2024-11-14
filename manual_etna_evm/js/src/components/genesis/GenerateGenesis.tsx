import Button from "../../lib/Button";
import { useWalletStore } from "../balance/walletStore";
import { useGenesisStore } from "./genesisStore";

export default function GenerateGenesis({ canRegenerate }: { canRegenerate: boolean }) {
    const generateGenesis = useGenesisStore(state => state.generateGenesis);
    const clearGenesis = useGenesisStore(state => state.clearGenesis);
    const cAddress = useWalletStore(state => state.cAddress);
    const genesis = useGenesisStore(state => state.genesis);

    if (genesis) {
        const blob = new Blob([genesis], { type: 'text/plain' });
        const downloadUrl = URL.createObjectURL(blob);

        return (<>
            <div className="mb-4">✅ Genesis generated!</div>
            <div className="flex gap-4">
                <a href={downloadUrl} className="text-blue-600 hover:text-blue-800 underline" download="genesis.json">Download Genesis</a>

                {canRegenerate && <a href="#" className="text-gray-500 hover:text-blue-600 underline" onClick={() => clearGenesis()}>Start over</a>}
            </div>
        </>
        );
    }

    return <>
        <Button onClick={() => generateGenesis({ userAddress: cAddress })}>Generate genesis</Button>
    </>
}
