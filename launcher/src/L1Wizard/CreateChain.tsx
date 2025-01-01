import { useState, useEffect } from 'react';
import { useWizardStore } from './store';
import NextPrev from "./ui/NextPrev";
import { createPublicClient, createWalletClient, custom, http, parseEther, formatEther } from 'viem';
import { avalancheFuji } from 'viem/chains';
import { newPrivateKey, getAddresses } from './wallet';

type Status = 'not_started' | 'in_progress' | 'error' | 'success';

interface StepStatus {
    status: Status;
    error?: string;
    data?: any;
}

export default function CreateChain() {
    const { nodesCount, setNodesCount, tempPrivateKeyHex, setTempPrivateKeyHex } = useWizardStore();
    const nodeCounts = [1, 3, 5];
    const [cChainBalance, setCChainBalance] = useState<bigint>(BigInt(0));
    const [transferring, setTransferring] = useState(false);

    // Mock states for each step
    const [chainStatus, setChainStatus] = useState<StepStatus>({ status: 'not_started' });
    const [walletStatus, setWalletStatus] = useState<StepStatus>({ status: 'not_started' });
    const [pChainStatus, setPChainStatus] = useState<StepStatus>({ status: 'not_started' });
    const [subnetStatus, setSubnetStatus] = useState<StepStatus>({ status: 'not_started' });
    const [createChainStatus, setCreateChainStatus] = useState<StepStatus>({ status: 'not_started' });

    useEffect(() => {
        if (!tempPrivateKeyHex) {
            setTempPrivateKeyHex(newPrivateKey());
        }
    }, [tempPrivateKeyHex, setTempPrivateKeyHex]);

    const renderStepIcon = (status: Status) => {
        switch (status) {
            case 'in_progress':
                return (
                    <svg className="animate-spin h-5 w-5 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                );
            case 'success':
                return (
                    <svg className="w-5 h-5 text-green-500" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 12">
                        <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M1 5.917 5.724 10.5 15 1.5" />
                    </svg>
                );
            case 'error':
                return (
                    <svg className="w-5 h-5 text-red-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                );
            default:
                return (
                    <div className="w-5 h-5">
                        <svg className="w-6 h-6 text-gray-800" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                            <path stroke="currentColor" strokeLinecap="round" strokeWidth="2" d="M6 12h.01m6 0h.01m5.99 0h.01" />
                        </svg>
                    </div>
                );
        }
    };

    const addresses = tempPrivateKeyHex ? getAddresses(tempPrivateKeyHex) : null;

    // Check if user is on the right chain
    const checkChain = async () => {
        if (!window.ethereum) {
            setChainStatus({ status: 'error', error: 'No wallet detected' });
            return;
        }

        try {
            const chainId = await window.ethereum.request({ method: 'eth_chainId' });
            if (chainId === '0xa869') { // 43113 in hex
                setChainStatus({ status: 'success' });
            } else {
                setChainStatus({ status: 'error', error: 'Wrong chain' });
            }
        } catch (error: any) {
            setChainStatus({ status: 'error', error: error.message });
        }
    };

    useEffect(() => {
        checkChain();
        // Listen for chain changes
        if (window.ethereum) {
            window.ethereum.on('chainChanged', checkChain);
            return () => {
                window.ethereum.removeListener('chainChanged', checkChain);
            };
        }
    }, []);

    const switchToFuji = async () => {
        if (!window.ethereum) {
            setChainStatus({ status: 'error', error: 'No wallet detected' });
            return;
        }

        try {
            await window.ethereum.request({
                method: 'wallet_switchEthereumChain',
                params: [{ chainId: '0xa869' }], // 43113 in hex
            });
        } catch (error: any) {
            // If the chain hasn't been added to MetaMask, add it
            if (error.code === 4902) {
                try {
                    await window.ethereum.request({
                        method: 'wallet_addEthereumChain',
                        params: [{
                            chainId: '0xa869',
                            chainName: 'Avalanche Fuji Testnet',
                            nativeCurrency: {
                                name: 'Avalanche',
                                symbol: 'AVAX',
                                decimals: 18
                            },
                            rpcUrls: ['https://api.avax-test.network/ext/bc/C/rpc'],
                            blockExplorerUrls: ['https://testnet.snowtrace.io/']
                        }]
                    });
                } catch (addError: any) {
                    setChainStatus({ status: 'error', error: addError.message });
                }
            } else {
                setChainStatus({ status: 'error', error: error.message });
            }
        }
    };

    // Check C-Chain balance
    const checkCChainBalance = async () => {
        if (!addresses?.C) return;

        const client = createPublicClient({
            chain: avalancheFuji,
            transport: http()
        });

        try {
            const balance = await client.getBalance({
                address: addresses.C as `0x${string}`
            });
            setCChainBalance(balance);
        } catch (error) {
            console.error('Failed to get C-Chain balance:', error);
        }
    };

    useEffect(() => {
        checkCChainBalance();
        const interval = setInterval(checkCChainBalance, 5000); // Check every 5 seconds
        return () => clearInterval(interval);
    }, [addresses?.C]);

    const handleTransfer = async () => {
        if (!window.ethereum || !addresses?.C) return;

        const requiredTotal = nodesCount + 0.5;
        const currentBalance = Number(formatEther(cChainBalance));
        const transferAmount = requiredTotal - currentBalance;

        if (transferAmount <= 0) {
            return; // Already have enough funds
        }

        setTransferring(true);
        try {
            const walletClient = createWalletClient({
                chain: avalancheFuji,
                transport: custom(window.ethereum)
            });

            const [account] = await walletClient.requestAddresses();

            await walletClient.sendTransaction({
                account,
                to: addresses.C as `0x${string}`,
                value: parseEther(transferAmount.toFixed(2))
            });

            // Balance will update via the interval
        } catch (error: any) {
            console.error('Transfer failed:', error);
        } finally {
            setTransferring(false);
        }
    };

    const handleCreate = async () => {
        // TODO: Implement the create chain logic
        console.log("Create chain functionality not implemented yet");
    };

    return (
        <div className="max-w-3xl mx-auto">
            <h1 className="text-2xl font-medium mb-6">Create a Chain</h1>

            {/* Node Count Selection */}
            <div className="mb-8">
                <h3 className="mb-4 font-medium">How many nodes do you want to run?</h3>
                <ul className="mb-4 items-center w-full text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-lg sm:flex">
                    {nodeCounts.map((count) => (
                        <li key={count} className="w-full border-b border-gray-200 sm:border-b-0 sm:border-r last:border-r-0">
                            <div className="flex items-center ps-3">
                                <input
                                    id={`nodes-${count}`}
                                    type="radio"
                                    checked={nodesCount === count}
                                    onChange={() => setNodesCount(count)}
                                    name="nodes-count"
                                    className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-0"
                                />
                                <label htmlFor={`nodes-${count}`} className="w-full py-3 ms-2 text-sm font-medium text-gray-900">
                                    {count} {count === 1 ? 'Node' : 'Nodes'}
                                    {count === 1 && <span className="ms-2 bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded-full">Dev</span>}
                                    {count === 3 && <span className="ms-2 bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded-full">Testnet</span>}
                                    {count === 5 && <span className="ms-2 bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full">Mainnet</span>}
                                </label>
                            </div>
                        </li>
                    ))}
                </ul>
            </div>

            {/* Chain Switch Section */}
            <div className="mb-8 p-4 border rounded-lg">
                <h3 className="font-medium mb-4">1. Network Check</h3>
                {chainStatus.status === 'success' ? (
                    <div className="text-green-600">You are on Avalanche Fuji Testnet</div>
                ) : (
                    <button
                        onClick={switchToFuji}
                        className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
                    >
                        Switch to Avalanche Fuji
                    </button>
                )}
            </div>

            {/* Wallet Funding Section */}
            <div className="mb-8 p-4 border rounded-lg">
                <h3 className="font-medium mb-4">2. Fund Temporary Wallet</h3>
                <div className="space-y-4">
                    <div className="bg-gray-50 p-4 rounded">
                        <div className="flex justify-between items-center mb-1">
                            <div className="text-sm text-gray-600">C-Chain Address:</div>
                            <div className="text-sm text-gray-600">Balance: {formatEther(cChainBalance)} AVAX</div>
                        </div>
                        <div className="font-mono text-sm break-all mb-3">{addresses?.C}</div>
                        {(nodesCount + 0.5 - Number(formatEther(cChainBalance))) > 0 && (
                            <button
                                onClick={handleTransfer}
                                disabled={transferring || chainStatus.status !== 'success'}
                                className={`w-full px-4 py-2 rounded text-white ${transferring || chainStatus.status !== 'success'
                                    ? 'bg-gray-400 cursor-not-allowed'
                                    : 'bg-blue-500 hover:bg-blue-600'
                                    }`}
                            >
                                {transferring ? 'Transferring...' : `Transfer ${(nodesCount + 0.5 - Number(formatEther(cChainBalance))).toFixed(2)} AVAX`}
                            </button>
                        )}
                    </div>
                    <div className="bg-gray-50 p-4 rounded">
                        <div className="flex justify-between items-center mb-1">
                            <div className="text-sm text-gray-600">P-Chain Address:</div>
                            <div className="text-sm text-gray-600">Balance: 0 AVAX</div>
                        </div>
                        <div className="font-mono text-sm break-all">{addresses?.P}</div>
                    </div>
                </div>
            </div>

            {/* P Chain Transfer Section */}
            <div className="mb-8 p-4 border rounded-lg">
                <h3 className="font-medium mb-4">3. Transfer & Create</h3>
                <div className="mb-4">
                    <div className="flex flex-col mb-4">
                        <div className="flex items-center gap-3">
                            {renderStepIcon(pChainStatus.status)}
                            <span className="text-gray-700">Transfer funds from C-Chain to P-Chain</span>
                        </div>
                        {pChainStatus.data && (
                            <p className="ml-8 mt-1 text-gray-600">
                                <code><a href={`https://subnets-test.avax.network/c-chain/tx/${pChainStatus.data}`} target="_blank" rel="noopener noreferrer" className="text-blue-500 underline hover:text-blue-700">{pChainStatus.data}</a></code>
                            </p>
                        )}
                    </div>

                    <div className="flex flex-col mb-4">
                        <div className="flex items-center gap-3">
                            {renderStepIcon(subnetStatus.status)}
                            <span className="text-gray-700">Create a Subnet</span>
                        </div>
                        {subnetStatus.data && (
                            <p className="ml-8 mt-1 text-gray-600">
                                <code><a href={`https://subnets-test.avax.network/p-chain/tx/${subnetStatus.data}`} target="_blank" rel="noopener noreferrer" className="text-blue-500 underline hover:text-blue-700">{subnetStatus.data}</a></code>
                            </p>
                        )}
                    </div>

                    <div className="flex flex-col mb-4">
                        <div className="flex items-center gap-3">
                            {renderStepIcon(createChainStatus.status)}
                            <span className="text-gray-700">Create a Chain</span>
                        </div>
                        {createChainStatus.data && (
                            <p className="ml-8 mt-1 text-gray-600">
                                <code><a href={`https://subnets-test.avax.network/p-chain/tx/${createChainStatus.data}`} target="_blank" rel="noopener noreferrer" className="text-blue-500 underline hover:text-blue-700">{createChainStatus.data}</a></code>
                            </p>
                        )}
                    </div>
                </div>

                <button
                    onClick={handleCreate}
                    disabled={pChainStatus.status === 'in_progress' || subnetStatus.status === 'in_progress' || createChainStatus.status === 'in_progress'}
                    className={`px-6 py-2 rounded-md ${pChainStatus.status === 'in_progress' || subnetStatus.status === 'in_progress' || createChainStatus.status === 'in_progress'
                        ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                        : 'bg-blue-500 text-white hover:bg-blue-600'
                        }`}
                >
                    {pChainStatus.status === 'in_progress' || subnetStatus.status === 'in_progress' || createChainStatus.status === 'in_progress'
                        ? 'Processing...'
                        : 'Start Process'}
                </button>
            </div>

            <NextPrev
                nextDisabled={createChainStatus.status !== 'success'}
                currentStepName="create-chain"
            />
        </div>
    );
}
