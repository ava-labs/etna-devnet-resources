import { useWizardStore } from "./store";
import { useState } from "react";
const validateNodePop = (json: string): boolean => {
    try {
        const parsed = JSON.parse(json);
        if (!parsed.result?.nodeID || !parsed.result?.nodePOP?.publicKey || !parsed.result?.nodePOP?.proofOfPossession) {
            return false;
        }

        // Validate nodeID is base58
        const base58Regex = /^NodeID-[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]+$/;
        if (!base58Regex.test(parsed.result.nodeID)) {
            return false;
        }

        // Validate publicKey and proofOfPossession are hex strings
        const hexRegex = /^0x[0-9a-fA-F]+$/;
        if (!hexRegex.test(parsed.result.nodePOP.publicKey) || !hexRegex.test(parsed.result.nodePOP.proofOfPossession)) {
            return false;
        }

        return true;
    } catch {
        return false;
    }
};

export default function CreateL1() {
    const { advanceFrom, nodePopJsons, setNodePopJsons, nodesCount } = useWizardStore();
    const [errors, setErrors] = useState<string[]>(Array(nodesCount).fill(''));

    const handleNodePopChange = (index: number, value: string) => {
        const newJsons = [...nodePopJsons];
        newJsons[index] = value;
        setNodePopJsons(newJsons);

        const newErrors = [...errors];
        if (value) {
            if (!validateNodePop(value)) {
                newErrors[index] = 'Invalid JSON format. Must contain nodeID and nodePOP fields';
            } else {
                newErrors[index] = '';
            }
        } else {
            newErrors[index] = '';
        }
        setErrors(newErrors);
    };

    const nextDisabled = nodePopJsons.length !== nodesCount ||
        nodePopJsons.some((json, i) => !json || errors[i]);

    return <>
        <h1 className="text-2xl font-medium mb-6">Create an L1</h1>

        <h3 className="mb-4 font-medium">Paste the node credentials for each node:</h3>
        {Array.from({ length: nodesCount }).map((_, index) => (
            <div key={index} className="mb-4">
                <label className="block mb-2">
                    Node {index + 1} Credentials:
                </label>
                <div className="relative">
                    <textarea
                        className={`w-full p-2 border rounded-md font-mono ${nodePopJsons[index] && !errors[index]
                            ? 'bg-green-50 border-green-200'
                            : 'bg-gray-100'
                            }`}
                        rows={8}
                        value={nodePopJsons[index] || ''}
                        onChange={(e) => handleNodePopChange(index, e.target.value)}
                        placeholder={`{"jsonrpc":"2.0","result":{"nodeID":"NodeID-....","nodePOP":{"publicKey":"0x...","proofOfPossession":"0x..."}},"id":1}`}
                    />
                    {nodePopJsons[index] && !errors[index] && (
                        <div className="absolute right-2 top-2 text-green-500">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                            </svg>
                        </div>
                    )}
                </div>
                {errors[index] && (
                    <p className="text-red-500 text-sm mt-1">{errors[index]}</p>
                )}
            </div>
        ))}

        <div className="flex justify-between">
            <button
                onClick={() => advanceFrom('create-l1', 'down')}
                className="px-4 py-2 rounded-md bg-gray-100 hover:bg-gray-200"
            >
                Previous
            </button>
            <button
                onClick={() => advanceFrom('create-l1')}
                disabled={nextDisabled}
                className={`px-4 py-2 rounded-md ${nextDisabled
                    ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                    : 'bg-blue-500 text-white hover:bg-blue-600'
                    }`}
            >
                Continue
            </button>
        </div>
    </>
}
