import { useState } from 'react';
import { getWalletAddress } from './wallet';
import { useWizardStore } from './store';

function isValidL1Name(name: string): boolean {
    return name.split('').every(char => {
        const code = char.charCodeAt(0);
        return code <= 127 && // MaxASCII check
            (char.match(/[a-zA-Z0-9 ]/) !== null); // only letters, numbers, spaces
    });
}

export default function Genesis() {
    const { ownerEthAddress, setOwnerEthAddress, advanceFrom, evmChainId, setEvmChainId, genesisString, regenerateGenesis, l1Name, setL1Name } = useWizardStore();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');
    const [isRegenerating, setIsRegenerating] = useState(false);

    const handleConnectWallet = async () => {
        setIsLoading(true);
        setError('');
        try {
            const walletAddress = await getWalletAddress();
            setOwnerEthAddress(walletAddress);
        } catch (err: any) {
            setError(err.message || 'Failed to connect wallet');
        } finally {
            setIsLoading(false);
        }
    };

    const handleContinue = () => {
        advanceFrom('genesis')
    }

    const handleInputChange = (field: 'evmChainId' | 'ownerEthAddress' | 'l1Name', value: any) => {
        if (field === 'evmChainId') {
            setEvmChainId(parseInt(value));
        } else if (field === 'l1Name') {
            setL1Name(value);
        } else {
            setOwnerEthAddress(value);
        }
    };

    const handleGenerateGenesis = async () => {
        setIsRegenerating(true);
        setError('');
        try {
            await regenerateGenesis();
        } catch (err: any) {
            setError(err.message || 'Failed to regenerate genesis');
        } finally {
            setIsRegenerating(false);
        }
    };

    return (
        <div className="">
            <h1 className="text-2xl font-medium mb-6">Genesis Settings</h1>

            {error && (
                <div className="mb-4 p-3 text-sm text-red-500 bg-red-50 rounded-md">
                    {error}
                </div>
            )}

            <div className="mb-6">
                <input
                    type="text"
                    value={l1Name}
                    onChange={(e) => handleInputChange('l1Name', e.target.value)}
                    placeholder="L1 Name"
                    className={`w-full p-2 border rounded-md ${l1Name && !isValidL1Name(l1Name)
                        ? 'border-red-500 text-red-500'
                        : 'border-gray-200'
                        }`}
                />
                <p className={`mt-2 text-sm ${l1Name && !isValidL1Name(l1Name)
                    ? 'text-red-500'
                    : 'text-gray-500'
                    }`}>
                    {l1Name && !isValidL1Name(l1Name)
                        ? 'Only letters, numbers, and spaces are allowed'
                        : 'The name of your L1 blockchain.'}
                </p>
            </div>

            <div className="mb-6">
                <input
                    type="number"
                    value={evmChainId}
                    onChange={(e) => handleInputChange('evmChainId', e.target.value)}
                    onBlur={() => handleInputChange('evmChainId', evmChainId)}
                    placeholder="Chain ID"
                    className="w-full p-2 border border-gray-200 rounded-md"
                />
                <p className="mt-2 text-sm text-gray-500">
                    Unique identifier for your blockchain network.  Check if it's unique <a href={`https://chainlist.org/?search=${evmChainId}`} target="_blank" rel="noopener noreferrer" className="text-blue-500 underline hover:text-blue-600">on chainlist.org</a>.
                </p>
            </div>

            <div className="mb-6">
                <div className="flex gap-3 items-start">
                    <div className="flex-grow">
                        <input
                            type="text"
                            value={ownerEthAddress}
                            onChange={(e) => handleInputChange('ownerEthAddress', e.target.value)}
                            onBlur={() => handleInputChange('ownerEthAddress', ownerEthAddress)}
                            placeholder="Wallet Address"
                            className="w-full p-2 border border-gray-200 rounded-md"
                        />
                    </div>
                    {window.ethereum && <button
                        onClick={handleConnectWallet}
                        disabled={isLoading}
                        className={`px-4 py-2 rounded-md ${isLoading
                            ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                            : 'bg-gray-100 hover:bg-gray-200'
                            }`}
                    >
                        {isLoading ? 'Loading...' : 'Fill from Wallet'}
                    </button>}
                </div>
                <p className="mt-2 text-sm text-gray-500">
                    This address will receive all tokens and control in case of Proof of Authority chain.
                </p>
            </div>

            <div className="mb-6 flex justify-between">
                <button
                    onClick={handleGenerateGenesis}
                    disabled={!evmChainId || !ownerEthAddress || isRegenerating}
                    className={`px-4 py-2 rounded-md ${!evmChainId || !ownerEthAddress || isRegenerating
                        ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                        : 'bg-blue-500 text-white hover:bg-blue-600'
                        }`}
                >
                    {isRegenerating ? 'Generating...' : 'Generate Genesis'}
                </button>
                <button
                    onClick={handleContinue}
                    disabled={!isValidEthereumAddress(ownerEthAddress) || !genesisString || !isValidL1Name(l1Name)}
                    className={`px-4 py-2 rounded-md ${!isValidEthereumAddress(ownerEthAddress) || !genesisString || !isValidL1Name(l1Name)
                        ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                        : 'bg-blue-500 text-white hover:bg-blue-600'
                        }`}
                >
                    Next
                </button>
            </div>

            {genesisString && <div className="mb-6">
                <label className="block text-sm text-gray-500 mb-2">
                    Genesis JSON:
                </label>
                <div className="bg-gray-50 overflow-x-auto text-sm font-mono border border-gray-200 rounded-md">
                    <pre className="w-full p-3 break-words overflow-auto">
                        {genesisString}
                    </pre>
                </div>
            </div>}


        </div>
    );
}

function isValidEthereumAddress(address: string) {
    return /^0x[a-fA-F0-9]{40}$/.test(address);
}
